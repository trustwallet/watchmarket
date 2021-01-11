package ratescontroller

import (
	"encoding/json"
	"errors"

	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
)

type Controller struct {
	database      db.Instance
	dataCache     cache.Provider
	ratesPriority []string
	configuration config.Configuration
}

func NewController(
	database db.Instance,
	cache cache.Provider,
	ratesPriority []string,
	configuration config.Configuration,
) Controller {
	return Controller{
		dataCache:     cache,
		database:      database,
		ratesPriority: ratesPriority,
		configuration: configuration,
	}
}

func (c Controller) HandleRatesRequest(r controllers.RateRequest) (controllers.RateResponse, error) {
	fromRate, err := c.getRateByCurrency(r.From)
	if err != nil {
		return controllers.RateResponse{}, err
	}
	toRate, err := c.getRateByCurrency(r.To)
	if err != nil {
		return controllers.RateResponse{}, err
	}
	fromAmountInUSD := r.Amount * fromRate.Rate
	if fromRate.Rate == 0 {
		return controllers.RateResponse{}, errors.New("from rate is zero")
	}
	result := fromAmountInUSD / toRate.Rate
	return controllers.RateResponse{Amount: result}, nil
}

func (c Controller) getRateByCurrency(currency string) (watchmarket.Rate, error) {
	rawResult, err := c.dataCache.Get(currency)
	if err != nil {
		return watchmarket.Rate{}, err
	}
	var result watchmarket.Rate
	if err = json.Unmarshal(rawResult, &result); err != nil {
		return watchmarket.Rate{}, err
	}
	return result, nil
}
