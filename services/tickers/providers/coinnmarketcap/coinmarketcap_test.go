package coinnmarketcap

import (
	"testing"
)

func Test_normalizeTickers(t *testing.T) {
	//mapping := coinnmarketcap.CmcMapping{
	//	666: {{
	//		Coin: 1023,
	//		Id:   666,
	//		Type: "coin",
	//	}},
	//}
	//type args struct {
	//	prices   coinnmarketcap.CoinPrices
	//	provider string
	//}
	//tests := []struct {
	//	name        string
	//	args        args
	//	wantTickers watchmarket.Tickers
	//}{
	//	{
	//		"test normalize coinnmarketcap quote",
	//		args{prices: coinnmarketcap.CoinPrices{Data: []coinnmarketcap.Data{
	//			{Coin: coinnmarketcap.Coin{Symbol: "BTC", Id: 0}, LastUpdated: time.Unix(111, 0), Quote: coinnmarketcap.Quote{
	//				USD: coinnmarketcap.USD{Price: 223.55, PercentChange24h: 10}}},
	//			{Coin: coinnmarketcap.Coin{Symbol: "ETH", Id: 60}, LastUpdated: time.Unix(333, 0), Quote: coinnmarketcap.Quote{
	//				USD: coinnmarketcap.USD{Price: 11.11, PercentChange24h: 20}}},
	//			{Coin: coinnmarketcap.Coin{Symbol: "SWP", Id: 6969}, LastUpdated: time.Unix(444, 0), Quote: coinnmarketcap.Quote{
	//				USD: coinnmarketcap.USD{Price: 463.22, PercentChange24h: -3}},
	//				Platform: &coinnmarketcap.Platform{Coin: coinnmarketcap.Coin{Symbol: "ETH"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}},
	//			{Coin: coinnmarketcap.Coin{Symbol: "ONE", Id: 666}, LastUpdated: time.Unix(555, 0), Quote: coinnmarketcap.Quote{
	//				USD: coinnmarketcap.USD{Price: 123.09, PercentChange24h: -1.4}},
	//				Platform: &coinnmarketcap.Platform{Coin: coinnmarketcap.Coin{Symbol: "BNB"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}}}},
	//			provider: "coinnmarketcap"},
	//		watchmarket.Tickers{
	//			&watchmarket.Ticker{CoinName: "BTC", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(111, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     223.55,
	//					Change24h: 10,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "coinnmarketcap",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ETH", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(333, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     11.11,
	//					Change24h: 20,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "coinnmarketcap",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ETH", TokenId: "0x8CE9137d39326AD0cD6491fb5CC0CbA0e089b6A9", CoinType: watchmarket.TypeToken, LastUpdate: time.Unix(444, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     463.22,
	//					Change24h: -3,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "coinnmarketcap",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ONE", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(555, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     123.09,
	//					Change24h: -1.4,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "coinnmarketcap",
	//				},
	//			},
	//		},
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		gotTickers := normalizeTickers(tt.args.prices, tt.args.provider, mapping)
	//		sort.SliceStable(gotTickers, func(i, j int) bool {
	//			return gotTickers[i].LastUpdate.Unix() < gotTickers[j].LastUpdate.Unix()
	//		})
	//		if !assert.Equal(t, len(tt.wantTickers), len(gotTickers)) {
	//			t.Fatal("invalid tickers length")
	//		}
	//		for i, obj := range tt.wantTickers {
	//			assert.Equal(t, obj, gotTickers[i])
	//		}
	//	})
	//}
}
