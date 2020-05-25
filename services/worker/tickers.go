package worker

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
	"go.elastic.co/apm"
	"sync"
)

func (w Worker) FetchAndSaveTickers() {
	tx := apm.DefaultTracer.StartTransaction("FetchAndSaveTickers", "app")
	ctx := apm.ContextWithTransaction(context.Background(), tx)
	defer tx.End()

	logger.Info("Fetching Tickers ...")
	fetchedTickers := fetchTickers(w.tickersApis, ctx)
	normalizedTickers := toTickersModel(fetchedTickers)

	if err := w.db.AddTickers(normalizedTickers, ctx); err != nil {
		logger.Error(err)
	}
}

func toTickersModel(tickers watchmarket.Tickers) []models.Ticker {
	result := make([]models.Ticker, 0, len(tickers))
	for _, t := range tickers {
		result = append(result, models.Ticker{
			Coin:        t.Coin,
			CoinName:    t.CoinName,
			CoinType:    string(t.CoinType),
			TokenId:     t.TokenId,
			Change24h:   t.Price.Change24h,
			Currency:    t.Price.Currency,
			Provider:    t.Price.Provider,
			Value:       t.Price.Value,
			Volume:      t.Volume,
			MarketCap:   t.MarketCap,
			LastUpdated: t.LastUpdate,
		})
	}

	return result
}

func fetchTickers(tickersApis markets.TickersAPIs, ctx context.Context) watchmarket.Tickers {
	wg := new(sync.WaitGroup)
	s := new(tickers)
	for _, t := range tickersApis {
		wg.Add(1)
		go fetchTickersByProvider(t, wg, s, ctx)
	}
	wg.Wait()

	return s.tickers
}

func fetchTickersByProvider(t markets.TickersAPI, wg *sync.WaitGroup, s *tickers, ctx context.Context) {
	defer wg.Done()

	tickers, err := t.GetTickers(ctx)
	if err != nil {
		logger.Error("Failed to fetch tickers", logger.Params{"provider": t.GetProvider(), "details": err})
	}

	logger.Info("Tickers fetching done", logger.Params{"provider": t.GetProvider(), "tickers": len(tickers)})

	s.Lock()
	s.tickers = append(s.tickers, tickers...)
	s.Unlock()
}
