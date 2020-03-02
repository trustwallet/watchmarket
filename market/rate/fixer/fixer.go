package fixer

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/watchmarket/market/rate"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/url"
)

const (
	id = "fixer"
)

type Fixer struct {
	rate.Rate
	APIKey string
	blockatlas.Request
}

func InitRate(api string, apiKey string, updateTime string) rate.RateProvider {
	return &Fixer{
		Rate: rate.Rate{
			Id:         id,
			UpdateTime: updateTime,
		},
		Request: blockatlas.InitClient(api),
		APIKey:  apiKey,
	}
}

func (f *Fixer) FetchLatestRates() (rates watchmarket.Rates, err error) {
	values := url.Values{
		"access_key": {f.APIKey},
		"base":       {watchmarket.DefaultCurrency}, // Base USD supported only in paid api
	}
	var latest Latest
	err = f.Get(&latest, "latest", values)
	if err != nil {
		return
	}
	rates = normalizeRates(latest, f.GetId())
	return
}

func normalizeRates(latest Latest, provider string) (rates watchmarket.Rates) {
	for currency, rate := range latest.Rates {
		rates = append(rates, watchmarket.Rate{Currency: currency, Rate: rate, Timestamp: latest.Timestamp, Provider: provider})
	}
	return
}
