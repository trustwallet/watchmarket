package worker

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Worker struct {
	ratesApis         markets.RatesAPIs
	ratesUpdateTime   string
	tickersApis       markets.TickersAPIs
	tickersUpdateTime string
	db                db.Instance
}

func Init(
	ratesApis markets.RatesAPIs,
	ratesUpdateTime string,
	tickersApis markets.TickersAPIs,
	tickersUpdateTime string,
	db db.Instance,
) Worker {
	return Worker{
		ratesApis,
		ratesUpdateTime,
		tickersApis,
		tickersUpdateTime,
		db,
	}
}

func (w Worker) addRatesOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if err := c.AddFunc(spec, w.fetchAndSaveRates); err != nil {
		logger.Fatal(err)
	}

	return c
}

func (w Worker) addTickersOperation(c *cron.Cron, updateTime string) *cron.Cron {
	spec := fmt.Sprintf("@every %s", updateTime)

	if err := c.AddFunc(spec, w.fetchAndSaveTickers); err != nil {
		logger.Fatal(err)
	}
	c.Entries()

	return c
}
