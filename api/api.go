package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/storage"
)

func Bootstrap(engine *gin.Engine, market storage.Market) {
	engine.GET("/", GetRoot)
	marketAPI := engine.Group("/v1/market")
	SetupMarketAPI(marketAPI, market)
}