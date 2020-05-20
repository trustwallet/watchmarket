package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
)

func (i *Instance) AddRates(rates []models.Rate) error {
	// TODO: Upsert
	db := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (currency,provider) DO NOTHING") //UPDATE SET rate = excluded.rate, percent_change24h = excluded.percent_change24h, timestamp = excluded.timestamp
	return BulkInsert(db, rates)
}

func (i *Instance) GetRates(currency string) ([]models.Rate, error) {
	var rates []models.Rate
	if err := i.Gorm.Where("currency = ?", currency).
		Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}
