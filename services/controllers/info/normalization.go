package infocontroller

import (
	"context"
	"errors"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"strconv"
)

type (
	detailsNormalizedRequest struct {
		Coin            uint
		Token, Currency string
	}
)

const info = "info"

func toDetailsRequestData(dr controllers.DetailsRequest) (detailsNormalizedRequest, error) {
	if len(dr.CoinQuery) == 0 {
		return detailsNormalizedRequest{}, errors.New("invalid arguments length")
	}

	coinId, err := strconv.Atoi(dr.CoinQuery)
	if err != nil {
		return detailsNormalizedRequest{}, err
	}

	if _, ok := coin.Coins[uint(coinId)]; !ok {
		return detailsNormalizedRequest{}, errors.New(watchmarket.ErrBadRequest)
	}

	currency := watchmarket.DefaultCurrency
	if dr.Currency != "" {
		currency = dr.Currency
	}

	return detailsNormalizedRequest{
		Coin:     uint(coinId),
		Token:    dr.Token,
		Currency: currency,
	}, nil
}

func (c Controller) getDetailsByPriority(data detailsNormalizedRequest, ctx context.Context) (watchmarket.CoinDetails, error) {
	availableProviders := c.coinInfoPriority

	for _, p := range availableProviders {
		data, err := c.api[p].GetCoinData(data.Coin, data.Token, data.Currency, ctx)
		if err == nil {
			return data, nil
		}
	}
	return watchmarket.CoinDetails{}, nil
}
