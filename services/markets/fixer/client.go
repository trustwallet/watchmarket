package fixer

import (
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	key, currency, api string
	r                  *req.Req
}

func NewClient(api, key, currency string) Client {
	return Client{r: req.New(), key: key, currency: currency, api: api}
}

func (c Client) FetchRates() (Rate, error) {
	var (
		values = req.Param{"access_key": c.key, "base": c.currency} // Base USD supported only in paid api}
		result Rate
	)
	resp, err := c.r.Get(c.api+"/latest", values)
	if err != nil {
		return Rate{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.WithFields(log.Fields{
			"url":      resp.Request().URL.String(),
			"status":   resp.Response().Status,
			"response": resp,
		}).Error("Fixer Fetch Rates: ", resp.Response().Status)
		return Rate{}, err
	}
	return result, nil
}
