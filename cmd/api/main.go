package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
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
	"github.com/trustwallet/watchmarket/services/controllers/rates"
	tickerscontroller "github.com/trustwallet/watchmarket/services/controllers/tickers"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/worker"
	"time"
)

const (
	defaultPort       = "8420"
	defaultConfigPath = "../../config.yml"
)

var (
	port, confPath string
	engine         *gin.Engine
	configuration  config.Configuration
	tickers        controllers.TickersController
	rates          controllers.RatesController
	charts         controllers.ChartsController
	info           controllers.InfoController
	w              worker.Worker
	c              *cron.Cron
	memoryCache    cache.Provider
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)

	configuration = internal.InitConfig(confPath)
	logger.InitLogger()
	port = configuration.RestAPI.Port
	chartsPriority := configuration.Markets.Priority.Charts
	ratesPriority := configuration.Markets.Priority.Rates
	tickerPriority := configuration.Markets.Priority.Tickers
	coinInfoPriority := configuration.Markets.Priority.CoinInfo
	a := assets.Init(configuration.Markets.Assets)

	m, err := markets.Init(configuration, a)
	if err != nil {
		logger.Fatal(err)
	}

	database, err := postgres.New(
		configuration.Storage.Postgres.Uri,
		configuration.Storage.Postgres.APM,
		configuration.Storage.Postgres.Logs,
	)
	if err != nil {
		logger.Fatal(err)
	}

	if configuration.RestAPI.UseMemoryCache {
		memoryCache = memory.Init()
		c = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
		w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, memoryCache, configuration)
	} else {
		go postgres.FatalWorker(time.Second*10, *database)
	}

	r := internal.InitRedis(configuration.Storage.Redis)
	redisCache := rediscache.Init(*r, configuration.RestAPI.Cache)

	charts = chartscontroller.NewController(redisCache, memoryCache, database, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	info = infocontroller.NewController(redisCache, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	tickers = tickerscontroller.NewController(database, memoryCache, ratesPriority, tickerPriority, configuration)
	rates = ratescontroller.NewController(database, memoryCache, ratesPriority, configuration)
	engine = internal.InitEngine(configuration.RestAPI.Mode)
}

func main() {
	if configuration.RestAPI.UseMemoryCache {
		w.SaveRatesToMemory()
		w.SaveTickersToMemory()

		w.AddOperation(c, configuration.RestAPI.UpdateTime.Rates, w.SaveRatesToMemory)
		w.AddOperation(c, configuration.RestAPI.UpdateTime.Tickers, w.SaveTickersToMemory)

		c.Start()

		if memoryCache.GetLenOfSavedItems() <= 0 {
			panic("no items in memory cache")
		}
	}

	internal.InitAPI(engine, tickers, rates, charts, info, configuration)
	internal.SetupGracefulShutdown(port, engine)
}
