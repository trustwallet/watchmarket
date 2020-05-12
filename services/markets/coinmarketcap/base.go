package coinmarketcap

import "github.com/trustwallet/watchmarket/services/assets"

const (
	id = "coinmarketcap"
)

type Provider struct {
	id, currency string
	client       Client
	info         assets.Client
}

func InitProvider(proApi, assetsApi, webApi, widgetApi, infoApi, key, currency string, info assets.Client) Provider {
	return Provider{id: id, currency: currency, client: NewClient(proApi, assetsApi, webApi, widgetApi, key), info: info}
}

func (p Provider) GetProvider() string {
	return p.id
}
