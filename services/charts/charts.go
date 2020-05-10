package charts

//
//import (
//	"github.com/trustwallet/watchmarket/services/charts/providers/coingecko"
//	"github.com/trustwallet/watchmarket/services/charts/providers/coinmarketcap"
//)
//
//type (
//	Provider interface {
//		GetChartData(coinID uint, token, currency string, timeStart int64) (Data, error)
//		GetCoinData(coinID uint, token, currency string) (CoinDetails, error)
//	}
//
//	Providers map[string]Provider
//)
//
//func InitProviders() Providers {
//	providers := make(map[string]Provider, 0)
//
//	cmc := coinmarketcap.InitProvider("", "", "", "")
//	cg := coingecko.InitProvider("", "")
//
//	providers[cmc.ID] = cmc
//	providers[cg.ID] = cg
//
//	return providers
//}
