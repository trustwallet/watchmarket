package storage

import (
	"net/http"
	"time"
)

const EntityCache = "MARKET_CACHE"

type (
	CacheResponse struct {
		Status int
		Header http.Header
		Data   []byte
	}

	CacheData struct {
		Response CacheResponse
		Expired  int64
	}
)

func (s *Storage) Set(key string, data CacheData) (SaveResult, error) {
	err := s.AddHM(EntityCache, key, &data)
	if err != nil {
		return SaveResultStorageFailure, err
	}
	return SaveResultSuccess, nil
}

func (s *Storage) Get(key string) (CacheData, error) {
	var cd CacheData
	err := s.GetHMValue(EntityCache, key, &cd)
	if err != nil {
		return CacheData{}, err
	}
	return cd, nil
}

func (s *Storage) Delete(key string) (SaveResult, error) {
	err := s.DeleteHM(EntityCache, key)
	if err != nil {
		return SaveResultStorageFailure, err
	}

	return SaveResultSuccess, nil
}

func (c *CacheData) IsExpired() bool {
	return c.Expired < time.Now().Unix()
}
