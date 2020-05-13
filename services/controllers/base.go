package controllers

import (
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
)

type Controller struct {
	cache            cache.Instance
	chartsPriority   priority.Controller
	coinInfoPriority priority.Controller
	ratesPriority    priority.Controller
	tickersPriority  priority.Controller
	api              markets.APIs
}

func NewController(cache cache.Instance, chartsPriority, coinInfoPriority, ratesPriority, tickersPriority priority.Controller, api markets.APIs) Controller {
	return Controller{cache, chartsPriority, coinInfoPriority, ratesPriority, tickersPriority, api}
}
