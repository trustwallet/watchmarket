package chartscontroller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"strconv"
	"strings"
	"time"
)

const charts = "charts"

func toChartsRequestData(cr controllers.ChartRequest) (chartsNormalizedRequest, error) {
	if len(cr.CoinQuery) == 0 {
		return chartsNormalizedRequest{}, errors.New("invalid arguments length")
	}

	coinId, err := strconv.Atoi(cr.CoinQuery)
	if err != nil {
		return chartsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return chartsNormalizedRequest{}, errors.New(watchmarket.ErrBadRequest)
	}
	var timeStart int64
	if cr.TimeStartRaw == "" {
		timeStart = time.Now().Unix() - 60*60*24
	} else {
		timeStart, err = strconv.ParseInt(cr.TimeStartRaw, 10, 64)
		if err != nil {
			return chartsNormalizedRequest{}, err
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

	return chartsNormalizedRequest{
		Coin:      uint(coinId),
		Token:     cr.Token,
		Currency:  currency,
		TimeStart: timeStart,
		MaxItems:  maxItems,
	}, nil
}

func (c Controller) checkTickersAvailability(coin uint, token string, ctx context.Context) ([]models.Ticker, error) {
	tr := []models.TickerQuery{{Coin: coin, TokenId: strings.ToLower(token)}}
	if c.configuration.RestAPI.UseMemoryCache {
		key := strings.ToLower(asset.BuildID(coin, token))
		rawResult, err := c.memoryCache.Get(key, ctx)
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

func (c Controller) getChartsByPriority(data chartsNormalizedRequest, ctx context.Context) (watchmarket.Chart, error) {
	availableProviders := c.chartsPriority
	for _, p := range availableProviders {
		price, err := c.api[p].GetChartData(data.Coin, data.Token, data.Currency, data.TimeStart, ctx)
		if len(price.Prices) > 0 && err == nil {
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
