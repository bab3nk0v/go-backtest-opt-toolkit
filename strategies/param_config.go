package strategies

import "github.com/c-bata/goptuna"

type ParamConfigurator struct {
	prms  OptParams
	trial *goptuna.Trial
	err   error
}

func NewParamConfigurator(trial *goptuna.Trial) *ParamConfigurator {
	p := ParamConfigurator{}
	p.trial = trial
	p.prms = make(OptParams)
	return &p
}

func (pcfg *ParamConfigurator) IntParam(name string, low int, high int) *ParamConfigurator {
	if pcfg.err != nil {
		return pcfg
	}
	var p int
	if pcfg.trial != nil {
		pSuggested, err := pcfg.trial.SuggestInt(name, low, high)
		if err != nil {
			pcfg.prms = nil
			pcfg.err = err
		}
		p = pSuggested
	} else {
		p = 42
	}
	pcfg.prms[name] = p
	return pcfg
}

func (pcfg *ParamConfigurator) FloatParam(name string, low float64, high float64) *ParamConfigurator {
	if pcfg.err != nil {
		return pcfg
	}
	var p float64
	if pcfg.trial != nil {
		pSuggested, err := pcfg.trial.SuggestFloat(name, low, high)
		if err != nil {
			pcfg.prms = nil
			pcfg.err = err
		}
		p = pSuggested
	} else {
		p = 42
	}
	pcfg.prms[name] = p
	return pcfg
}

func (pcfg *ParamConfigurator) OptParams() (OptParams, error) {
	return pcfg.prms, pcfg.err
}
