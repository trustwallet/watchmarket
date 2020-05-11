package coingecko

import (
	"github.com/trustwallet/watchmarket/services/markets"
	"strings"
)

func (p Provider) GetRates() (markets.Rates, error) {
	coins, err := p.client.fetchCoins()
	if err != nil {
		return markets.Rates{}, err
	}
	prices := p.client.fetchRates(coins)

	return normalizeRates(prices, p.ID), nil
}

func normalizeRates(prices CoinPrices, provider string) markets.Rates {
	var result markets.Rates

	for _, price := range prices {
		result = append(result, markets.Rate{
			Currency:  strings.ToUpper(price.Symbol),
			Rate:      price.CurrentPrice,
			Timestamp: price.LastUpdated.Unix(),
			Provider:  provider,
		})
	}
	return result
}
