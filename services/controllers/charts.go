package controllers

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"strconv"
)

func (c Controller) HandleChartsRequest(cr ChartRequest) (watchmarket.Chart, error) {
	verifiedData, err := verifyChartsRequestData(cr)
	if err != nil {
		return watchmarket.Chart{}, err
	}

	key := cache.GenerateKey(cr.coinQuery + cr.token + cr.currency + cr.maxItems)
	cachedChart, err := c.cache.GetCharts(key, verifiedData.timeStart)
	if err == nil {
		return cachedChart, nil
	}

	price, err := c.getDataWithProviders(verifiedData)
	if err != nil {
		return watchmarket.Chart{}, err
	}
	if err = c.cache.SaveCharts(key, price, verifiedData.timeStart); err != nil {
		logger.Error("Failed to save cache", logger.Params{"err": err})
	}
	return price, nil
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
