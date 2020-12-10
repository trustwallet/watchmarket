package worker

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
)

func (w Worker) FetchAndSaveRates() {
	log.Info("Fetching Rates ...")
	fetchedRates := fetchRates(w.ratesApis)
	normalizedRates := toRatesModel(fetchedRates)

	if err := w.db.AddRates(normalizedRates, w.configuration.Worker.BatchLimit); err != nil {
		log.Error(err)
	}
}

func fetchRates(ratesApis markets.RatesAPIs) watchmarket.Rates {
	wg := new(sync.WaitGroup)
	s := new(rates)
	for _, t := range ratesApis {
		wg.Add(1)
		go fetchRatesByProvider(t, wg, s)
	}
	wg.Wait()

	return s.rates
}

func fetchRatesByProvider(r markets.RatesAPI, wg *sync.WaitGroup, s *rates) {
	defer wg.Done()

	rates, err := r.GetRates()
	if err != nil {
		log.WithFields(log.Fields{"provider": r.GetProvider(), "details": err}).Error("Failed to fetch rates")
		return
	}

	log.WithFields(log.Fields{"provider": r.GetProvider(), "rates": len(rates)}).Info("Rates fetching done")

	s.Add(rates)
}
