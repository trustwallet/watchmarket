package coinmarketcap

import (
	"github.com/trustwallet/watchmarket/services/rates"
	tickersClient "github.com/trustwallet/watchmarket/services/tickers/providers/coinmarketcap"
	"math/big"
)

const (
	id = "coinmarketcap"
)

type Provider struct {
	ID       string
	client   tickersClient.Client
	currency string
}

func InitProvider(proApi, assetsApi, key, currency string) Provider {
	return Provider{
		ID:       id,
		client:   tickersClient.NewClient(proApi, assetsApi, key),
		currency: currency,
	}
}

func (p Provider) GetData() (rates rates.Rates, err error) {
	prices, err := p.client.FetchPrices(p.currency)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.ID)
	return
}

func normalizeRates(prices tickersClient.CoinPrices, provider string) rates.Rates {
	var (
		res           rates.Rates
		emptyPlatform tickersClient.Platform
	)

	for _, price := range prices.Data {
		if price.Platform != emptyPlatform {
			continue
		}
		res = append(res, rates.Rate{
			Currency:         price.Symbol,
			Rate:             price.Quote.USD.Price,
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: *big.NewFloat(price.Quote.USD.PercentChange24h),
			Provider:         provider,
		})
	}
	return res
}
