package strategies

import (
	"trade-optimizer/dataloader"
	"trade-optimizer/structs"
)

type OptParams map[string]interface{}

type Strategy interface {
	NumPrevRequired() int
	SuggestParams(pcfg *ParamConfigurator)
	SetOptParams(params OptParams)
	OptParams() OptParams
	CalcTimeseriesFeatures(c map[string]dataloader.Candles) map[string]map[string][]float64
	CalcFeatures(market TradingState, timeseriesFeatures map[string]map[string][]float64)
	TakeAction(timeseriesFeatures map[string]map[string][]float64) []Action
	Clone() interface{}
	DegreesOfFreedom() int
}

type Action struct {
	Pair       string
	ActionType int
	Confidence float64
	Stoploss   float64
}

const (
	ActionBuy  = iota
	ActionSell = iota
)

type TradingState struct {
	TradablePairs []string
	OpenTrades    map[string]*structs.Trade
	TsNow         uint64
}
