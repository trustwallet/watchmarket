package markets

import (
	"context"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/markets/binancedex"
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
		GetRates(ctx context.Context) (watchmarket.Rates, error)
	}

	TickersAPI interface {
		Provider
		GetTickers(ctx context.Context) (watchmarket.Tickers, error)
	}

	ChartsAPI interface {
		Provider
		GetChartData(coinID uint, token, currency string, timeStart int64, ctx context.Context) (watchmarket.Chart, error)
		GetCoinData(coinID uint, token, currency string, ctx context.Context) (watchmarket.CoinDetails, error)
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
	b := binancedex.InitProvider(config.Markets.BinanceDex.API)
	cmc := coinmarketcap.InitProvider(
		config.Markets.Coinmarketcap.API,
		config.Markets.Coinmarketcap.WebAPI,
		config.Markets.Coinmarketcap.WidgetAPI,
		config.Markets.Coinmarketcap.Key,
		config.Markets.Coinmarketcap.Currency,
		assets)
	cg := coingecko.InitProvider(config.Markets.Coingecko.API, config.Markets.Coingecko.Currency, assets)
	f := fixer.InitProvider(config.Markets.Fixer.API, config.Markets.Fixer.Key, config.Markets.Fixer.Currency)

	ps := make(Providers, 4)

	ps[b.GetProvider()] = b
	ps[cmc.GetProvider()] = cmc
	ps[cg.GetProvider()] = cg
	ps[f.GetProvider()] = f

	return ps
}
