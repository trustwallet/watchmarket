package tickers

import (
	"time"
)

type (
	CoinType string

	Response struct {
		Currency string  `json:"currency"`
		Docs     Tickers `json:"docs"`
	}

	Ticker struct {
		Coin       uint      `json:"coin"`
		CoinName   string    `json:"coin_name,omitempty"`
		CoinType   CoinType  `json:"type,omitempty"`
		Error      string    `json:"error,omitempty"`
		LastUpdate time.Time `json:"last_update,omitempty"`
		Price      Price     `json:"price,omitempty"`
		TokenId    string    `json:"token_id,omitempty"`
	}

	Price struct {
		Change24h float64 `json:"change_24h"`
		Currency  string  `json:"currency,omitempty"`
		Provider  string  `json:"provider,omitempty"`
		Value     float64 `json:"value"`
	}

	Tickers []Ticker
)

const (
	Coin  CoinType = "coin"
	Token CoinType = "token"
)
