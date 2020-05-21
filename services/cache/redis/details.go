package rediscache

import (
	"encoding/json"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

func (i Instance) SaveCoinDetails(key string, data watchmarket.CoinDetails, timeStart int64) error {
	if data.IsEmpty() {
		return errors.E("data is empty")
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cachingKey := i.GenerateKey(key + strconv.Itoa(int(timeStart)))
	interval := CachedInterval{
		Timestamp: timeStart,
		Duration:  int64(watchmarket.DurationToUnix(i.detailsCaching)),
		Key:       cachingKey,
	}

	err = i.updateInterval(key, interval)
	if err != nil {
		return err
	}

	err = i.redis.Set(cachingKey, rawData, i.detailsCaching)
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) GetCoinDetails(key string, timeStart int64) (watchmarket.CoinDetails, error) {
	var (
		keyInterval string
		data        watchmarket.CoinDetails
	)
	keyInterval, err := i.getIntervalKey(key, timeStart)
	if err != nil {
		return data, err
	}

	cacheData, err := i.redis.Get(keyInterval)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(cacheData, &data)

	if err == nil && !data.IsEmpty() {
		return data, nil
	}

	err = i.redis.Delete(keyInterval)
	if err != nil {
		return watchmarket.CoinDetails{}, errors.E("invalid cache is not deleted")
	}

	return watchmarket.CoinDetails{}, errors.E("cache is not valid")
}
