package worker

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/markets"
)

type (
	Worker struct {
		ratesApis     markets.RatesAPIs
		tickersApis   markets.TickersAPIs
		db            db.Instance
		configuration config.Configuration
	}
)

func Init(
	ratesApis markets.RatesAPIs,
	tickersApis markets.TickersAPIs,
	db db.Instance,
	configuration config.Configuration,
) Worker {
	return Worker{
		ratesApis,
		tickersApis,
		db,
		configuration,
	}
}

func (w Worker) AddRatesOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if _, err := c.AddFunc(spec, w.FetchAndSaveRates); err != nil {
		logger.Fatal(err)
	}

	return c
}

func (w Worker) AddTickersOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if _, err := c.AddFunc(spec, w.FetchAndSaveTickers); err != nil {
		logger.Fatal(err)
	}

	return c
}
