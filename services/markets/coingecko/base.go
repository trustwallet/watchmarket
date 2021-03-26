package coingecko

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

const (
	id            = watchmarket.CoinGecko
	bucketSize    = 250
	chartDataSize = 2
)

type Provider struct {
	id       string
	currency string
	client   Client
	info     assets.Client
}

func InitProvider(api, currency string, info assets.Client) Provider {
	return Provider{id: id, currency: currency, client: NewClient(api, bucketSize), info: info}
}

func (p Provider) GetProvider() string {
	return p.id
}
