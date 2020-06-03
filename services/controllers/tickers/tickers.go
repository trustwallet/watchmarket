package tickerscontroller

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
)

type (
	sortedTickersResponse struct {
		sync.Mutex
		tickers []models.Ticker
	}
)

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery, ctx context.Context) (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetTickersByQueries(tickerQueries, ctx)
	if err != nil {
		logger.Error(err, "getTickersByPriority")
		return nil, err
	}
	providers := c.tickersPriority

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res, c.configuration)
	}

	wg.Wait()

	sortedTickers := res.tickers

	result := make(watchmarket.Tickers, len(sortedTickers))

	for i, sr := range sortedTickers {
		result[i] = watchmarket.Ticker{
			Coin:       sr.Coin,
			CoinName:   sr.CoinName,
			CoinType:   watchmarket.CoinType(sr.CoinType),
			LastUpdate: sr.LastUpdated,
			Price: watchmarket.Price{
				Change24h: sr.Change24h,
				Currency:  sr.Currency,
				Provider:  sr.Provider,
				Value:     sr.Value,
			},
			TokenId: sr.TokenId,
		}
	}

	return result, nil
}
