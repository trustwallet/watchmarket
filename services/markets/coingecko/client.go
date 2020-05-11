package coingecko

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	blockatlas.Request
	currency   string
	bucketSize int
}

func NewClient(api, currency string, bucketSize int) Client {
	return Client{Request: blockatlas.InitClient(api), currency: currency, bucketSize: bucketSize}
}

func (c Client) fetchCharts(id, currency string, timeStart, timeEnd int64) (Charts, error) {
	var (
		result Charts
		values = url.Values{
			"vs_currency": {currency},
			"from":        {strconv.FormatInt(timeStart, 10)},
			"to":          {strconv.FormatInt(timeEnd, 10)},
		}
	)
	err := c.Get(&result, fmt.Sprintf("v3/coins/%s/market_chart/range", id), values)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) fetchRates(coins Coins) (prices CoinPrices) {
	ci := coins.coinIds()

	i := 0
	prChan := make(chan CoinPrices)
	var wg sync.WaitGroup
	for i < len(ci) {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var end = len(ci)
			if len(ci) > i+c.bucketSize {
				end = i + c.bucketSize
			}
			bucket := ci[i:end]
			ids := strings.Join(bucket[:], ",")

			cp, err := c.fetchMarkets(ids)
			if err != nil {
				logger.Error(err)
				return
			}
			prChan <- cp
		}(i)

		i += c.bucketSize
	}

	go func() {
		wg.Wait()
		close(prChan)
	}()

	for bucket := range prChan {
		prices = append(prices, bucket...)
	}

	return
}

func (c Client) fetchMarkets(ids string) (CoinPrices, error) {
	var (
		result CoinPrices
		values = url.Values{"vs_currency": {c.currency}, "sparkline": {"false"}, "ids": {ids}}
	)

	err := c.Get(&result, "v3/coins/markets", values)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) fetchCoins() (Coins, error) {
	var result Coins
	err := c.Get(&result, "v3/coins/list", url.Values{"include_platform": {"true"}})
	if err != nil {
		return nil, err
	}
	return result, nil
}
