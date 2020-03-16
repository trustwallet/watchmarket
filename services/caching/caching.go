package caching

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
	"strconv"
)

type Provider struct {
	DB storage.Caching
}

var (
	ChartsCachingDuration     int64
	ChartsCachingInfoDuration int64
)

func SetChartsCachingDuration(duration int64) {
	if duration >= 0 {
		logger.Info("Setting charts caching duration (seconds)", logger.Params{"duration": duration})
		ChartsCachingDuration = duration
	}
}

func SetChartsCachingInfoDuration(duration int64) {
	if duration >= 0 {
		logger.Info("Setting charts caching INFO duration (seconds)", logger.Params{"duration": duration})
		ChartsCachingInfoDuration = duration
	}
}

func InitCaching(db *storage.Storage) *Provider {
	if ChartsCachingDuration == 0 {
		logger.Warn("Caching only the absolutely same response", logger.Params{"caching_duration": 0})
	}
	if ChartsCachingInfoDuration == 0 {
		logger.Warn("Caching INFO only the absolutely same response", logger.Params{"caching_duration": 0})
	}
	return &Provider{DB: db}
}

func (p *Provider) GenerateKey(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (p *Provider) SaveCoinInfoCache(key string, data watchmarket.ChartCoinInfo, timeStart int64) error {
	if data.IsEmpty() {
		return errors.New("data is empty")
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cachingKey := p.GenerateKey(key + strconv.Itoa(int(timeStart)))
	interval := storage.CachedInterval{
		Timestamp: timeStart,
		Duration:  ChartsCachingInfoDuration,
		Key:       cachingKey,
	}

	err = p.DB.UpdateInterval(key, interval)
	if err != nil {
		return err
	}

	err = p.DB.Set(cachingKey, rawData)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GetCoinInfoCache(key string, timeStart int64) (watchmarket.ChartCoinInfo, error) {
	var (
		keyInterval string
		data        watchmarket.ChartCoinInfo
	)
	keyInterval, err := p.DB.GetIntervalKey(key, timeStart)
	if err != nil {
		return data, err
	}

	cacheData, err := p.DB.Get(keyInterval)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(cacheData, &data)

	if err == nil && !data.IsEmpty() {
		return data, nil
	}

	err = p.DB.Delete(keyInterval)
	if err != nil {
		return watchmarket.ChartCoinInfo{}, errors.New("invalid cache is not deleted")
	}

	return watchmarket.ChartCoinInfo{}, errors.New("cache is not valid")
}

func (p *Provider) SaveChartsCache(key string, data watchmarket.ChartData, timeStart int64) error {
	if data.IsEmpty() {
		return errors.New("data is empty")
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cachingKey := p.GenerateKey(key + strconv.Itoa(int(timeStart)))
	interval := storage.CachedInterval{
		Timestamp: timeStart,
		Duration:  ChartsCachingDuration,
		Key:       cachingKey,
	}

	err = p.DB.UpdateInterval(key, interval)
	if err != nil {
		return err
	}

	err = p.DB.Set(cachingKey, rawData)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GetChartsCache(key string, timeStart int64) (watchmarket.ChartData, error) {
	var (
		keyInterval string
		data        watchmarket.ChartData
	)
	keyInterval, err := p.DB.GetIntervalKey(key, timeStart)
	if err != nil {
		return data, err
	}

	cacheData, err := p.DB.Get(keyInterval)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(cacheData, &data)

	if err == nil && !data.IsEmpty() {
		return data, nil
	}

	err = p.DB.Delete(keyInterval)
	if err != nil {
		return watchmarket.ChartData{}, errors.New("invalid cache is not deleted")
	}

	return watchmarket.ChartData{}, errors.New("cache is not valid")
}
