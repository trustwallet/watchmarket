package coinmarketcap

import (
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/services/assets"
	"io/ioutil"
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

func InitProvider(proApi, assetsApi, webApi, widgetApi, key, currency, mappingPath string, info assets.Client) Provider {
	cm, err := setupCoinMap(mappingPath)
	if err != nil {
		logger.Error("Init provider coin map: " + err.Error())
	}
	return Provider{
		id:       id,
		currency: currency,
		client:   NewClient(proApi, assetsApi, webApi, widgetApi, key),
		info:     info,
		cm:       cm,
	}
}

func (p Provider) GetProvider() string {
	return p.id
}

func setupCoinMap(path string) ([]CoinMap, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result []CoinMap
	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
