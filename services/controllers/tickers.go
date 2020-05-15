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
	watchMarketRate, err := c.getRateByPriority(strings.ToUpper(tr.Currency))
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

type tickersRes struct {
	sync.Mutex
	tickers []models.Ticker
}

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery) (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetTickersByQueries(tickerQueries)
	if err != nil {
		return nil, err
	}
	providers := c.tickersPriority.GetAllProviders()

	res := new(tickersRes)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res)
	}

	wg.Wait()

	sortedTickers := res.tickers
	result := make(watchmarket.Tickers, len(sortedTickers))

	for _, sr := range sortedTickers {
		result = append(result, watchmarket.Ticker{
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
		})
	}

	return result, nil
}

///


func findBestProviderForQuery(coin uint, token string, sliceToFind []models.Ticker, providers []string, wg *sync.WaitGroup, res *tickersRes) {
	for _, p := range providers {
	ProvidersLoop:
		for _, t := range sliceToFind {
			if coin == t.Coin && strings.ToLower(token) == t.TokenId && p == t.Provider {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				break ProvidersLoop
			}
		}
	}
	wg.Done()
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

// пройтись по tickers
// для них получить rates
// normalize
// return

// пройтись по TickerRequest
// формируем rates req to db & tickers req to db
// db req
// data
// map tickers
// normalize
// return

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

//rate, err := storage.GetRate(strings.ToUpper(md.Currency))
//		if err == watchmarket.ErrNotFound {
//			c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E(fmt.Sprintf("Currency %s not found", md.Currency))))
//			logger.Warn(fmt.Sprintf("Currency %s not found", md.Currency))
//			return
//		} else if err != nil {
//			logger.Error(err, "Failed to retrieve rate", logger.Params{"currency": md.Currency})
//			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E(fmt.Sprintf("Failed to get rate for %s", md.Currency))))
//			return
//		}
//
//		tickers := make(watchmarket.Tickers, 0)
//		for _, coinRequest := range md.Assets {
//			exchangeRate := rate.Rate
//			percentChange := rate.PercentChange24h
//
//			coinObj, err := getCoinObj(coinRequest.Coin)
//			if err != nil {
//				logger.Warn("Requested coin does not exist", logger.Params{"coin": coinRequest.Coin})
//				continue
//			}
//
//			r, err := storage.GetTicker(coinObj.Symbol, strings.ToUpper(coinRequest.TokenId))
//			if err != nil {
//				if err == watchmarket.ErrNotFound {
//					logger.Warn("Ticker not found", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
//				} else {
//					logger.Error(err, "Failed to retrieve ticker", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
//					c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Failed to retrieve tickers")))
//					return
//				}
//				continue
//			}
//			if r.Price.Currency != watchmarket.DefaultCurrency {
//				newRate, err := storage.GetRate(strings.ToUpper(r.Price.Currency))
//				if err == nil {
//					exchangeRate *= newRate.Rate
//					percentChange = newRate.PercentChange24h
//				} else {
//					tickerRate, err := storage.GetTicker(strings.ToUpper(r.Price.Currency), "")
//					if err == nil {
//						exchangeRate *= tickerRate.Price.Value
//						percentChange = big.NewFloat(tickerRate.Price.Change24h)
//					}
//				}
//			}
//
//			r.ApplyRate(md.Currency, exchangeRate, percentChange)
//			r.SetCoinId(coinRequest.Coin)
//			tickers = append(tickers, r)
//		}
