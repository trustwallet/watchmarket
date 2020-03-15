package storage

import (
	"errors"
)

type CachedInterval struct {
	Timestamp int64
	Duration  int64
	Key       string
}

const EntityInterval = "MARKET_INTERVAL"

func (s *Storage) GetIntervalKey(key string, time int64) (string, error) {
	var currentIntervals []CachedInterval
	err := s.GetHMValue(EntityInterval, key, &currentIntervals)
	if err != nil {
		return "", err
	}

	var results = make([]string, 0)
	var counter int
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

func (s *Storage) UpdateInterval(key string, interval CachedInterval) error {
	var currentIntervals []CachedInterval
	err := s.GetHMValue(EntityInterval, key, &currentIntervals)
	if err != nil && err.Error() != "record does not exist" {
		return err
	}

	var newCurrentIntervals []CachedInterval
	for i, iv := range currentIntervals {
		if iv.Timestamp+iv.Duration != interval.Timestamp {
			newCurrentIntervals = append(newCurrentIntervals, currentIntervals[i])
		}
	}

	newCurrentIntervals = append(newCurrentIntervals, interval)

	err = s.AddHM(EntityInterval, key, &newCurrentIntervals)
	if err != nil {
		return err
	}
	return nil
}
