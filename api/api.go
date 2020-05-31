package api

import (
	"context"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/api/middleware"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm"
	"net/http"
	"time"
)

func SetupMarketAPI(engine *gin.Engine, controller controllers.Controller) {
	engine.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, `Watchmarket API`) })
	engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	engine.POST("v1/market/ticker",
		middleware.CacheControl(time.Minute, getTickersHandler(controller)))
	engine.GET("v1/market/charts",
		middleware.CacheControl(time.Minute*10, getChartsHandler(controller)))
	engine.GET("v1/market/info",
		middleware.CacheControl(time.Minute*10, getCoinInfoHandler(controller)))
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
func getTickersHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v1/market/ticker", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
			return
		}
		response, err := controller.HandleTickersRequest(request, ctx)
		if err != nil || len(response.Tickers) == 0 {
			handleTickersError(c, request)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get charts data for a specific coin
// @Id get_charts_data
// @Description Get the charts data from an market and coin/token
// @Accept json
// @Produce json
// @Tags Market
// @Param coin query int true "Coin id" default(60)
// @Param token query string false "Token id"
// @Param time_start query int false "Start timestamp" default(1574483028)
// @Param max_items query int false "Max number of items in result prices array" default(64)
// @Param currency query string false "The currency to show charts" default(USD)
// @Success 200 {object} watchmarket.ChartData
// @Router /v1/market/charts [get]
func getChartsHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("GET /v1/market/charts", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.ChartRequest{
			CoinQuery:    c.Query("coin"),
			Token:        c.Query("token"),
			Currency:     c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			TimeStartRaw: c.Query("time_start"),
			MaxItems:     c.Query("max_items"),
		}

		response, err := controller.HandleChartsRequest(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get charts coin assets data for a specific coin
// @Id get_charts_coin_info
// @Description Get the charts coin assets data from an market and coin/contract
// @Accept json
// @Produce json
// @Tags Market
// @Param coin query int true "Coin id" default(60)
// @Param token query string false "Token id"
// @Param currency query string false "The currency to show coin assets in" default(USD)
// @Success 200 {object} watchmarket.ChartCoinInfo
// @Router /v1/market/info [get]
func getCoinInfoHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		tx := apm.DefaultTracer.StartTransaction("POST /v1/market/info", "request")
		ctx := apm.ContextWithTransaction(context.Background(), tx)
		defer tx.End()

		request := controllers.DetailsRequest{
			CoinQuery: c.Query("coin"),
			Token:     c.Query("token"),
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, err := controller.HandleDetailsRequest(request, ctx)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func handleError(c *gin.Context, err error) {
	switch err.Error() {
	case controllers.ErrInternal:
		c.JSON(http.StatusInternalServerError, model.CreateErrorResponse(model.InternalFail, errors.E("Internal Fail")))
		return
	case controllers.ErrBadRequest:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
		return
	case controllers.ErrNotFound:
		c.JSON(http.StatusNotFound, model.CreateErrorResponse(model.RequestedDataNotFound, errors.E("Not found")))
		return
	default:
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
	}
}

func handleTickersError(c *gin.Context, req controllers.TickerRequest) {
	if len(req.Assets) == 0 || req.Assets == nil {
		c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
		return
	}
	emptyResponse := controllers.TickerResponse{
		Currency: req.Currency,
	}
	tickers := make(watchmarket.Tickers, 0, len(req.Assets))
	for _, t := range req.Assets {
		tickers = append(tickers, watchmarket.Ticker{
			Coin:     t.Coin,
			TokenId:  t.TokenId,
			CoinType: t.CoinType,
		})
	}
	emptyResponse.Tickers = tickers
	c.JSON(http.StatusOK, emptyResponse)
}
