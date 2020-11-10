package postgres

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
)

func (i *Instance) GetAssetsToShow(currency string, ctx context.Context) ([]string, error) {
	return nil, nil
}

func (i *Instance) GetAssetsFromAlerts(interval models.Interval, ctx context.Context) ([]string, error) {
	return nil, nil
}

func (i Instance) AddNewAlerts(alerts []models.Alert, ctx context.Context) error {
	return nil
}
