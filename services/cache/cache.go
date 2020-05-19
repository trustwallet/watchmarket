package cache

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/cache/redis"
	"time"
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

	Providers    map[string]Provider
	RatesCache   map[string]Rates
	TickersCache map[string]Tickers
	ChartsCache  map[string]Charts

	Cache struct {
		RatesCache
		TickersCache
		ChartsCache
	}
)

func Init(redis redis.Redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching time.Duration) Cache {
	var (
		ratesCache   = make(RatesCache, 0)
		tickersCache = make(TickersCache, 0)
		chartsCache  = make(ChartsCache, 0)
		providers    = setupProviders(redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching)
	)

	for id, p := range providers {
		if t, ok := p.(Rates); ok {
			ratesCache[id] = t
		}
		if t, ok := p.(Tickers); ok {
			tickersCache[id] = t
		}
		if t, ok := p.(Charts); ok {
			chartsCache[id] = t
		}
	}

	return Cache{ratesCache, tickersCache, chartsCache}
}

func setupProviders(redis redis.Redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching time.Duration) Providers {
	ri := rediscache.Init(redis, chartsCaching, tickersCaching, ratesCaching, detailsCaching)
	p := make(Providers, 1)
	p[ri.GetID()] = ri
	return p
}
