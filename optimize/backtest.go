package optimize

import (
	"math"
	"sort"
	"trade-optimizer/dataloader"
	"trade-optimizer/strategies"
	"trade-optimizer/structs"
	"trade-optimizer/util"
)

func CreateSliceOfValues(size int, val float64) []float64 {
	r := make([]float64, size)
	for i, _ := range r {
		r[i] = val
	}
	return r
}

func PrePad(tsFeat map[string][]float64, numCandles int) map[string][]float64 {
	copyFeat := make(map[string][]float64)
	for k, ts := range tsFeat {
		copyFeat[k] = CreateSliceOfValues(numCandles, math.NaN())
		j := numCandles - 1
		for i := len(ts) - 1; i >= 0; i-- {
			copyFeat[k][j] = ts[i]
			j--
		}
	}
	return copyFeat
}

func FirstFullFilled(tsFeat map[string][]float64) int {
	r := 0
	for _, ts := range tsFeat {
		for i, x := range ts {
			if !math.IsNaN(x) {
				r = util.Max(r, i)
				break
			}
		}
	}
	return r
}

func TakeBatch(data map[string]map[string][]float64, i, j int) map[string]map[string][]float64 {
	r := make(map[string]map[string][]float64)
	for pairname, featuredict := range data {
		r[pairname] = make(map[string][]float64)
		for feature, ts := range featuredict {
			r[pairname][feature] = ts[i:j]
		}
	}
	return r
}

func Backtest(optCtx *OptimizeContext) []*structs.Trade {
	exc := NewExchange()
	pairCandleMap := make(map[string]dataloader.Candles)
	pairCandleMap[optCtx.Pairname] = optCtx.Candles
	tsFeat := optCtx.Strategy.CalcTimeseriesFeatures(pairCandleMap)
	for pairname, featdict := range tsFeat {
		tsFeat[pairname] = PrePad(featdict, len(optCtx.Candles))
	}
	npr := optCtx.Strategy.NumPrevRequired()
	firstIndex := FirstFullFilled(tsFeat[optCtx.Pairname])
	wc := firstIndex + (npr - 1) //iterate not over candles but over ts features
	for i := wc; i < len(optCtx.Candles); i++ {
		candle := optCtx.Candles[i]
		if candle.Volume == 0 {
			continue
		}
		leftBatchIndex := i - npr + 1
		rightBatchIndex := i + 1
		historyCandles := optCtx.Candles[leftBatchIndex:rightBatchIndex]
		tsFeatBatch := TakeBatch(tsFeat, leftBatchIndex, rightBatchIndex)
		exc.UpdateState(optCtx.Pairname, &candle)
		exc.FireStoplosses(&candle)
		pairMap := make(map[string]dataloader.Candles)
		pairMap[optCtx.Pairname] = historyCandles
		openTradesMap := make(map[string]*structs.Trade)
		if ot, ok := exc.openTrades[optCtx.Pairname]; ok {
			openTradesMap[optCtx.Pairname] = ot
		}
		market := strategies.TradingState{
			TradablePairs: []string{optCtx.Pairname},
			OpenTrades:    openTradesMap,
			TsNow:         candle.Time,
		}
		optCtx.Strategy.CalcFeatures(market, tsFeat)
		actions := optCtx.Strategy.TakeAction(tsFeatBatch)
		sort.Slice(actions, func(i, j int) bool {
			return actions[i].Confidence > actions[j].Confidence
		})
		for _, a := range actions {
			if a.ActionType == strategies.ActionBuy {
				if val, ok := exc.coolDown[a.Pair]; !ok || (ok && val <= 0) {
					exc.Buy(a.Pair, optCtx.Fee, a.Stoploss, &optCtx.Candles[i], optCtx.TradeSize, ReasonStrategy)
				}
			}
			if a.ActionType == strategies.ActionSell {
				if val, ok := exc.coolDown[a.Pair]; !ok || (ok && val <= 0) {
					exc.Sell(a.Pair, optCtx.Fee, a.Stoploss, &optCtx.Candles[i], optCtx.TradeSize, ReasonStrategy)
				}
			}
		}
		for p, cd := range exc.coolDown {
			exc.coolDown[p] = util.Max(0, cd-1)
		}
	}
	return exc.closedTrades
}
