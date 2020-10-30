package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
)

type CachedInterval struct {
	Timestamp int64
	Duration  int64
	Key       string
}

func (i Instance) GetWithTime(key string, time int64, ctx context.Context) ([]byte, error) {
	var (
		keyInterval string
		data        []byte
	)
	keyInterval, err := i.getIntervalKey(key, time, ctx)
	if err != nil {
		return data, err
	}

	cacheData, err := i.redis.Get(keyInterval, ctx)
	if err == nil {
		return cacheData, err
	}

	err = i.redis.Delete(keyInterval, ctx)
	if err != nil {
		return data, errors.New("invalid cache is not deleted")
	}

	return data, errors.New("cache is not valid")
}

func (i Instance) SetWithTime(key string, data []byte, time int64, ctx context.Context) error {
	if data == nil {
		return errors.New("data is empty")
	}

	cachingKey := i.GenerateKey(key + strconv.Itoa(int(time)))
	interval := CachedInterval{
		Timestamp: time,
		Duration:  int64(watchmarket.DurationToUnix(i.cachingPeriod)),
		Key:       cachingKey,
	}

	err := i.updateInterval(key, interval, ctx)
	if err != nil {
		return err
	}

	err = i.redis.Set(cachingKey, data, i.cachingPeriod, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) getIntervalKey(key string, time int64, ctx context.Context) (string, error) {
	var currentIntervals []CachedInterval

	rawIntervals, err := i.redis.Get(key, ctx)
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
		return "", errors.New("no suitable intervals")
	}
	return results[0], nil
}

func (i Instance) updateInterval(key string, interval CachedInterval, ctx context.Context) error {
	var currentIntervals []CachedInterval

	rawIntervals, err := i.redis.Get(key, ctx)
	if err != nil && err.Error() != "not found" {
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

	err = i.redis.Set(key, rawNewIntervalsRaw, i.cachingPeriod, ctx)
	if err != nil {
		return err
	}
	return nil
}
