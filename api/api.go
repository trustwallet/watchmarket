package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/model"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/api/middleware"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
	"time"
)

func SetupMarketAPI(router gin.IRouter, controller controllers.Controller) {
	router.POST("/ticker",
		middleware.CacheControl(time.Minute, getTickersHandler(controller)))
	router.GET("/charts",
		middleware.CacheControl(time.Minute*10, getChartsHandler(controller)))
	router.GET("/info",
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
		request := controllers.TickerRequest{Currency: watchmarket.DefaultCurrency}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, model.CreateErrorResponse(model.InvalidQuery, errors.E("Invalid request payload")))
			return
		}
		response, err := controller.HandleTickersRequest(request)
		if err != nil {
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
		request := controllers.ChartRequest{
			CoinQuery:    c.Query("coin"),
			Token:        c.Query("token"),
			Currency:     c.DefaultQuery("currency", watchmarket.DefaultCurrency),
			TimeStartRaw: c.Query("time_start"),
			MaxItems:     c.Query("max_items"),
		}
		response, _ := controller.HandleChartsRequest(request)
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
// @Router /v1/market/assets [get]
func getCoinInfoHandler(controller controllers.Controller) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := controllers.DetailsRequest{
			CoinQuery: c.Query("coin"),
			Token:     c.Query("token"),
			Currency:  c.DefaultQuery("currency", watchmarket.DefaultCurrency),
		}
		response, _ := controller.HandleDetailsRequest(request)
		c.JSON(http.StatusOK, response)
	}
}
