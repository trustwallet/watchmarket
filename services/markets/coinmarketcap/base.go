package coinmarketcap

import "github.com/trustwallet/watchmarket/services/charts/info"

const (
	id = "coinmarketcap"
)

type Provider struct {
	ID, currency string
	client       Client
	info         info.Client
}

func InitProvider(proApi, assetsApi, webApi, widgetApi, infoApi, key, currency string) Provider {
	return Provider{ID: id, currency: currency, client: NewClient(proApi, assetsApi, webApi, widgetApi, key), info: info.NewClient(infoApi)}
}
