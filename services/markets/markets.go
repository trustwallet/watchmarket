package markets

import "github.com/trustwallet/watchmarket/pkg/watchmarket"

type (
	Provider interface {
		GetProvider() string
	}

	RatesAPI interface {
		Provider
		GetRates() watchmarket.Rates
	}

	TickersAPI interface {
		Provider
		GetTickers() watchmarket.Tickers
	}

	ChartsAPI interface {
		Provider
		GetChartData(coinID uint, token, currency string, timeStart int64) (watchmarket.Chart, error)
		GetCoinData(coinID uint, token, currency string) (watchmarket.CoinDetails, error)
	}

	RatesAPIs   map[string]RatesAPI
	TickersAPIs map[string]TickersAPI
	ChartsAPIs  map[string]ChartsAPI

	APIs struct {
		RatesAPIs
		TickersAPIs
		ChartsAPIs
	}
)

func Init() (APIs, error) {
	a := APIs{}

	return a, nil
}
