package db

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
)

type (
	Instance interface {
		GetRates(currency string, ctx context.Context) ([]models.Rate, error)
		GetAllRates(ctx context.Context) ([]models.Rate, error)
		AddRates(rates []models.Rate, batchLimit uint, ctx context.Context) error

		AddTickers(tickers []models.Ticker, batchLimit uint, ctx context.Context) error
		GetTickers(coin uint, tokenId string, ctx context.Context) ([]models.Ticker, error)
		GetAllTickers(ctx context.Context) ([]models.Ticker, error)
		GetTickersByQueries(tickerQueries []models.TickerQuery, ctx context.Context) ([]models.Ticker, error)
	}
)
