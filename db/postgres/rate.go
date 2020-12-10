package postgres

import (
	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/clause"
)

func (i *Instance) AddRates(rates []models.Rate, batchLimit uint) error {
	normalizedRates := normalizeRates(rates)
	batch := toRatesBatch(normalizedRates, batchLimit)
	for _, b := range batch {
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
		}).Create(&b).Error
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

func toRatesBatch(rates []models.Rate, sizeUint uint) [][]models.Rate {
	size := int(sizeUint)
	resultLength := (len(rates) + size - 1) / size
	result := make([][]models.Rate, resultLength)
	lo, hi := 0, size
	for i := range result {
		if hi > len(rates) {
			hi = len(rates)
		}
		result[i] = rates[lo:hi:hi]
		lo, hi = hi, hi+size
	}
	return result
}
