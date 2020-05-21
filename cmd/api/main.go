package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/api"
	"github.com/trustwallet/watchmarket/db/postgres"
	_ "github.com/trustwallet/watchmarket/docs"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/services/assets"
	rediscache "github.com/trustwallet/watchmarket/services/cache/redis"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
	"time"
)

const (
	defaultPort       = "8421"
	defaultConfigPath = "config.yml"
)

var (
	port, confPath string
	engine         *gin.Engine
	controller     controllers.Controller
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)

	configuration := internal.InitConfig(confPath)
	logger.InitLogger()

	chartsPriority, err := priority.Init(configuration.Markets.Priority.Charts)
	if err != nil {
		logger.Fatal(err)
	}

	ratesPriority, err := priority.Init(configuration.Markets.Priority.Rates)
	if err != nil {
		logger.Fatal(err)
	}

	tickerPriority, err := priority.Init(configuration.Markets.Priority.Tickers)
	if err != nil {
		logger.Fatal(err)
	}

	coinInfoPriority, err := priority.Init(configuration.Markets.Priority.CoinInfo)
	if err != nil {
		logger.Fatal(err)
	}

	a := assets.Init(configuration.Markets.Assets)

	m, err := markets.Init(configuration, a)
	if err != nil {
		logger.Fatal(err)
	}

	database, err := postgres.New(configuration.Storage.Postgres)
	if err != nil {
		logger.Fatal(err)
	}

	r := internal.InitRedis(configuration.Storage.Redis)

	cache := rediscache.Init(*r, time.Minute, time.Minute, time.Minute, time.Minute)

	controller = controllers.NewController(cache, database, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m)

	engine = internal.InitEngine(configuration.RestAPI.Mode)
}

func main() {
	api.Bootstrap(engine, controller)
	internal.SetupGracefulShutdown(port, engine)
}
