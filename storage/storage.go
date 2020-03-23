package storage

import (
	"errors"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type SaveResult string

const (
	SaveResultSuccess                      SaveResult = "Success"
	SaveResultStorageFailure               SaveResult = "StorageFailure"
	SaveResultSkippedLowPriority           SaveResult = "SkippedLowPriority"
	SaveResultSkippedLowPriorityOrOutdated SaveResult = "SkippedLowPriorityOrOutdated"
)

var ErrNotExist = errors.New("record does not exist")

type Storage struct {
	DB
}

type DB interface {
	GetAllHM(entity string) (map[string]string, error)
	DeleteHM(entity, key string) error
	GetHMValue(entity, key string, value interface{}) error
	AddHM(entity, key string, value interface{}) error
	Init(host string) error
	InitCluster(host []string) error
}

type Market interface {
	SaveTicker(coin *watchmarket.Ticker, pl ProviderList) (SaveResult, error)
	GetTicker(coin, token string) (*watchmarket.Ticker, error)
	SaveRates(rates watchmarket.Rates, pl ProviderList) map[SaveResult]int
	GetRate(currency string) (*watchmarket.Rate, error)
}

type Caching interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	UpdateInterval(key string, interval CachedInterval) error
	GetIntervalKey(key string, time int64) (string, error)
}
