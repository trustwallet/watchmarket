package coingecko

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	blockatlas.Request
}

func NewClient(api string) *Client {
	c := Client{
		Request: blockatlas.InitClient(api),
	}
	return &c
}

func (c *Client) FetchLatestRates(coins GeckoCoins, currency string, bucketSize int) (prices CoinPrices) {
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

			cp, err := c.FetchCoinsMarkets(currency, ids)
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

func (c *Client) FetchCoinsMarkets(currency, ids string) (cp CoinPrices, err error) {
	values := url.Values{
		"vs_currency": {currency},
		"sparkline":   {"false"},
		"ids":         {ids},
	}

	err = c.Get(&cp, "v3/coins/markets", values)
	return
}

func (c *Client) FetchCoinsList() (coins GeckoCoins, err error) {
	values := url.Values{
		"include_platform": {"true"},
	}
	err = c.GetWithCache(&coins, "v3/coins/list", values, time.Hour)
	return
}
