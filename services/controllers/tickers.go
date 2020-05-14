package controllers

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"strings"
)

func (c Controller) HandleTickersRequest(tr TickerRequest) (watchmarket.Tickers, error) {

	rates, _ := c.database.GetRates(strings.ToUpper(tr.Currency))
	rate := rates[0]

	for _, coinRequest := range tr.Assets {
		exchangeRate := rate.Rate
		percentChange := rate.PercentChange24h

		coin := coinRequest.Coin

	}

	return watchmarket.Tickers{}, nil
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

func buildTickersGroup(coins []Coin) map[string]string {
	tickersMap := make(map[string]string, len(coins))

	for _, c := range coins {
		rawCoin := strconv.Itoa(int(c.Coin))
		tickersMap[rawCoin+c.TokenId] = c.TokenId
	}
	return tickersMap
}

func (c Controller) getRate() {

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
