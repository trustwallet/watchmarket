package coingecko

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/url"
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
		return result, err
	}
	return result, nil
}

func (c Client) FetchCoins() (Coins, error) {
	var result Coins
	err := c.Get(&result, "v3/coins/list", url.Values{"include_platform": {"true"}})
	if err != nil {
		return result, err
	}
	return result, nil
}
