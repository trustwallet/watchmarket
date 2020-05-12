package cache

import (
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (i Instance) GetTickers(key string) (watchmarket.Tickers, error) {
	raw, err := i.redis.Get(key)
	if err != nil {
		return nil, err
	}
	var result watchmarket.Tickers
	err = json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i Instance) SaveTickers(key string, tickers watchmarket.Tickers) error {
	if len(tickers) == 0 {
		return errors.E("Tickers are empty")
	}
	return nil

}
