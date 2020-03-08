package chart

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type ChartProvider interface {
	GetId() string
	GetChartData(coin uint, token string, currency string, timeStart int64) (watchmarket.ChartData, error)
	GetCoinData(coin uint, token string, currency string) (watchmarket.ChartCoinInfo, error)
}

type ChartProviders map[int]ChartProvider

func (ps ChartProviders) GetPriority(providerId string) int {
	for priority, provider := range ps {
		if provider.GetId() == providerId {
			return priority
		}
	}
	return -1
}
