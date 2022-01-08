package formulas

import (
	"trade-optimizer/dataloader"
	"trade-optimizer/structs"
)

func AmountForBid(bid float64, candle *dataloader.Candle, fee float64) float64 {
	return bid / candle.Open
}

func CalcProfit(t *structs.Trade, exitPrice float64, fee float64) float64 {
	priceDiff := exitPrice - t.EnterPrice
	profit := t.Amount * priceDiff
	if t.TradeType == structs.TradeShort {
		profit = -profit
	}
	return profit - ((t.ExitPrice * t.Amount) * fee)
}
