package main

import (
	"github.com/robfig/cron"
	"github.com/trustwallet/blockatlas/pkg/logger"
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
	_, confPath := internal.ParseArgs("", defaultConfigPath)

	configuration = internal.InitConfig(confPath)

	assets := internal.InitAssets(configuration.Markets.Assets)

	m, err := markets.Init(configuration, assets)
	if err != nil {
		logger.Fatal(err)
	}

	database, err := postgres.New(configuration.Storage.Postgres)
	if err != nil {
		logger.Fatal(err)
	}

	w = worker.Init(m.RatesAPIs, m.TickersAPIs, database)

	c = cron.New()

	logger.InitLogger()
}

func main() {
	c = w.AddRatesOperation(c, "5m")
	c = w.AddTickersOperation(c, "5m")
	c.Start()
}
