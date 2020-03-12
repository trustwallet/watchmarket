package market

import (
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
)

var tickerProviders ticker.Providers

func InitTickers(storage storage.Market, providers *ticker.Providers) *cron.Cron {
	tickerProviders = *providers
	return scheduleTickers(storage, tickerProviders)
}

func scheduleTickers(storage storage.Market, ps ticker.Providers) *cron.Cron {
	c := cron.New()
	for _, p := range ps {
		scheduleTasks(storage, p, c)
	}
	return c
}

func runTicker(storage storage.Market, p ticker.TickerProvider) {
	data, err := p.GetData()
	if err != nil {
		logger.Error("Failed to fetch market data", logger.Params{"error": err, "provider": p.GetId()})
		return
	}

	ch := make(chan string)
	save := func(ticker *watchmarket.Ticker, ch chan string) {
		res, err := storage.SaveTicker(ticker, tickerProviders)
		if err != nil {
			logger.Error("Failed to save ticker", logger.Params{"error": err, "provider": p.GetId(), "result": ticker})
		}
		ch <- string(res)
	}

	for _, result := range data {
		go save(result, ch)
	}

	results := make(map[string]int)
	for i := 0; i < len(data); i++ {
		results[<-ch]++
	}

	logger.Info("Market data result", logger.Params{"markets": len(data), "provider": p.GetId(), "results": results})
}
