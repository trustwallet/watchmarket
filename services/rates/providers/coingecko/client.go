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
	currency   string
	bucketSize int
}

func NewClient(api, currency string, bucketSize int) Client {
	return Client{
		Request:    blockatlas.InitClient(api),
		currency:   currency,
		bucketSize: bucketSize,
	}
}

func (c Client) fetchCoinsList() (Coins, error) {
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

func (c Client) fetchLatestRates(coins Coins, currency string) Prices {
	var (
		ci     = coins.getCoinsID()
		i      = 0
		prices Prices
	)

	prChan := make(chan Prices)
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

			cp, err := c.fetchCoinsMarkets(currency, ids)
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

	return prices
}

func (c Client) fetchCoinsMarkets(currency, ids string) (Prices, error) {
	var (
		values = url.Values{"vs_currency": {currency}, "sparkline": {"false"}, "ids": {ids}}
		result Prices
	)

	err := c.Get(&result, "v3/coins/markets", values)
	if err != nil {
		return result, err
	}
	return result, nil
}
