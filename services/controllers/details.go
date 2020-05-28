package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

const info = "info"

func (c Controller) HandleDetailsRequest(dr DetailsRequest, ctx context.Context) (watchmarket.CoinDetails, error) {
	var cd watchmarket.CoinDetails

	req, err := toDetailsRequestData(dr)
	if err != nil {
		return cd, errors.New(ErrBadRequest)
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
		return watchmarket.CoinDetails{}, errors.New(ErrInternal)
	}

	if result.Info != nil && result.IsEmpty() {
		result.Info = nil
	}

	newCache, err := json.Marshal(result)
	if err != nil {
		logger.Error(err)
	}

	if result.Info != nil {
		err = c.dataCache.Set(key, newCache, ctx)
		if err != nil {
			logger.Error("failed to save cache", logger.Params{"err": err})
		}
	}

	return result, nil
}

func toDetailsRequestData(dr DetailsRequest) (DetailsNormalizedRequest, error) {
	if len(dr.CoinQuery) == 0 {
		return DetailsNormalizedRequest{}, errors.New("invalid arguments length")
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

func (c Controller) getDetailsByPriority(data DetailsNormalizedRequest, ctx context.Context) (watchmarket.CoinDetails, error) {
	availableProviders := c.coinInfoPriority

	for _, p := range availableProviders {
		data, err := c.api[p].GetCoinData(data.Coin, data.Token, data.Currency, ctx)
		if err == nil {
			return data, nil
		}
	}
	return watchmarket.CoinDetails{}, nil
}
