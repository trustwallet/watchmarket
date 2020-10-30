package worker

import (
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
)

type (
	Worker struct {
		ratesApis     markets.RatesAPIs
		tickersApis   markets.TickersAPIs
		db            db.Instance
		cache         cache.Provider
		configuration config.Configuration
	}
)

func Init(
	ratesApis markets.RatesAPIs,
	tickersApis markets.TickersAPIs,
	db db.Instance,
	cache cache.Provider,
	configuration config.Configuration,
) Worker {
	return Worker{
		ratesApis,
		tickersApis,
		db,
		cache,
		configuration,
	}
}

func (w Worker) AddOperation(c *cron.Cron, updateTime string, f func()) {
	spec := fmt.Sprintf("@every %s", updateTime)

	if _, err := c.AddFunc(spec, f); err != nil {
		log.Fatal(err)
	}
}
