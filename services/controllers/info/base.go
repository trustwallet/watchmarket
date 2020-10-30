package infocontroller

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Controller struct {
	dataCache        cache.Provider
	chartsPriority   []string
	coinInfoPriority []string
	ratesPriority    []string
	tickersPriority  []string
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	cache cache.Provider,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		cache,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
		configuration,
	}
}

func (c Controller) HandleInfoRequest(dr controllers.DetailsRequest, ctx context.Context) (watchmarket.CoinDetails, error) {
	var cd watchmarket.CoinDetails

	req, err := toDetailsRequestData(dr)
	if err != nil {
		return cd, errors.New(watchmarket.ErrBadRequest)
	}

	key := c.dataCache.GenerateKey(info + dr.CoinQuery + dr.Token + dr.Currency)

	cachedDetails, err := c.dataCache.Get(key, ctx)
	if err == nil && len(cachedDetails) > 0 {
		if json.Unmarshal(cachedDetails, &cd) == nil {
			return cd, nil
		}
	}

	result, err := c.getDetailsByPriority(req, ctx)
	if err != nil {
		return watchmarket.CoinDetails{}, errors.New(watchmarket.ErrInternal)
	}

	if result.Info != nil && result.IsEmpty() {
		result.Info = nil
	}

	newCache, err := json.Marshal(result)
	if err != nil {
		log.Error(err)
	}

	if result.Info != nil {
		err = c.dataCache.Set(key, newCache, ctx)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("failed to save cache")
		}
	}

	return result, nil
}
