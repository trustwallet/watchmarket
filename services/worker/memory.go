package worker

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"go.elastic.co/apm"
)

func (w Worker) SaveTickersToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveTickersToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	log.Info("---------------------------------------")
	log.Info("Memory Cache: request to DB for tickers ...")

	allTickers, err := w.db.GetAllTickers(ctx)
	if err != nil {
		log.Warn("Failed to get tickers: ", err.Error())
		return
	}

	log.WithFields(log.Fields{"len": len(allTickers)}).Info("Memory Cache: got tickers From DB")
	tickersMap := createTickersMap(allTickers, w.configuration)
	for key, val := range tickersMap {
		rawVal, err := json.Marshal(val)
		if err != nil {
			log.Error(err)
			continue
		}
		if err = w.cache.Set(key, rawVal, ctx); err != nil {
			log.Error(err)
			continue
		}
	}

	log.WithFields(log.Fields{"len": len(tickersMap)}).Info("Memory Cache: tickers saved to the cache")
	log.Info("---------------------------------------")
}

func (w Worker) SaveRatesToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveRatesToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	log.Info("---------------------------------------")
	log.Info("Memory Cache: request to DB for rates ...")

	allRates, err := w.db.GetAllRates(ctx)
	if err != nil {
		log.Warn("Failed to get rates: ", err.Error())
		return
	}

	log.WithFields(log.Fields{"len": len(allRates)}).Info("Memory Cache: got rates From DB")

	ratesMap := createRatesMap(allRates, w.configuration)
	for key, val := range ratesMap {
		rawVal, err := json.Marshal(val)
		if err != nil {
			log.Error(err)
			continue
		}
		if err = w.cache.Set(key, rawVal, ctx); err != nil {
			log.Error(err)
			continue
		}
	}

	log.WithFields(log.Fields{"len": len(ratesMap)}).Info("Memory Cache: rates saved to the cache")
	log.Info("---------------------------------------")
}

func createTickersMap(allTickers []models.Ticker, configuration config.Configuration) map[string]watchmarket.Ticker {
	m := make(map[string]watchmarket.Ticker, len(allTickers))
	for _, ticker := range allTickers {
		if ticker.ShowOption == models.NeverShow {
			continue
		}
		key := ticker.ID
		if ticker.ShowOption == models.AlwaysShow {
			m[key] = fromModelToTicker(ticker)
			continue
		}
		baseCheck :=
			(watchmarket.IsRespectableValue(ticker.MarketCap, configuration.RestAPI.Tickers.RespsectableMarketCap) || ticker.Provider != "coingecko") &&
				(watchmarket.IsRespectableValue(ticker.Volume, configuration.RestAPI.Tickers.RespsectableVolume) || ticker.Provider != "coingecko") &&
				watchmarket.IsSuitableUpdateTime(ticker.LastUpdated, configuration.RestAPI.Tickers.RespectableUpdateTime)

		result, ok := m[key]
		if ok {
			if isHigherPriority(configuration.Markets.Priority.Tickers, result.Price.Provider, ticker.Provider) &&
				baseCheck && result.ShowOption != models.AlwaysShow {
				m[key] = fromModelToTicker(ticker)
			}
			continue
		} else if baseCheck {
			m[key] = fromModelToTicker(ticker)
		}
	}
	return m
}

func createRatesMap(allRates []models.Rate, configuration config.Configuration) map[string]watchmarket.Rate {
	m := make(map[string]watchmarket.Rate, len(allRates))
	for _, rate := range allRates {
		if rate.ShowOption == models.NeverShow {
			continue
		}
		key := rate.Currency
		if rate.ShowOption == models.AlwaysShow {
			m[key] = fromModelToRate(rate)
			continue
		}
		if rate.Provider == "fixer" {
			if !watchmarket.IsFiatRate(key) {
				continue
			}
		}
		result, ok := m[key]
		if ok {
			if isHigherPriority(configuration.Markets.Priority.Rates, result.Provider, rate.Provider) && result.ShowOption != models.AlwaysShow {
				m[key] = fromModelToRate(rate)
			}
		} else {
			m[key] = fromModelToRate(rate)
		}
	}
	return m
}
