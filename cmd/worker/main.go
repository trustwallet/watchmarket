package main

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
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
	w.AddOperation(c, configuration.Worker.Tickers, w.AlertsIndexer)

	go c.Start()
	go w.FetchAndSaveRates()
	go w.FetchAndSaveTickers()
	go w.AlertsIndexer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown worker gracefully...")
	ctx := c.Stop()
	<-ctx.Done()
}
