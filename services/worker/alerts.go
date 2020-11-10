package worker

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"time"
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

	updatedAlerts, err := w.updateAlerts(currentPrices, alerts)
	if err != nil {
		log.Error(err)
		return
	}

	err = w.db.AddNewAlerts(updatedAlerts, context.Background())
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
			intervalsToInit = append(intervalsToInit, interval)
		}
	}

	if len(intervalsToInit) == 0 {
		return nil
	}

	allTickers, err := w.db.GetTickersByQueries([]models.TickerQuery{{Coin: 714}, {Coin: 0}, {Coin: 60}}, ctx)
	if err != nil {
		return err
	}

	for _, interval := range intervalsToInit {
		var alerts []models.Alert
		for _, ticker := range allTickers {
			// will need to resolve multiple providers later
			if ticker.Provider == "coingecko" {
				continue
			}
			a := models.Alert{
				AssetID:    ticker.ID,
				Interval:   interval,
				Difference: 0,
				Price:      ticker.Value,
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
	var queries []models.TickerQuery
	for _, a := range alerts {
		c, t, err := asset.ParseID(a.AssetID)
		if err != nil {
			continue
		}
		q := models.TickerQuery{
			Coin:    c,
			TokenId: t,
		}
		queries = append(queries, q)
	}
	tickers, err := w.db.GetTickersByQueries(queries, context.Background())
	if err != nil {
		return nil, err
	}
	result := make(map[string]float64)
	for _, t := range tickers {
		// will need to resolve multiple providers later
		if t.Provider == "coingecko" {
			continue
		}
		result[t.ID] = t.Value
	}
	return result, nil
}

func (w Worker) updateAlerts(currentPrices map[string]float64, alerts []models.Alert) ([]models.Alert, error) {
	result := make([]models.Alert, 0, len(alerts))
	for _, a := range alerts {
		oldPrice := a.Price
		newPrice, ok := currentPrices[a.AssetID]
		if !ok {
			continue
		}
		difference := calculatePriceDifference(oldPrice, newPrice)
		a.Difference = watchmarket.TruncateWithPrecision(difference, 2)

		if time.Now().Unix()-a.UpdatedAt.Unix() >= intervalToUnix("1") {
			a.Price = newPrice
		}
		result = append(result, a)
	}
	return result, nil
}

func intervalToUnix(interval models.Interval) int64 {
	switch interval {
	case models.Day:
		return 24 * 60 * 60
	case models.Week:
		return 7 * 24 * 60 * 60
	case models.Hour:
		return 60 * 60
	default:
		return 0
	}
}

func calculatePriceDifference(old, new float64) float64 {
	difference := old - new
	if difference < 0 {
		return -(difference / 100)
	} else {
		return difference / 100
	}
}
