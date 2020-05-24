package controllers

import (
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (c Controller) HandleChartsRequest(cr ChartRequest) (watchmarket.Chart, error) {
	var ch watchmarket.Chart

	verifiedData, err := toChartsRequestData(cr)
	if err != nil {
		return ch, errors.New(ErrBadRequest)
	}

	key := c.dataCache.GenerateKey(cr.CoinQuery + cr.Token + cr.Currency + cr.MaxItems)

	cachedChartRaw, err := c.dataCache.GetWithTime(key, verifiedData.TimeStart)
	if err == nil && len(cachedChartRaw) > 0 {
		err = json.Unmarshal(cachedChartRaw, &ch)
		if err == nil {
			return ch, nil
		}
	}

	rawChart, err := c.getChartsByPriority(verifiedData)
	if err != nil {
		return watchmarket.Chart{}, errors.New(ErrInternal)
	}

	chart := normalizeChart(rawChart, verifiedData.MaxItems)

	chartRaw, err := json.Marshal(&chart)
	if err = c.dataCache.SetWithTime(key, chartRaw, verifiedData.TimeStart); err != nil {
		logger.Error("Failed to save cache", logger.Params{"err": err})
	}
	return chart, nil
}

func toChartsRequestData(cr ChartRequest) (ChartsNormalizedRequest, error) {
	if len(cr.TimeStartRaw) == 0 || len(cr.CoinQuery) == 0 {
		return ChartsNormalizedRequest{}, errors.New("Invalid arguments length")
	}

	coinId, err := strconv.Atoi(cr.CoinQuery)
	if err != nil {
		return ChartsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return ChartsNormalizedRequest{}, errors.New(ErrBadRequest)
	}

	timeStart, err := strconv.ParseInt(cr.TimeStartRaw, 10, 64)
	if err != nil {
		return ChartsNormalizedRequest{}, err
	}

	maxItems, err := strconv.Atoi(cr.MaxItems)
	if err != nil || maxItems <= 0 {
		maxItems = watchmarket.DefaultMaxChartItems
	}

	currency := watchmarket.DefaultCurrency
	if cr.Currency != "" {
		currency = cr.Currency
	}

	return ChartsNormalizedRequest{
		Coin:      uint(coinId),
		Token:     cr.Token,
		Currency:  currency,
		TimeStart: timeStart,
		MaxItems:  maxItems,
	}, nil
}

func (c Controller) getChartsByPriority(data ChartsNormalizedRequest) (watchmarket.Chart, error) {
	availableProviders := c.chartsPriority.GetAllProviders()

	for _, p := range availableProviders {
		price, err := c.api.ChartsAPIs[p].GetChartData(data.Coin, data.Token, data.Currency, data.TimeStart)
		if err == nil {
			return price, nil
		}
	}
	return watchmarket.Chart{}, nil
}

func normalizeChart(chart watchmarket.Chart, maxItems int) watchmarket.Chart {
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
