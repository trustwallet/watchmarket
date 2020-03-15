package market

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/market/rate"
	"github.com/trustwallet/watchmarket/market/ticker"
	"github.com/trustwallet/watchmarket/storage"
)

type Provider interface {
	Init(storage.Market) error
	GetId() string
	GetLogType() string
	GetUpdateTime() string
}

func scheduleTasks(storage storage.Market, p Provider, c *cron.Cron) {
	err := p.Init(storage)
	if err != nil {
		logger.Error(err, "Init Market Error", logger.Params{"Type": p.GetLogType(), "Market": p.GetId()})
		return
	}
	t := p.GetUpdateTime()
	spec := fmt.Sprintf("@every %s", t)
	logger.Info("Scheduling market data task", logger.Params{
		"Type":     p.GetLogType(),
		"Market":   p.GetId(),
		"Interval": spec,
	})
	_, err = c.AddFunc(spec, func() { go run(storage, p) })
	go run(storage, p)
	if err != nil {
		logger.Error(err, "AddFunc")
	}
}

func run(storage storage.Market, md Provider) {
	logger.Info("Starting market data task...", logger.Params{"Type": md.GetLogType(), "Market": md.GetId()})
	switch m := md.(type) {
	case ticker.TickerProvider:
		runTicker(storage, m)
	case rate.RateProvider:
		runRate(storage, m)
	default:
		logger.Error("Invalid market interface provided")
	}
}
