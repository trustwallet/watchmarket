package caching

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
)

type Provider struct {
	DB storage.Caching
}

var ChartsCachingDuration int64

func SetChartsCachingDuration(duration int64) {
	if duration >= 0 {
		logger.Info("Setting charts caching duration (seconds)", logger.Params{"duration": duration})
		ChartsCachingDuration = duration
	}
}

func InitCaching(db *storage.Storage) *Provider {
	if ChartsCachingDuration == 0 {
		logger.Warn("Caching only the absolutely same response", logger.Params{"caching_duration": 0})
	}
	return &Provider{DB: db}
}

func (p *Provider) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (p *Provider) SaveChartsCache(key string, data watchmarket.ChartData, timeStart int64) error {
	if data.IsEmpty() {
		return errors.New("data is empty")
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.DB.Set(key, storage.CacheData{
		RawData:      rawData,
		WasSavedTime: timeStart,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GetChartsCache(key string, timeStart int64) (watchmarket.ChartData, error) {
	var data watchmarket.ChartData

	cacheData, err := p.DB.Get(key)
	if err != nil {
		return data, err
	}

	if cacheData.Validate(timeStart, ChartsCachingDuration) {
		err = json.Unmarshal(cacheData.RawData, &data)
	}

	if err == nil && !data.IsEmpty() {
		return data, nil
	}

	err = p.DB.Delete(key)
	if err != nil {
		return watchmarket.ChartData{}, errors.New("invalid cache is not deleted")
	}

	return data, errors.New("cache is not valid")
}
