package market

import (
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/storage"
)

var tickerProviders ticker.Providers

func InitTickers(storage storage.Market, providers *ticker.Providers) *cron.Cron {
	tickerProviders = *providers
	return scheduleTickers(storage, tickerProviders)
}

func runTicker(storage storage.Market, p ticker.TickerProvider) {
	data, err := p.GetData()
	if err != nil {
		logger.Error("Failed to fetch market data", logger.Params{"error": err, "provider": p.GetId()})
		return
	}
	results := make(map[string]int)
	for _, result := range data {
		res, err := storage.SaveTicker(result, tickerProviders)
		results[string(res)]++
		if err != nil {
			logger.Error("Failed to save ticker", logger.Params{"error": err, "provider": p.GetId(), "result": result})
		}
	}

	logger.Info("Market data result", logger.Params{"markets": len(data), "provider": p.GetId(), "results": results})
}

func scheduleTickers(storage storage.Market, ps ticker.Providers) *cron.Cron {
	c := cron.New()
	for _, p := range ps {
		scheduleTasks(storage, p, c)
	}
	return c
}
