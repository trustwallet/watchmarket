package db

import (
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
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

func (i *Instance) GetRates(coin, token string) ([]watchmarket.Rate, error) {
	var ticker []watchmarket.Rate

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
