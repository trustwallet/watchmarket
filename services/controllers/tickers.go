package controllers

import (
	"errors"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"strings"
)

func (c Controller) HandleTickersRequest(tr TickerRequest) (watchmarket.Tickers, error) {
	watchMarketRate, err := c.getRateByPriority(strings.ToUpper(tr.Currency))
	if err != nil {
		return nil, err
	}

	tickers, err := c.getTickersByPriority(makeTickerQueries(tr.Assets), *watchMarketRate)
	if err != nil {
		return nil, err
	}

	// normalize

	return nil, nil
}

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery) (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetTickersByQueries(tickerQueries)
	if err != nil {
		return nil, err
	}
	providers := c.tickersPriority.GetAllProviders()

	tickersMap := make(map[string][]models.Ticker)

	for _, dbTicker := range dbTickers {
		rawCoin := strconv.Itoa(int(dbTicker.Coin))
		key := rawCoin + dbTicker.TokenId
		tickersMap[key] = append(tickersMap[key], dbTicker)
	}

	for _, p := range providers {
		for k, v := range tickersMap {

		}
	}

	return result, nil
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
