package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/postgres"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache/memory"
	rediscache "github.com/trustwallet/watchmarket/services/cache/redis"
	"github.com/trustwallet/watchmarket/services/controllers"
	chartscontroller "github.com/trustwallet/watchmarket/services/controllers/charts"
	infocontroller "github.com/trustwallet/watchmarket/services/controllers/info"
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
	charts         controllers.ChartsController
	info           controllers.InfoController
	w              worker.Worker
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
		configuration.Storage.Postgres.Env,
		configuration.Storage.Postgres.Logs,
	)
	if err != nil {
		logger.Fatal(err)
	}

	memoryCache := memory.Init()
	w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, memoryCache, configuration)

	r := internal.InitRedis(configuration.Storage.Redis)
	cache := rediscache.Init(*r, configuration.RestAPI.Cache)

	charts = chartscontroller.NewController(cache, memoryCache, database, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	info = infocontroller.NewController(cache, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m.ChartsAPIs, configuration)
	tickers = tickerscontroller.NewController(database, memoryCache, ratesPriority, tickerPriority, configuration)
	engine = internal.InitEngine(configuration.RestAPI.Mode)

	go postgres.FatalWorker(time.Second*10, *database)
}

func main() {
	w.SaveRatesToMemory()
	w.SaveTickersToMemory()
	if err := internal.InitAPI(engine, tickers, charts, info, configuration); err != nil {
		panic(err)
	}
	internal.SetupGracefulShutdown(port, engine)
}
