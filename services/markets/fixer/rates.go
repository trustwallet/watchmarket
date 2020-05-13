package fixer

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	id = "fixer"
)

func (p Provider) GetRates() (watchmarket.Rates, error) {
	var result watchmarket.Rates
	rawRates, err := p.client.FetchRates()
	if err != nil {
		return result, err
	}
	result = normalizeRates(rawRates, p.id)
	return result, nil
}

func normalizeRates(rawRate Rate, provider string) watchmarket.Rates {
	var result watchmarket.Rates
	for currency, rate := range rawRate.Rates {
		result = append(result, watchmarket.Rate{Currency: currency, Rate: rate, Timestamp: rawRate.Timestamp, Provider: provider})
	}
	return result
}
