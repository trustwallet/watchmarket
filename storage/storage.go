package storage

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/storage/redis"
)

type Storage struct {
	redis.Redis
}

func New() *Storage {
	s := new(Storage)
	return s
}

type Market interface {
	SaveTicker(coin *blockatlas.Ticker, pl ProviderList) error
	GetTicker(coin, token string) (*blockatlas.Ticker, error)
	SaveRates(rates blockatlas.Rates, pl ProviderList)
	GetRate(currency string) (*blockatlas.Rate, error)
}
