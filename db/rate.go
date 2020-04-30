package db

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (i *Instance) AddRates(rates []watchmarket.Rate, provider market.Provider) error {
	var errorsList []error

	for _, rate := range rates {
		r := models.Rate{
			Rate:     rate,
			Provider: provider.GetId(),
		}

		err := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (id) DO NOTHING").Create(&r).Error
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}

	if len(errorsList) > 0 {
		return ConvertToError(errorsList)
	}
	return nil
}

func (i *Instance) GetRates(coin uint, token string) ([]watchmarket.Rate, error) {
	var ticker []watchmarket.Rate

	err := i.Gorm.Where("coin = ? AND token = ?", coin, token).Find(&ticker).Error
	if err != nil {
		return ticker, err
	}

	return ticker, nil
}
