package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/blockatlas/pkg/ginutils/gincache"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/storage"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultMaxChartItems = 64
	tokenTWT             = "TWT-8C2"
)

type TickerRequest struct {
	Currency string `json:"currency"`
	Assets   []Coin `json:"assets"`
}

type Coin struct {
	Coin     uint                `json:"coin"`
	CoinType blockatlas.CoinType `json:"type"`
	TokenId  string              `json:"token_id,omitempty"`
}

func SetupMarketAPI(router gin.IRouter, db storage.Market, charts *market.Charts, ac assets.AssetClient) {
	router.Use(ginutils.TokenAuthMiddleware(viper.GetString("market.auth")))
	// Ticker
	router.POST("/ticker", getTickersHandler(db))
	// Charts
	router.GET("/charts", gincache.CacheMiddleware(time.Minute*10, getChartsHandler(charts)))
	router.GET("/info", gincache.CacheMiddleware(time.Minute*5, getCoinInfoHandler(charts, ac)))
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
			ginutils.RenderError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		rate, err := storage.GetRate(strings.ToUpper(md.Currency))
		if err == watchmarket.ErrNotFound {
			ginutils.RenderError(c, http.StatusNotFound, fmt.Sprintf("Currency %s not found", md.Currency))
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve rate", logger.Params{"currency": md.Currency})
			ginutils.RenderError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get rate for %s", md.Currency))
			return
		}

		type getTickerResult struct {
			Error error
			Ticker *watchmarket.Ticker
		}

		getTicker := func(coinRequest Coin, c chan getTickerResult) {
			exchangeRate := rate.Rate
			percentChange := rate.PercentChange24h

			coinObj, ok := coin.Coins[coinRequest.Coin]
			if !ok {
				logger.Warn("Requested coin does not exist", logger.Params{"coin": coinRequest.Coin})
				c <- getTickerResult{Error: watchmarket.ErrNotFound}
				return
			}

			r, err := storage.GetTicker(coinObj.Symbol, strings.ToUpper(coinRequest.TokenId))
			if err != nil {
				if err == watchmarket.ErrNotFound {
					logger.Warn("Ticker not found", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
				} else if err != nil {
					logger.Error(err, "Failed to retrieve ticker", logger.Params{"coin": coinObj.Symbol, "token": coinRequest.TokenId})
				}
				c <- getTickerResult{Error: err}
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

			c <- getTickerResult{Ticker: r}
		}

		ch := make(chan getTickerResult)
		for _, coinRequest := range md.Assets { go getTicker(coinRequest, ch) }
		tickers := make(watchmarket.Tickers, 0)
		for i := 0; i < len(md.Assets); i++ {
			res := <-ch
			if res.Error != nil && res.Error != blockatlas.ErrNotFound {
				ginutils.RenderError(c, http.StatusInternalServerError, "Failed to retrieve tickers")
				return
			}
			tickers = append(tickers, res.Ticker)
		}

		ginutils.RenderSuccess(c, watchmarket.TickerResponse{Currency: md.Currency, Docs: tickers})
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
func getChartsHandler(charts *market.Charts) func(c *gin.Context) {
	return func(c *gin.Context) {
		coinQuery := c.Query("coin")
		if len(coinQuery) == 0 {
			ginutils.RenderError(c, http.StatusBadRequest, "No coin provided")
			return
		}
		if len(c.Query("time_start")) == 0 {
			ginutils.RenderError(c, http.StatusBadRequest, "No time_start provided")
			return
		}

		coinId, err := strconv.Atoi(coinQuery)
		if err != nil {
			ginutils.RenderError(c, http.StatusBadRequest, "Invalid coin provided")
			return
		}
		token := c.Query("token")

		timeStart, err := strconv.ParseInt(c.Query("time_start"), 10, 64)
		if err != nil {
			ginutils.RenderError(c, http.StatusBadRequest, "Invalid time_start provided")
			return
		}
		maxItems, err := strconv.Atoi(c.Query("max_items"))
		if err != nil || maxItems <= 0 {
			maxItems = defaultMaxChartItems
		}

		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)
		chart, err := charts.GetChartData(uint(coinId), token, currency, timeStart, maxItems)
		if err == watchmarket.ErrNotFound {
			ginutils.RenderError(c, http.StatusNotFound, "Chart data not found")
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve chart", logger.Params{"coin": coinId, "currency": currency, "token": token, "time_start": timeStart})
			ginutils.RenderError(c, http.StatusInternalServerError, "Failed to retrieve chart")
			return
		}

		ginutils.RenderSuccess(c, chart)
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
// @Param time_start query int false "Start timestamp" default(1574483028)
// @Param currency query string false "The currency to show coin info in" default(USD)
// @Success 200 {object} watchmarket.ChartCoinInfo
// @Router /v1/market/info [get]
func getCoinInfoHandler(charts *market.Charts, ac assets.AssetClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		coinQuery := c.Query("coin")
		if len(coinQuery) == 0 {
			ginutils.RenderError(c, http.StatusBadRequest, "No coin provided")
			return
		}

		coinId, err := strconv.Atoi(coinQuery)
		if err != nil {
			ginutils.RenderError(c, http.StatusBadRequest, "Invalid coin provided")
			return
		}

		token := c.Query("token")
		currency := c.DefaultQuery("currency", watchmarket.DefaultCurrency)
		chart, err := charts.GetCoinInfo(uint(coinId), token, currency)
		// TODO: cover special casing of TWT token in tests
		if err == watchmarket.ErrNotFound && token != tokenTWT {
			ginutils.RenderError(c, http.StatusNotFound, fmt.Sprintf("Coin info for coin id %d (token: %s, currency: %s) not found", coinId, token, currency))
			return
		} else if err != nil && token != tokenTWT {
			logger.Error(err, "Failed to retrieve coin info", logger.Params{"coinId": coinId, "token": token, "currency": currency})
			ginutils.RenderError(c, http.StatusInternalServerError, "Failed to retrieve coin info")
			return
		}

		chart.Info, err = ac.GetCoinInfo(coinId, token)
		if err == watchmarket.ErrNotFound {
			logger.Warn(err, fmt.Sprintf("Coin assets for coin id %d (token: %s) not found", coinId, token))
			ginutils.RenderSuccess(c, chart)
			return
		} else if err != nil {
			logger.Error(err, "Failed to retrieve coin assets", logger.Params{"coinId": coinId, "token": token})
			ginutils.RenderError(c, http.StatusInternalServerError, "Failed to retrieve coin assets")
			return
		}
		ginutils.RenderSuccess(c, chart)
	}
}
