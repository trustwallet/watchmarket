package coingecko

import (
	"context"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

func (p Provider) GetRates(ctx context.Context) (watchmarket.Rates, error) {
	coins, err := p.client.fetchCoins(ctx)
	if err != nil {
		return watchmarket.Rates{}, err
	}
	prices := p.client.fetchRates(coins, p.currency, ctx)

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
