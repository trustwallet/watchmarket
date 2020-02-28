package storage

import (
	"github.com/trustwallet/blockatlas/pkg/storage/redis"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Storage struct {
	redis.Redis
}

func New() *Storage {
	s := new(Storage)
	return s
}

type Market interface {
	SaveTicker(coin *watchmarket.Ticker, pl ProviderList) error
	GetTicker(coin, token string) (*watchmarket.Ticker, error)
	SaveRates(rates watchmarket.Rates, pl ProviderList)
	GetRate(currency string) (*watchmarket.Rate, error)
}
