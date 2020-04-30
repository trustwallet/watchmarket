package db

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (i *Instance) AddTickers(tickers []watchmarket.Ticker, provider string) []error {
	var errorsList []error

	for _, ticker := range tickers {
		t := models.Ticker{
			Ticker:   ticker,
			Provider: provider,
		}

		err := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (id) DO NOTHING").Create(&t).Error
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}

	if len(errorsList) > 0 {
		return errorsList
	}
	return nil
}

func (i *Instance) GetTickers(coin uint, token string) ([]watchmarket.Ticker, error) {
	var ticker []watchmarket.Ticker

	err := i.Gorm.Where("coin = ? AND token = ?", coin, token).Find(&ticker).Error
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}
