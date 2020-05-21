package controllers

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
)

const (
	ErrNotFound = "not found"
	ErrInternal = "internal"
)

type (
	ChartRequest struct {
		CoinQuery, Token, Currency, TimeStartRaw, MaxItems string
	}
	TickerRequest struct {
		Currency string `json:"Currency"`
		Assets   []Coin `json:"assets"`
	}

	TickerResponse struct {
		Currency string              `json:"Currency"`
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
