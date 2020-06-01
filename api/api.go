package api

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/watchmarket/api/endpoint"
	"github.com/trustwallet/watchmarket/api/middleware"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
	"time"
)

func SetupMarketAPI(engine *gin.Engine, controller controllers.Controller) {
	engine.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, `Watchmarket API`) })
	engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	engine.GET("v2/market/ticker/:id",
		middleware.CacheControl(time.Minute, endpoint.GetTickersHandlerV2(controller)))
	engine.GET("v2/market/charts/:id",
		middleware.CacheControl(time.Minute, endpoint.GetChartsHandlerV2(controller)))
	engine.GET("v2/market/info/:id",
		middleware.CacheControl(time.Minute, endpoint.GetCoinInfoHandlerV2(controller)))

	engine.POST("v1/market/ticker",
		middleware.CacheControl(time.Minute, endpoint.GetTickersHandler(controller)))
	engine.GET("v1/market/charts",
		middleware.CacheControl(time.Minute*10, endpoint.GetChartsHandler(controller)))
	engine.GET("v1/market/info",
		middleware.CacheControl(time.Minute*10, endpoint.GetCoinInfoHandler(controller)))
}
