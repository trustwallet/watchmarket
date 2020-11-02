package fixer

import (
	"context"
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

func (c Client) FetchRates(ctx context.Context) (Rate, error) {
	var (
		values = req.Param{"access_key": c.key, "base": c.currency} // Base USD supported only in paid api}
		result Rate
	)
	resp, err := c.r.Get(c.api+"/latest", values, ctx)
	if err != nil {
		return Rate{}, err
	}
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return Rate{}, err
	}
	return result, nil
}
