package storage

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type SaveResult string

const (
	SaveResultSuccess                      SaveResult = "Success"
	SaveResultStorageFailure               SaveResult = "StorageFailure"
	SaveResultSkippedLowPriority           SaveResult = "SkippedLowPriority"
	SaveResultSkippedLowPriorityOrOutdated SaveResult = "SkippedLowPriorityOrOutdated"
)

type Storage struct {
	DB
}

type DB interface {
	DeleteHM(entity, key string) error
	GetHMValue(entity, key string, value interface{}) error
	AddHM(entity, key string, value interface{}) error
	Init(host string) error
}

type Market interface {
	SaveTicker(coin *watchmarket.Ticker, pl ProviderList) (SaveResult, error)
	GetTicker(coin, token string) (*watchmarket.Ticker, error)
	SaveRates(rates watchmarket.Rates, pl ProviderList) map[SaveResult]int
	GetRate(currency string) (*watchmarket.Rate, error)
}

type Middleware interface {
	Set(key string, data CacheData) (SaveResult, error)
	Get(key string) (CacheData, error)
	Delete(key string) (SaveResult, error)
}
