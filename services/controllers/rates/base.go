package ratescontroller

import (
	"context"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
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

func (c Controller) HandleRatesRequest(r controllers.RateRequest, ctx context.Context) (controllers.RateResponse, error) {
	fromRate, err := c.getRateByCurrency(r.From, ctx)
	if err != nil {
		return controllers.RateResponse{}, err
	}
	toRate, err := c.getRateByCurrency(r.To, ctx)
	if err != nil {
		return controllers.RateResponse{}, err
	}
	fromAmountInUSD := r.Amount * fromRate.Rate
	result := fromAmountInUSD / toRate.Rate
	return controllers.RateResponse{Amount: result}, nil
}

func (c Controller) getRateByCurrency(currency string, ctx context.Context) (models.Rate, error) {
	emptyRate := models.Rate{}
	rates, err := c.database.GetRates(currency, ctx)
	if err != nil {
		logger.Error(err, "getRateByPriority")
		return emptyRate, err
	}

	providers := c.ratesPriority
	var result models.Rate
ProvidersLoop:
	for _, p := range providers {
		for _, r := range rates {
			if p == r.Provider {
				result = r
				break ProvidersLoop
			}
		}
	}

	if result == emptyRate {
		return emptyRate, errors.New(watchmarket.ErrNotFound)
	}

	return result, nil
}
