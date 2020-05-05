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

func (c Client) GetPrices() ([]CoinPrice, error) {
	var prices []CoinPrice
	err := c.Get(&prices, "v1/ticker/24hr", url.Values{"limit": {"1000"}})
	if err != nil {
		return nil, err
	}
	return prices, nil
}
