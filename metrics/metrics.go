package metrics

import (
	"time"
	"trade-optimizer/dataloader"
	"trade-optimizer/formulas"
	"trade-optimizer/objectives"
	"trade-optimizer/structs"
	"trade-optimizer/util"
)

func MakeBuyAndHoldTrade(candles []dataloader.Candle, tradeSize float64, fee float64) *structs.Trade {
	firstCandle := candles[0]
	lastCandle := candles[len(candles)-1]
	t := &structs.Trade{
		TradeSize:  tradeSize,
		Fee:        fee,
		Amount:     tradeSize / firstCandle.Close,
		EnterPrice: firstCandle.Close,
		ExitPrice:  lastCandle.Close,
		TsEnter:    firstCandle.Time,
		TsExit:     lastCandle.Time,
		TradeType:  structs.TradeLong,
	}
	t.Profit = formulas.CalcProfit(t, t.ExitPrice, fee)
	return t
}

func CalcDailyBHProfit(candles []dataloader.Candle, tradeSize float64) []DayStructure {
	eod := util.Eod(time.Unix(int64(candles[0].Time), 0))
	var daily []DayStructure
	dayOpen := candles[0].Close
	prevPrice := 0.0
	amount := tradeSize / candles[0].Close
	for _, c := range candles {
		if int64(c.Time) > eod.Unix() {
			profit := (prevPrice - dayOpen) * amount
			dailystr := DayStructure{
				Day:    util.Bod(eod).Unix(),
				Profit: profit,
			}
			daily = append(daily, dailystr)
			dayOpen = c.Close
			eod = util.Eod(time.Unix(int64(c.Time), 0))
		}
		prevPrice = c.Close
	}
	profit := (prevPrice - dayOpen) * amount
	dailystr := DayStructure{
		Day:    util.Bod(eod).Unix(),
		Profit: profit,
	}
	daily = append(daily, dailystr)
	return daily
}

func CalcMetrics(t structs.Trades, c []dataloader.Candle, tradeSize float64, fee float64) TradeMetrics {
	tsStart := c[0].Time
	tsEnd := c[len(c)-1].Time
	buyAndHoldTrade := MakeBuyAndHoldTrade(c, tradeSize, fee)
	tm := TradeMetrics{
		Trades:                 t,
		TsStart:                tsStart,
		TsEnd:                  tsEnd,
		BuyAndHoldProfit:       buyAndHoldTrade.Profit,
		ProfitFactor:           objectives.ProfitFactor(t, tsStart, tsEnd),
		ProfitFactorRobust:     objectives.ProfitFactorRobust(t, tsStart, tsEnd),
		TotalProfit:            objectives.ProfitSum(t),
		AvgReturnByAvgDrawdown: objectives.AvgReturnByAvgDrawdownLoss(t, tsStart, tsEnd),
		SharpeRobust:           objectives.SharpeRobust(t, tsStart, tsEnd),
		AvgDrawdown:            util.AvgDrawdown(t.Profits()),
		MaxDrawdown:            util.MaxDrawdownFloat64(t.Profits()),
		NumTrades:              len(t),
		BuyAndHoldDaily:        CalcDailyBHProfit(c, tradeSize),
		OverfittingScore:       objectives.CalcOverfittingScore(t),
	}
	return tm
}
