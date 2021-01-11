package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/network/middleware"
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
	unixFile       string
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)
	flag.StringVar(&unixFile, "u", "", "unix file address for api")

	configuration = internal.InitConfig(confPath)
	port = configuration.Port
	chartsPriority := configuration.Markets.Priority.Charts
	ratesPriority := configuration.Markets.Priority.Rates
	tickerPriority := configuration.Markets.Priority.Tickers
	coinInfoPriority := configuration.Markets.Priority.CoinInfo

	err := middleware.SetupSentry(configuration.Sentry.DSN)
	if err != nil {
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

	memoryCache = memory.Init()
	c = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, memoryCache, configuration)

	r := internal.InitRedis(configuration.Storage.Redis.Url)
	redisCache := rediscache.Init(*r, configuration.RestAPI.Cache)

	charts = chartscontroller.NewController(redisCache, memoryCache, database, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	info = infocontroller.NewController(database, memoryCache, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	tickers = tickerscontroller.NewController(database, memoryCache, ratesPriority, tickerPriority, configuration)
	rates = ratescontroller.NewController(database, memoryCache, ratesPriority, configuration)
	engine = internal.InitEngine(configuration.Mode)
}

func main() {
	w.SaveRatesToMemory()
	w.SaveTickersToMemory()

	w.AddOperation(c, configuration.RestAPI.UpdateTime.Rates, w.SaveRatesToMemory)
	w.AddOperation(c, configuration.RestAPI.UpdateTime.Tickers, w.SaveTickersToMemory)

	c.Start()

	internal.InitAPI(engine, tickers, charts, info)

	if len(unixFile) > 0 {
		// heroku specific to start nginx
		os.Create("/tmp/app-initialized")
		if err := engine.RunUnix(unixFile); err != nil {
			log.Fatal(err, "Application failed")
		}
		log.WithFields(log.Fields{"unixFile": ":" + unixFile}).Info("Running application")
	} else {
		if err := engine.Run(":" + port); err != nil {
			log.Fatal(err, "Application failed")
		}
		log.WithFields(log.Fields{"port": ":" + port}).Info("Running application")
	}
}
