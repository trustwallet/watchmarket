package controllers

import (
	"errors"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"strings"
	"sync"
)

func (c Controller) HandleTickersRequest(tr TickerRequest) (watchmarket.Tickers, error) {
	_, err := c.getRateByPriority(strings.ToUpper(tr.Currency))
	if err != nil {
		return nil, err
	}

	tickers, err := c.getTickersByPriority(makeTickerQueries(tr.Assets))
	if err != nil {
		return nil, err
	}

	// normalize

	return tickers, nil
}

func (c Controller) getRateByPriority(currency string) (*watchmarket.Rate, error) {
	rates, err := c.database.GetRates(currency)
	if err != nil {
		return nil, err
	}

	providers := c.tickersPriority.GetAllProviders()

	result := new(models.Rate)
ProvidersLoop:
	for _, p := range providers {
		for _, r := range rates {
			if p == r.Provider {
				result = &r
				break ProvidersLoop
			}
		}
	}
	if result == nil {
		return nil, errors.New("Not found")
	}

	return normalizeRate(*result), nil
}

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery) (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetTickersByQueries(tickerQueries)
	if err != nil {
		return nil, err
	}
	providers := c.tickersPriority.GetAllProviders()

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res)
	}

	wg.Wait()

	sortedTickers := res.tickers

	result := make(watchmarket.Tickers, len(sortedTickers))

	for i, sr := range sortedTickers {
		result[i] = watchmarket.Ticker{
			Coin:       sr.Coin,
			CoinName:   sr.CoinName,
			CoinType:   watchmarket.CoinType(sr.CoinType),
			LastUpdate: sr.UpdatedAt,
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

func findBestProviderForQuery(coin uint, token string, sliceToFind []models.Ticker, providers []string, wg *sync.WaitGroup, res *sortedTickersResponse) {
	for _, p := range providers {
		for _, t := range sliceToFind {
			if coin == t.Coin && strings.ToLower(token) == t.TokenId && p == t.Provider {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				wg.Done()
				return
			}
		}
	}
	wg.Done()
}

func normalizeRate(r models.Rate) *watchmarket.Rate {
	rateStr := strconv.FormatFloat(r.Rate, 'f', 10, 64)
	return &watchmarket.Rate{
		Currency:         rateStr,
		PercentChange24h: r.PercentChange24h,
		Provider:         r.Provider,
		Rate:             r.Rate,
		Timestamp:        r.Timestamp,
	}
}

func makeTickerQueries(coins []Coin) []models.TickerQuery {
	tickerQueries := make([]models.TickerQuery, 0, len(coins))
	for _, c := range coins {
		tickerQueries = append(tickerQueries, models.TickerQuery{
			Coin:    c.Coin,
			TokenId: c.TokenId,
		})
	}
	return tickerQueries
}
