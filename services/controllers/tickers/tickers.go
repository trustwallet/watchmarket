package tickerscontroller

import (
	"encoding/json"
	"github.com/trustwallet/golibs/asset"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	sortedTickersResponse struct {
		sync.Mutex
		tickers []models.Ticker
	}
)

func (c Controller) getBaseTickers() (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetBaseTickers()
	if err != nil {
		log.Error(err, "getBaseTickers")
		return nil, err
	}
	return c.mapTickers(dbTickers)
}

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery) (watchmarket.Tickers, error) {
	if c.configuration.RestAPI.UseMemoryCache {
		var results watchmarket.Tickers
		for _, tr := range tickerQueries {
			key := strings.ToLower(asset.BuildID(tr.Coin, tr.TokenId))
			rawResult, err := c.cache.Get(key)
			if err != nil {
				continue
			}
			var result watchmarket.Ticker
			if err = json.Unmarshal(rawResult, &result); err != nil {
				continue
			}
			results = append(results, result)
		}
		return results, nil
	}

	dbTickers, err := c.database.GetTickersByQueries(tickerQueries)
	if err != nil {
		log.Error(err, "getTickersByPriority")
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
	return c.mapTickers(sortedTickers)
}

func (c Controller) mapTickers(tickers []models.Ticker) (watchmarket.Tickers, error) {
	result := make(watchmarket.Tickers, len(tickers))

	for i, sr := range tickers {
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
