package coinmarketcap

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestProvider_GetTickers(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, "USD", assets.Init("assets.api"))
	data, err := provider.GetTickers(context.Background())
	assert.Nil(t, err)
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedTickers, string(rawData))
}

func Test_normalizeTickers(t *testing.T) {
	coinMap := CoinMap{
		Coin:    1023,
		Id:      666,
		Type:    "coin",
		TokenId: "",
	}
	type args struct {
		prices   CoinPrices
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantTickers watchmarket.Tickers
	}{
		{
			"test normalize coinmarketcap quote",
			args{prices: CoinPrices{Data: []Data{
				{Coin: Coin{Symbol: "BTC", Id: 0}, LastUpdated: time.Unix(111, 0), Quote: Quote{
					USD: USD{Price: 223.55, PercentChange24h: 10}}},
				{Coin: Coin{Symbol: "ETH", Id: 60}, LastUpdated: time.Unix(333, 0), Quote: Quote{
					USD: USD{Price: 11.11, PercentChange24h: 20}}},
				{Coin: Coin{Symbol: "SWP", Id: 6969}, LastUpdated: time.Unix(444, 0), Quote: Quote{
					USD: USD{Price: 463.22, PercentChange24h: -3}},
					Platform: Platform{Coin: Coin{Symbol: "ETH"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}},
				{Coin: Coin{Symbol: "ONE", Id: 666}, LastUpdated: time.Unix(555, 0), Quote: Quote{
					USD: USD{Price: 123.09, PercentChange24h: -1.4}},
					Platform: Platform{Coin: Coin{Symbol: "BNB"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}},
			}}, provider: "coinmarketcap"},
			watchmarket.Tickers{
				watchmarket.Ticker{Coin: watchmarket.UnknownCoinID, CoinName: "BTC", CoinType: watchmarket.Coin, LastUpdate: time.Unix(111, 0),
					Price: watchmarket.Price{
						Value:     223.55,
						Change24h: 10,
						Currency:  "USD",
						Provider:  "coinmarketcap",
					},
				},
				watchmarket.Ticker{Coin: watchmarket.UnknownCoinID, CoinName: "ETH", CoinType: watchmarket.Coin, LastUpdate: time.Unix(333, 0),
					Price: watchmarket.Price{
						Value:     11.11,
						Change24h: 20,
						Currency:  "USD",
						Provider:  "coinmarketcap",
					},
				},
				watchmarket.Ticker{Coin: watchmarket.UnknownCoinID, CoinName: "ETH", TokenId: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9", CoinType: watchmarket.Token, LastUpdate: time.Unix(444, 0),
					Price: watchmarket.Price{
						Value:     463.22,
						Change24h: -3,
						Currency:  "USD",
						Provider:  "coinmarketcap",
					},
				},
				watchmarket.Ticker{Coin: 1023, CoinName: "ONE", CoinType: watchmarket.Coin, LastUpdate: time.Unix(555, 0),
					Price: watchmarket.Price{
						Value:     123.09,
						Change24h: -1.4,
						Currency:  "USD",
						Provider:  "coinmarketcap",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTickers := normalizeTickers(tt.args.prices, []CoinMap{coinMap}, tt.args.provider, "USD")
			sort.SliceStable(gotTickers, func(i, j int) bool {
				return gotTickers[i].LastUpdate.Unix() < gotTickers[j].LastUpdate.Unix()
			})
			if !assert.Equal(t, len(tt.wantTickers), len(gotTickers)) {
				t.Fatal("invalid tickers length")
			}
			for i, obj := range tt.wantTickers {
				assert.Equal(t, obj, gotTickers[i])
			}
		})
	}
}
