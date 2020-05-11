package coingecko

import (
	"github.com/trustwallet/watchmarket/services/rates"
	"strings"
)

func (p Provider) GetRates() (rates.Rates, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return rates.Rates{}, err
	}
	prices := p.client.fetchRates(coins)

	return normalizeRates(prices, p.ID), nil
}

func normalizeRates(prices CoinPrices, provider string) rates.Rates {
	var result rates.Rates

	for _, price := range prices {
		result = append(result, rates.Rate{
			Currency:  strings.ToUpper(price.Symbol),
			Rate:      price.CurrentPrice,
			Timestamp: price.LastUpdated.Unix(),
			Provider:  provider,
		})
	}
	return result
}
