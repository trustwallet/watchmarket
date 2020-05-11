package fixer

import (
	"github.com/trustwallet/watchmarket/services/markets"
)

const (
	id = "fixer"
)

func (p Provider) GetRates() (markets.Rates, error) {
	var result markets.Rates
	rawRates, err := p.client.FetchRates()
	if err != nil {
		return result, err
	}
	result = normalizeRates(rawRates, p.ID)
	return result, nil
}

func normalizeRates(rawRate Rate, provider string) markets.Rates {
	var result markets.Rates
	for currency, rate := range rawRate.Rates {
		result = append(result, markets.Rate{Currency: currency, Rate: rate, Timestamp: rawRate.Timestamp, Provider: provider})
	}
	return result
}
