package coingecko

import (
	"fmt"
	"time"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	Charts struct {
		Prices     []Volume `json:"prices"`
		MarketCaps []Volume `json:"market_caps"`
		Volumes    []Volume `json:"total_volumes"`
	}

	Volume []float64

	CoinResult struct {
		Symbol          string
		TokenId         string
		CoinType        watchmarket.CoinType
		PotentialCoinID uint
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

	Coins []Coin

	Coin struct {
		Id        string    `json:"id"`
		Symbol    string    `json:"symbol"`
		Name      string    `json:"name"`
		Platforms Platforms `json:"platforms"`
	}

	Platforms map[string]string
)

func (coins Coins) coinIds() []string {
	coinIds := make([]string, 0)
	for _, coin := range coins {
		coinIds = append(coinIds, coin.Id)
	}
	return coinIds
}

func (cp CoinPrice) getUrl() string {
	return fmt.Sprintf("https://www.coingecko.com/en/coins/%s", cp.Id)
}
