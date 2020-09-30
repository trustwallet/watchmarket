package chartscontroller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Controller struct {
	redisCache       cache.Provider
	memoryCache      cache.Provider
	database         db.Instance
	chartsPriority   []string
	coinInfoPriority []string
	ratesPriority    []string
	tickersPriority  []string
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	redisCache cache.Provider,
	memoryCache cache.Provider,
	database db.Instance,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		redisCache,
		memoryCache,
		database,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
		configuration,
	}
}

func (c Controller) HandleChartsRequest(cr controllers.ChartRequest, ctx context.Context) (watchmarket.Chart, error) {
	var ch watchmarket.Chart

	verifiedData, err := toChartsRequestData(cr)
	if err != nil {
		return ch, errors.New(watchmarket.ErrBadRequest)
	}

	key := c.redisCache.GenerateKey(charts + cr.CoinQuery + cr.Token + cr.Currency + cr.MaxItems)

	cachedChartRaw, err := c.redisCache.GetWithTime(key, verifiedData.TimeStart, ctx)
	if err == nil && len(cachedChartRaw) > 0 {
		err = json.Unmarshal(cachedChartRaw, &ch)
		if err == nil && len(ch.Prices) > 0 {
			return ch, nil
		}
	}

	res, err := c.checkTickersAvailability(verifiedData.Coin, verifiedData.Token, ctx)
	if err != nil || len(res) == 0 {
		return ch, err
	}

	rawChart, err := c.getChartsByPriority(verifiedData, ctx)
	if err != nil {
		return watchmarket.Chart{}, errors.New(watchmarket.ErrInternal)
	}

	if len(rawChart.Prices) < 1 {
		return watchmarket.Chart{}, errors.New(watchmarket.ErrNotFound)
	}

	chart := normalizeChart(rawChart, verifiedData.MaxItems)

	chartRaw, err := json.Marshal(&chart)
	if err != nil {
		logger.Error(err)
	}

	if err == nil && len(chart.Prices) > 0 {
		err = c.redisCache.SetWithTime(key, chartRaw, verifiedData.TimeStart, ctx)
		if err != nil {
			logger.Error("failed to save cache", logger.Params{"err": err})
		}
	}

	return chart, nil
}
