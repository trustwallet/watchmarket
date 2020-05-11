package charts

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/markets"
)

type Resolver struct {
	providers markets.Providers
	priority  map[uint]string
}

func Init(providers markets.Providers, priority map[uint]string) Resolver {
	return Resolver{providers: providers, priority: priority}
}

func (r Resolver) HandleChartsRequest(coinID uint, token, currency string, timeStart int64) watchmarket.Data {
	p, _ := r.getProvider(0)
	result, err := p.GetChartData(coinID, token, currency, timeStart)
	if err != nil {
		return result
	}

	return result
}

func (r Resolver) getProvider(currentProvider uint) (markets.Provider, error) {
	bestProvider, _ := r.providers[r.priority[currentProvider]]
	return bestProvider, nil
}
