package coingecko

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	blockatlas.Request
}

func NewClient(api string) Client {
	return Client{
		Request: blockatlas.InitClient(api),
	}
}

func (c Client) fetchCoins() (Coins, error) {
	var (
		result Coins
		values = url.Values{"include_platform": {"true"}}
	)
	err := c.GetWithCache(&result, "v3/coins/list", values, time.Hour)
	if err != nil {
		return result, err
	}
	return result, nil
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

func (c Client) fetchRates(coins Coins, currency string, bucketSize int) (prices CoinPrices) {
	ci := coins.coinIds()

	i := 0
	prChan := make(chan CoinPrices)
	var wg sync.WaitGroup
	for i < len(ci) {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var end = len(ci)
			if len(ci) > i+bucketSize {
				end = i + bucketSize
			}
			bucket := ci[i:end]
			ids := strings.Join(bucket[:], ",")

			cp, err := c.fetchMarkets(currency, ids)
			if err != nil {
				logger.Error(err)
				return
			}
			prChan <- cp
		}(i)

		i += bucketSize
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

func (c Client) fetchMarkets(currency, ids string) (cp CoinPrices, err error) {
	values := url.Values{
		"vs_currency": {currency},
		"sparkline":   {"false"},
		"ids":         {ids},
	}

	err = c.Get(&cp, "v3/coins/markets", values)
	return
}
