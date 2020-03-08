package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/storage"
)

func Bootstrap(engine *gin.Engine, market storage.Market, charts *market.Charts, ac assets.AssetClient) {
	engine.GET("/", GetRoot)
	marketAPI := engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, market, charts, ac)
}