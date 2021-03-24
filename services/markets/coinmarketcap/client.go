package coinmarketcap

import (
	"net/url"
	"strconv"

	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/network/middleware"
)

type Client struct {
	proApi    client.Request
	webApi    client.Request
	widgetApi client.Request
}

func NewClient(proApi, webApi, widgetApi, key string) Client {
	c := Client{
		proApi:    client.InitClient(proApi, middleware.SentryErrorHandler),
		webApi:    client.InitClient(webApi, middleware.SentryErrorHandler),
		widgetApi: client.InitClient(widgetApi, middleware.SentryErrorHandler),
	}
	c.proApi.AddHeader("X-CMC_PRO_API_KEY", key)
	return c
}

func (c Client) fetchPrices(currency, cryptocurrencyType string) (result CoinPrices, err error) {
	params := url.Values{"limit": {"5000"}, "convert": {currency}, "cryptocurrency_type": {cryptocurrencyType}}
	err = c.proApi.Get(&result, "/v1/cryptocurrency/listings/latest", params)
	return
}

func (c Client) fetchChartsData(id uint, currency string, timeStart int64, timeEnd int64, interval string) (result Charts, err error) {
	values := url.Values{
		"convert":    {currency},
		"format":     {"chart_crypto_details"},
		"id":         {strconv.FormatInt(int64(id), 10)},
		"time_start": {strconv.FormatInt(timeStart, 10)},
		"time_end":   {strconv.FormatInt(timeEnd, 10)},
		"interval":   {interval},
	}
	err = c.webApi.Get(&result, "/v1/cryptocurrency/quotes/historical", values)
	return
}

func (c Client) fetchCoinData(id uint, currency string) (result ChartInfo, err error) {
	values := url.Values{
		"id":      {strconv.FormatInt(int64(id), 10)},
		"convert": {currency},
	}
	err = c.widgetApi.Get(&result, "/v1/cryptocurrency/widget", values)
	return
}
