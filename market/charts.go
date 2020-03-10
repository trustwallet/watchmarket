package market

import (
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/numbers"
	"github.com/trustwallet/watchmarket/market/chart"
	"github.com/trustwallet/watchmarket/market/chart/coingecko"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"math"
	"sort"
)

const (
	minUnixTime = 1000000000
)

type Charts struct {
	ChartProviders chart.ChartProviders
}

func InitCharts() *Charts {
	return &Charts{chart.ChartProviders{
		0: coingecko.InitChart(
			viper.GetString("market.coingecko.api"),
		),
	}}
}

func (c *Charts) GetChartData(coin uint, token string, currency string, timeStart int64, maxItems int) (watchmarket.ChartData, error) {
	chartsData := watchmarket.ChartData{}
	timeStart = numbers.Max(timeStart, minUnixTime)
	for i := 0; i < len(c.ChartProviders); i++ {
		c := c.ChartProviders[i]
		charts, err := c.GetChartData(coin, token, currency, timeStart)
		if err != nil {
			continue
		}
		charts.Prices = normalizePrices(charts.Prices, maxItems)
		return charts, nil
	}

	return chartsData, watchmarket.ErrNotFound
}

func (c *Charts) GetCoinInfo(coin uint, token string, currency string) (watchmarket.ChartCoinInfo, error) {
	coinInfoData := watchmarket.ChartCoinInfo{}
	for i := 0; i < len(c.ChartProviders); i++ {
		c := c.ChartProviders[i]
		info, err := c.GetCoinData(coin, token, currency)
		if err != nil {
			continue
		}
		return info, nil
	}

	return coinInfoData, watchmarket.ErrNotFound
}

func normalizePrices(prices []watchmarket.ChartPrice, maxItems int) (result []watchmarket.ChartPrice) {
	sort.Slice(prices, func(p, q int) bool {
		return prices[p].Date < prices[q].Date
	})
	if len(prices) > maxItems && maxItems > 0 {
		skip := int(math.Ceil(float64(len(prices) / maxItems)))
		i := 0
		for i < len(prices) {
			result = append(result, prices[i])
			i += skip + 1
		}
		lastPrice := prices[len(prices)-1]
		if len(result) > 0 && lastPrice.Date != result[len(result)-1].Date {
			result = append(result, lastPrice)
		}
	} else {
		result = prices
	}
	return
}
