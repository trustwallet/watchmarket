package controllers

import (
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (c Controller) HandleDetailsRequest(dr DetailsRequest) (watchmarket.CoinDetails, error) {
	var cd watchmarket.CoinDetails

	req, err := toDetailsRequestData(dr)
	if err != nil {
		return cd, errors.New(ErrBadRequest)
	}

	key := c.dataCache.GenerateKey(dr.CoinQuery + dr.Token + dr.Currency)

	cachedDetails, err := c.dataCache.Get(key)
	if err == nil && len(cachedDetails) > 0 {
		if json.Unmarshal(cachedDetails, &cd) == nil {
			return cd, nil
		}
	}

	result, err := c.getDetailsByPriority(req)
	if err != nil {
		return watchmarket.CoinDetails{}, errors.New(ErrInternal)
	}
	newCache, err := json.Marshal(result)
	if err != nil {
		logger.Error(err)
	}
	err = c.dataCache.Set(key, newCache)
	if err != nil {
		logger.Error(err)
	}

	return result, nil
}

func toDetailsRequestData(dr DetailsRequest) (DetailsNormalizedRequest, error) {
	if len(dr.CoinQuery) == 0 {
		return DetailsNormalizedRequest{}, errors.New("Invalid arguments length")
	}

	coinId, err := strconv.Atoi(dr.CoinQuery)
	if err != nil {
		return DetailsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return DetailsNormalizedRequest{}, errors.New(ErrBadRequest)
	}

	currency := watchmarket.DefaultCurrency
	if dr.Currency != "" {
		currency = dr.Currency
	}

	return DetailsNormalizedRequest{
		Coin:     uint(coinId),
		Token:    dr.Token,
		Currency: currency,
	}, nil
}

func (c Controller) getDetailsByPriority(data DetailsNormalizedRequest) (watchmarket.CoinDetails, error) {
	availableProviders := c.coinInfoPriority.GetAllProviders()

	for _, p := range availableProviders {
		data, err := c.api.ChartsAPIs[p].GetCoinData(data.Coin, data.Token, data.Currency)
		if err == nil {
			return data, nil
		}
	}
	return watchmarket.CoinDetails{}, nil
}
