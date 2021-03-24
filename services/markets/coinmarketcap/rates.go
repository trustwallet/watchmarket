package coinmarketcap

import (
	"strings"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (p Provider) GetRates() (rates watchmarket.Rates, err error) {
	prices, err := p.client.fetchPrices(p.currency, "coins")
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.id)
	return
}

func normalizeRates(prices CoinPrices, provider string) watchmarket.Rates {
	var (
		result watchmarket.Rates
	)

	for _, price := range prices.Data {
		result = append(result, watchmarket.Rate{
			Currency:         strings.ToUpper(price.Symbol),
			Rate:             watchmarket.TruncateWithPrecision(price.Quote.USD.Price, watchmarket.DefaultPrecision),
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: price.Quote.USD.PercentChange24h,
			Provider:         provider,
		})
	}
	return result
}
