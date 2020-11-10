package worker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db/models"
)

func (w Worker) AlertsIndexer() {
	err := w.initAssetsListForDB()
	if err != nil {
		log.Error(err)
		return
	}

	assets, err := w.getAssetsToUpdate()
	if err != nil {
		log.Error(err)
		return
	}

	currentPrices, err := w.getCurrentPrices()
	if err != nil {
		log.Error(err)
		return
	}

	oldPrices, err := w.getOldPrices()
	if err != nil {
		log.Error(err)
		return
	}

	priceDifference, err := w.getPriceDifference()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(assets)
	log.Info(currentPrices)
	log.Info(oldPrices)
	log.Info(priceDifference)

	err = w.updateAlerts()
	if err != nil {
		log.Error(err)
		return
	}
}

func (w Worker) initAssetsListForDB() error {
	ctx := context.Background()
	intervals := []models.Interval{models.Hour, models.Day, models.Week}

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

func (w Worker) getAssetsToUpdate() ([]string, error) {
	return nil, nil
}

func (w Worker) getCurrentPrices() (map[string]float64, error) {
	return nil, nil
}

func (w Worker) getOldPrices() (map[string]float64, error) {
	return nil, nil
}

func (w Worker) getPriceDifference() (map[string]float64, error) {
	return nil, nil
}

func (w Worker) updateAlerts() error {
	return nil
}
