package coinmarketcap

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
)

type Client struct {
	api    blockatlas.Request
	assets blockatlas.Request
}

func NewClient(proApi, assetsApi, key string) Client {
	c := Client{
		api:    blockatlas.InitClient(proApi),
		assets: blockatlas.InitClient(assetsApi),
	}
	c.api.Headers["X-CMC_PRO_API_KEY"] = key
	return c
}

func (c Client) FetchPrices(currency string) (CoinPrices, error) {
	var (
		result CoinPrices
		path   = "v1/cryptocurrency/listings/latest"
	)

	err := c.api.Get(&result, path, url.Values{"limit": {"5000"}, "convert": {currency}})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c Client) FetchCoinMap() ([]CoinMap, error) {
	var (
		result []CoinMap
		path   = "mapping.json"
	)

	err := c.assets.Get(&result, path, nil)
	if err != nil {
		return result, err
	}
	return result, nil
}
