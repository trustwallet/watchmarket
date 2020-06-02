package controllers

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
)

const (
	ErrNotFound   = "not found"
	ErrBadRequest = "bad request"
	ErrInternal   = "internal"
)

type (
	ChartRequest struct {
		CoinQuery, Token, Currency, TimeStartRaw, MaxItems string
	}

	ChartsNormalizedRequest struct {
		Coin            uint
		Token, Currency string
		TimeStart       int64
		MaxItems        int
	}

	DetailsRequest struct {
		CoinQuery, Token, Currency string
	}

	DetailsNormalizedRequest struct {
		Coin            uint
		Token, Currency string
	}

	TickerRequest struct {
		Currency string `json:"Currency"`
		Assets   []Coin `json:"assets"`
	}

	TickerRequestV2 struct {
		Currency string   `json:"Currency"`
		Ids      []string `json:"ids"`
	}

	Coin struct {
		Coin     uint                 `json:"Coin"`
		CoinType watchmarket.CoinType `json:"type"`
		TokenId  string               `json:"token_id,omitempty"`
	}

	TickerResponse struct {
		Currency string              `json:"currency"`
		Tickers  watchmarket.Tickers `json:"docs"`
	}

	TickerResponseV2 struct {
		Currency string        `json:"currency"`
		Tickers  []TickerPrice `json:"tickers"`
	}

	TickerPrice struct {
		Change24h float64 `json:"change_24h"`
		Provider  string  `json:"provider"`
		Price     float64 `json:"value"`
		ID        string
	}

	sortedTickersResponse struct {
		sync.Mutex
		tickers []models.Ticker
	}
)
