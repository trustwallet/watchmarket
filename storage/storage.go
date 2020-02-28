package storage

import (
	"github.com/trustwallet/blockatlas/pkg/storage/redis"
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
	redis.Redis
}

func New() *Storage {
	s := new(Storage)
	return s
}

type Market interface {
	SaveTicker(coin *watchmarket.Ticker, pl ProviderList) (SaveResult, error)
	GetTicker(coin, token string) (*watchmarket.Ticker, error)
	SaveRates(rates watchmarket.Rates, pl ProviderList) map[SaveResult]int
	GetRate(currency string) (*watchmarket.Rate, error)
}
