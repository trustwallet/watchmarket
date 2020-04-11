package internal

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/api/middleware"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
	"net/http"
	"path/filepath"
	"time"
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

func InitRedis(host string) *storage.Storage {
	cache := &storage.Storage{DB: &redis.Redis{}}
	err := cache.Init(host)
	if err != nil {
		logger.Fatal(err)
	}
	return cache
}

func InitConfig(confPath string) {
	confPath, err := filepath.Abs(confPath)
	if err != nil {
		logger.Fatal(err)
	}

	config.LoadConfig(confPath)
}

func InitEngine(handler *gin.HandlerFunc, ginMode string) *gin.Engine {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(middleware.CheckReverseProxy, *handler)
	engine.Use(middleware.CORSMiddleware())
	engine.Use(gin.Logger())

	engine.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status": true,
		})
	})

	return engine
}

func InitCaching(db *storage.Storage, chartsDuration string, chartsInfoDuration string) *caching.Provider {
	chartsCachingDuration, err := time.ParseDuration(chartsDuration)
	if err != nil {
		logger.Warn("Failed to parse charts duration from config, using default value")
	} else {
		caching.SetChartsCachingDuration(int64(chartsCachingDuration.Seconds()))
	}
	chartsInfoCachingDuration, err := time.ParseDuration(chartsInfoDuration)
	if err != nil {
		logger.Warn("Failed to parse charts INFO duration from config, using default value")
	} else {
		caching.SetChartsCachingInfoDuration(int64(chartsInfoCachingDuration.Seconds()))
	}
	return caching.InitCaching(db)
}
