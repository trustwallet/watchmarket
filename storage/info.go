package storage

import "time"

const (
	EntityInfo              = "ATLAS_INFO"
	defaultCoinCacheTimeout = 600
)

func (s Storage) GetInfo(key string) (*CoinInfo, error) {
	var info CoinInfo
	err := s.GetHMValue(EntityInfo, key, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (s Storage) SaveInfo(key string, info *CoinInfo) (SaveResult, error) {
	err := s.AddHM(EntityInfo, key, &info)
	if err != nil {
		return SaveResultAddHMFailure, err
	}
	return SaveResultSuccess, err
}

func (info *CoinInfo) IsOutdated() bool {
	timeNow := time.Now().Unix()
	return timeNow-info.Timestamp > defaultCoinCacheTimeout
}
