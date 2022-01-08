package strategies

import (
	"github.com/dmitryikh/leaves"
	"trade-optimizer/dataloader"
)

// MlStategy Заготовка для стратегии, использующей lgbm в качестве предиктора
// Не было использовано в финальной версии, на будущее
type MlStategy struct {
	params OptParams
}

func (m MlStategy) NumPrevRequired() int {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) SuggestParams(pcfg *ParamConfigurator) {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) SetOptParams(params OptParams) {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) OptParams() OptParams {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) CalcTimeseriesFeatures(c map[string]dataloader.Candles) map[string]map[string][]float64 {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) CalcFeatures(market TradingState, timeseriesFeatures map[string]map[string][]float64) {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) TakeAction(timeseriesFeatures map[string]map[string][]float64) []Action {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) Clone() interface{} {
	//TODO implement me
	panic("implement me")
}

func (m MlStategy) DegreesOfFreedom() int {
	//TODO implement me
	panic("implement me")
}

func LoadModel(path string) (*leaves.Ensemble, error) {
	// 1. Read model
	model, err := leaves.LGEnsembleFromFile(path, true)
	if err != nil {
		return nil, err
	}
	return model, nil
}
