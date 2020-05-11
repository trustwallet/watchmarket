package coinmarketcap

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"math/big"
)

func (p Provider) GetRates() (rates watchmarket.Rates, err error) {
	prices, err := p.client.fetchPrices(p.currency)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.id)
	return
}

func normalizeRates(prices CoinPrices, provider string) watchmarket.Rates {
	var (
		res           watchmarket.Rates
		emptyPlatform Platform
	)

	for _, price := range prices.Data {
		if price.Platform != emptyPlatform {
			continue
		}
		res = append(res, watchmarket.Rate{
			Currency:         price.Symbol,
			Rate:             price.Quote.USD.Price,
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: *big.NewFloat(price.Quote.USD.PercentChange24h),
			Provider:         provider,
		})
	}
	return res
}
