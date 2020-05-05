package binance_dex

//func Test_normalizeTickers(t *testing.T) {
//	type args struct {
//		prices   []*CoinPrice
//		provider string
//	}
//	tests := []struct {
//		name        string
//		args        args
//		wantTickers watchmarket.Tickers
//	}{
//		{
//			"test normalize binance_dex quote",
//			args{prices: []*CoinPrice{
//				{
//					BaseAssetName:      "RAVEN-F66",
//					QuoteAssetName:     "BNB",
//					LastPrice:          "0.00001082",
//					PriceChangePercent: "-2.2500",
//				},
//				{
//					BaseAssetName:      "SLV-986",
//					QuoteAssetName:     "BNB",
//					LastPrice:          "0.04494510",
//					PriceChangePercent: "-5.3700",
//				},
//				{
//					BaseAssetName:      "CBIX-3C9",
//					QuoteAssetName:     "TAUD-888",
//					LastPrice:          "0.00100235",
//					PriceChangePercent: "5.2700",
//				},
//			},
//				provider: "binance_dex"},
//			watchmarket.Tickers{
//				&watchmarket.Ticker{CoinName: "BNB", TokenId: "RAVEN-F66", CoinType: watchmarket.TypeToken, LastUpdate: time.Now(),
//					Price: watchmarket.TickerPrice{
//						Value:     0.00001082,
//						Change24h: -2.2500,
//						Currency:  "BNB",
//						Provider:  "binance_dex",
//					},
//				},
//				&watchmarket.Ticker{CoinName: "BNB", TokenId: "SLV-986", CoinType: watchmarket.TypeToken, LastUpdate: time.Now(),
//					Price: watchmarket.TickerPrice{
//						Value:     0.0449451,
//						Change24h: -5.3700,
//						Currency:  "BNB",
//						Provider:  "binance_dex",
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotTickers := normalizeTickers(tt.args.prices, tt.args.provider)
//			now := time.Now()
//			sort.Slice(gotTickers, func(i, j int) bool {
//				gotTickers[i].LastUpdate = now
//				gotTickers[j].LastUpdate = now
//				return gotTickers[i].Coin > gotTickers[j].Coin
//			})
//			sort.Slice(tt.wantTickers, func(i, j int) bool {
//				tt.wantTickers[i].LastUpdate = now
//				tt.wantTickers[j].LastUpdate = now
//				return tt.wantTickers[i].Coin > tt.wantTickers[j].Coin
//			})
//			assert.Equal(t, tt.wantTickers, gotTickers)
//		})
//	}
//}
