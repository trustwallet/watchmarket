package worker

import (
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
	"sync"
)

func (w Worker) FetchAndSaveTickers() {
	fetchedTickers := fetchTickers(w.tickersApis)
	normalizedTickers := toTickersModel(fetchedTickers)

	if err := w.db.AddTickers(normalizedTickers); err != nil {
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

func fetchTickers(tickersApis markets.TickersAPIs) watchmarket.Tickers {
	wg := new(sync.WaitGroup)
	s := new(tickers)
	for _, t := range tickersApis {
		wg.Add(1)
		go fetchTickersByProvider(t, wg, s)
	}
	wg.Wait()

	return s.tickers
}

func fetchTickersByProvider(t markets.TickersAPI, wg *sync.WaitGroup, s *tickers) {
	defer wg.Done()

	tickers, err := t.GetTickers()
	if err != nil {
		logger.Error(err)
	}

	s.Lock()
	s.tickers = append(s.tickers, tickers...)
	s.Unlock()
}
