package binancedex

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestProvider_GetTickers(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	provider := InitProvider(server.URL)
	data, err := provider.GetTickers()
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "BNB", data[0].CoinName)
	assert.Equal(t, uint(714), data[0].Coin)
	assert.Equal(t, watchmarket.Price{Value: 123, Change24h: 10, Currency: "BNB", Provider: "binancedex"}, data[0].Price)
	assert.Equal(t, watchmarket.CoinType("token"), data[0].CoinType)
	assert.Equal(t, "", data[0].Error)
	assert.LessOrEqual(t, data[0].LastUpdate.Unix(), time.Now().Unix())
}

func Test_normalizeTickers(t *testing.T) {
	type args struct {
		prices   []CoinPrice
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantTickers watchmarket.Tickers
	}{
		{
			"test normalize binancedex quote",
			args{prices: []CoinPrice{
				{
					BaseAssetName:      "RAVEN-F66",
					QuoteAssetName:     "BNB",
					LastPrice:          "0.00001082",
					PriceChangePercent: "-2.2500",
				},
				{
					BaseAssetName:      "SLV-986",
					QuoteAssetName:     "BNB",
					LastPrice:          "0.04494510",
					PriceChangePercent: "-5.3700",
				},
				{
					BaseAssetName:      "CBIX-3C9",
					QuoteAssetName:     "TAUD-888",
					LastPrice:          "0.00100235",
					PriceChangePercent: "5.2700",
				},
			},
				provider: "binancedex"},
			watchmarket.Tickers{
				watchmarket.Ticker{Coin: uint(714), CoinName: "BNB", TokenId: "RAVEN-F66", CoinType: watchmarket.Token, LastUpdate: time.Now(),
					Price: watchmarket.Price{
						Value:     0.00001082,
						Change24h: -2.2500,
						Currency:  "BNB",
						Provider:  "binancedex",
					},
				},
				watchmarket.Ticker{Coin: uint(714), CoinName: "BNB", TokenId: "SLV-986", CoinType: watchmarket.Token, LastUpdate: time.Now(),
					Price: watchmarket.Price{
						Value:     0.0449451,
						Change24h: -5.3700,
						Currency:  "BNB",
						Provider:  "binancedex",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTickers := normalizeTickers(tt.args.prices, tt.args.provider)
			now := time.Now()
			sort.Slice(gotTickers, func(i, j int) bool {
				gotTickers[i].LastUpdate = now
				gotTickers[j].LastUpdate = now
				return gotTickers[i].Coin < gotTickers[j].Coin
			})
			sort.Slice(tt.wantTickers, func(i, j int) bool {
				tt.wantTickers[i].LastUpdate = now
				tt.wantTickers[j].LastUpdate = now
				return tt.wantTickers[i].Coin < tt.wantTickers[j].Coin
			})
			assert.Equal(t, tt.wantTickers, gotTickers)
		})
	}
}
