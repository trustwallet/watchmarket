package worker

import (
	"context"
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"go.elastic.co/apm"
)

func (w Worker) SaveTickersToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveTickersToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	logger.Info("---------------------------------------")
	logger.Info("Memory Cache: request to DB for tickers ...")
	allTickers, err := w.db.GetAllTickers(ctx)
	if err != nil {
		panic(err)
	}

	logger.Info("Memory Cache: got tickers From DB", logger.Params{"len": len(allTickers)})
	tickersMap := createTickersMap(allTickers, w.configuration)
	for key, val := range tickersMap {
		rawVal, err := json.Marshal(val)
		if err != nil {
			logger.Error(err)
			continue
		}
		if err = w.cache.Set(key, rawVal, ctx); err != nil {
			logger.Error(err)
			continue
		}
	}
	logger.Info("Memory Cache: tickers saved to the cache", logger.Params{"len": len(tickersMap)})
	logger.Info("---------------------------------------")
}

func (w Worker) SaveRatesToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveRatesToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()
	logger.Info("---------------------------------------")
	logger.Info("Memory Cache: request to DB for rates ...")
	allRates, err := w.db.GetAllRates(ctx)
	if err != nil {
		panic(err)
	}

	logger.Info("Memory Cache: got rates From DB", logger.Params{"len": len(allRates)})
	ratesMap := createRatesMap(allRates, w.configuration)
	for key, val := range ratesMap {
		rawVal, err := json.Marshal(val)
		if err != nil {
			logger.Error(err)
			continue
		}
		if err = w.cache.Set(key, rawVal, ctx); err != nil {
			logger.Error(err)
			continue
		}
	}
	logger.Info("Memory Cache: rates saved to the cache", logger.Params{"len": len(ratesMap)})
	logger.Info("---------------------------------------")
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
			if isHigherPriority(configuration.Markets.Priority.Tickers, result.Price.Provider, ticker.Provider) && baseCheck {
				m[key] = fromModelToTicker(ticker)
			}
			continue
		}
		if baseCheck {
			m[key] = fromModelToTicker(ticker)
		}
	}
	return m
}

func createRatesMap(allRates []models.Rate, configuration config.Configuration) map[string]watchmarket.Rate {
	m := make(map[string]watchmarket.Rate, len(allRates))
	for _, rate := range allRates {
		key := rate.Currency
		if rate.Provider != "fixer" {
			if !watchmarket.IsFiatRate(key) {
				continue
			}
		}
		result, ok := m[key]
		if ok {
			if isHigherPriority(configuration.Markets.Priority.Rates, result.Provider, rate.Provider) {
				m[key] = fromModelToRate(rate)
			}
		}
		m[key] = fromModelToRate(rate)
	}
	return m
}

func fromModelToTicker(m models.Ticker) watchmarket.Ticker {
	return watchmarket.Ticker{
		Coin:       m.Coin,
		CoinName:   m.CoinName,
		CoinType:   watchmarket.CoinType(m.CoinType),
		LastUpdate: m.LastUpdated,
		Price: watchmarket.Price{
			Change24h: m.Change24h,
			Currency:  m.Currency,
			Provider:  m.Provider,
			Value:     m.Value,
		},
		TokenId: m.TokenId,
	}
}

func fromModelToRate(m models.Rate) watchmarket.Rate {
	return watchmarket.Rate{
		Currency:         m.Currency,
		PercentChange24h: m.PercentChange24h,
		Provider:         m.Provider,
		Rate:             m.Rate,
		Timestamp:        m.LastUpdated.Unix(),
	}
}
