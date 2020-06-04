package controllers

import "github.com/trustwallet/watchmarket/pkg/watchmarket"

type (
	TickerRequest struct {
		Currency string `json:"Currency"`
		Assets   []Coin `json:"assets"`
	}

	TickerRequestV2 struct {
		Currency string   `json:"currency"`
		Ids      []string `json:"assets"`
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
		Price     float64 `json:"price"`
		ID        string  `json:"id"`
	}

	ChartRequest struct {
		CoinQuery, Token, Currency, TimeStartRaw, MaxItems string
	}

	DetailsRequest struct {
		CoinQuery, Token, Currency string
	}
)
