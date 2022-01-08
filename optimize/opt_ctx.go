package optimize

import (
	"github.com/c-bata/goptuna"
	"github.com/pkg/errors"
	"trade-optimizer/dataloader"
	"trade-optimizer/strategies"
	"trade-optimizer/structs"
)

type Loss func(trades structs.Trades, startTs uint64, endTs uint64) float64

var errInvalidStrategyClone = errors.New("Cannot use type as strategy")

type OptimizationResult struct {
	TrainTrades  structs.Trades
	TrainStart   uint64
	TrainEnd     uint64
	ObjectiveVal float64
	BestParams   map[string]interface{}
	StudyName    string
}

type OptimizeContext struct {
	Candles     []dataloader.Candle
	loss           Loss
	TradeSize      float64
	Fee            float64
	Pairname       string //POTENTIALLY STORE MULTIPLE PAIRS HERE
	nItersNoChange int
	Strategy       strategies.Strategy
}

func NewOptimizeContext(candles []dataloader.Candle, loss Loss, tradeSize float64, fee float64, nItersNoChange int, pairname string, strategy strategies.Strategy) *OptimizeContext {
	return &OptimizeContext{
		Candles:        candles,
		loss:           loss,
		TradeSize:      tradeSize,
		Fee:            fee,
		nItersNoChange: nItersNoChange,
		Pairname:       pairname,
		Strategy:       strategy,
	}
}

func (optCtx *OptimizeContext) Clone() (*OptimizeContext, error) {
	clone := *optCtx
	if strategyClone, ok := optCtx.Strategy.Clone().(strategies.Strategy); ok {
		clone.Strategy = strategyClone
	} else {
		return nil, errors.Wrapf(errInvalidStrategyClone, "cannot use %s as Strategy", strategyClone)
	}
	return &clone, nil
}

func (optCtx *OptimizeContext) computeObjective(trial goptuna.Trial) (float64, error) {
	pcfg := strategies.NewParamConfigurator(&trial)
	optCtx.Strategy.SuggestParams(pcfg)
	paramsMap, err := pcfg.OptParams()
	if err != nil {
		return 0.0, errors.Wrap(err, "error when building trial params")
	}
	optCtx.Strategy.SetOptParams(paramsMap)
	trades := Backtest(optCtx)

	startTs := optCtx.Candles[0].Time
	endTs := optCtx.Candles[len(optCtx.Candles)-1].Time
	trainLoss := optCtx.loss(trades, startTs, endTs)
	return trainLoss, nil
}
