package main

import (
	"github.com/gin-contrib/cors"
	"github.com/trustwallet/golibs/network/middleware"
	"github.com/trustwallet/watchmarket/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/postgres"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/cache/memory"
	rediscache "github.com/trustwallet/watchmarket/services/cache/redis"
	"github.com/trustwallet/watchmarket/services/controllers"
	chartscontroller "github.com/trustwallet/watchmarket/services/controllers/charts"
	infocontroller "github.com/trustwallet/watchmarket/services/controllers/info"
	ratescontroller "github.com/trustwallet/watchmarket/services/controllers/rates"
	tickerscontroller "github.com/trustwallet/watchmarket/services/controllers/tickers"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/worker"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	port          string
	engine        *gin.Engine
	configuration config.Configuration
	tickers       controllers.TickersController
	rates         controllers.RatesController
	charts        controllers.ChartsController
	info          controllers.InfoController
	w             worker.Worker
	c             *cron.Cron
	memoryCache   cache.Provider
)

func init() {
	var err error
	confPath := internal.GetConfigPath(defaultConfigPath)

	configuration, err = config.Init(confPath)
	if err != nil {
		log.Panic("Config read error: ", err)
	}
	port = configuration.RestAPI.Port
	chartsPriority := configuration.Markets.Priority.Charts
	ratesPriority := configuration.Markets.Priority.Rates
	tickerPriority := configuration.Markets.Priority.Tickers
	coinInfoPriority := configuration.Markets.Priority.CoinInfo

	if err = middleware.SetupSentry(configuration.Sentry.DSN); err != nil {
		log.Error(err)
	}

	a := assets.Init(configuration.Markets.Assets)

	m, err := markets.Init(configuration, a)
	if err != nil {
		log.Fatal(err)
	}

	database, err := postgres.New(
		configuration.Storage.Postgres.Url,
		configuration.Storage.Postgres.Logs,
	)
	if err != nil {
		log.Fatal(err)
	}

	if configuration.RestAPI.UseMemoryCache {
		memoryCache = memory.Init()
		c = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
		w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, memoryCache, configuration)
	} else {
		go postgres.FatalWorker(time.Second*10, *database)
	}

	redisCache, err := rediscache.Init(configuration.Storage.Redis.Url, configuration.RestAPI.Cache)
	if err != nil {
		log.Fatal(err)
	}

	charts = chartscontroller.NewController(redisCache, memoryCache, database, chartsPriority, m.ChartsAPIs, configuration)
	info = infocontroller.NewController(database, memoryCache, coinInfoPriority, ratesPriority, m.ChartsAPIs)
	tickers = tickerscontroller.NewController(database, memoryCache, ratesPriority, tickerPriority, configuration)
	rates = ratescontroller.NewController(database, memoryCache, ratesPriority, configuration)
}

func main() {
	if configuration.RestAPI.UseMemoryCache {
		w.SaveRatesToMemory()
		w.SaveTickersToMemory()

		w.AddOperation(c, configuration.RestAPI.UpdateTime.Rates, w.SaveRatesToMemory)
		w.AddOperation(c, configuration.RestAPI.UpdateTime.Tickers, w.SaveTickersToMemory)

		c.Start()

		log.Info("No items in memory cache")
	}

	gin.SetMode(configuration.RestAPI.Mode)
	engine = gin.New()
	engine.Use(cors.Default())
	engine.Use(middleware.Logger())

	api.SetupBasicAPI(engine)
	api.SetupTickersAPI(engine, tickers, configuration.RestAPI.Tickers.CacheControl)
	api.SetupChartsAPI(engine, charts, configuration.RestAPI.Charts.CacheControl)
	api.SetupInfoAPI(engine, info, configuration.RestAPI.Info.CacheControl)
	api.SetupRatesAPI(engine, rates)
	api.SetupSwaggerAPI(engine)
	internal.SetupGracefulShutdown(port, engine)
}
