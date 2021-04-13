package db

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/services/controllers"
)

type (
	Instance interface {
		GetRates(currency string) ([]models.Rate, error)
		GetAllRates() ([]models.Rate, error)
		GetRatesByProvider(provider string) ([]models.Rate, error)
		AddRates(rates []models.Rate) error

		AddTickers(tickers []models.Ticker) error
		GetTickers(assets []controllers.Asset) ([]models.Ticker, error)
		GetAllTickers() ([]models.Ticker, error)
	}
)
