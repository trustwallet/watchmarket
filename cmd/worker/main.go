package main

import (
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/postgres"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/worker"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	_, confPath := internal.ParseArgs("", defaultConfigPath)
	configuration = internal.InitConfig(confPath)
	assets := internal.InitAssets(configuration.Markets.Assets)

	m, err := markets.Init(configuration, assets)
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

	w = worker.Init(m.RatesAPIs, m.TickersAPIs, database, nil, configuration)
	c = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	logger.InitLogger()

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
	logger.Info("Shutdown worker gracefully...")
	ctx := c.Stop()
	<-ctx.Done()
}
