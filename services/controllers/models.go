package controllers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"time"
)

type (
	TickerRequest struct {
		Currency string  `json:"Currency"`
		Assets   []Asset `json:"assets"`
	}

	TickerRequestV2 struct {
		Currency string   `json:"currency"`
		Ids      []string `json:"assets"`
	}

	Asset struct {
		CoinId   uint                 `json:"Coin"`
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

	FiatRate struct {
		Currency string  `json:"currency"`
		Rate     float64 `json:"rate"`
	}

	FiatRates []FiatRate

	TickerPrice struct {
		Change24h float64 `json:"change_24h"`
		Provider  string  `json:"provider"`
		Price     float64 `json:"price"`
		ID        string  `json:"id"`
	}

	ChartRequest struct {
		Asset     Asset
		Currency  string
		TimeStart int64
		MaxItems  int
	}

	DetailsRequest struct {
		Asset    Asset
		Currency string
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
)

func GetCoinId(rawCoinId string) (uint, error) {
	coinId, err := strconv.Atoi(rawCoinId)
	if err != nil {
		return 0, err
	}
	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return 0, errors.New(fmt.Sprintf("invalid coin Id: %d", coinId))
	}
	return uint(coinId), nil
}

func GetCurrency(rawCurrency string) string {
	currency := rawCurrency
	if currency == "" {
		currency = watchmarket.DefaultCurrency
	}
	return currency
}

func GetTimeStart(rawTime string) int64 {
	timeStart, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		timeStart = time.Now().Unix() - int64(time.Hour)*24
	}
	return timeStart
}

func GetMaxItems(rawMax string) int {
	maxItems, err := strconv.Atoi(rawMax)
	if err != nil || maxItems <= 0 {
		maxItems = watchmarket.DefaultMaxChartItems
	}
	return maxItems
}
