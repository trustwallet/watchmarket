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

		GetAssetsFromAlerts(interval models.Interval, ctx context.Context) ([]string, error)
		GetAlertsByInterval(interval models.Interval, ctx context.Context) ([]models.Alert, error)
		GetAlertsByIntervalWithDifference(interval models.Interval,
			difference float64, ctx context.Context) ([]models.Alert, error)
		AddNewAlerts(alerts []models.Alert, ctx context.Context) error
	}
)
