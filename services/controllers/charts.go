package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"strings"
	"time"
)

const charts = "charts"

func (c Controller) HandleChartsRequest(cr ChartRequest, ctx context.Context) (watchmarket.Chart, error) {
	var ch watchmarket.Chart

	verifiedData, err := toChartsRequestData(cr)
	if err != nil {
		return ch, errors.New(ErrBadRequest)
	}

	key := c.dataCache.GenerateKey(charts + cr.CoinQuery + cr.Token + cr.Currency + cr.MaxItems)

	cachedChartRaw, err := c.dataCache.GetWithTime(key, verifiedData.TimeStart, ctx)
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
		return watchmarket.Chart{}, errors.New(ErrInternal)
	}

	if len(rawChart.Prices) < 1 {
		return watchmarket.Chart{}, errors.New(ErrNotFound)
	}

	chart := normalizeChart(rawChart, verifiedData.MaxItems)

	chartRaw, err := json.Marshal(&chart)
	if err != nil {
		logger.Error(err)
	}

	if err == nil && len(chart.Prices) > 0 {
		err = c.dataCache.SetWithTime(key, chartRaw, verifiedData.TimeStart, ctx)
		if err != nil {
			logger.Error("failed to save cache", logger.Params{"err": err})
		}
	}

	return chart, nil
}

func toChartsRequestData(cr ChartRequest) (ChartsNormalizedRequest, error) {
	if len(cr.CoinQuery) == 0 {
		return ChartsNormalizedRequest{}, errors.New("invalid arguments length")
	}

	coinId, err := strconv.Atoi(cr.CoinQuery)
	if err != nil {
		return ChartsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return ChartsNormalizedRequest{}, errors.New(ErrBadRequest)
	}
	var timeStart int64
	if cr.TimeStartRaw == "" {
		timeStart = time.Now().Unix() - 60*60*24
	} else {
		timeStart, err = strconv.ParseInt(cr.TimeStartRaw, 10, 64)
		if err != nil {
			return ChartsNormalizedRequest{}, err
		}
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

func (c Controller) checkTickersAvailability(coin uint, token string, ctx context.Context) ([]models.Ticker, error) {
	tr := []models.TickerQuery{{Coin: coin, TokenId: strings.ToLower(token)}}
	dbTickers, err := c.database.GetTickersByQueries(tr, ctx)
	if err != nil {
		return nil, err
	}
	res := make([]models.Ticker, 0, len(dbTickers))
	for _, t := range dbTickers {
		if t.ShowOption != 2 {
			res = append(res, t)
		}
	}
	return res, nil
}

func (c Controller) getChartsByPriority(data ChartsNormalizedRequest, ctx context.Context) (watchmarket.Chart, error) {
	availableProviders := c.chartsPriority

	for _, p := range availableProviders {
		price, err := c.api[p].GetChartData(data.Coin, data.Token, data.Currency, data.TimeStart, ctx)
		a := len(price.Prices) > 0
		b := err == nil
		if b && a {
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
