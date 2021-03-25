package infocontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
	"strings"
)

const info = "info"

type Controller struct {
	cache            cache.Provider
	coinInfoPriority []string
	ratesPriority    []string
	api              markets.ChartsAPIs
}

func NewController(
	cache cache.Provider,
	coinInfoPriority []string,
	ratesPriority []string,
	api markets.ChartsAPIs,
) Controller {
	return Controller{
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
		return controllers.InfoResponse{}, errors.New(watchmarket.ErrNotFound)
	}
	var infoResponse controllers.InfoResponse
	err = json.Unmarshal(cachedDetails, &infoResponse)
	return infoResponse, err
}

func (c Controller) getDetailsByPriority(request controllers.DetailsRequest) (controllers.InfoResponse, error) {
	key := strings.ToLower(asset.BuildID(request.Asset.CoinId, request.Asset.TokenId))
	rawResult, err := c.cache.Get(key)
	if err != nil {
		return controllers.InfoResponse{}, errors.New(watchmarket.ErrNotFound)
	}
	var ticker watchmarket.Ticker
	if err = json.Unmarshal(rawResult, &ticker); err != nil {
		return controllers.InfoResponse{}, errors.New(watchmarket.ErrNotFound)
	}
	result := c.getCoinDataFromApi(request.Asset, request.Currency)
	result.CirculatingSupply = ticker.CirculatingSupply
	result.MarketCap = ticker.MarketCap
	result.Vol24 = ticker.Volume
	result.TotalSupply = ticker.TotalSupply

	if request.Currency != watchmarket.DefaultCurrency {
		rateRaw, err := c.cache.Get(request.Currency)
		var rate watchmarket.Rate
		if err != nil {
			return controllers.InfoResponse{}, errors.New(watchmarket.ErrNotFound)
		}
		if err = json.Unmarshal(rateRaw, &rate); err != nil {
			return result, err
		}
		result.MarketCap *= 1 / rate.Rate
		result.Vol24 *= 1 / rate.Rate
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
