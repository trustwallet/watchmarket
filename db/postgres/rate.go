package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/clause"
)

func (i *Instance) AddRates(rates []models.Rate) error {
	normalizedRates := normalizeRates(rates)
	for _, b := range normalizedRates {
		err := i.Gorm.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "currency",
				},
				{
					Name: "provider",
				},
			},
			DoUpdates: clause.AssignmentColumns([]string{"rate", "percent_change24h", "last_updated", "updated_at"}),
		}).CreateInBatches(&b, 1000).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) GetRates(currency string) ([]models.Rate, error) {
	var rates []models.Rate
	if err := i.Gorm.Where("currency = ?", currency).
		Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (i *Instance) GetAllRates() ([]models.Rate, error) {
	var rates []models.Rate
	if err := i.Gorm.Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (i *Instance) GetRatesByProvider(provider string) ([]models.Rate, error) {
	var rates []models.Rate
	if err := i.Gorm.Where("provider = ?", provider).Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func normalizeRates(rates []models.Rate) []models.Rate {
	ratesMap := make(map[string]models.Rate)
	for _, rate := range rates {
		key := rate.Currency + rate.Provider
		if _, ok := ratesMap[key]; !ok {
			ratesMap[key] = rate
		}
	}
	result := make([]models.Rate, 0, len(ratesMap))
	for _, rate := range ratesMap {
		result = append(result, rate)
	}
	return result
}
