package internal

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/middleware"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
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
		logger.Fatal(err)
	}
	return &c
}

func InitAPI(
	engine *gin.Engine,
	tickers controllers.TickersController,
	charts controllers.ChartsController,
	info controllers.InfoController,
	configuration config.Configuration,
) error {
	var counter int
	for _, a := range configuration.RestAPI.APIs {
		switch a {
		case "base":
			logger.Info("Running base api")
			api.SetupBasicAPI(engine)
			counter++
		case "tickers":
			logger.Info("Running tickers api")
			api.SetupTickersAPI(engine, tickers, configuration.RestAPI.Tickers.CacheControl)
			counter++
		case "charts":
			logger.Info("Running charts api")
			api.SetupChartsAPI(engine, charts, configuration.RestAPI.Charts.CacheControl)
			counter++
		case "info":
			logger.Info("Running info api")
			api.SetupInfoAPI(engine, info, configuration.RestAPI.Info.CacheControl)
			counter++
		case "swagger":
			logger.Info("Running swagger api")
			api.SetupSwaggerAPI(engine)
			counter++
		default:
			continue
		}
	}
	if counter == 0 {
		return errors.E("no apis provided")
	}
	return nil
}

func InitAssets(assetsHost string) assets.Client {
	return assets.Init(assetsHost)
}

func InitConfig(confPath string) config.Configuration {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		logger.Fatal(err)
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
