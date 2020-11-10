package worker

import log "github.com/sirupsen/logrus"

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
