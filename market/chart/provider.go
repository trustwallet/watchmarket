package chart

import (
	watchmarket "github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Provider interface {
	GetId() string
	GetChartData(coin uint, token string, currency string, timeStart int64) (watchmarket.ChartData, error)
	GetCoinData(coin uint, token string, currency string) (watchmarket.ChartCoinInfo, error)
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
