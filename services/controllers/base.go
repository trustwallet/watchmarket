package controllers

import (
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
)

type Controller struct {
	chartsCache      cache.Charts
	database         db.Instance
	chartsPriority   priority.Controller
	coinInfoPriority priority.Controller
	ratesPriority    priority.Controller
	tickersPriority  priority.Controller
	api              markets.APIs
}

func NewController(
	cache cache.Charts,
	database db.Instance,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority priority.Controller,
	api markets.APIs,
) Controller {
	return Controller{
		cache,
		database,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
	}
}
