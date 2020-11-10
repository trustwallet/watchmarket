package internal

import (
	"flag"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs-networking/middleware"
	"github.com/trustwallet/watchmarket/api"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/controllers"
	"go.elastic.co/apm/module/apmgin"
	"path/filepath"
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
	rates controllers.RatesController,
	charts controllers.ChartsController,
	info controllers.InfoController,
	alerts controllers.AlertsController,
	configuration config.Configuration,
) {
	api.SetupBasicAPI(engine)
	api.SetupTickersAPI(engine, tickers, configuration.RestAPI.Tickers.CacheControl)
	api.SetupChartsAPI(engine, charts, configuration.RestAPI.Charts.CacheControl)
	api.SetupInfoAPI(engine, info, configuration.RestAPI.Info.CacheControl)
	api.SetupRatesAPI(engine, rates)
	api.SetupAlertsAPI(engine, alerts)
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
	engine.Use(middleware.CORSMiddleware())
	engine.Use(gin.Logger())
	engine.Use(middleware.Prometheus())
	engine.Use(apmgin.Middleware(engine))

	return engine
}
