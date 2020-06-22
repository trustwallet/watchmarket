package tickerscontroller

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"strings"
	"sync"
	"time"
)

func createResponse(tr controllers.TickerRequest, tickers watchmarket.Tickers) controllers.TickerResponse {
	mergedTickers := make(watchmarket.Tickers, 0, len(tickers))
	for _, t := range tickers {
		newTicker, ok := findTickerInAssets(tr.Assets, t)
		if !ok {
			continue
		}
		mergedTickers = append(mergedTickers, newTicker)
	}
	return controllers.TickerResponse{Currency: tr.Currency, Tickers: mergedTickers}
}

func createResponseV2(tr controllers.TickerRequestV2, tickers watchmarket.Tickers) controllers.TickerResponseV2 {
	result := controllers.TickerResponseV2{
		Currency: tr.Currency,
	}
	tickersPrices := make([]controllers.TickerPrice, 0, len(tickers))
	for _, ticker := range tickers {
		id, ok := findIDInRequest(tr, watchmarket.BuildID(ticker.Coin, ticker.TokenId))
		if !ok {
			logger.Error("Cannot find ID in request")
		}
		tp := controllers.TickerPrice{
			Change24h: ticker.Price.Change24h,
			Provider:  ticker.Price.Provider,
			Price:     ticker.Price.Value,
			ID:        id,
		}
		tickersPrices = append(tickersPrices, tp)
	}
	result.Tickers = tickersPrices
	return result
}

func makeTickerQueries(coins []controllers.Coin) []models.TickerQuery {
	tickerQueries := make([]models.TickerQuery, 0, len(coins))
	for _, c := range coins {
		tickerQueries = append(tickerQueries, models.TickerQuery{
			Coin:    c.Coin,
			TokenId: strings.ToLower(c.TokenId),
		})
	}
	return tickerQueries
}

func makeTickerQueriesV2(ids []string) []models.TickerQuery {
	tickerQueries := make([]models.TickerQuery, 0, len(ids))
	for _, id := range ids {
		coin, token, err := watchmarket.ParseID(id)
		if err != nil {
			continue
		}
		tickerQueries = append(tickerQueries, models.TickerQuery{
			Coin:    coin,
			TokenId: strings.ToLower(token),
		})
	}
	return tickerQueries
}

func (c Controller) normalizeTickers(tickers watchmarket.Tickers, rate watchmarket.Rate, ctx context.Context) watchmarket.Tickers {
	result := make(watchmarket.Tickers, 0, len(tickers))
	for _, t := range tickers {
		r, ok := c.rateToDefaultCurrency(t, rate, ctx)
		if !ok {
			continue
		}
		if !c.isSuitableUpdateTime(t) {
			continue
		}
		result = append(result, applyRateToTicker(t, r))
	}
	return result
}

func (c Controller) isSuitableUpdateTime(ticker watchmarket.Ticker) bool {
	if ticker.Price.Provider == "coinmarketcap" {
		return true
	}
	now := time.Now().Unix()
	last := ticker.LastUpdate.Unix()
	if now < last {
		return true
	}
	diff := now - last
	if diff < 0 {
		return true
	}
	respectableTime := watchmarket.DurationToUnix(c.configuration.RestAPI.Tickers.RespectableUpdateTime)
	return uint(diff) <= respectableTime
}

func isRespectableMarketCap(provider string, marketCap float64, configuration config.Configuration) bool {
	if provider != "coingecko" {
		return true
	}
	return marketCap >= configuration.RestAPI.Tickers.RespsectableMarketCap
}

func isRespectableVolume(provider string, volume float64, configuration config.Configuration) bool {
	if provider != "coingecko" {
		return true
	}
	return volume >= configuration.RestAPI.Tickers.RespsectableVolume
}

func findIDInRequest(request controllers.TickerRequestV2, id string) (string, bool) {
	for _, i := range request.Ids {
		givenCoin, givenToken, err := watchmarket.ParseID(i)
		if err != nil {
			continue
		}
		coin, token, err := watchmarket.ParseID(id)
		if err != nil {
			continue
		}

		if givenCoin == coin && strings.EqualFold(givenToken, token) {
			return i, true
		}
	}
	return id, false
}

func findTickerInAssets(assets []controllers.Coin, t watchmarket.Ticker) (watchmarket.Ticker, bool) {
	for _, c := range assets {
		if c.Coin == t.Coin && strings.ToLower(c.TokenId) == t.TokenId {
			t.TokenId = c.TokenId
			return t, true
		}
	}
	return watchmarket.Ticker{}, false
}

func findBestProviderForQuery(coin uint, token string, sliceToFind []models.Ticker, providers []string, wg *sync.WaitGroup, res *sortedTickersResponse, configuration config.Configuration) {
	defer wg.Done()
	for _, p := range providers {
		for _, t := range sliceToFind {
			baseCheck := coin == t.Coin && strings.ToLower(token) == t.TokenId

			if baseCheck && t.ShowOption == models.AlwaysShow {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				return
			}
			if baseCheck && p == t.Provider && t.ShowOption != models.NeverShow &&
				isRespectableMarketCap(t.Provider, t.MarketCap, configuration) && isRespectableVolume(t.Provider, t.Volume, configuration) {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				return
			}
		}
	}
}
