package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/trustwallet/watchmarket/services/assets"

	"github.com/trustwallet/golibs/network/middleware"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/postgres"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/worker"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	w             worker.Worker
	configuration config.Configuration
	c             *cron.Cron
)

func init() {
	var err error
	confPath := internal.GetConfigPath(defaultConfigPath)
	configuration, err = config.Init(confPath)
	if err != nil {
		log.Panic("Config read error: ", err)
	}

	err = middleware.SetupSentry(configuration.Sentry.DSN)
	if err != nil {
		log.Error(err)
	}

	assets := assets.Init(configuration.Markets.Assets)

	m, err := markets.Init(configuration, assets)
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

	w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, nil, configuration)
	c = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))

	go postgres.FatalWorker(time.Second*10, *database)
}

func main() {
	w.AddOperation(c, configuration.Worker.Rates, w.FetchAndSaveRates)
	w.AddOperation(c, configuration.Worker.Tickers, w.FetchAndSaveTickers)

	go c.Start()
	go w.FetchAndSaveRates()
	go w.FetchAndSaveTickers()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown worker gracefully...")
	ctx := c.Stop()
	<-ctx.Done()
}
