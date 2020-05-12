package coingecko

import "github.com/trustwallet/watchmarket/services/assets"

const (
	id            = "coingecko"
	bucketSize    = 500
	unknownCoinID = 111111
	chartDataSize = 2
)

type Provider struct {
	id, currency string
	client       Client
	info         assets.Client
}

func InitProvider(api, infoApi, currency string, info assets.Client) Provider {
	return Provider{id: id, currency: currency, client: NewClient(api, bucketSize), info: info}
}

func (p Provider) GetProvider() string {
	return p.id
}
