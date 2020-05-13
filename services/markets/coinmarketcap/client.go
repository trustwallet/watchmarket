package coinmarketcap

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
	"strconv"
)

type Client struct {
	api    blockatlas.Request
	assets blockatlas.Request
	web    blockatlas.Request
	widget blockatlas.Request
}

func NewClient(proApi, assetsApi, webApi, widgetApi, key string) Client {
	c := Client{
		api:    blockatlas.InitClient(proApi),
		assets: blockatlas.InitClient(assetsApi),
		web:    blockatlas.InitClient(webApi),
		widget: blockatlas.InitClient(widgetApi),
	}
	c.api.Headers["X-CMC_PRO_API_KEY"] = key
	return c
}

func (c Client) fetchPrices(currency string) (CoinPrices, error) {
	var (
		result CoinPrices
		path   = "v1/cryptocurrency/listings/latest"
	)

	err := c.api.Get(&result, path, url.Values{"limit": {"5000"}, "convert": {currency}})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) fetchCoinMap() ([]CoinMap, error) {
	var (
		result []CoinMap
		path   = "mapping.json"
	)

	err := c.assets.Get(&result, path, nil)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) fetchChartsData(id uint, currency string, timeStart int64, timeEnd int64, interval string) (charts Charts, err error) {
	values := url.Values{
		"convert":    {currency},
		"format":     {"chart_crypto_details"},
		"id":         {strconv.FormatInt(int64(id), 10)},
		"time_start": {strconv.FormatInt(timeStart, 10)},
		"time_end":   {strconv.FormatInt(timeEnd, 10)},
		"interval":   {interval},
	}
	err = c.web.Get(&charts, "v1/cryptocurrency/quotes/historical", values)
	return
}

func (c Client) fetchCoinData(id uint, currency string) (charts ChartInfo, err error) {
	values := url.Values{
		"convert": {currency},
		"ref":     {"widget"},
	}
	err = c.widget.Get(&charts, fmt.Sprintf("v2/ticker/%d", id), values)
	return
}
