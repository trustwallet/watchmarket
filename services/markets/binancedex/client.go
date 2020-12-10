package binancedex

import (
	"github.com/imroc/req"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseURL string
	r       *req.Req
}

func NewClient(api string) Client {
	return Client{
		baseURL: api,
		r:       req.New(),
	}
}

func (c Client) fetchPrices() ([]CoinPrice, error) {
	resp, err := c.r.Get(c.baseURL+"/v1/ticker/24hr", req.Param{"limit": "1000"})
	if err != nil {
		return nil, err
	}
	var result []CoinPrice
	err = resp.ToJSON(&result)
	if err != nil {
		log.Error("URL: " + resp.Request().URL.String())
		log.Error("Status code: " + resp.Response().Status)
		return nil, err
	}
	return result, nil
}
