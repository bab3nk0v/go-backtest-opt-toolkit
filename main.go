package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"strconv"
	"time"
	"trade-optimizer/dataloader"
	"trade-optimizer/metrics"
	"trade-optimizer/objectives"
	"trade-optimizer/optimize"
	"trade-optimizer/strategies"
	"trade-optimizer/util"
)

type WalkForwardFold struct {
	Params      strategies.OptParams
	InSample    metrics.TradeMetrics
	OutOfSample metrics.TradeMetrics
	StudyName   string
}

type WalkForwardResult struct {
	Windows map[int][]WalkForwardFold
}

func main() {
	tradeSize := flag.Float64("tradeSize", 100.0, "trade size")
	fee := flag.Float64("fee", 0.0, "fee")
	nItersNoChange := flag.Int("nItersNoChange", 3000, "number of iterations without improvement to stop")
	timeMillis := flag.Bool("timeInMillseconds", true, "true=time in milliseconds false=time in seconds")
	pairname := flag.String("pairname", "XMR_USD", "pair name")
	runName := "optim-" + *pairname + "-" + fmt.Sprintf("%d", time.Now().Unix())
	path := flag.String("dataPath", "datasets/XMR_1m.csv", "path to OHLCV dataset")
	resultFolder := flag.String("resultFolder", "opt_results", "where to dump optimization results")

	dbLogin := flag.String("dbLogin", "goptuna", "db username")
	dbPass := flag.String("dbPass", "password", "db pass")
	dbHost := flag.String("dbHost", "localhost", "db host")
	dbPort := flag.Int("dbPort", 3306, "db port")
	dbName := flag.String("dbName", "goptuna", "database name")
	nThreads := flag.Int("parallelism", 1, "number of parallel optimisation processes")

	flag.Parse()

	candles, err := dataloader.ReadOHLCV(*path, *timeMillis)
	if err != nil {
		panic(err)
	}

	strat := strategies.NewSmaStrategy()

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", *dbLogin, *dbPass, *dbHost, *dbPort, *dbName)

	//db block
	db, _ := gorm.Open(mysql.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDb.Close()
	//db block

	stepRes := WalkForwardResult{}
	stepRes.Windows = make(map[int][]WalkForwardFold)

	for _, nPhases := range []int{1, 3, 10} {
		println("Optimizing " + strconv.Itoa(nPhases))
		walkForwardData := util.WalkForwardSplit(candles, strat.DegreesOfFreedom(), 1, nPhases)

		var res []WalkForwardFold

		for i, wfd := range walkForwardData {
			println(fmt.Sprintf("Optimizing fold %d", i+1))
			inSample, outOfSample := wfd[0], wfd[1]
			optCtx := optimize.NewOptimizeContext(
				inSample,
				objectives.SharpeRobust,
				*tradeSize,
				*fee,
				*nItersNoChange,
				*pairname,
				strat,
			)
			or, err := optimize.Fit(runName, optCtx, db, *nThreads)
			if err != nil {
				panic(err)
			}
			testCtx, err := optCtx.Clone()
			if err != nil {
				panic(err)
			}
			testCtx.Candles = outOfSample
			testTrades := optimize.Backtest(testCtx)
			res = append(res, WalkForwardFold{
				Params:      or.BestParams,
				InSample:    metrics.CalcMetrics(or.TrainTrades, inSample, optCtx.TradeSize, optCtx.Fee),
				OutOfSample: metrics.CalcMetrics(testTrades, outOfSample, optCtx.TradeSize, optCtx.Fee),
				StudyName:   or.StudyName,
			})
		}
		stepRes.Windows[nPhases] = res
	}

	dataMarshalled, err := json.Marshal(stepRes)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*resultFolder+"/"+runName+".json", dataMarshalled, 0644)
	if err != nil {
		panic(err)
	}
	println(runName + " finished")
}
