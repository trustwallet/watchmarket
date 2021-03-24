package coinmarketcap

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

const (
	id = watchmarket.CoinMarketCap
)

type Provider struct {
	id, currency string
	client       Client
	info         assets.Client
	Cm           []CoinMap
}

func InitProvider(proApi, webApi, widgetApi, key, currency string, info assets.Client) Provider {
	cm, err := setupCoinMap(Mapping)
	if err != nil {
		log.Error("Init provider coin map: " + err.Error())
	}
	return Provider{
		id:       id,
		currency: currency,
		client:   NewClient(proApi, webApi, widgetApi, key),
		info:     info,
		Cm:       cm,
	}
}

func (p Provider) GetProvider() string {
	return p.id
}

func setupCoinMap(mapping string) ([]CoinMap, error) {
	var result []CoinMap
	err := json.Unmarshal([]byte(mapping), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
