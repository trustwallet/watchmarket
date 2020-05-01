package cmc

import (
	"testing"
)

func Test_normalizeTickers(t *testing.T) {
	//mapping := cmc.CmcMapping{
	//	666: {{
	//		Coin: 1023,
	//		Id:   666,
	//		Type: "coin",
	//	}},
	//}
	//type args struct {
	//	prices   cmc.CoinPrices
	//	provider string
	//}
	//tests := []struct {
	//	name        string
	//	args        args
	//	wantTickers watchmarket.Tickers
	//}{
	//	{
	//		"test normalize cmc quote",
	//		args{prices: cmc.CoinPrices{Data: []cmc.Data{
	//			{Coin: cmc.Coin{Symbol: "BTC", Id: 0}, LastUpdated: time.Unix(111, 0), Quote: cmc.Quote{
	//				USD: cmc.USD{Price: 223.55, PercentChange24h: 10}}},
	//			{Coin: cmc.Coin{Symbol: "ETH", Id: 60}, LastUpdated: time.Unix(333, 0), Quote: cmc.Quote{
	//				USD: cmc.USD{Price: 11.11, PercentChange24h: 20}}},
	//			{Coin: cmc.Coin{Symbol: "SWP", Id: 6969}, LastUpdated: time.Unix(444, 0), Quote: cmc.Quote{
	//				USD: cmc.USD{Price: 463.22, PercentChange24h: -3}},
	//				Platform: &cmc.Platform{Coin: cmc.Coin{Symbol: "ETH"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}},
	//			{Coin: cmc.Coin{Symbol: "ONE", Id: 666}, LastUpdated: time.Unix(555, 0), Quote: cmc.Quote{
	//				USD: cmc.USD{Price: 123.09, PercentChange24h: -1.4}},
	//				Platform: &cmc.Platform{Coin: cmc.Coin{Symbol: "BNB"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}}}},
	//			provider: "cmc"},
	//		watchmarket.Tickers{
	//			&watchmarket.Ticker{CoinName: "BTC", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(111, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     223.55,
	//					Change24h: 10,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "cmc",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ETH", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(333, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     11.11,
	//					Change24h: 20,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "cmc",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ETH", TokenId: "0x8CE9137d39326AD0cD6491fb5CC0CbA0e089b6A9", CoinType: watchmarket.TypeToken, LastUpdate: time.Unix(444, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     463.22,
	//					Change24h: -3,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "cmc",
	//				},
	//			},
	//			&watchmarket.Ticker{CoinName: "ONE", CoinType: watchmarket.TypeCoin, LastUpdate: time.Unix(555, 0),
	//				Price: watchmarket.TickerPrice{
	//					Value:     123.09,
	//					Change24h: -1.4,
	//					Currency:  watchmarket.DefaultCurrency,
	//					Provider:  "cmc",
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
