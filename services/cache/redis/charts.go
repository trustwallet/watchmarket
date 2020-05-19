package rediscache

import (
	"encoding/json"
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"time"
)

func (i Instance) GetCharts(key string, timeStart int64) (watchmarket.Chart, error) {
	var (
		keyInterval string
		data        watchmarket.Chart
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
		fmt.Println("cached")
		return data, nil
	}

	err = i.redis.Delete(keyInterval)
	if err != nil {
		return watchmarket.Chart{}, errors.E("invalid cache is not deleted")
	}

	return watchmarket.Chart{}, errors.E("cache is not valid")
}

func (i Instance) SaveCharts(key string, data watchmarket.Chart, timeStart int64) error {
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
		Duration:  int64(DurationToUnix(i.chartsCaching)),
		Key:       cachingKey,
	}

	err = i.updateInterval(key, interval)
	if err != nil {
		return err
	}

	err = i.redis.Set(cachingKey, rawData, i.chartsCaching)
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) getIntervalKey(key string, time int64) (string, error) {
	var currentIntervals []CachedInterval

	rawIntervals, err := i.redis.Get(key)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(rawIntervals, &currentIntervals)
	if err != nil {
		return "", err
	}

	var (
		results = make([]string, 0)
		counter int
	)
	for _, interval := range currentIntervals {
		if time >= interval.Timestamp && time <= interval.Timestamp+interval.Duration {
			results = append(results, interval.Key)
			counter++
		}
	}

	if len(results) == 0 {
		return "", errors.E("no suitable intervals")
	}
	return results[0], nil
}

func (i Instance) updateInterval(key string, interval CachedInterval) error {
	var currentIntervals []CachedInterval

	rawIntervals, err := i.redis.Get(key)
	if err != nil && err.Error() != "Not found" {
		return err
	}

	if err == nil {
		err = json.Unmarshal(rawIntervals, &currentIntervals)
		if err != nil {
			return err
		}
	}

	var newCurrentIntervals []CachedInterval
	for i, iv := range currentIntervals {
		if iv.Timestamp+iv.Duration != interval.Timestamp {
			newCurrentIntervals = append(newCurrentIntervals, currentIntervals[i])
		}
	}

	newCurrentIntervals = append(newCurrentIntervals, interval)

	rawNewIntervalsRaw, err := json.Marshal(newCurrentIntervals)
	if err != nil {
		return err
	}

	err = i.redis.Set(key, rawNewIntervalsRaw, i.chartsCaching)
	if err != nil {
		return err
	}
	return nil
}

func UnixToDuration(unixTime uint) time.Duration {
	return time.Duration(unixTime * 1000000000)
}

func DurationToUnix(duration time.Duration) uint {
	return uint(duration.Seconds())
}
