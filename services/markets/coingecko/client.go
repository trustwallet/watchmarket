package coingecko

import (
	"context"
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
	bucketSize int
}

func NewClient(api string, bucketSize int) Client {
	c := Client{Request: blockatlas.InitClient(api), bucketSize: bucketSize}
	c.SetTimeout(time.Minute)
	return c
}

func (c Client) fetchCharts(id, currency string, timeStart, timeEnd int64, ctx context.Context) (Charts, error) {
	var (
		result Charts
		values = url.Values{
			"vs_currency": {currency},
			"from":        {strconv.FormatInt(timeStart, 10)},
			"to":          {strconv.FormatInt(timeEnd, 10)},
		}
	)
	err := c.GetWithContext(&result, fmt.Sprintf("v3/coins/%s/market_chart/range", id), values, ctx)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) fetchRates(coins Coins, currency string, ctx context.Context) (prices CoinPrices) {
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

			cp, err := c.fetchMarkets(ids, currency, ctx)
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

func (c Client) fetchMarkets(ids, currency string, ctx context.Context) (CoinPrices, error) {
	var (
		result CoinPrices
		values = url.Values{"vs_currency": {currency}, "sparkline": {"false"}, "ids": {ids}}
	)

	err := c.GetWithContext(&result, "v3/coins/markets", values, ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) fetchCoins(ctx context.Context) (Coins, error) {
	var result Coins
	err := c.GetWithCacheAndContext(&result, "v3/coins/list", url.Values{"include_platform": {"true"}}, time.Minute*10, ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
