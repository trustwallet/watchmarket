package storage

const EntityCache = "MARKET_CACHE"

type (
	CacheData struct {
		RawData      []byte
		WasSavedTime int64
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

func (c *CacheData) Validate(time, duration int64) bool {
	// must not be expired
	// must not be before caching
	isExpired := c.WasSavedTime+duration < time
	isBeforeCaching := c.WasSavedTime-time > 0

	return !isExpired && !isBeforeCaching
}
