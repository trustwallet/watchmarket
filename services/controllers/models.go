package controllers

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
)

type (
	ChartRequest struct {
		coinQuery, token, currency, timeStartRaw, maxItems string
	}
	TickerRequest struct {
		Currency string `json:"currency"`
		Assets   []Coin `json:"assets"`
	}

	TickerResponse struct {
		Currency string              `json:"currency"`
		Tickers  watchmarket.Tickers `json:"docs"`
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

	sortedTickersResponse struct {
		sync.Mutex
		tickers []models.Ticker
	}
)
