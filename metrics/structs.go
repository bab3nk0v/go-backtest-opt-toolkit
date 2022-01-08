package metrics

import "trade-optimizer/structs"

type DayStructure struct {
	Day    int64
	Profit float64
}

type TradeMetrics struct {
	Trades                 structs.Trades
	TsStart                uint64
	TsEnd                  uint64
	ProfitFactor           float64
	ProfitFactorRobust     float64
	AvgReturnByAvgDrawdown float64
	SharpeRobust           float64
	TotalProfit            float64
	MaxDrawdown            float64
	BuyAndHoldProfit       float64
	AvgDrawdown            float64
	NumTrades              int
	BuyAndHoldDaily        []DayStructure
	OverfittingScore       float64
}
