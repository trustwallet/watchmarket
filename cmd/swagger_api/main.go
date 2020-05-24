package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/trustwallet/blockatlas/pkg/logger"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"net/http"
)

const (
	defaultPort       = "8423"
	defaultConfigPath = "../../config.yml"
)

var (
	port, confPath string
	engine         *gin.Engine
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)
	configuration := internal.InitConfig(confPath)
	logger.InitLogger()
	engine = internal.InitEngine(configuration.RestAPI.Mode)
}

func main() {
	logger.Info("Loading Swagger API")
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "swagger/index.html")
	})
	engine.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	internal.SetupGracefulShutdown(port, engine)
}
