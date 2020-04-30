package model

import (
	"math/big"
	"time"
)

type (
	CoinType string

	TickerResponse struct {
		Currency string  `json:"currency"`
		Docs     Tickers `json:"docs"`
	}

	Ticker struct {
		Coin       uint        `json:"coin"`
		CoinName   string      `json:"coin_name,omitempty"`
		TokenId    string      `json:"token_id,omitempty"`
		CoinType   CoinType    `json:"type,omitempty"`
		Price      TickerPrice `json:"price,omitempty"`
		LastUpdate time.Time   `json:"last_update,omitempty"`
		Error      string      `json:"error,omitempty"`
	}

	TickerPrice struct {
		Value     float64 `json:"value"`
		Change24h float64 `json:"change_24h"`
		Currency  string  `json:"currency,omitempty"`
		Provider  string  `json:"provider,omitempty"`
	}

	Tickers []*Ticker
)

const (
	TypeCoin  CoinType = "coin"
	TypeToken CoinType = "token"
)

func (t *Ticker) SetCoinId(coinId uint) {
	t.Coin = coinId
	t.CoinName = ""
	t.Price.Provider = ""
	t.Price.Currency = ""
}

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
