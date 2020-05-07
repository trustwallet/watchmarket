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
}

func NewClient(api string) Client {
	return Client{Request: blockatlas.InitClient(api)}
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

func (c Client) fetchMarkets(currency, ids string) (CoinPrices, error) {
	var (
		result CoinPrices
		values = url.Values{"vs_currency": {currency}, "sparkline": {"false"}, "ids": {ids}}
	)

	err := c.Get(&result, "v3/coins/markets", values)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) fetchCoins() (coins Coins, err error) {
	values := url.Values{
		"include_platform": {"true"},
	}
	err = c.Get(&coins, "v3/coins/list", values)
	return
}
