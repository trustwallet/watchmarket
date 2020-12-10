package coinmarketcap

import (
	"fmt"
	"strconv"

	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	proApiURL    string
	webApiURL    string
	widgetApiURL string
	key          string
	r            *req.Req
}

func NewClient(proApi, webApi, widgetApi, key string) Client {
	return Client{
		r:            req.New(),
		proApiURL:    proApi,
		webApiURL:    webApi,
		widgetApiURL: widgetApi,
		key:          key,
	}
}

func (c Client) fetchPrices(currency string) (CoinPrices, error) {
	var (
		result CoinPrices
		path   = c.proApiURL + "/v1/cryptocurrency/listings/latest"
		header = req.Header{"X-CMC_PRO_API_KEY": c.key}
	)

	resp, err := c.r.Get(path, req.Param{"limit": "5000", "convert": currency}, header)
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

func (c Client) fetchChartsData(id uint, currency string, timeStart int64, timeEnd int64, interval string) (Charts, error) {
	values := req.Param{
		"convert":    currency,
		"format":     "chart_crypto_details",
		"id":         strconv.FormatInt(int64(id), 10),
		"time_start": strconv.FormatInt(timeStart, 10),
		"time_end":   strconv.FormatInt(timeEnd, 10),
		"interval":   interval,
	}
	var result Charts
	resp, err := c.r.Get(c.webApiURL+"/v1/cryptocurrency/quotes/historical", values)
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

func (c Client) fetchCoinData(id uint, currency string) (ChartInfo, error) {
	values := req.Param{
		"convert": currency,
		"ref":     "widget",
	}
	var result ChartInfo
	resp, err := c.r.Get(c.widgetApiURL+fmt.Sprintf("/v2/ticker/%d", id), values)
	if err != nil {
		return ChartInfo{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return ChartInfo{}, err
	}
	return result, nil
}
