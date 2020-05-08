package coingecko

import (
	"github.com/trustwallet/watchmarket/services/rates"
	"strings"
)

const (
	id         = "coingecko"
	bucketSize = 500
)

type Provider struct {
	ID       string
	client   Client
	currency string
}

func InitProvider(api, currency string) Provider {
	return Provider{
		ID:       id,
		client:   NewClient(api, currency, bucketSize),
		currency: currency,
	}
}

func (p Provider) GetData() (rates.Rates, error) {
	coins, err := p.client.fetchCoinsList()
	if err != nil {
		return rates.Rates{}, err
	}
	prices := p.client.fetchLatestRates(coins, p.currency)

	return normalizeRates(prices, p.ID), nil
}

func normalizeRates(prices Prices, provider string) rates.Rates {
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
