package api

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
	"net/http"
)

type BootstrapProviders struct {
	Engine *gin.Engine
	Market storage.Market
	Charts *market.Charts
	Ac     assets.AssetClient
	Cache  *caching.Provider
}

func Bootstrap(providers BootstrapProviders) {
	providers.Engine.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, `Watchmarket API`) })
	providers.Engine.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	marketAPI := providers.Engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, providers)
}
