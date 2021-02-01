package coinmarketcap

import (
	"errors"
	"fmt"
	"github.com/trustwallet/watchmarket/services/controllers"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

const (
	chartDataSize = 3
)

func (p Provider) GetChartData(asset controllers.Asset, currency string, timeStart int64) (watchmarket.Chart, error) {
	chartsData := watchmarket.Chart{}
	coinsFromCmcMap := CmcSlice(p.Cm).coinToCmcMap()
	coinObj, err := coinsFromCmcMap.getCoinByContract(asset)
	if err != nil {
		return chartsData, err
	}
	if timeStart < 1000000000 {
		timeStart = 1000000000
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

func (p Provider) GetCoinData(asset controllers.Asset, currency string) (watchmarket.CoinDetails, error) {
	details := watchmarket.CoinDetails{}
	coinsFromCmcMap := CmcSlice(p.Cm).coinToCmcMap()
	coinObj, err := coinsFromCmcMap.getCoinByContract(asset)
	if err != nil {
		return details, err
	}
	priceData, err := p.client.fetchCoinData(coinObj.Id, currency)
	if err != nil {
		return details, err
	}
	assetsData, err := p.info.GetCoinInfo(asset)
	if err != nil {
		log.WithFields(log.Fields{"coinID": asset.CoinId, "token": asset.TokenId}).Warn("No assets assets about that coinID")
	}

	return normalizeInfo(priceData, &assetsData)
}

func normalizeCharts(currency string, c Charts) watchmarket.Chart {
	chartsData := watchmarket.Chart{}
	prices := make([]watchmarket.ChartPrice, 0)
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
		prices = append(prices, watchmarket.ChartPrice{
			Price: quote[0],
			Date:  date.Unix(),
		})
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Date < prices[j].Date
	})
	chartsData.Prices = prices
	chartsData.Provider = id
	return chartsData
}

func normalizeInfo(priceData ChartInfo, assetsData *watchmarket.Info) (watchmarket.CoinDetails, error) {
	var chartInfoData ChartInfoData
	for _, chartInfoData = range priceData.Data {
		break
	}
	return watchmarket.CoinDetails{
		Provider:    id,
		ProviderURL: getUrl(chartInfoData.Slug),
		Info:        assetsData,
	}, nil
}

func (c CmcSlice) coinToCmcMap() (m CoinMapping) {
	m = make(map[string]CoinMap)
	for _, cm := range c {
		m[createID(cm.Coin, cm.TokenId)] = cm
	}
	return
}

func (cm CoinMapping) getCoinByContract(asset controllers.Asset) (c CoinMap, err error) {
	c, ok := cm[createID(asset.CoinId, asset.TokenId)]
	if !ok {
		err = errors.New("No coin found")
	}

	return
}

func getUrl(slug string) string {
	return fmt.Sprintf("https://coinmarketcap.com/currencies/%s/", slug)
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
