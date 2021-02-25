package db

import (
	"github.com/trustwallet/watchmarket/db/models"
)

type (
	Instance interface {
		GetRates(currency string) ([]models.Rate, error)
		GetAllRates() ([]models.Rate, error)
		AddRates(rates []models.Rate) error

		AddTickers(tickers []models.Ticker) error
		GetTickers(coin uint, tokenId string) ([]models.Ticker, error)
		GetAllTickers() ([]models.Ticker, error)
		GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error)
	}
)
