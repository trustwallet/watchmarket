package fixer

import (
	"github.com/trustwallet/watchmarket/services/rates"
)

const (
	id = "fixer"
)

type Parser struct {
	ID, currency string
	client       Client
}

func InitParser(api, key, currency string) Parser {
	return Parser{
		ID:       id,
		currency: currency,
		client:   NewClient(api, key, currency),
	}
}

func (p Parser) GetData() (rates.Rates, error) {
	var result rates.Rates
	rawRates, err := p.client.FetchRates()
	if err != nil {
		return result, err
	}
	result = normalizeRates(rawRates, p.ID)
	return result, nil
}

func normalizeRates(rawRate Rate, provider string) rates.Rates {
	var result rates.Rates
	for currency, rate := range rawRate.Rates {
		result = append(result, rates.Rate{Currency: currency, Rate: rate, Timestamp: rawRate.Timestamp, Provider: provider})
	}
	return result
}
