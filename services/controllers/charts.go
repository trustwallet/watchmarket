package controllers

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (c Controller) HandleChartsRequest(cr ChartRequest) (watchmarket.Chart, error) {
	verifiedData, err := verifyChartsRequestData(cr)
	if err != nil {
		return watchmarket.Chart{}, err
	}

	key := c.chartsCache.GenerateKey(cr.coinQuery + cr.token + cr.currency + cr.maxItems)
	cachedChart, err := c.chartsCache.GetCharts(key, verifiedData.timeStart)
	if err == nil {
		return cachedChart, nil
	}

	rawChart, err := c.getDataWithProviders(verifiedData)
	if err != nil {
		return watchmarket.Chart{}, err
	}

	chart := normalizeChart(rawChart, verifiedData.maxItems)

	if err = c.chartsCache.SaveCharts(key, chart, verifiedData.timeStart); err != nil {
		logger.Error("Failed to save cache", logger.Params{"err": err})
	}
	return chart, nil
}

func verifyChartsRequestData(cr ChartRequest) (ChartsNormalizedRequest, error) {
	if len(cr.timeStartRaw) == 0 || len(cr.coinQuery) == 0 {
		return ChartsNormalizedRequest{}, errors.E("Invalid arguments length")
	}

	coinId, err := strconv.Atoi(cr.coinQuery)
	if err != nil {
		return ChartsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return ChartsNormalizedRequest{}, err
	}

	timeStart, err := strconv.ParseInt(cr.timeStartRaw, 10, 64)
	if err != nil {
		return ChartsNormalizedRequest{}, err
	}

	maxItems, err := strconv.Atoi(cr.maxItems)
	if err != nil || maxItems <= 0 {
		maxItems = watchmarket.DefaultMaxChartItems
	}

	currency := watchmarket.DefaultCurrency
	if cr.currency != "" {
		currency = cr.currency
	}

	return ChartsNormalizedRequest{
		coin:      uint(coinId),
		token:     cr.token,
		currency:  currency,
		timeStart: timeStart,
		maxItems:  maxItems,
	}, nil
}

func (c Controller) getDataWithProviders(data ChartsNormalizedRequest) (watchmarket.Chart, error) {
	availableProviders := c.chartsPriority.GetAllProviders()

	for _, p := range availableProviders {
		price, err := c.api.ChartsAPIs[p].GetChartData(data.coin, data.token, data.currency, data.timeStart)
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
