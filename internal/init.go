package internal

import (
	"flag"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/network/middleware"
	"github.com/trustwallet/watchmarket/api"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/controllers"
)

func ParseArgs(defaultPort, defaultConfigPath string) (string, string) {
	var (
		port, confPath string
	)

	flag.StringVar(&port, "p", defaultPort, "port for api")
	flag.StringVar(&confPath, "c", defaultConfigPath, "config file for api")
	flag.Parse()

	return port, confPath
}

func InitRedis(host string) *redis.Redis {
	c, err := redis.Init(host)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func InitAPI(
	engine *gin.Engine,
	tickers controllers.TickersController,
	charts controllers.ChartsController,
	info controllers.InfoController,
) {
	api.SetupBasicAPI(engine)
	api.SetupTickersAPI(engine, tickers)
	api.SetupChartsAPI(engine, charts)
	api.SetupInfoAPI(engine, info)
	api.SetupSwaggerAPI(engine)
}

func InitAssets(assetsHost string) assets.Client {
	return assets.Init(assetsHost)
}

func InitConfig(confPath string) config.Configuration {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		log.Fatal(err)
	}

	return config.Init(confPath)
}

func InitEngine(ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(middleware.Logger())

	return engine
}
