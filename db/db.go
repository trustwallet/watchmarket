package db

import "github.com/trustwallet/watchmarket/db/models"

type (
	Instance interface {
		GetRates(currency, provider string) ([]models.Rate, error)
		AddRates(rates []models.Rate) error

		AddTickers(tickers []models.Ticker) error
		GetTickers(coin uint, tokenId string) ([]models.Ticker, error)
		GetTickersByMap(tickersMap map[string]string) ([]models.Ticker, error)
	}
)
