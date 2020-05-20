package worker

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/markets"
)

type (
	Worker struct {
		ratesApis   markets.RatesAPIs
		tickersApis markets.TickersAPIs
		db          db.Instance
	}
)

func Init(
	ratesApis markets.RatesAPIs,
	tickersApis markets.TickersAPIs,
	db db.Instance,
) Worker {
	return Worker{
		ratesApis,
		tickersApis,
		db,
	}
}

func (w Worker) AddRatesOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if err := c.AddFunc(spec, w.fetchAndSaveRates); err != nil {
		logger.Fatal(err)
	}

	return c
}

func (w Worker) AddTickersOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if err := c.AddFunc(spec, w.fetchAndSaveTickers); err != nil {
		logger.Fatal(err)
	}
	c.Entries()

	return c
}
