package chartscontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

const charts = "charts"

type Controller struct {
	redisCache         cache.Provider
	memoryCache        cache.Provider
	database           db.Instance
	availableProviders []string
	api                markets.ChartsAPIs
	useMemoryCache     bool
}

func NewController(
	redisCache cache.Provider,
	memoryCache cache.Provider,
	database db.Instance,
	chartsPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		redisCache,
		memoryCache,
		database,
		chartsPriority,
		api,
		configuration.RestAPI.UseMemoryCache,
	}
}

// ChartsController interface implementation
func (c Controller) HandleChartsRequest(request controllers.ChartRequest) (chart watchmarket.Chart, err error) {

	if !c.hasTickers(request.Asset) {
		return chart, nil
	}

	chart, err = c.getChartFromRedis(request)
	if err == nil && len(chart.Prices) > 0 {
		return chart, nil
	}

	rawChart, err := c.getChartsFromApi(request)
	if err != nil {
		return watchmarket.Chart{}, errors.New(watchmarket.ErrInternal)
	}

	if len(rawChart.Prices) < 1 {
		return watchmarket.Chart{}, errors.New(watchmarket.ErrNotFound)
	}

	chart = calculateChartByMaxItems(rawChart, request.MaxItems)
	c.putChartsToRedis(chart, request)
	return chart, nil
}

func (c Controller) hasTickers(assetData controllers.Asset) bool {
	var tickers []models.Ticker
	var err error

	if c.useMemoryCache {
		if tickers, err = c.getChartsFromMemory(assetData); err != nil {
			return false
		}
	} else {
		dbTickers, err := c.database.GetTickers([]controllers.Asset{assetData})
		if err != nil {
			return false
		}
		for _, t := range dbTickers {
			if t.ShowOption != 2 { // TODO: 2 to constants
				tickers = append(tickers, t)
			}
		}
	}
	return len(tickers) > 0
}

func (c Controller) getChartsFromApi(data controllers.ChartRequest) (ch watchmarket.Chart, err error) {
	for _, p := range c.availableProviders {
		price, err := c.api[p].GetChartData(data.Asset, data.Currency, data.TimeStart)
		if err == nil && len(price.Prices) > 0 {
			return price, nil
		}
	}
	return watchmarket.Chart{}, nil
}

func (c Controller) getRedisKey(request controllers.ChartRequest) string {
	return c.redisCache.GenerateKey(fmt.Sprintf("%s%d%s%s%d", charts, request.Asset.CoinId, request.Asset.TokenId, request.Currency, request.MaxItems))
}

func (c Controller) getChartFromRedis(request controllers.ChartRequest) (ch watchmarket.Chart, err error) {
	key := c.getRedisKey(request)
	cachedChartRaw, err := c.redisCache.GetWithTime(key, request.TimeStart)
	if err != nil || len(cachedChartRaw) <= 0 {
		return ch, err
	}
	err = json.Unmarshal(cachedChartRaw, &ch)
	return ch, err
}

func (c Controller) putChartsToRedis(chart watchmarket.Chart, request controllers.ChartRequest) {
	key := c.getRedisKey(request)
	chartRaw, err := json.Marshal(&chart)
	if err != nil {
		log.Error(err)
	}

	if err == nil && len(chart.Prices) > 0 {
		err = c.redisCache.SetWithTime(key, chartRaw, request.TimeStart)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("failed to save cache")
		}
	}
}

func (c Controller) getChartsFromMemory(assetData controllers.Asset) ([]models.Ticker, error) {
	key := strings.ToLower(asset.BuildID(assetData.CoinId, assetData.TokenId))
	rawResult, err := c.memoryCache.Get(key)
	if err != nil {
		return nil, err
	}
	var t watchmarket.Ticker
	if err = json.Unmarshal(rawResult, &t); err != nil {
		return nil, err
	}
	result := models.Ticker{
		Coin:        t.Coin,
		CoinName:    t.CoinName,
		CoinType:    string(t.CoinType),
		TokenId:     t.TokenId,
		Currency:    t.Price.Currency,
		Provider:    t.Price.Provider,
		Change24h:   t.Price.Change24h,
		Value:       t.Price.Value,
		LastUpdated: t.LastUpdate,
	}
	return []models.Ticker{result}, nil
}

func calculateChartByMaxItems(chart watchmarket.Chart, maxItems int) watchmarket.Chart {
	var newPrices []watchmarket.ChartPrice
	if len(chart.Prices) > maxItems && maxItems > 0 {
		skip := int(float64(len(chart.Prices) / maxItems))
		i := 0
		for i < len(chart.Prices) {
			newPrices = append(newPrices, chart.Prices[i])
			i += skip + 1
		}
		lastPrice := chart.Prices[len(chart.Prices)-1]
		if len(newPrices) > 0 && lastPrice.Date != newPrices[len(newPrices)-1].Date {
			newPrices = append(newPrices, lastPrice)
		}
	} else {
		newPrices = chart.Prices
	}

	chart.Prices = newPrices
	return chart
}
