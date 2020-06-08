package worker

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
	"go.elastic.co/apm"
	"sync"
	"time"
)

func (w Worker) FetchAndSaveRates() {
	tx := apm.DefaultTracer.StartTransaction("FetchAndSaveRates", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	logger.Info("Fetching Rates ...")
	fetchedRates := fetchRates(w.ratesApis, ctx)
	normalizedRates := toRatesModel(fetchedRates)

	if err := w.db.AddRates(normalizedRates, w.configuration.Worker.BatchLimit, ctx); err != nil {
		logger.Error(err)
	}
}

func toRatesModel(rates watchmarket.Rates) []models.Rate {
	result := make([]models.Rate, 0, len(rates))
	for _, r := range rates {
		result = append(result, models.Rate{
			Currency:         r.Currency,
			PercentChange24h: r.PercentChange24h,
			Provider:         r.Provider,
			Rate:             r.Rate,
			LastUpdated:      time.Unix(r.Timestamp, 0),
		})
	}
	return result
}

func fetchRates(ratesApis markets.RatesAPIs, ctx context.Context) watchmarket.Rates {
	wg := new(sync.WaitGroup)
	s := new(rates)
	for _, t := range ratesApis {
		wg.Add(1)
		go fetchRatesByProvider(t, wg, s, ctx)
	}
	wg.Wait()

	return s.rates
}

func fetchRatesByProvider(r markets.RatesAPI, wg *sync.WaitGroup, s *rates, ctx context.Context) {
	defer wg.Done()

	rates, err := r.GetRates(ctx)
	if err != nil {
		logger.Error("Failed to fetch rates", logger.Params{"provider": r.GetProvider(), "details": err})
		return
	}

	logger.Info("Rates fetching done", logger.Params{"provider": r.GetProvider(), "rates": len(rates)})

	s.Add(rates)
}
