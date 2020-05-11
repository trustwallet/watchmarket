package coinmarketcap

import (
	"github.com/trustwallet/watchmarket/services/markets"
	"math/big"
)

func (p Provider) GetRates() (rates markets.Rates, err error) {
	prices, err := p.client.fetchPrices(p.currency)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.ID)
	return
}

func normalizeRates(prices CoinPrices, provider string) markets.Rates {
	var (
		res           markets.Rates
		emptyPlatform Platform
	)

	for _, price := range prices.Data {
		if price.Platform != emptyPlatform {
			continue
		}
		res = append(res, markets.Rate{
			Currency:         price.Symbol,
			Rate:             price.Quote.USD.Price,
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: *big.NewFloat(price.Quote.USD.PercentChange24h),
			Provider:         provider,
		})
	}
	return res
}
