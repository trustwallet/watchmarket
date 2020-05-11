package coinmarketcap

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/services/markets"
	"sort"
	"strings"
	"time"
)

const chartDataSize = 3

func (p Provider) GetChartData(coinID uint, token, currency string, timeStart int64) (markets.Data, error) {
	chartsData := markets.Data{}
	coinsFromCmc, err := p.client.fetchCoinMap()
	if err != nil {
		return chartsData, err
	}
	coinsFromCmcMap := CmcSlice(coinsFromCmc).coinToCmcMap()
	coinObj, err := coinsFromCmcMap.getCoinByContract(coinID, token)
	if err != nil {
		return chartsData, err
	}
	timeStartDate := time.Unix(timeStart, 0)
	days := int(time.Since(timeStartDate).Hours() / 24)
	timeEnd := time.Now().Unix()
	c, err := p.client.fetchChartsData(coinObj.Id, currency, timeStart, timeEnd, getInterval(days))
	if err != nil {
		return chartsData, err
	}
	return normalizeCharts(currency, c), nil
}


func (p Provider) GetCoinData(coin uint, token, currency string) (charts.CoinDetails, error) {
	info := charts.CoinDetails{}


	

	details := charts.CoinDetails{}

func (p Provider) GetCoinData(coinID uint, token, currency string) (markets.CoinDetails, error) {
	details := markets.CoinDetails{}
	coinsFromCmc, err := p.client.fetchCoinMap()

	if err != nil {
		return details, err
	}

	coinsFromCmcMap := coinsFromCmc.coinToCmcMap()
	coinObj, err := coinsFromCmcMap.getCoinByContract(coin, token)

	coinsFromCmcMap := CmcSlice(coinsFromCmc).coinToCmcMap()
	coinObj, err := coinsFromCmcMap.getCoinByContract(coinID, token)

	if err != nil {
		return details, err
	}
	priceData, err := p.client.fetchCoinData(coinObj.Id, currency)
	if err != nil {
		return details, err
	}
	assetsData, err := p.info.GetCoinInfo(coin, token)
	if err != nil {
		logger.Warn("No assets info about that coin", logger.Params{"coin": coin, "token": token})
	}

	return normalizeInfo(currency, coinObj.Id, priceData, assetsData)
}



func normalizeCharts(currency string, c Charts) charts.Data {
	chartsData := charts.Data{}
	prices := make([]charts.Price, 0)

func normalizeCharts(currency string, c Charts) markets.Data {
	chartsData := markets.Data{}
	prices := make([]markets.ChartsPrice, 0)

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
		prices = append(prices, markets.ChartsPrice{
			Price: quote[0],
			Date:  date.Unix(),
		})
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date < prices[j].Date
	})
	chartsData.Prices = prices
	return chartsData
}

func normalizeInfo(currency string, cmcCoin uint, priceData ChartInfo, assetsData markets.Info) (markets.CoinDetails, error) {
	details := markets.CoinDetails{}
	quote, ok := priceData.Data.Quotes[currency]
	if !ok {
		return details, errors.E("Cant get coin details", errors.Params{"cmcCoin": cmcCoin, "currency": currency})
	}
	return markets.CoinDetails{
		Vol24:             quote.Volume24,
		MarketCap:         quote.MarketCap,
		CirculatingSupply: data.Data.CirculatingSupply,
		TotalSupply:       data.Data.TotalSupply,
		Provider:          id,
		CirculatingSupply: priceData.Data.CirculatingSupply,
		TotalSupply:       priceData.Data.TotalSupply,
		Info:              assetsData,
	}, nil
}

func (c CmcSlice) coinToCmcMap() (m CoinMapping) {
	m = make(map[string]CoinMap)
	for _, cm := range c {
		m[createID(cm.Coin, cm.TokenId)] = cm
	}
	return
}

func (cm CoinMapping) getCoinByContract(coinId uint, contract string) (c CoinMap, err error) {
	c, ok := cm[createID(coinId, contract)]
	if !ok {
		err = errors.E("No coin found", errors.Params{"coin": coinId, "token": contract})
	}

	return
}

func createID(id uint, token string) string {
	return strings.ToLower(fmt.Sprintf("%d:%s", id, token))
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