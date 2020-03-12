package ticker

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
)

type TickerProvider interface {
	Init(storage.Market) error
	GetId() string
	GetUpdateTime() string
	GetData() (watchmarket.Tickers, error)
	GetLogType() string
}

type Providers map[int]TickerProvider

func (ps Providers) GetPriority(providerId string) int {
	for priority, provider := range ps {
		if provider.GetId() == providerId {
			return priority
		}
	}
	return -1
}
