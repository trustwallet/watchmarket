package coingecko

import (
	"github.com/trustwallet/watchmarket/services/charts"
	"time"
)

type (
	Platforms map[string]string

	Coin struct {
		Id        string    `json:"id"`
		Symbol    string    `json:"symbol"`
		Name      string    `json:"name"`
		Platforms Platforms `json:"platforms"`
	}
	Coins []Coin

	CoinResult struct {
		Symbol   string
		TokenId  string
		CoinType charts.CoinType
	}

	CoinPrices []CoinPrice

	CoinPrice struct {
		Id                           string    `json:"id"`
		Symbol                       string    `json:"symbol"`
		Name                         string    `json:"name"`
		CurrentPrice                 float64   `json:"current_price"`
		PriceChange24h               float64   `json:"price_change_24h"`
		PriceChangePercentage24h     float64   `json:"price_change_percentage_24h"`
		MarketCapChange24h           float64   `json:"market_cap_change_24h"`
		MarketCapChangePercentage24h float64   `json:"market_cap_change_percentage_24h"`
		MarketCap                    float64   `json:"market_cap"`
		TotalVolume                  float64   `json:"total_volume"`
		CirculatingSupply            float64   `json:"circulating_supply"`
		TotalSupply                  float64   `json:"total_supply"`
		LastUpdated                  time.Time `json:"last_updated"`
	}

	Charts struct {
		Prices     []Volume `json:"prices"`
		MarketCaps []Volume `json:"market_caps"`
		Volumes    []Volume `json:"total_volumes"`
	}

	Volume []float64
)

func (coins Coins) coinIds() []string {
	coinIds := make([]string, 0)
	for _, coin := range coins {
		coinIds = append(coinIds, coin.Id)
	}
	return coinIds
}