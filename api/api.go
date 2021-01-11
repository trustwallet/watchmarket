package api

import (
	"net/http"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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

func SetupInfoAPI(engine *gin.Engine, info controllers.InfoController) {
	engine.GET("v2/market/info/:id", endpoint.GetCoinInfoHandlerV2(info))
	engine.GET("v1/market/info", endpoint.GetCoinInfoHandler(info))
}

func SetupChartsAPI(engine *gin.Engine, charts controllers.ChartsController) {
	engine.GET("v2/market/charts/:id", endpoint.GetChartsHandlerV2(charts))
	engine.GET("v1/market/charts", endpoint.GetChartsHandler(charts))
}

func SetupTickersAPI(engine *gin.Engine, tickers controllers.TickersController) {
	engine.POST("v2/market/tickers", endpoint.PostTickersHandlerV2(tickers))
	engine.GET("v2/market/ticker/:id", endpoint.GetTickerHandlerV2(tickers))
	engine.GET("v2/market/tickers/:assets", endpoint.GetTickersHandlerV2(tickers))
	engine.POST("v1/market/ticker", endpoint.GetTickersHandler(tickers))
}
