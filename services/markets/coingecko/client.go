package coingecko

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL    string
	bucketSize int
	r          *req.Req
}

func NewClient(api string, bucketSize int) Client {
	r := req.New()
	c := Client{r: r, bucketSize: bucketSize, baseURL: api}
	c.r.SetTimeout(time.Minute)
	return c
}

func (c Client) fetchCharts(id, currency string, timeStart, timeEnd int64, ctx context.Context) (Charts, error) {
	var (
		result Charts
		values = req.Param{
			"vs_currency": currency,
			"from":        strconv.FormatInt(timeStart, 10),
			"to":          strconv.FormatInt(timeEnd, 10),
		}
	)
	resp, err := c.r.Get(c.baseURL+fmt.Sprintf("/v3/coins/%s/market_chart/range", id), values, ctx)
	if err != nil {
		return Charts{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return Charts{}, err
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
				log.Error(err)
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

	resp, err := c.r.Get(c.baseURL+"/v3/coins/markets", values, ctx)
	if err != nil {
		return CoinPrices{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return CoinPrices{}, err
	}
	return result, nil
}

func (c Client) fetchCoins(ctx context.Context) (Coins, error) {
	var result Coins
	resp, err := c.r.Get(c.baseURL+"/v3/coins/list", req.Param{"include_platform": "true"}, ctx)
	if err != nil {
		return Coins{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return Coins{}, err
	}
	return result, nil
}
