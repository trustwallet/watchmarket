package main

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/api"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
)

const (
	defaultPort       = "8421"
	defaultConfigPath = "../../config.yml"
)

var (
	port, confPath        string
	db                    *storage.Storage
	engine                *gin.Engine
	chartsCachingDuration int64
	cache                 *caching.Provider
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)

	internal.InitConfig(confPath)
	logger.InitLogger()
	tmp := sentrygin.New(sentrygin.Options{})
	sg := &tmp

	redisHost := viper.GetString("storage.redis")
	db = internal.InitRedis(redisHost)
	engine = internal.InitEngine(sg, viper.GetString("gin.mode"))
	cache = internal.InitCaching(db)

}

func main() {
	api.Bootstrap(engine, db, market.InitCharts(), &assets.HttpAssetClient{HttpClient: resty.New()}, cache)
	internal.SetupGracefulShutdown(port, engine)
}
