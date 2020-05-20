package worker

import (
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
	"sync"
)

func (w Worker) FetchAndSaveRates() {
	fetchedRates := fetchRates(w.ratesApis)
	normalizedRates := toRatesModel(fetchedRates)

	if err := w.db.AddRates(normalizedRates); err != nil {
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
			Timestamp:        r.Timestamp,
		})
	}
	return result
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
		logger.Error(err)
	}

	s.Lock()
	s.rates = append(s.rates, rates...)
	s.Unlock()
}
