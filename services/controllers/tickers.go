package controllers

import (
	"context"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

func (c Controller) HandleTickersRequestV2(tr TickerRequestV2, ctx context.Context) (TickerResponseV2, error) {
	if tr.Ids == nil {
		return TickerResponseV2{}, errors.New(ErrBadRequest)
	}

	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency), ctx)
	if err != nil {
		return TickerResponseV2{}, errors.New(ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueriesV2(tr.Ids), ctx)
	if err != nil {
		return TickerResponseV2{}, errors.New(ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate, ctx)

	return createResponseV2(tr, tickers), nil
}

func createResponseV2(tr TickerRequestV2, tickers watchmarket.Tickers) TickerResponseV2 {
	result := TickerResponseV2{
		Currency: tr.Currency,
	}
	tickersPrices := make([]TickerPrice, 0, len(tickers))

	for _, ticker := range tickers {
		id, ok := foundIDInRequest(tr, BuildID(ticker.Coin, ticker.TokenId))
		if !ok {
			logger.Error("Cannot find ID in request")
		}
		tp := TickerPrice{
			Change24h: ticker.Price.Change24h,
			Provider:  ticker.Price.Provider,
			Price:     ticker.Price.Value,
			ID:        id,
		}
		tickersPrices = append(tickersPrices, tp)
	}

	result.Tickers = tickersPrices
	return result
}

func makeTickerQueriesV2(ids []string) []models.TickerQuery {
	tickerQueries := make([]models.TickerQuery, 0, len(ids))
	for _, id := range ids {
		coin, token, _, err := ParseID(id)
		if err != nil {
			logger.Error(err.Error() + " " + "makeTickerQueriesV2")
		}
		tickerQueries = append(tickerQueries, models.TickerQuery{
			Coin:    coin,
			TokenId: strings.ToLower(token),
		})
	}

	return tickerQueries
}
