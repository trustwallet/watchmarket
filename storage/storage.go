package storage

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type SaveResult string
const(
	SaveResultSuccess                      SaveResult = "Success"
	SaveResultStorageFailure               SaveResult = "StorageFailure"
	SaveResultSkippedLowPriority           SaveResult = "SkippedLowPriority"
	SaveResultSkippedLowPriorityOrOutdated SaveResult = "SkippedLowPriorityOrOutdated"
)

type Storage struct {
	DB
}

type DB interface {
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
