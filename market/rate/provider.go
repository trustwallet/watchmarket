package rate

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
)

type Provider interface {
	Init(storage.Market) error
	FetchLatestRates() (watchmarket.Rates, error)
	GetUpdateTime() string
	GetId() string
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
