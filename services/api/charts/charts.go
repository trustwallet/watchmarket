package charts

import (
	"github.com/trustwallet/watchmarket/services/charts"
)

type Resolver struct {
	providers charts.Providers
	priority  map[uint]string
}

func Init(providers charts.Providers, priority map[uint]string) Resolver {
	return Resolver{providers: providers, priority: priority}
}

func (r Resolver) HandleChartsRequest(coinID uint, token, currency string, timeStart int64) charts.Data {
	p, _ := r.getProvider(0)
	result, err := p.GetChartData(coinID, token, currency, timeStart)
	if err != nil {
		return result
	}

	return result
}

func (r Resolver) getProvider(currentProvider uint) (charts.Provider, error) {
	bestProvider, _ := r.providers[r.priority[currentProvider]]
	return bestProvider, nil
}
