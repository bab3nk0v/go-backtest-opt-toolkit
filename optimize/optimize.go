package optimize

import (
	"context"
	"fmt"
	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb.v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"math"
	"time"
	"trade-optimizer/datascanner"
)

func createStudy(studyName string, db *gorm.DB) error {
	err := rdb.RunAutoMigrate(db)
	if err != nil {
		return err
	}
	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy(studyName)
	if err != nil {
		return err
	}
	studyName, err = storage.GetStudyNameFromID(studyID)
	if err != nil {
		return err
	}
	err = storage.SetStudyDirection(studyID, goptuna.StudyDirectionMinimize)
	if err != nil {
		return err
	}
	return nil
}

func earlyStoppingWatcher(trialchan chan goptuna.FrozenTrial, optCtx *OptimizeContext) {
	curIterNoChange := 0
	curBest := math.Inf(1)
	for t := range trialchan {
		trialRuntime := t.DatetimeComplete.Sub(t.DatetimeStart)
		if t.Value < curBest {
			println(fmt.Sprintf("%dms, Found new best value %f on trial %d", trialRuntime.Milliseconds(), t.Value, t.ID))
			curIterNoChange = 0
			curBest = t.Value
		} else {
			curIterNoChange += 1
			println(fmt.Sprintf("%dms, No improvement on trial %d. The best is still %f. %d iters left", trialRuntime.Milliseconds(), t.ID, curBest, optCtx.nItersNoChange-curIterNoChange))
		}
		if curIterNoChange >= optCtx.nItersNoChange {
			println("Optimizing study %d is done, stopping workers", t.StudyID)
			close(trialchan)
			break
		}
	}
}

func Fit(baseRunName string, optCtx *OptimizeContext, db *gorm.DB, nThreads int) (*OptimizationResult, error) {
	if ok, err := datascanner.AnalyzeCandles(optCtx.Candles); !ok {
		return nil, errors.Wrap(err, "cannot run optimizer because of datascanner error")
	}
	storage := rdb.NewStorage(db)
	runName := baseRunName + "&" + "study" + "-" + fmt.Sprintf("%d", time.Now().Unix())
	err := createStudy(runName, db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create study")
	}
	trialchan := make(chan goptuna.FrozenTrial, 256)
	study, err := goptuna.LoadStudy(
		runName,
		goptuna.StudyOptionStorage(storage),
		goptuna.StudyOptionSetTrialNotifyChannel(trialchan),
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load the created study")
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	study.WithContext(ctx)
	for i := 0; i < nThreads; i++ {
		cloneCtx, err := optCtx.Clone()
		if err != nil {
			panic(errors.Wrap(err, "cannot clone context"))
		}
		w := newWorker(i, cloneCtx, study)
		eg.Go(w.run)
	}
	go earlyStoppingWatcher(trialchan, optCtx)
	err = eg.Wait()
	if err != nil {
		return nil, err
	}
	bestValue, err := study.GetBestValue()
	if err != nil {
		return nil, err
	}
	params, err := study.GetBestParams()
	if err != nil {
		return nil, err
	}
	optCtx.Strategy.SetOptParams(params)
	trainTrades := Backtest(optCtx)
	return &OptimizationResult{
		TrainTrades:  trainTrades,
		ObjectiveVal: bestValue,
		BestParams:   params,
		StudyName:    runName,
	}, nil
}
