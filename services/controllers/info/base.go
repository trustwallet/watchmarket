package infocontroller

import (
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Controller struct {
	database         db.Instance
	cache            cache.Provider
	chartsPriority   []string
	coinInfoPriority []string
	ratesPriority    []string
	tickersPriority  []string
	api              markets.ChartsAPIs
	configuration    config.Configuration
}

func NewController(
	database db.Instance,
	cache cache.Provider,
	chartsPriority, coinInfoPriority, ratesPriority, tickersPriority []string,
	api markets.ChartsAPIs,
	configuration config.Configuration,
) Controller {
	return Controller{
		database,
		cache,
		chartsPriority,
		coinInfoPriority,
		ratesPriority,
		tickersPriority,
		api,
		configuration,
	}
}

func (c Controller) HandleInfoRequest(dr controllers.DetailsRequest) (controllers.InfoResponse, error) {
	var cd controllers.InfoResponse

	req, err := toDetailsRequestData(dr)
	if err != nil {
		return cd, errors.New(watchmarket.ErrBadRequest)
	}

	key := c.cache.GenerateKey(info + dr.CoinQuery + dr.Token + dr.Currency)

	cachedDetails, err := c.cache.Get(key)
	if err == nil && len(cachedDetails) > 0 {
		if json.Unmarshal(cachedDetails, &cd) == nil {
			return cd, nil
		}
	}

	result, err := c.getDetailsByPriority(req)
	if err != nil {
		return controllers.InfoResponse{}, errors.New(watchmarket.ErrInternal)
	}

	if result.Info != nil && result.Vol24 == 0 && result.TotalSupply == 0 && result.CirculatingSupply == 0 {
		result.Info = nil
	}

	newCache, err := json.Marshal(result)
	if err != nil {
		log.Error(err)
	}

	if result.Info != nil {
		err = c.cache.Set(key, newCache)
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Error("failed to save cache")
		}
	}

	return result, nil
}
