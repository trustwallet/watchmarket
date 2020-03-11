package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/storage"
)

func Bootstrap(engine *gin.Engine, market storage.Market, charts *market.Charts, ac assets.AssetClient) {
	engine.GET("/", func(c *gin.Context) { ginutils.RenderSuccess(c, `Watchmarket API`) })
	marketAPI := engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, market, charts, ac)
}
