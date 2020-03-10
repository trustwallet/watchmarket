package market

import (
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market/rate"
	"github.com/trustwallet/watchmarket/storage"
)

var rateProviders rate.Providers

func InitRates(storage storage.Market, providers *rate.Providers) *cron.Cron {
	rateProviders = *providers
	return scheduleRates(storage, rateProviders)
}

func scheduleRates(storage storage.Market, rates rate.Providers) *cron.Cron {
	c := cron.New()
	for _, r := range rates {
		scheduleTasks(storage, r, c)
	}
	return c
}

func runRate(storage storage.Market, p rate.RateProvider) {
	rates, err := p.FetchLatestRates()
	if err != nil {
		logger.Error("Failed to fetch rates", logger.Params{"error": err, "provider": p.GetId()})
		return
	}
	results := storage.SaveRates(rates, rateProviders)
	logger.Info("Done fetching latest rates", logger.Params{"numFetchedRates": len(rates), "provider": p.GetId(), "results": results})
}
