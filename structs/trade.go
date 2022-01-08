package structs

const (
	TradeLong  = iota
	TradeShort = iota
)

type Trade struct {
	Profit     float64 `json:"profit"`
	MaxProfit  float64 `json:"max_profit"`
	TradeSize  float64 `json:"trade_size"`
	Fee        float64 `json:"fee"`
	Amount     float64 `json:"amount"`
	EnterPrice float64 `json:"enter_price"`
	ExitPrice  float64 `json:"exit_price"`
	TsEnter    uint64  `json:"ts_buy"`
	TsExit     uint64  `json:"ts_sell"`
	TradeType  int     `json:"trade_type"`
	Stoploss   float64 `json:"stoploss"`
	ExitReason int     `json:"exit_reason"`
}

type Trades []*Trade

func (t Trades) Profits() []float64 {
	r := make([]float64, 0, len(t))
	for _, x := range t {
		r = append(r, x.Profit)
	}
	return r
}

func (t Trades) ProfitsInPercents() []float64 {
	r := make([]float64, 0, len(t))
	for _, x := range t {
		r = append(r, (x.Profit/x.TradeSize)*100)
	}
	return r
}
