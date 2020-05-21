package api

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
)

func Bootstrap(engine *gin.Engine, controller controllers.Controller) {
	engine.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, `Watchmarket API`) })
	engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	marketAPI := engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, controller)
}
