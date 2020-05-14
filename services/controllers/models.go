package controllers

import "github.com/trustwallet/watchmarket/pkg/watchmarket"

type (
	ChartRequest struct {
		coinQuery, token, currency, timeStartRaw, maxItems string
	}
	TickerRequest struct {
		Currency string `json:"currency"`
		Assets   []Coin `json:"assets"`
	}
	Coin struct {
		Coin     uint                 `json:"coin"`
		CoinType watchmarket.CoinType `json:"type"`
		TokenId  string               `json:"token_id,omitempty"`
	}
	ChartsNormalizedRequest struct {
		coin            uint
		token, currency string
		timeStart       int64
		maxItems        int
	}
)
