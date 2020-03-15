package setup

import (
	"github.com/trustwallet/watchmarket/market/chart"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	id = "cmc"
)

type Chart struct {
	chart.Chart
}

func InitChartProviders() *chart.ChartProviders {
	return &chart.ChartProviders{
		0: InitChart(),
	}
}

func InitChart() chart.ChartProvider {
	m := &Chart{
		Chart: chart.Chart{
			Id: id,
		},
	}
	return m
}

func (c *Chart) GetId() string {
	return c.Id
}

func (c *Chart) GetCoinData(coin uint, token string, currency string) (watchmarket.ChartCoinInfo, error) {
	coinInfoData := watchmarket.ChartCoinInfo{}
	return coinInfoData, nil
}

func (c *Chart) GetChartData(coin uint, token string, currency string, timeStart int64) (watchmarket.ChartData, error) {
	price := watchmarket.ChartPrice{
		Price: 10,
		Date:  0,
	}

	prices := make([]watchmarket.ChartPrice, 0)
	prices = append(prices, price)
	prices = append(prices, price)

	data := watchmarket.ChartData{
		Prices: prices,
		Error:  "",
	}

	return data, nil
}

func (c *Chart) GetCoinInfo(coin uint, token string, currency string) (watchmarket.ChartCoinInfo, error) {
	coinInfoData := watchmarket.ChartCoinInfo{}
	return coinInfoData, nil
}
