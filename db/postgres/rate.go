package postgres

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/trustwallet/watchmarket/db/models"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgorm"
	"strings"
	"time"
)

const (
	rawBulkRatesInsert = `INSERT INTO rates(updated_at,created_at,currency,percent_change24h,provider,rate,last_updated,show_option) VALUES %s ON CONFLICT ON CONSTRAINT rates_pkey DO UPDATE SET rate = excluded.rate, percent_change24h = excluded.percent_change24h, updated_at = excluded.updated_at, last_updated = excluded.last_updated`
)

func (i *Instance) AddRates(rates []models.Rate, batchLimit uint, ctx context.Context) error {
	g := apmgorm.WithContext(ctx, i.Gorm)
	span, _ := apm.StartSpan(ctx, "AddRates", "postgresql")
	defer span.End()
	normalizedRates := normalizeRates(rates)
	batch := toRatesBatch(normalizedRates, batchLimit)
	for _, b := range batch {
		err := bulkCreateRate(g, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) GetRates(currency string, ctx context.Context) ([]models.Rate, error) {
	g := apmgorm.WithContext(ctx, i.Gorm)
	var rates []models.Rate
	if err := g.Where("currency = ?", currency).
		Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (i *Instance) GetAllRates(ctx context.Context) ([]models.Rate, error) {
	g := apmgorm.WithContext(ctx, i.Gorm)
	var rates []models.Rate
	if err := g.Find(&rates).Error; err != nil {
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

func bulkCreateRate(db *gorm.DB, dataList []models.Rate) error {
	var (
		valueStrings []string
		valueArgs    []interface{}
	)

	for _, d := range dataList {
		valueStrings = append(valueStrings, "(?, ? ,?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, d.Currency)
		valueArgs = append(valueArgs, d.PercentChange24h)
		valueArgs = append(valueArgs, d.Provider)
		valueArgs = append(valueArgs, d.Rate)
		valueArgs = append(valueArgs, d.LastUpdated)
		valueArgs = append(valueArgs, d.ShowOption)
	}

	smt := fmt.Sprintf(rawBulkRatesInsert, strings.Join(valueStrings, ","))
	if err := db.Exec(smt, valueArgs...).Error; err != nil {
		return err
	}

	return nil
}
