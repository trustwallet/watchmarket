package coingecko

import (
	"github.com/trustwallet/watchmarket/services/rates"
	"strings"
)

const (
	id         = "coingecko"
	bucketSize = 500
)

type Parser struct {
	ID       string
	client   Client
	currency string
}

func InitParser(api, key, currency string) Parser {
	return Parser{
		ID:       id,
		client:   NewClient(api, key, bucketSize),
		currency: currency,
	}
}

func (p Parser) FetchLatestRates() (rates.Rates, error) {
	coins, err := p.client.FetchCoinsList()
	if err != nil {
		return rates.Rates{}, err
	}
	prices := p.client.FetchLatestRates(coins, p.currency)

	return normalizeRates(prices, p.ID), nil
}

func normalizeRates(prices Prices, provider string) rates.Rates {
	var result rates.Rates

	for _, price := range prices {
		r := 0.0
		if price.CurrentPrice != 0 {
			r = 1.0 / price.CurrentPrice
		}
		result = append(result, rates.Rate{
			Currency:  strings.ToUpper(price.Symbol),
			Rate:      r,
			Timestamp: price.LastUpdated.Unix(),
			Provider:  provider,
		})
	}
	return result
}
