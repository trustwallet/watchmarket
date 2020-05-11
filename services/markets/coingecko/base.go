package coingecko

import "github.com/trustwallet/watchmarket/services/charts/info"

const (
	id            = "coingecko"
	bucketSize    = 500
	unknownCoinID = 111111
	chartDataSize = 2
)

type Provider struct {
	ID, currency string
	client       Client
	info         info.Client
}

func InitProvider(api, infoApi, currency string) Provider {
	return Provider{ID: id, currency: currency, client: NewClient(api, currency, bucketSize), info: info.NewClient(infoApi)}
}
