package main

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/api"
	"github.com/trustwallet/watchmarket/api/middleware"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/storage"
)

const (
	defaultPort       = "8421"
	defaultConfigPath = "../../config.yml"
)

var (
	port, confPath string
	cache          *storage.Storage
	engine         *gin.Engine
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)

	internal.InitConfig(confPath)
	logger.InitLogger()
	tmp := sentrygin.New(sentrygin.Options{})
	sg := &tmp

	redisHost := viper.GetString("storage.redis")
	cache = internal.InitRedis(redisHost)
	middleware.InitCachingMiddleware(cache)
	engine = internal.InitEngine(sg, viper.GetString("gin.mode"))

}

func main() {
	api.Bootstrap(engine, cache, market.InitCharts(), &assets.HttpAssetClient{HttpClient: resty.New()})
	internal.SetupGracefulShutdown(port, engine)
}
