package optimize

import (
	"math"
	"trade-optimizer/dataloader"
	"trade-optimizer/formulas"
	"trade-optimizer/structs"
)

const (
	ReasonStrategy = iota
	ReasonStoploss = iota
)

type Exchange struct {
	openTrades   map[string]*structs.Trade
	closedTrades structs.Trades
	bid          float64
	coolDown     map[string]int
}

func NewExchange() *Exchange {
	return &Exchange{
		openTrades: make(map[string]*structs.Trade),
		coolDown:   make(map[string]int),
	}
}

func newTrade(bid float64, fee float64, stoploss float64, candle *dataloader.Candle, tType int) *structs.Trade {
	t := &structs.Trade{
		TradeSize:  bid,
		TradeType:  tType,
		EnterPrice: candle.Close,
		ExitPrice:  -1.0,
		TsEnter:    candle.Time,
		TsExit:     0,
		Amount:     formulas.AmountForBid(bid, candle, fee),
		Fee:        fee,
		Stoploss:   stoploss,
		ExitReason: -1,
	}
	if stoploss > 0 {
		t.Stoploss = stoploss
	}
	return t
}

func (e *Exchange) Buy(pairname string, fee float64, stoploss float64, candle *dataloader.Candle, tradeSize float64, reason int) {
	if _, ok := e.openTrades[pairname]; !ok {
		//enterLong
		e.openTrades[pairname] = newTrade(tradeSize, fee, stoploss, candle, structs.TradeLong)
	} else if e.openTrades[pairname].TradeType == structs.TradeShort {
		//exitShort
		trade := e.openTrades[pairname]
		delete(e.openTrades, pairname)
		trade.TsExit = candle.Time
		trade.ExitPrice = candle.Close
		newProfit := formulas.CalcProfit(trade, candle.Close, trade.Fee)
		trade.MaxProfit = math.Max(trade.Profit, newProfit)
		trade.Profit = newProfit
		trade.ExitReason = reason
		e.closedTrades = append(e.closedTrades, trade)
		e.coolDown[pairname] = 1
	}
}

func (e *Exchange) Sell(pairname string, fee float64, stoploss float64, candle *dataloader.Candle, tradeSize float64, reason int) {
	if _, ok := e.openTrades[pairname]; !ok {
		//enter short
		e.openTrades[pairname] = newTrade(tradeSize, fee, stoploss, candle, structs.TradeShort)
	} else if e.openTrades[pairname].TradeType == structs.TradeLong {
		//exitLong
		trade := e.openTrades[pairname]
		delete(e.openTrades, pairname)
		trade.TsExit = candle.Time
		trade.ExitPrice = candle.Close
		newProfit := formulas.CalcProfit(trade, candle.Close, trade.Fee)
		trade.MaxProfit = math.Max(trade.Profit, newProfit)
		trade.Profit = newProfit
		trade.ExitReason = reason
		e.closedTrades = append(e.closedTrades, trade)
		e.coolDown[pairname] = 1
	}
}

func (e *Exchange) UpdateState(pairname string, c *dataloader.Candle) {
	if trade, ok := e.openTrades[pairname]; ok {
		newProfit := formulas.CalcProfit(trade, c.Close, trade.Fee)
		trade.MaxProfit = math.Max(newProfit, trade.Profit)
		trade.Profit = newProfit
		e.openTrades[pairname] = trade
	}
}

func (e *Exchange) FireStoplosses(c *dataloader.Candle) {
	for pairname, t := range e.openTrades {
		if t.Stoploss > 0 {
			panic("Positive stoploss")
		}
		if t.Profit/t.Amount/t.EnterPrice < t.Stoploss {
			if t.TradeType == structs.TradeLong {
				e.Sell(pairname, t.Fee, t.Stoploss, c, t.TradeSize, ReasonStoploss)
			} else if t.TradeType == structs.TradeShort {
				e.Buy(pairname, t.Fee, t.Stoploss, c, t.TradeSize, ReasonStoploss)
			}
		}
	}
}
