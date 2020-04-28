package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/api/middleware"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultMaxChartItems = 64
	day                  = 24 * 60 * 60
)

type (
	CoinType string

	TickerRequest struct {
		Currency string `json:"currency"`
		Assets   []Coin `json:"assets"`
	}

	Coin struct {
		Coin     uint     `json:"coin"`
		CoinType CoinType `json:"type"`
		TokenId  string   `json:"token_id,omitempty"`
	}
)

func SetupMarketAPI(router gin.IRouter, provider BootstrapProviders) {
	router.POST("/ticker",
		middleware.CacheControl(time.Minute, getTickersHandler(provider.Market)))
	router.GET("/charts",
		middleware.CacheControl(time.Minute*10, getChartsHandler(provider.Charts, provider.Cache, provider.Market)))
	router.GET("/info",
		middleware.CacheControl(time.Minute*10, getCoinInfoHandler(provider.Charts, provider.Ac, provider.Cache)))
}

// @Summary Get ticker values for a specific market
// @Id get_tickers
// @Description Get the ticker values from many market and coin/token
// @Accept json
// @Produce json
// @Tags Market
// @Param tickers body api.TickerRequest true "Ticker"
// @Success 200 {object} watchmarket.Tickers
// @Router /v1/market/ticker [post]
func getTickersHandler(storage storage.Market) func(c *gin.Context) {
	if storage == nil {
		return nil
	}
	return func(c *gin.Context) {
		md := TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&md); err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
			return
		}

		rate, err := storage.GetRate(strings.ToUpper(md.Currency))
		if err == watchmarket.ErrNotFound {
			c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E(fmt.Sprintf("Currency %s not found", md.Currency))))
			logger.Warn(fmt.Sprintf("Currency %s not found", md.Currency))
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve rate", logger.Params{"currency": md.Currency})
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E(fmt.Sprintf("Failed to get rate for %s", md.Currency))))
			return
		}

		tickers := make(watchmarket.Tickers, 0)
		for _, coinRequest := range md.Assets {
			exchangeRate := rate.Rate
			percentChange := rate.PercentChange24h

			coinObj, ok := coin.Coins[coinRequest.Coin]
			if !ok {
				logger.Warn("Requested coin does not exist", logger.Params{"coin": coinRequest.Coin})
				continue
			}
			r, err := storage.GetTicker(coinObj.Symbol, strings.ToUpper(coinRequest.TokenId))
			if err != nil {
				if err == watchmarket.ErrNotFound {
					logger.Warn("Ticker not found", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
				} else {
					logger.Error(err, "Failed to retrieve ticker", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
					c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Failed to retrieve tickers")))
					return
				}
				continue
			}
			if r.Price.Currency != watchmarket.DefaultCurrency {
				newRate, err := storage.GetRate(strings.ToUpper(r.Price.Currency))
				if err == nil {
					exchangeRate *= newRate.Rate
					percentChange = newRate.PercentChange24h
				} else {
					tickerRate, err := storage.GetTicker(strings.ToUpper(r.Price.Currency), "")
					if err == nil {
						exchangeRate *= tickerRate.Price.Value
						percentChange = big.NewFloat(tickerRate.Price.Change24h)
					}
				}
			}

			r.ApplyRate(md.Currency, exchangeRate, percentChange)
			r.SetCoinId(coinRequest.Coin)
			tickers = append(tickers, r)
		}
		c.JSON(http.StatusOK, watchmarket.TickerResponse{Currency: md.Currency, Docs: tickers})
	}
}

// @Summary Get charts data for a specific coin
// @Id get_charts_data
// @Description Get the charts data from an market and coin/token
// @Accept json
// @Produce json
// @Tags Market
// @Param coin query int true "Coin ID" default(60)
// @Param token query string false "Token ID"
// @Param time_start query int false "Start timestamp" default(1574483028)
// @Param max_items query int false "Max number of items in result prices array" default(64)
// @Param currency query string false "The currency to show charts" default(USD)
// @Success 200 {object} watchmarket.ChartData
// @Router /v1/market/charts [get]
func getChartsHandler(charts *market.Charts, cache *caching.Provider, db storage.Market) func(c *gin.Context) {
	return func(c *gin.Context) {
		coinQuery := c.Query("coin")
		if len(coinQuery) == 0 {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("No coin provided")))
			return
		}

		coinId, err := strconv.Atoi(coinQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid coin provided")))
			return
		}

		timeStart := time.Now().Unix() - day
		if len(c.Query("time_start")) != 0 {
			timeStart, err = strconv.ParseInt(c.Query("time_start"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid time_start provided")))
				return
			}
		}
		token := c.Query("token")

		coinObj, ok := coin.Coins[uint(coinId)]
		if !ok {
			c.JSON(http.StatusOK, watchmarket.ChartData{})
			return
		}

		r, err := db.GetTicker(coinObj.Symbol, strings.ToUpper(token))
		if r == nil || r.Price.Value == 0 || err != nil {
			c.JSON(http.StatusOK, watchmarket.ChartData{})
			return
		}

		maxItemsRaw := c.Query("max_items")
		maxItems, err := strconv.Atoi(maxItemsRaw)
		if err != nil || maxItems <= 0 {
			maxItems = defaultMaxChartItems
		}

		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)

		key := cache.GenerateKey(coinQuery + token + maxItemsRaw + currency)

		chart, err := cache.GetChartsCache(key, timeStart)
		if err == nil {
			c.JSON(http.StatusOK, chart)
			return
		}

		chart, err = charts.GetChartData(uint(coinId), token, currency, timeStart, maxItems)
		if err == watchmarket.ErrNotFound {
			c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E("Chart data not found")))
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve chart", logger.Params{"coin": coinId, "currency": currency, "token": token, "time_start": timeStart})
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Failed to retrieve chart")))
			return
		}

		err = cache.SaveChartsCache(key, chart, timeStart)
		if err != nil {
			logger.Error(err, "Failed to save cache chart", logger.Params{"coin": coinId, "currency": currency, "token": token, "time_start": timeStart, "key": key, "err": err})
		}

		c.JSON(http.StatusOK, chart)
	}
}

// @Summary Get charts coin info data for a specific coin
// @Id get_charts_coin_info
// @Description Get the charts coin info data from an market and coin/contract
// @Accept json
// @Produce json
// @Tags Market
// @Param coin query int true "Coin ID" default(60)
// @Param token query string false "Token ID"
// @Param currency query string false "The currency to show coin info in" default(USD)
// @Success 200 {object} watchmarket.ChartCoinInfo
// @Router /v1/market/info [get]
func getCoinInfoHandler(charts *market.Charts, ac assets.AssetClient, cache *caching.Provider) func(c *gin.Context) {
	return func(c *gin.Context) {
		coinQuery := c.Query("coin")
		if len(coinQuery) == 0 {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("No coin provided")))
			return
		}

		coinId, err := strconv.Atoi(coinQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid coin provided")))
			return
		}

		token := c.Query("token")
		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)

		var chart watchmarket.ChartCoinInfo
		key := cache.GenerateKey(coinQuery + token + currency)
		timeStart := time.Now().Unix() - day

		chart, err = cache.GetCoinInfoCache(key, timeStart)
		if err == nil {
			c.JSON(http.StatusOK, chart)
			return
		}

		chart, err = charts.GetCoinInfo(uint(coinId), token, currency)
		if err == watchmarket.ErrNotFound {
			logger.Info(fmt.Sprintf("Coin info for coin id %d (token: %s, currency: %s) not found", coinId, token, currency))
		} else if err != nil {
			logger.Info(err, "Failed to retrieve coin info", logger.Params{"coinId": coinId, "token": token, "currency": currency})
		}

		chart.Info, err = ac.GetCoinInfo(coinId, token)
		if err == watchmarket.ErrNotFound {
			logger.Warn(err, fmt.Sprintf("Coin assets for coin id %d (token: %s) not found", coinId, token))
			c.JSON(http.StatusOK, chart)
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve coin assets", logger.Params{"coinId": coinId, "token": token})
			c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InvalidQuery, errors.E("Failed to retrieve coin assets")))
			return
		}

		err = cache.SaveCoinInfoCache(key, chart, timeStart)
		if err != nil {
			logger.Error(err, "Failed to save cache info chart", logger.Params{"coin": coinId, "currency": currency, "token": token, "time_start": timeStart, "key": key, "err": err})
		}
		c.JSON(http.StatusOK, chart)
	}
}
