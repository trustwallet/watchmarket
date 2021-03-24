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

	allowCryptoCurrency := make(map[string]bool)
	for _, r := range w.configuration.Markets.Priority.RatesAllow {
		allowCryptoCurrency[r] = true
	}

	rates := FilterRates(fetchedRates, allowCryptoCurrency)

	normalizedRates := toRatesModel(rates)

	if err := w.db.AddRates(normalizedRates); err != nil {
		log.Error(err)
	}
}

func FilterRates(rates []watchmarket.Rate, cryptoCurrency map[string]bool) []watchmarket.Rate {
	result := make([]watchmarket.Rate, 0)
	for _, rate := range rates {
		if rate.Provider == watchmarket.Fixer || cryptoCurrency[rate.Currency] {
			result = append(result, rate)
		}
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
		log.WithFields(log.Fields{"provider": r.GetProvider(), "details": err}).Error("Failed to fetch rates")
		return
	}

	log.WithFields(log.Fields{"provider": r.GetProvider(), "rates": len(rates)}).Info("Rates fetching done")

	s.Add(rates)
}
