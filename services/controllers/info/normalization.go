package infocontroller

import (
	"context"
	"errors"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/db/models"
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

func (c Controller) getDetailsByPriority(data detailsNormalizedRequest, ctx context.Context) (controllers.InfoResponse, error) {
	availableTickerProviders := c.coinInfoPriority
	availableRateProviders := c.configuration.Markets.Priority.Rates

	var (
		result  controllers.InfoResponse
		details watchmarket.CoinDetails
	)
	for _, p := range availableTickerProviders {
		data, err := c.api[p].GetCoinData(data.Coin, data.Token, data.Currency, ctx)
		if err == nil {
			details = data
			break
		}
	}
	result.Info = details.Info
	result.Provider = details.Provider
	result.ProviderURL = details.ProviderURL

	dbTickers, err := c.database.GetTickersByQueries([]models.TickerQuery{{Coin: data.Coin, TokenId: data.Token}}, ctx)
	if err != nil {
		return controllers.InfoResponse{}, err
	}
	if len(dbTickers) == 0 {
		return controllers.InfoResponse{}, errors.New("empty db ticker")
	}
	tickerData, err := getTickerDataAccordingToPriority(availableTickerProviders, dbTickers)
	if err != nil {
		return controllers.InfoResponse{}, err
	}
	result.CirculatingSupply = tickerData.CirculatingSupply
	result.MarketCap = tickerData.MarketCap
	result.Vol24 = tickerData.Vol24
	result.TotalSupply = tickerData.TotalSupply

	if data.Currency != watchmarket.DefaultCurrency {
		rates, err := c.database.GetRates(data.Currency, ctx)
		if err != nil {
			return controllers.InfoResponse{}, err
		}
		if len(rates) == 0 {
			return controllers.InfoResponse{}, errors.New("empty db rate")
		}
		rate, err := getRateDataAccordingToPriority(availableRateProviders, rates)
		if err != nil {
			return controllers.InfoResponse{}, err
		}
		result.CirculatingSupply *= 1 / rate
		result.MarketCap *= 1 / rate
		result.Vol24 *= 1 / rate
		result.TotalSupply *= 1 / rate
	}
	return result, nil
}

func getTickerDataAccordingToPriority(availableProviders []string, tickers []models.Ticker) (tickerData, error) {
	for _, p := range availableProviders {
		for _, t := range tickers {
			if t.Provider == p && t.ShowOption != models.NeverShow {
				return tickerData{
					Vol24:             t.Volume,
					MarketCap:         t.MarketCap,
					CirculatingSupply: t.CirculatingSupply,
					TotalSupply:       t.TotalSupply,
				}, nil
			}
		}
	}
	return tickerData{}, errors.New("bad ticker or providers")
}

func getRateDataAccordingToPriority(availableProviders []string, rates []models.Rate) (float64, error) {
	for _, p := range availableProviders {
		for _, r := range rates {
			if r.Provider == p && r.ShowOption != models.NeverShow {
				return r.Rate, nil
			}
		}
	}
	return 0, errors.New("bad ticker or providers")
}
