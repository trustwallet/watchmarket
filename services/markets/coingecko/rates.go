package coingecko

import (
	"strings"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (p Provider) GetRates() (watchmarket.Rates, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return watchmarket.Rates{}, err
	}
	prices := p.client.fetchRates(coins, p.currency)

	return normalizeRates(prices, p.id), nil
}

func normalizeRates(prices CoinPrices, provider string) watchmarket.Rates {
	var result watchmarket.Rates

	for _, price := range prices {
		result = append(result, watchmarket.Rate{
			Currency:  strings.ToUpper(price.Symbol),
			Rate:      watchmarket.TruncateWithPrecision(price.CurrentPrice, watchmarket.DefaultPrecision),
			Timestamp: price.LastUpdated.Unix(),
			Provider:  provider,
		})
	}
	return result
}
