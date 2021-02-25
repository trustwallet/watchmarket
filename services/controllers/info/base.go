package infocontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

const info = "info"

type Controller struct {
	database         db.Instance
	cache            cache.Provider
	coinInfoPriority []string
	ratesPriority    []string
	api              markets.ChartsAPIs
}

func NewController(
	database db.Instance,
	cache cache.Provider,
	coinInfoPriority []string,
	ratesPriority []string,
	api markets.ChartsAPIs,
) Controller {
	return Controller{
		database,
		cache,
		coinInfoPriority,
		ratesPriority,
		api,
	}
}

func (c Controller) HandleInfoRequest(request controllers.DetailsRequest) (controllers.InfoResponse, error) {
	result, err := c.getFromCache(request)
	if err == nil {
		return result, nil
	}

	result, err = c.getDetailsByPriority(request)
	if err != nil {
		return controllers.InfoResponse{}, errors.New(watchmarket.ErrInternal)
	}

	if result.Info != nil && result.Vol24 == 0 && result.TotalSupply == 0 && result.CirculatingSupply == 0 {
		result.Info = nil
	}

	c.putToCache(request, result)

	return result, nil
}

func (c Controller) putToCache(request controllers.DetailsRequest, result controllers.InfoResponse) {
	if result.Info == nil {
		return
	}
	key := c.getCacheKey(request)
	newCache, err := json.Marshal(result)
	if err != nil {
		return
	}

	err = c.cache.Set(key, newCache)
	if err != nil {
		log.Error("failed to save cache")
	}
}

func (c Controller) getCacheKey(request controllers.DetailsRequest) string {
	return c.cache.GenerateKey(fmt.Sprintf("%s%d%s%s", info, request.Asset.CoinId, request.Asset.TokenId, request.Currency))
}

func (c Controller) getFromCache(request controllers.DetailsRequest) (controllers.InfoResponse, error) {
	key := c.getCacheKey(request)

	cachedDetails, err := c.cache.Get(key)
	if err != nil || len(cachedDetails) <= 0 {
		return controllers.InfoResponse{}, errors.New("cache is empty")
	}
	var infoResponse controllers.InfoResponse
	err = json.Unmarshal(cachedDetails, &infoResponse)
	return infoResponse, err
}

func (c Controller) getDetailsByPriority(request controllers.DetailsRequest) (controllers.InfoResponse, error) {
	dbTickers, err := c.database.GetTickers([]controllers.Asset{request.Asset})

	if err != nil || len(dbTickers) == 0 {
		return controllers.InfoResponse{}, fmt.Errorf("no tickers in db or db error: %w", err)
	}

	ticker, err := c.getTickerDataAccordingToPriority(dbTickers)
	if err != nil {
		return controllers.InfoResponse{}, err
	}
	result := c.getCoinDataFromApi(request.Asset, request.Currency)
	result.CirculatingSupply = ticker.CirculatingSupply
	result.MarketCap = ticker.MarketCap
	result.Vol24 = ticker.Volume
	result.TotalSupply = ticker.TotalSupply

	if request.Currency != watchmarket.DefaultCurrency {
		rates, err := c.database.GetRates(request.Currency)
		if err != nil || len(rates) == 0 {
			return controllers.InfoResponse{}, fmt.Errorf("empty db rate or db error: %w", err)
		}
		rate, err := c.getRateDataAccordingToPriority(rates)
		if err != nil {
			return controllers.InfoResponse{}, err
		}
		result.MarketCap *= 1 / rate
		result.Vol24 *= 1 / rate
	}
	return result, nil
}

func (c Controller) getCoinDataFromApi(assetData controllers.Asset, currency string) controllers.InfoResponse {
	var result controllers.InfoResponse

	for _, p := range c.coinInfoPriority {
		if data, err := c.api[p].GetCoinData(assetData, currency); err == nil {
			result.Info = data.Info
			result.Provider = data.Provider
			result.ProviderURL = data.ProviderURL
			break
		}
	}
	return result
}

func (c Controller) getTickerDataAccordingToPriority(tickers []models.Ticker) (models.Ticker, error) {
	for _, p := range c.coinInfoPriority {
		for _, t := range tickers {
			if t.Provider == p && t.ShowOption != models.NeverShow {
				return t, nil
			}
		}
	}
	return models.Ticker{}, errors.New("bad ticker or providers")
}

func (c Controller) getRateDataAccordingToPriority(rates []models.Rate) (float64, error) {
	for _, p := range c.ratesPriority {
		for _, r := range rates {
			if r.Provider == p && r.ShowOption != models.NeverShow {
				return r.Rate, nil
			}
		}
	}
	return 0, errors.New("bad ticker or providers")
}
