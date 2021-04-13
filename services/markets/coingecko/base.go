package coingecko

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

const (
	id            = watchmarket.CoinGecko
	bucketSize    = 200
	chartDataSize = 2
)

type Provider struct {
	id       string
	key      string
	currency string
	client   Client
	info     assets.Client
}

func InitProvider(api, key, currency string, info assets.Client) Provider {
	return Provider{id: id, key: key, currency: currency, client: NewClient(api, key, bucketSize), info: info}
}

func (p Provider) GetProvider() string {
	return p.id
}
