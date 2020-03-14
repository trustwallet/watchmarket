package watchmarket

import (
	"math/big"
	"time"
)

const (
	TypeCoin  CoinType = "coin"
	TypeToken CoinType = "token"

	DefaultCurrency = "USD"
)

type CoinType string

type TickerResponse struct {
	Currency string  `json:"currency"`
	Docs     Tickers `json:"docs"`
}

type Ticker struct {
	Coin       uint        `json:"coin"`
	CoinName   string      `json:"coin_name,omitempty"`
	TokenId    string      `json:"token_id,omitempty"`
	CoinType   CoinType    `json:"type,omitempty"`
	Price      TickerPrice `json:"price,omitempty"`
	LastUpdate time.Time   `json:"last_update,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type ChartData struct {
	Prices []ChartPrice `json:"prices,omitempty"`
	Error  string       `json:"error,omitempty"`
}

type ChartPrice struct {
	Price float64 `json:"price"`
	Date  int64   `json:"date"`
}

type ChartCoinInfo struct {
	Vol24             float64   `json:"volume_24"`
	MarketCap         float64   `json:"market_cap"`
	CirculatingSupply float64   `json:"circulating_supply"`
	TotalSupply       float64   `json:"total_supply"`
	Info              *CoinInfo `json:"info,omitempty"`
}

type CoinInfo struct {
	Name             string       `json:"name,omitempty"`
	Website          string       `json:"website,omitempty"`
	SourceCode       string       `json:"source_code,omitempty"`
	WhitePaper       string       `json:"white_paper,omitempty"`
	Description      string       `json:"description,omitempty"`
	ShortDescription string       `json:"short_description,omitempty"`
	Explorers        []Link       `json:"explorers,omitempty"`
	Socials          []SocialLink `json:"socials,omitempty"`
	DataSource       string       `json:"data_source,omitempty"`
}

type SocialLink struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Handle string `json:"handle"`
}

type Link struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (t *Ticker) SetCoinId(coinId uint) {
	t.Coin = coinId
	t.CoinName = ""
	t.Price.Provider = ""
	t.Price.Currency = ""
}

type TickerPrice struct {
	Value     float64 `json:"value"`
	Change24h float64 `json:"change_24h"`
	Currency  string  `json:"currency,omitempty"`
	Provider  string  `json:"provider,omitempty"`
}

type Rate struct {
	Currency         string     `json:"currency"`
	Rate             float64    `json:"rate"`
	Timestamp        int64      `json:"timestamp"`
	PercentChange24h *big.Float `json:"percent_change_24h,omitempty"`
	Provider         string     `json:"provider,omitempty"`
}

type Rates []Rate
type Tickers []*Ticker

func (ts Tickers) ApplyRate(currency string, rate float64, percentChange24h *big.Float) {
	for _, t := range ts {
		t.ApplyRate(currency, rate, percentChange24h)
	}
}

func (t *Ticker) ApplyRate(currency string, rate float64, percentChange24h *big.Float) {
	if t.Price.Currency == currency {
		return
	}
	t.Price.Value *= rate
	t.Price.Currency = currency

	if percentChange24h != nil {
		change24h, _ := percentChange24h.Float64()
		t.Price.Change24h -= change24h
	}
}

func (d ChartData) IsEmpty() bool {
	if len(d.Prices) == 0 {
		return true
	}
	return false
}
