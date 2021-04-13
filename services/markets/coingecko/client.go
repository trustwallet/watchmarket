package coingecko

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/network/middleware"
)

type Client struct {
	client     client.Request
	key        string
	bucketSize int
}

func NewClient(api string, key string, bucketSize int) Client {
	c := Client{client: client.InitClient(api, middleware.SentryErrorHandler), key: key, bucketSize: bucketSize}
	c.client.HttpClient = &http.Client{
		Timeout: time.Minute,
	}
	return c
}

func (c Client) Get(result interface{}, path string, values url.Values) error {
	values.Add("x_cg_pro_api_key", c.key)
	return c.client.Get(&result, path, values)
}

func (c Client) fetchCharts(id, currency string, timeStart, timeEnd int64) (charts Charts, err error) {

	values := url.Values{
		"vs_currency": {currency},
		"from":        {strconv.FormatInt(timeStart, 10)},
		"to":          {strconv.FormatInt(timeEnd, 10)},
	}

	err = c.Get(&charts, fmt.Sprintf("/v3/coins/%s/market_chart/range", id), values)
	return
}

func (c Client) fetchRates(coins Coins, currency string) (prices CoinPrices) {
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

			cp, err := c.fetchMarkets(ids, currency)
			if err != nil {
				log.WithFields(log.Fields{
					"ids":      ids,
					"currency": currency,
				}).Error("CoinGecko Fetch Rates")
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

func (c Client) fetchMarkets(ids, currency string) (result CoinPrices, err error) {
	values := url.Values{"vs_currency": {currency}, "sparkline": {"false"}, "ids": {ids}}
	err = c.Get(&result, "/v3/coins/markets", values)
	return
}

func (c Client) fetchCoins() (result Coins, err error) {
	err = c.Get(&result, "/v3/coins/list", url.Values{"include_platform": {"true"}})
	return
}
