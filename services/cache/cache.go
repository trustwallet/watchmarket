package cache

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	Provider interface {
		GetID() string
		GenerateKey(data string) string
	}

	Charts interface {
		Provider
		GetCharts(key string, timeStart int64) (watchmarket.Chart, error)
		SaveCharts(key string, data watchmarket.Chart, timeStart int64) error
		SaveCoinDetails(key string, data watchmarket.CoinDetails, timeStart int64) error
		GetCoinDetails(key string, timeStart int64) (watchmarket.CoinDetails, error)
	}

	Tickers interface {
		Provider
		GetTickers(key string) (watchmarket.Tickers, error)
		SaveTickers(key string, tickers watchmarket.Tickers) error
	}

	Rates interface {
		Provider
		GetRates(key string) (watchmarket.Rates, error)
		SaveRates(key string, tickers watchmarket.Rates) error
	}
)
