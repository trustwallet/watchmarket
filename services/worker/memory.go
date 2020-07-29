package worker

import (
	"context"
	"fmt"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"go.elastic.co/apm"
)

func (w Worker) SaveTickersToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveTickersToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	allTickers, err := w.db.GetAllTickers(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(allTickers))
	tickersMap := createTickersMap(allTickers, w.configuration)
	fmt.Println(len(tickersMap))
}

func (w Worker) SaveRatesToMemory() {
	tx := apm.DefaultTracer.StartTransaction("SaveRatesToMemory", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	allRates, err := w.db.GetAllRates(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(allRates))
	ratesMap := createRatesMap(allRates, w.configuration)
	fmt.Println(len(ratesMap))
}

func createTickersMap(allTickers []models.Ticker, configuration config.Configuration) map[string]watchmarket.Ticker {
	m := make(map[string]watchmarket.Ticker, len(allTickers))
	for _, ticker := range allTickers {
		if ticker.ShowOption == models.NeverShow {
			continue
		}
		key := watchmarket.BuildID(ticker.Coin, ticker.TokenId)
		if ticker.ShowOption == models.AlwaysShow {
			m[key] = fromModelToTicker(ticker)
			continue
		}
		baseCheck := watchmarket.IsRespectable(ticker.Provider, ticker.MarketCap, configuration.RestAPI.Tickers.RespsectableMarketCap) &&
			watchmarket.IsRespectable(ticker.Provider, ticker.Volume, configuration.RestAPI.Tickers.RespsectableVolume) &&
			watchmarket.IsSuitableUpdateTime(ticker.LastUpdated, configuration.RestAPI.Tickers.RespectableUpdateTime)

		result, ok := m[key]
		if ok {
			if isHigherPriority(configuration.Markets.Priority.Tickers, result.Price.Provider, ticker.Provider) && baseCheck {
				m[key] = fromModelToTicker(ticker)
			}
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
