package controllers

import (
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Controller struct {
	dataCache        cache.Provider
	database         db.Instance
	chartsPriority   []string
	coinInfoPriority []string
	ratesPriority    []string
	tickersPriority  []string
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	cache cache.Provider,
	database db.Instance,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		cache,
		database,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
		configuration,
	}
}
