package coingecko

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/market/rate"
	watchmarket "github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

const (
	id = "coingecko"
)

type Coingecko struct {
	client *coingecko.Client
	rate.Rate
}

func InitRate(api string, updateTime string) rate.Provider {
	return &Coingecko{
		client: coingecko.NewClient(api),
		Rate: rate.Rate{
			Id:         id,
			UpdateTime: updateTime,
		},
	}
}

func (c *Coingecko) FetchLatestRates() (rates watchmarket.Rates, err error) {
	coins, err := c.client.FetchCoinsList()
	if err != nil {
		return
	}
	prices := c.client.FetchLatestRates(coins, blockatlas.DefaultCurrency)

	rates = normalizeRates(prices, c.GetId())
	return
}

func normalizeRates(coinPrices coingecko.CoinPrices, provider string) (rates watchmarket.Rates) {
	for _, price := range coinPrices {
		rates = append(rates, watchmarket.Rate{
			Currency:  strings.ToUpper(price.Symbol),
			Rate:      1.0 / price.CurrentPrice,
			Timestamp: price.LastUpdated.Unix(),
			Provider:  provider,
		})
	}
	return
}
