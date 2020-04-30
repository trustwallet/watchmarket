package coingecko

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/market/chart"
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"time"
)

const (
	id            = "coingecko"
	chartDataSize = 2
)

type Chart struct {
	chart.Chart
	client *coingecko.Client
}

func InitChart(api string) chart.ChartProvider {
	m := &Chart{
		Chart: chart.Chart{
			Id: id,
		},
		client: coingecko.NewClient(api),
	}
	return m
}

func (c *Chart) GetChartData(coinId uint, token string, currency string, timeStart int64) (watchmarket.ChartData, error) {
	chartsData := watchmarket.ChartData{}
	coins, err := c.client.FetchCoinsList()
	if err != nil {
		return chartsData, err
	}
	cache := coingecko.NewSymbolsCache(coins)

	coinResult, err := getCoinObj(cache, coinId, token)
	if err != nil {
		return chartsData, err
	}

	timeEndDate := time.Now().Unix()
	charts, err := c.client.GetChartsData(coinResult.Id, currency, timeStart, timeEndDate)
	if err != nil {
		return chartsData, err
	}

	return normalizeCharts(charts), nil
}

func (c *Chart) GetCoinData(coinId uint, token string, currency string) (watchmarket.ChartCoinInfo, error) {
	coins, err := c.client.FetchCoinsList()
	if err != nil {
		return watchmarket.ChartCoinInfo{}, err
	}
	cache := coingecko.NewSymbolsCache(coins)

	coinResult, err := getCoinObj(cache, coinId, token)
	if err != nil {
		return watchmarket.ChartCoinInfo{}, err
	}

	data := c.client.FetchLatestRates(coingecko.GeckoCoins{coinResult}, currency)
	if len(data) == 0 {
		return watchmarket.ChartCoinInfo{}, errors.E("No rates found", errors.Params{"id": coinResult.Id})
	}
	return normalizeInfo(data[0]), nil
}

func getCoinObj(cache *coingecko.SymbolsCache, coinId uint, token string) (coingecko.GeckoCoin, error) {
	c := coingecko.GeckoCoin{}
	coinObj, ok := coin.Coins[coinId]
	if !ok {
		return c, errors.E("Coin not found", errors.Params{"coindId": coinId})
	}

	c, err := cache.GetCoinsBySymbol(coinObj.Symbol, token)
	if err != nil {
		return c, err
	}

	return c, nil
}

func normalizeCharts(charts coingecko.Charts) watchmarket.ChartData {
	chartsData := watchmarket.ChartData{}
	prices := make([]watchmarket.ChartPrice, 0)
	for _, quote := range charts.Prices {
		if len(quote) != chartDataSize {
			continue
		}

		date := time.Unix(int64(quote[0])/1000, 0)
		prices = append(prices, watchmarket.ChartPrice{
			Price: quote[1],
			Date:  date.Unix(),
		})
	}

	chartsData.Prices = prices

	return chartsData
}

func normalizeInfo(data coingecko.CoinPrice) watchmarket.ChartCoinInfo {
	return watchmarket.ChartCoinInfo{
		Vol24:             data.TotalVolume,
		MarketCap:         data.MarketCap,
		CirculatingSupply: data.CirculatingSupply,
		TotalSupply:       data.TotalSupply,
	}
}
