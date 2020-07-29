package worker

import (
	"context"
	"fmt"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"go.elastic.co/apm"
	"time"
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
	//logger.Info("Fetching Tickers ...")
	//fetchedTickers := fetchTickers(w.tickersApis, ctx)
	//normalizedTickers := toTickersModel(fetchedTickers)
	//
	//if err := w.db.AddTickers(normalizedTickers, w.configuration.Worker.BatchLimit, ctx); err != nil {
	//	logger.Error(err)
	//}
	time.Sleep(time.Minute)
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
