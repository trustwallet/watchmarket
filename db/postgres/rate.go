package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
)

func (i *Instance) AddRates(rates []models.Rate) error {
	// TODO: Upsert
	db := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (subscription_id) DO UPDATE SET subscription_id = excluded.subscription_id")
	return BulkInsert(db, rates)
}

func (i *Instance) GetRates(currency, provider string) ([]models.Rate, error) {
	var rates []models.Rate
	if err := i.Gorm.Where("currency = ? AND provider = ?", currency, provider).
		Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}
