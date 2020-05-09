package coinmarketcap

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/charts"
	"time"
)

const (
	id            = "coinmarketcap"
	chartDataSize = 3
)

type Provider struct {
	ID     string
	client Client
}

func InitProvider(webApi string, widgetApi string, mapApi string) Provider {
	return Provider{ID: id, client: NewClient(webApi, widgetApi, mapApi)}
}

func (p Provider) GetChartData(coin uint, token string, currency string, timeStart int64) (charts.Data, error) {
	chartsData := charts.Data{}
	cmap, err := p.client.GetCoinMap()
	if err != nil {
		return chartsData, err
	}
	coinObj, err := cmap.GetCoinByContract(coin, token)
	if err != nil {
		return chartsData, err
	}

	timeStartDate := time.Unix(timeStart, 0)
	days := int(time.Since(timeStartDate).Hours() / 24)
	timeEnd := time.Now().Unix()
	c, err := p.client.GetChartsData(coinObj.Id, currency, timeStart, timeEnd, getInterval(days))
	if err != nil {
		return chartsData, err
	}

	return normalizeCharts(currency, c), nil
}

func (p Provider) GetCoinData(coin uint, token, currency string) (charts.CoinDetails, error) {
	info := charts.CoinDetails{}

	cmap, err := p.client.GetCoinMap()
	if err != nil {
		return info, err
	}
	coinObj, err := cmap.GetCoinByContract(coin, token)
	if err != nil {
		return info, err
	}

	data, err := p.client.GetCoinData(coinObj.Id, currency)
	if err != nil {
		return info, err
	}

	return normalizeInfo(currency, coinObj.Id, data)
}

func normalizeCharts(currency string, c Charts) charts.Data {
	chartsData := charts.Data{}
	prices := make([]charts.Price, 0)
	for dateSrt, q := range c.Data {
		date, err := time.Parse(time.RFC3339, dateSrt)
		if err != nil {
			continue
		}

		quote, ok := q[currency]
		if !ok {
			continue
		}

		if len(quote) < chartDataSize {
			continue
		}
		prices = append(prices, charts.Price{
			Price: quote[0],
			Date:  date.Unix(),
		})
	}

	chartsData.Prices = prices

	return chartsData
}

func normalizeInfo(currency string, cmcCoin uint, data ChartInfo) (charts.CoinDetails, error) {
	info := charts.CoinDetails{}
	quote, ok := data.Data.Quotes[currency]
	if !ok {
		return info, errors.E("Cant get coin info", errors.Params{"cmcCoin": cmcCoin, "currency": currency})
	}

	return charts.CoinDetails{
		Vol24:             quote.Volume24,
		MarketCap:         quote.MarketCap,
		CirculatingSupply: data.Data.CirculatingSupply,
		TotalSupply:       data.Data.TotalSupply,
	}, nil
}

func getInterval(days int) string {
	switch d := days; {
	case d >= 360:
		return "1d"
	case d >= 90:
		return "2h"
	case d >= 30:
		return "1h"
	case d >= 7:
		return "15m"
	default:
		return "5m"
	}
}
