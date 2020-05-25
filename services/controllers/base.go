package controllers

import (
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
)

type Controller struct {
	dataCache        cache.Provider
	database         db.Instance
	chartsPriority   priority.Controller
	coinInfoPriority priority.Controller
	ratesPriority    priority.Controller
	tickersPriority  priority.Controller
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	cache cache.Provider,
	database db.Instance,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority priority.Controller,
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
