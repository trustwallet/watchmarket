package worker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
	"go.elastic.co/apm"
	"sync"
)

func (w Worker) FetchAndSaveRates() {
	tx := apm.DefaultTracer.StartTransaction("FetchAndSaveRates", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	log.Info("Fetching Rates ...")
	fetchedRates := fetchRates(w.ratesApis, ctx)
	normalizedRates := toRatesModel(fetchedRates)

	if err := w.db.AddRates(normalizedRates, w.configuration.Worker.BatchLimit, ctx); err != nil {
		log.Error(err)
	}
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
		log.WithFields(log.Fields{"provider": r.GetProvider(), "details": err}).Error("Failed to fetch rates")
		return
	}

	log.WithFields(log.Fields{"provider": r.GetProvider(), "rates": len(rates)}).Info("Rates fetching done")

	s.Add(rates)
}
