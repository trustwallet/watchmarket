package internal

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/postgres"
	"github.com/trustwallet/watchmarket/services/cache"
	rediscache "github.com/trustwallet/watchmarket/services/cache/redis"
	"time"

	"github.com/trustwallet/blockatlas/api/middleware"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/assets"
	"net/http"
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

func InitAssets(assetsHost string) assets.Client {
	return assets.Init(assetsHost)
}

func InitCache(
	r redis.Redis,
	chartsCaching,
	tickersCaching,
	ratesCaching,
	detailsCaching time.Duration,
) (
	cache.Charts,
	cache.Tickers,
	cache.Rates,
) {
	i := rediscache.Init(r, chartsCaching, tickersCaching, ratesCaching, detailsCaching)
	return cache.Charts(i),
		cache.Tickers(i),
		cache.Rates(i)
}

func InitDB(uri string) db.Instance {
	pg, err := postgres.New(uri)
	if err != nil {
		logger.Fatal(err)
	}
	return pg
}

func InitConfig(confPath string) config.Configuration {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		logger.Fatal(err)
	}

	return config.Init(confPath)
}

func InitEngine(handler *gin.HandlerFunc, ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(middleware.CheckReverseProxy, *handler)
	engine.Use(middleware.CORSMiddleware())
	engine.Use(gin.Logger())

	engine.Use(middleware.Prometheus())

	engine.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
		})
	})

	return engine
}

//func InitCaching(db *storage.Storage, chartsDuration string, chartsInfoDuration string) *cache.Provider {
//	chartsCachingDuration, err := time.ParseDuration(chartsDuration)
//	if err != nil {
//		logger.Warn("Failed to parse charts duration from config, using default value")
//	} else {
//		cache.SetChartsCachingDuration(int64(chartsCachingDuration.Seconds()))
//	}
//	chartsInfoCachingDuration, err := time.ParseDuration(chartsInfoDuration)
//	if err != nil {
//		logger.Warn("Failed to parse charts INFO duration from config, using default value")
//	} else {
//		cache.SetChartsCachingInfoDuration(int64(chartsInfoCachingDuration.Seconds()))
//	}
//	return cache.InitCaching(db)
//}
