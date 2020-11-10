package worker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db/models"
)

func (w Worker) AlertsIndexer() {
	intervals := []models.Interval{models.Hour, models.Day, models.Week}

	err := w.initAssetsListForDB(intervals)
	if err != nil {
		log.Error(err)
		return
	}

	alerts, err := w.getAlertsToUpdate(intervals)
	if err != nil {
		log.Error(err)
		return
	}

	currentPrices, err := w.getCurrentPrices(alerts)
	if err != nil {
		log.Error(err)
		return
	}

	oldPrices, err := w.getOldPrices(alerts)
	if err != nil {
		log.Error(err)
		return
	}

	priceDifference, err := w.getPriceDifference()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(alerts)
	log.Info(currentPrices)
	log.Info(oldPrices)
	log.Info(priceDifference)

	err = w.updateAlerts()
	if err != nil {
		log.Error(err)
		return
	}
}

func (w Worker) initAssetsListForDB(intervals []models.Interval) error {
	ctx := context.Background()
	var intervalsToInit []models.Interval

	for _, interval := range intervals {
		assets, err := w.db.GetAssetsFromAlerts(interval, ctx)
		if err != nil {
			log.Error(err)
			continue
		}
		if len(assets) == 0 {
			continue
		}
		intervalsToInit = append(intervalsToInit, interval)
	}

	if len(intervalsToInit) == 0 {
		return nil
	}

	allTickers, err := w.db.GetAllTickers(ctx)
	if err != nil {
		return err
	}

	for _, interval := range intervalsToInit {
		var alerts []models.Alert
		for _, ticker := range allTickers {
			a := models.Alert{
				AssetID:    ticker.ID,
				Interval:   interval,
				Difference: 0,
				Display:    false,
			}
			alerts = append(alerts, a)
		}
		err = w.db.AddNewAlerts(alerts, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w Worker) getAlertsToUpdate(intervals []models.Interval) ([]models.Alert, error) {
	// get all assets where now - updated_at >= interval
	var result []models.Alert
	for _, interval := range intervals {
		a, err := w.db.GetAlertsByIntervalToUpdate(interval, context.Background())
		if err != nil {
			return nil, err
		}
		result = append(result, a...)
	}
	return result, nil
}

func (w Worker) getCurrentPrices(alerts []models.Alert) (map[string]float64, error) {
	return nil, nil
}

func (w Worker) getOldPrices(alerts []models.Alert) (map[string]float64, error) {
	return nil, nil
}

func (w Worker) getPriceDifference() (map[string]float64, error) {
	return nil, nil
}

func (w Worker) updateAlerts() error {
	return nil
}
