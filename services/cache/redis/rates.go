package rediscache

import (
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (i Instance) GetRates(key string) (watchmarket.Rates, error) {
	raw, err := i.redis.Get(key)
	if err != nil {
		return nil, err
	}
	var result watchmarket.Rates
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i Instance) SaveRates(key string, tickers watchmarket.Rates) error {
	if len(tickers) == 0 {
		return errors.E("Rates are empty")
	}

	raw, err := json.Marshal(tickers)
	if err != nil {
		return err
	}

	err = i.redis.Set(key, raw, i.ratesCaching)
	if err != nil {
		return err
	}
	return nil
}
