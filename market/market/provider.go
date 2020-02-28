package market

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
)

type Provider interface {
	Init(storage.Market) error
	GetId() string
	GetUpdateTime() string
	GetData() (watchmarket.Tickers, error)
	GetLogType() string
}

type Providers map[int]Provider

func (ps Providers) GetPriority(providerId string) int {
	for priority, provider := range ps {
		if provider.GetId() == providerId {
			return priority
		}
	}
	return -1
}
