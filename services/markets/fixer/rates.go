package fixer

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
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
	result = normalizeRates(rawRates, p.id, p.currency)
	return result, nil
}

func normalizeRates(rawRate Rate, provider, baseCurrency string) watchmarket.Rates {
	var result watchmarket.Rates
	for currency, rate := range rawRate.Rates {
		result = append(result, watchmarket.Rate{
			Currency:  strings.ToUpper(currency),
			Rate:      watchmarket.TruncateWithPrecision(1/rate, watchmarket.DefaultPrecision),
			Timestamp: rawRate.Timestamp,
			Provider:  provider,
		})
	}
	return result
}
