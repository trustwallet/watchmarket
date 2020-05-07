package coinnmarketcap

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
)

type Client struct {
	blockatlas.Request
}

func NewClient(api, key string) Client {
	c := Client{
		Request: blockatlas.InitClient(api),
	}
	c.Headers["X-CMC_PRO_API_KEY"] = key

	return c
}

func (c *Client) FetchPrices(currency string) (CoinPrices, error) {
	var (
		prices CoinPrices
		path   = "v1/cryptocurrency/listings/latest"
	)

	request := blockatlas.Request{
		BaseUrl:      c.BaseUrl,
		Headers:      c.Headers,
		HttpClient:   blockatlas.DefaultClient,
		ErrorHandler: blockatlas.DefaultErrorHandler,
	}
	err := request.Get(&prices, path, url.Values{"limit": {"5000"}, "convert": {currency}})
	if err != nil {
		return prices, err
	}
	return prices, nil
}
