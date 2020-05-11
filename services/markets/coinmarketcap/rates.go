package coinmarketcap

import (
	"github.com/trustwallet/watchmarket/services/rates"
	"math/big"
)

func (p Provider) GetRates() (rates rates.Rates, err error) {
	prices, err := p.client.fetchPrices(p.currency)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.ID)
	return
}

func normalizeRates(prices CoinPrices, provider string) rates.Rates {
	var (
		res           rates.Rates
		emptyPlatform Platform
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
