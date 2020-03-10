package market

import (
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market/market"
	"github.com/trustwallet/watchmarket/storage"
)

var marketProviders market.Providers

func InitMarkets(storage storage.Market, providers *market.Providers) *cron.Cron {
	marketProviders = *providers
	return scheduleMarkets(storage, marketProviders)
}

func scheduleMarkets(storage storage.Market, ps market.Providers) *cron.Cron {
	c := cron.New()
	for _, p := range ps {
		scheduleTasks(storage, p, c)
	}
	return c
}

func runMarket(storage storage.Market, p market.MarketProvider) {
	data, err := p.GetData()
	if err != nil {
		logger.Error("Failed to fetch market data", logger.Params{"error": err, "provider": p.GetId()})
		return
	}
	results := make(map[string]int)
	for _, result := range data {
		res, err := storage.SaveTicker(result, marketProviders)
		results[string(res)]++
		if err != nil {
			logger.Error("Failed to save ticker", logger.Params{"error": err, "provider": p.GetId(), "result": result})
		}
	}

	logger.Info("Market data result", logger.Params{"markets": len(data), "provider": p.GetId(), "results": results})
}
