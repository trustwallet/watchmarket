package coinmarketcap

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	web    blockatlas.Request
	widget blockatlas.Request
	assets blockatlas.Request
}

func NewClient(webApi, widgetApi, assetsApi string) Client {
	return Client{
		web:    blockatlas.InitClient(webApi),
		widget: blockatlas.InitClient(widgetApi),
		assets: blockatlas.InitClient(assetsApi),
	}
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

func (c Client) fetchCoinMap() (CmcSlice, error) {
	var results CmcSlice
	err := c.assets.GetWithCache(&results, "mapping.json", nil, time.Hour*1)
	if err != nil {
		return nil, err
	}
	return results, nil
}
