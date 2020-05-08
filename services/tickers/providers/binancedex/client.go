package binancedex

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
)

type Client struct {
	blockatlas.Request
}

func NewClient(api string) Client {
	return Client{
		blockatlas.InitClient(api),
	}
}

func (c Client) fetchPrices() ([]CoinPrice, error) {
	var result []CoinPrice
	err := c.Get(&result, "v1/ticker/24hr", url.Values{"limit": {"1000"}})
	if err != nil {
		return nil, err
	}
	return result, nil
}
