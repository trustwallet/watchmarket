package worker

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/logger"
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

	if err := w.db.AddTickers(normalizedTickers, w.configuration.Worker.BatchLimit, ctx); err != nil {
		logger.Error(err)
	}
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
		return
	}

	logger.Info("Tickers fetching done", logger.Params{"provider": t.GetProvider(), "tickers": len(tickers)})

	s.Add(tickers)
}
