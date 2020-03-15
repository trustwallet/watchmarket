package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
)

type BootstrapProviders struct {
	Engine *gin.Engine
	Market storage.Market
	Charts *market.Charts
	Ac     assets.AssetClient
	Cache  *caching.Provider
}

func Bootstrap(providers BootstrapProviders) {
	providers.Engine.GET("/", func(c *gin.Context) { ginutils.RenderSuccess(c, `Watchmarket API`) })
	marketAPI := providers.Engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, providers)
}
