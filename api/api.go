package api

import (
	"net/http"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/trustwallet/golibs/network/middleware"
	"github.com/trustwallet/watchmarket/api/endpoint"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/services/controllers"
)

func SetupBasicAPI(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, "Watchmarket API") })
	engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
}

func SetupSwaggerAPI(engine *gin.Engine) {
	engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func SetupInfoAPI(engine *gin.Engine, info controllers.InfoController, d time.Duration) {
	engine.GET("v2/market/info/:id",
		middleware.CacheControl(d, endpoint.GetCoinInfoHandlerV2(info)))

	engine.GET("v1/market/info",
		middleware.CacheControl(d, endpoint.GetCoinInfoHandler(info)))
}

func SetupChartsAPI(engine *gin.Engine, charts controllers.ChartsController, d time.Duration) {
	engine.GET("v2/market/charts/:id",
		middleware.CacheControl(d, endpoint.GetChartsHandlerV2(charts)))

	engine.GET("v1/market/charts",
		middleware.CacheControl(d, endpoint.GetChartsHandler(charts)))
}

func SetupTickersAPI(engine *gin.Engine, tickers controllers.TickersController, d time.Duration) {
	engine.POST("v2/market/tickers",
		middleware.CacheControl(d, endpoint.PostTickersHandlerV2(tickers)))

	engine.GET("v2/market/ticker/:id",
		middleware.CacheControl(d, endpoint.GetTickerHandlerV2(tickers)))

	engine.GET("v2/market/tickers/:assets",
		middleware.CacheControl(d, endpoint.GetTickersHandlerV2(tickers)))

	engine.POST("v1/market/ticker",
		middleware.CacheControl(d, endpoint.GetTickersHandler(tickers)))
}

func SetupRatesAPI(engine *gin.Engine, rates controllers.RatesController) {
	engine.GET("/v1/market/rate", endpoint.GetRate(rates))
	engine.GET("/v1/fiat_rates", endpoint.GetFiatRates(rates))
}
