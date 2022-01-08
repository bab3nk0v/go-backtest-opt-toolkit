package strategies

import (
	"github.com/frankrap/talib"
	"math"
	"trade-optimizer/dataloader"
	"trade-optimizer/util"
)

//const stoploss = -0.05

type SmaStrategyState struct {
	isNowTrading  bool
	currentProfit float64
	tradeLen      uint64
}

type SmaStrategy struct {
	params OptParams
	feat   map[string]SmaStrategyState
}

func NewSmaStrategy() *SmaStrategy {
	return &SmaStrategy{
		feat:   make(map[string]SmaStrategyState),
		params: make(OptParams),
	}
}

func (s *SmaStrategy) SetOptParams(params OptParams) {
	s.params = params
}

func (s *SmaStrategy) OptParams() OptParams {
	return s.params
}

func (s *SmaStrategy) Clone() interface{} {
	clone := NewSmaStrategy()
	for k, v := range s.feat {
		clone.feat[k] = v
	}
	for k, v := range s.params {
		clone.params[k] = v
	}
	return clone
}

func (s *SmaStrategy) NumPrevRequired() int {
	return 2
}

// SuggestParams DSL Для описания параметров и диапазонов параметров стратегии
func (s *SmaStrategy) SuggestParams(pcfg *ParamConfigurator) {
	pcfg.
		IntParam("n1", 8, 21).
		IntParam("n2_offset", 8, 21).
		IntParam("n3_offset", 8, 21).
		FloatParam("w1", 0.0, 50).
		FloatParam("w2", 0.0, 50).
		FloatParam("w3", 0.0, 50).
		FloatParam("w4", 0.0, 50).
		//FloatParam("w_bb", -10, 10).
		FloatParam("stoploss_m", -15, 0).
		FloatParam("tp_base", 0.85, 1).
		FloatParam("tp_scale", 2, 10).
		FloatParam("atr_thr", 0.01, 0.15)
}

func (s *SmaStrategy) CalcTimeseriesFeatures(c map[string]dataloader.Candles) map[string]map[string][]float64 {
	f := make(map[string]map[string][]float64)
	for pairname, candles := range c {
		f[pairname] = make(map[string][]float64)
		cloze := candles.Close()
		f[pairname]["close"] = cloze
		n1 := s.params["n1"].(int)
		sma1 :=  talib.Sma(cloze, int32(n1))
		n2 := n1 + s.params["n2_offset"].(int)
		sma2 := talib.Sma(cloze, int32(n2))
		n3 := n2 + s.params["n3_offset"].(int)
		sma3 := talib.Sma(cloze, int32(n3))
		atr20 := talib.Atr(candles.High(), candles.Low(), cloze, 20)
		f[pairname]["sma1"] = sma1
		f[pairname]["sma2"] = sma2
		f[pairname]["sma3"] = sma3
		f[pairname]["atr20"] = atr20
	}
	return f
}

func (s *SmaStrategy) CalcFeatures(market TradingState, timeseriesFeatures map[string]map[string][]float64) {
	for _, pairname := range market.TradablePairs {
		ftrs := SmaStrategyState{
			isNowTrading: false,
		}
		if val, ok := market.OpenTrades[pairname]; ok {
			ftrs.isNowTrading = ok
			ftrs.currentProfit = val.Profit / val.Amount / val.EnterPrice
			ftrs.tradeLen = market.TsNow - val.TsEnter
		}
		s.feat[pairname] = ftrs
	}
}

func (s *SmaStrategy) TakeAction(timeseriesFeatures map[string]map[string][]float64) []Action {
	var result []Action
	for pairname, pairFeatures := range s.feat {
		if pairFeatures.isNowTrading {
			tp_base := s.params["tp_base"].(float64)
			tp_scale := s.params["tp_scale"].(float64)
			tp := math.Pow(tp_base, float64(pairFeatures.tradeLen)/60) * tp_scale * 0.01
			if pairFeatures.currentProfit > tp && tp > 0 && (float64(pairFeatures.tradeLen)/60) > 30 {
				result = append(result, Action{
					Pair:       pairname,
					ActionType: ActionSell,
					Confidence: 1.0,
					Stoploss:   100500.0,
				})
				continue
			}
		}
		longhaulTrend := 0.0
		longhaulTrendDown := 0.0
		crossAbove := 0.0
		crossBelow := 0.0
		sma1 := timeseriesFeatures[pairname]["sma1"]
		sma2 := timeseriesFeatures[pairname]["sma2"]
		sma2Last := sma2[len(sma2)-1]
		sma3 := timeseriesFeatures[pairname]["sma3"]
		sma3Last := sma3[len(sma3)-1]
		smaDiff := (sma2Last - sma3Last) / sma2Last
		if sma2Last > sma3Last {
			longhaulTrend = 1.0
		} else {
			longhaulTrendDown = 1.0
		}
		if util.CrossedAbove(sma1, sma2) {
			crossAbove = 1.0
		}
		if util.CrossedBelow(sma1, sma2) {
			crossBelow = -1.0
		}
		w1 := s.params["w1"].(float64)
		w2 := s.params["w2"].(float64)
		atrThr := s.params["atr_thr"].(float64)
		w3 := s.params["w3"].(float64)
		w4 := s.params["w4"].(float64)
		cloze := timeseriesFeatures[pairname]["close"]
		atr20 := timeseriesFeatures[pairname]["atr20"]
		atr := atr20[len(atr20)-1]
		atrPercent := atr / cloze[len(cloze)-1]
		if atrPercent <= 0 {
			panic(atrPercent)
		}
		stoplossMul := s.params["stoploss_m"].(float64)
		stoploss := stoplossMul * atrPercent
		score := w1 * crossAbove
		score += w2 * longhaulTrend * smaDiff
		score += w3 * crossBelow
		score += w4 * longhaulTrendDown * smaDiff
		score = 1 / (1 + math.Exp(-score))
		if score > 0.75 && atrPercent < atrThr {
			result = append(result, Action{
				Pair:       pairname,
				ActionType: ActionBuy,
				Confidence: score,
				Stoploss:   stoploss,
			})
		}
		if score < 0.25 {
			result = append(result, Action{
				Pair:       pairname,
				ActionType: ActionSell,
				Confidence: 1 - score,
				Stoploss:   stoploss,
			})
		}
	}
	return result
}

func (s *SmaStrategy) DegreesOfFreedom() int {
	pcfg := NewParamConfigurator(nil)
	s.SuggestParams(pcfg)
	return len(pcfg.prms)
}
