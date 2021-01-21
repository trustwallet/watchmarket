package chartscontroller

import (
	"errors"
	"strconv"
	"time"

	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
)

const charts = "charts"

func normalizeRequest(cr controllers.ChartRequest) (chartsNormalizedRequest, error) {
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
