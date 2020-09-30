package fixer

import (
	"context"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
)

type Client struct {
	blockatlas.Request
	key, currency string
}

func NewClient(api, key, currency string) Client {
	return Client{Request: blockatlas.InitClient(api), key: key, currency: currency}
}

func (c Client) FetchRates(ctx context.Context) (Rate, error) {
	var (
		values  = url.Values{"access_key": {c.key}, "base": {c.currency}} // Base USD supported only in paid api}
		rawRate Rate
	)
	err := c.GetWithContext(&rawRate, "latest", values, ctx)
	if err != nil {
		return rawRate, err
	}
	return rawRate, nil
}
