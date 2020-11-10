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

	RateRequest struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	RateResponse struct {
		Amount float64 `json:"amount"`
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

	InfoResponse struct {
		Provider          string            `json:"provider,omitempty"`
		ProviderURL       string            `json:"provider_url,omitempty"`
		Vol24             float64           `json:"volume_24"`
		MarketCap         float64           `json:"market_cap"`
		CirculatingSupply float64           `json:"circulating_supply"`
		TotalSupply       float64           `json:"total_supply"`
		Info              *watchmarket.Info `json:"info,omitempty"`
	}

	AlertsRequest struct {
		Interval string `json:"interval"`
	}

	AlertsResponse struct {
		Assets map[string]AlertsDetails `json:"assets"`
	}

	AlertsDetails struct {
		PriceDifference float64 `json:"price_difference"`
		UpdatedAt       int64   `json:"updated_at"`
	}
)
