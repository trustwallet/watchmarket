package controllers

import (
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
)

type Controller struct {
	//database db.Instance
	chartsPriority   priority.Controller
	coinInfoPriority priority.Controller
	ratesPriority    priority.Controller
	tickersPriority  priority.Controller
	api              markets.APIs
}

func NewController(chartsPriority, coinInfoPriority, ratesPriority, tickersPriority priority.Controller, api markets.APIs) Controller {
	return Controller{chartsPriority, coinInfoPriority, ratesPriority, tickersPriority, api}
}
