package binancedex

import (
	"net/url"

	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/network/middleware"
)

type Client struct {
	client.Request
}

func NewClient(api string) Client {
	return Client{client.InitClient(api, middleware.SentryErrorHandler)}
}

func (c Client) fetchPrices() (result []CoinPrice, err error) {
	err = c.Get(&result, "/v1/ticker/24hr", url.Values{"limit": {"1000"}})
	return
}
