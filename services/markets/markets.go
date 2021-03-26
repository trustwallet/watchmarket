package markets

import (
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets/coingecko"
	"github.com/trustwallet/watchmarket/services/markets/coinmarketcap"
	"github.com/trustwallet/watchmarket/services/markets/fixer"
)

type (
	Provider interface {
		GetProvider() string
	}

	RatesAPI interface {
		Provider
		GetRates() (watchmarket.Rates, error)
	}

	TickersAPI interface {
		Provider
		GetTickers() (watchmarket.Tickers, error)
	}

	ChartsAPI interface {
		Provider
		GetChartData(asset controllers.Asset, currency string, timeStart int64) (watchmarket.Chart, error)
		GetCoinData(asset controllers.Asset, currency string) (watchmarket.CoinDetails, error)
	}

	Providers   map[string]Provider
	RatesAPIs   map[string]RatesAPI
	TickersAPIs map[string]TickersAPI
	ChartsAPIs  map[string]ChartsAPI

	APIs struct {
		RatesAPIs
		TickersAPIs
		ChartsAPIs
	}
)

func Init(config config.Configuration, assets assets.Client) (APIs, error) {
	var (
		ratesAPIs   = make(RatesAPIs)
		tickersAPIs = make(TickersAPIs)
		chartsAPIs  = make(ChartsAPIs)
		providers   = setupProviders(config, assets)
	)

	for id, p := range providers {
		if t, ok := p.(RatesAPI); ok {
			ratesAPIs[id] = t
		}
		if t, ok := p.(TickersAPI); ok {
			tickersAPIs[id] = t
		}
		if t, ok := p.(ChartsAPI); ok {
			chartsAPIs[id] = t
		}
	}

	return APIs{ratesAPIs, tickersAPIs, chartsAPIs}, nil
}

func setupProviders(config config.Configuration, assets assets.Client) Providers {
	coinmarketcapPriveder := coinmarketcap.InitProvider(
		config.Markets.Coinmarketcap.API,
		config.Markets.Coinmarketcap.WebAPI,
		config.Markets.Coinmarketcap.WidgetAPI,
		config.Markets.Coinmarketcap.Key,
		config.Markets.Coinmarketcap.Currency,
		assets)
	coingeckoProvider := coingecko.InitProvider(
		config.Markets.Coingecko.API,
		config.Markets.Coingecko.Key,
		config.Markets.Coingecko.Currency,
		assets)
	fixerProvider := fixer.InitProvider(config.Markets.Fixer.API, config.Markets.Fixer.Key, config.Markets.Fixer.Currency)

	providers := make(Providers, 4)

	providers[coinmarketcapPriveder.GetProvider()] = coinmarketcapPriveder
	providers[coingeckoProvider.GetProvider()] = coingeckoProvider
	providers[fixerProvider.GetProvider()] = fixerProvider

	return providers
}
