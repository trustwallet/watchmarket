package storage

import (
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

const (
	EntityRates  = "ATLAS_MARKET_RATES"
	EntityQuotes = "ATLAS_MARKET_QUOTES"
)

type ProviderList interface {
	GetPriority(providerId string) int
}

func (s *Storage) SaveTicker(coin *watchmarket.Ticker, pl ProviderList) (SaveResult, error) {
	cd, err := s.GetTicker(coin.CoinName, coin.TokenId)
	if err != nil && err != watchmarket.ErrNotFound {
		return SaveResultStorageFailure, err
	} else if err == nil {
		op := pl.GetPriority(cd.Price.Provider)
		np := pl.GetPriority(coin.Price.Provider)
		if op != -1 && np > op {
			logger.Debug("Skipping new ticker as its priority is lower than the existing record", logger.Params{
				"oldTickerPriority": op,
				"newTickerPriority": np,
			})
			return SaveResultSkippedLowPriority, nil
		}

		if cd.LastUpdate.After(coin.LastUpdate) && op >= np {
			logger.Debug("Skipping new ticker as its priority is lower than the existing record or its Timestamp is older", logger.Params{
				"oldTickerTime":     cd.LastUpdate,
				"newTickerTime":     coin.LastUpdate,
				"oldTickerPriority": op,
				"newTickerPriority": np,
			})
			return SaveResultSkippedLowPriorityOrOutdated, nil
		}
	}

	hm := createHashMap(coin.CoinName, coin.TokenId)
	err = s.AddHM(EntityQuotes, hm, coin)
	if err != nil {
		logger.Error(err, "SaveTicker")
		return SaveResultStorageFailure, err
	}

	return SaveResultSuccess, nil
}

func (s *Storage) GetTicker(coin, token string) (*watchmarket.Ticker, error) {
	hm := createHashMap(coin, token)
	var cd *watchmarket.Ticker
	err := s.GetHMValue(EntityQuotes, hm, &cd)
	if err != nil {
		return nil, err
	}
	return cd, nil
}

func (s *Storage) SaveRates(rates watchmarket.Rates, pl ProviderList) map[SaveResult]int {
	results := make(map[SaveResult]int)
	for _, rate := range rates {
		r, err := s.GetRate(rate.Currency)
		if err != nil && err != watchmarket.ErrNotFound {
			logger.Error(err, "SaveRates")
			results[SaveResultStorageFailure]++
			continue
		}
		if err == nil {
			op := pl.GetPriority(r.Provider)
			np := pl.GetPriority(rate.Provider)
			if op != -1 && np > op {
				results[SaveResultSkippedLowPriority]++
				continue
			}

			if rate.Timestamp < r.Timestamp && op >= np {
				results[SaveResultSkippedLowPriorityOrOutdated]++
				continue
			}
		}
		err = s.AddHM(EntityRates, rate.Currency, &rate)
		if err != nil {
			logger.Error(err, "SaveRates")
			results[SaveResultStorageFailure]++
			continue
		}
		results[SaveResultSuccess]++
	}
	return results
}

func (s *Storage) GetRate(currency string) (rate *watchmarket.Rate, err error) {
	err = s.GetHMValue(EntityRates, currency, &rate)
	return
}

func createHashMap(coin, token string) string {
	if len(token) == 0 {
		return strings.ToUpper(coin)
	}
	return strings.ToUpper(strings.Join([]string{coin, token}, "_"))
}
