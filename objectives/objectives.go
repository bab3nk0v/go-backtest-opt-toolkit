package objectives

import (
	"github.com/montanaflynn/stats"
	"math"
	"sort"
	"trade-optimizer/structs"
	"trade-optimizer/util"
)

func CalcOverfittingScore(t structs.Trades) float64 {
	if len(t) == 0 {
		return 1.00
	}
	var newTrades structs.Trades
	newTrades = append(make(structs.Trades, 0, len(t)), t...)
	sort.Slice(newTrades, func(i, j int) bool {
		return math.Abs(newTrades[i].Profit) > math.Abs(newTrades[j].Profit)
	})
	percentileIndex := int(math.Round(float64(len(newTrades)) * 0.01))
	percentileData := newTrades[0 : percentileIndex+1]
	profitPercent := 0.0
	for _, t := range percentileData {
		profitPercent += math.Abs(t.Profit)
	}
	profitWhole := 0.0
	for _, t := range newTrades {
		profitWhole += math.Abs(t.Profit)

	}
	if profitWhole == 0 {
		return 1.0
	}
	return profitPercent / profitWhole
}

func AvgReturnByAvgDrawdownLoss(trades structs.Trades, startTs uint64, endTs uint64) float64 {
	if len(trades) == 0 {
		return 0.0
	}
	os := CalcOverfittingScore(trades)
	numTradesPenalty := float64(len(trades)) / (float64(len(trades)) + 300)
	var newTrades structs.Trades
	newTrades = append(make(structs.Trades, 0, len(trades)), trades...)
	sort.Slice(newTrades, func(i, j int) bool {
		return math.Abs(newTrades[i].Profit) > math.Abs(newTrades[j].Profit)
	})
	percentileIndex := int(math.Round(float64(len(newTrades)) * 0.01))
	newTrades = newTrades[percentileIndex+1:]
	sort.Slice(newTrades, func(i, j int) bool {
		return newTrades[i].TsExit < newTrades[j].TsExit
	})
	equity := util.CumSumFloat64(newTrades.Profits())
	avgDd := util.AvgDrawdown(equity)
	avgReturn, err := stats.Mean(trades.ProfitsInPercents())
	if err != nil {
		panic(err)
	}
	if avgDd == 0 {
		return 100.0
	}
	ratio := avgReturn / avgDd
	loss := -ratio * (1 - os) * numTradesPenalty
	return loss
}

func SharpeRobust(trades structs.Trades, startTs uint64, endTs uint64) float64 {
	if len(trades) <= 2 {
		return 0.0
	}
	os := CalcOverfittingScore(trades)
	numTradesPenalty := float64(len(trades)) / (float64(len(trades)) + 300)
	var newTrades structs.Trades
	newTrades = append(make(structs.Trades, 0, len(trades)), trades...)
	sort.Slice(newTrades, func(i, j int) bool {
		return math.Abs(newTrades[i].Profit) > math.Abs(newTrades[j].Profit)
	})
	percentileIndex := int(math.Round(float64(len(newTrades)) * 0.01))
	newTrades = newTrades[percentileIndex+1:]
	ratio, err := stats.Mean(newTrades.ProfitsInPercents())
	if err != nil {
		panic(err)
	}
	dnm, err := stats.StandardDeviation(newTrades.ProfitsInPercents())
	if err != nil {
		panic(err)
	}
	ratio /= dnm
	loss := -ratio * (1 - os) * numTradesPenalty
	return loss
}

func ProfitFactorRobust(trades structs.Trades, startTs uint64, endTs uint64) float64 {

	os := CalcOverfittingScore(trades)

	if len(trades) == 0 {
		return 0.0
	}

	var newTrades structs.Trades
	newTrades = append(make(structs.Trades, 0, len(trades)), trades...)
	sort.Slice(newTrades, func(i, j int) bool {
		return math.Abs(newTrades[i].Profit) > math.Abs(newTrades[j].Profit)
	})
	percentileIndex := int(math.Round(float64(len(newTrades)) * 0.01))
	newTrades = newTrades[percentileIndex+1:]

	n := float64(len(newTrades))
	numTradesPenalty := n / (n + 300)
	pf := ProfitFactor(newTrades, startTs, endTs)
	pf *= numTradesPenalty
	pf *= 1 - os
	return pf
}

func ProfitFactor(trades structs.Trades, startTs uint64, endTs uint64) float64 {
	profit := 0.0
	loss := 0.0
	for _, t := range trades {
		if t.Profit > 0 {
			profit += t.Profit / t.TradeSize
		} else {
			loss += t.Profit / t.TradeSize
		}
	}
	if loss == 0 {
		return 0.0
	}
	return profit / loss
}

func ProfitSum(trades []*structs.Trade) float64 {
	s := 0.0
	for _, t := range trades {
		s += t.Profit
	}
	return s
}

// ExpectedProfit winrate * median of profits + (1 - winrate) * median of losses
func ExpectedProfit(trades structs.Trades, startTs uint64, endTs uint64, bhTrade *structs.Trade) float64 {
	n := float64(len(trades))
	if n == 0 {
		return 0.0
	}
	cntWins := 0
	var losses []float64
	var profits []float64
	for _, t := range trades {
		if t.Profit > 0 {
			profits = append(profits, t.Profit)
			cntWins++
		} else {
			losses = append(losses, t.Profit)
		}
	}
	if cntWins == 0 {
		return 0.0
	}
	p := float64(cntWins) / n
	avgWin, err := stats.Median(profits)
	if err != nil {
		avgWin = 0.0
	}
	avgLoss, err := stats.Median(losses)
	if err != nil {
		avgLoss = 0.0
	}
	profitExp := p * avgWin
	lossExp := (1 - p) * avgLoss
	l := profitExp + lossExp
	return -l
}
