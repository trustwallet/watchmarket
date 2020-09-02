package coinmarketcap

import (
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/services/assets"
)

const (
	id = "coinmarketcap"
)

type Provider struct {
	id, currency string
	client       Client
	info         assets.Client
	cm           []CoinMap
}

func InitProvider(proApi, webApi, widgetApi, key, currency string, info assets.Client) Provider {
	cm, err := setupCoinMap()
	if err != nil {
		logger.Error("Init provider coin map: " + err.Error())
	}
	return Provider{
		id:       id,
		currency: currency,
		client:   NewClient(proApi, webApi, widgetApi, key),
		info:     info,
		cm:       cm,
	}
}

func (p Provider) GetProvider() string {
	return p.id
}

func setupCoinMap() ([]CoinMap, error) {
	var result []CoinMap
	err := json.Unmarshal([]byte(mapping), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
