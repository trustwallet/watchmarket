package fixer

import (
	"net/url"

	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/network/middleware"
)

type Client struct {
	client   client.Request
	key      string
	currency string
}

func NewClient(api, key, currency string) Client {
	return Client{
		client:   client.InitClient(api, middleware.SentryErrorHandler),
		key:      key,
		currency: currency,
	}
}

func (c Client) FetchRates() (rate Rate, err error) {
	values := url.Values{"access_key": {c.key}, "base": {c.currency}} // Base USD supported only in paid api}
	err = c.client.Get(&rate, "/latest", values)
	return
}
