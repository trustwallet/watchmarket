package postgres

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/clause"
)

func (i *Instance) GetAlertsToShow(currency string, ctx context.Context) ([]string, error) {
	return nil, nil
}

func (i *Instance) GetAssetsFromAlerts(interval models.Interval, ctx context.Context) ([]string, error) {
	var alerts []models.Alert
	err := i.Gorm.Model(&models.Alert{}).Where("interval = ?", interval).Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(alerts))
	for _, a := range alerts {
		result = append(result, a.AssetID)
	}
	return result, nil
}

func (i Instance) AddNewAlerts(alerts []models.Alert, ctx context.Context) error {
	err := i.Gorm.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "asset_id",
			},
			{
				Name: "interval",
			},
		},
		DoUpdates: clause.AssignmentColumns([]string{"price", "difference", "updated_at"}),
	}).Create(&alerts).Error
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) GetAlertsByIntervalToUpdate(interval models.Interval, ctx context.Context) ([]models.Alert, error) {
	var alerts []models.Alert
	err := i.Gorm.Model(&models.Alert{}).Where("interval = ?", interval).Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}
