package cmc

import (
	"github.com/trustwallet/watchmarket/market/clients/cmc"
	"github.com/trustwallet/watchmarket/market/rate"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"math/big"
)

const (
	id = "cmc"
)

type Cmc struct {
	rate.Rate
	mapApi string
	client *cmc.Client
}

func InitRate(api string, apiKey string, mapApi string, updateTime string) rate.Provider {
	cmc := &Cmc{
		Rate: rate.Rate{
			Id:         id,
			UpdateTime: updateTime,
		},
		mapApi: mapApi,
		client: cmc.NewClient(api, apiKey),
	}
	return cmc
}

func (c *Cmc) FetchLatestRates() (rates watchmarket.Rates, err error) {
	prices, err := c.client.GetData()
	if err != nil {
		return
	}
	rates = normalizeRates(prices, c.GetId())
	return
}

func normalizeRates(prices cmc.CoinPrices, provider string) (rates watchmarket.Rates) {
	for _, price := range prices.Data {
		if price.Platform != nil {
			continue
		}
		rates = append(rates, watchmarket.Rate{
			Currency:         price.Symbol,
			Rate:             1.0 / price.Quote.USD.Price,
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: big.NewFloat(price.Quote.USD.PercentChange24h),
			Provider:         provider,
		})
	}
	return
}
