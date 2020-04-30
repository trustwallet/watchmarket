package db

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (i *Instance) AddTickers(tickers []watchmarket.Ticker, provider market.Provider) error {
	var errorsList []error

	for _, ticker := range tickers {
		t := models.Ticker{
			Ticker:   ticker,
			Provider: provider.GetId(),
		}

		err := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (id) DO NOTHING").Create(&t).Error
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}

	if len(errorsList) > 0 {
		return ConvertToError(errorsList)
	}
	return nil
}

func (i *Instance) GetTickers(coin, token string) ([]watchmarket.Ticker, error) {
	var ticker []watchmarket.Ticker

	coinID, err := strconv.Atoi(coin)
	if err != nil {
		return ticker, err
	}

	err = i.Gorm.Where("coin = ? AND token = ?", coinID, token).Find(&ticker).Error
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}
