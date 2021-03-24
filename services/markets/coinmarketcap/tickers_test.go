package coinmarketcap

import (
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

func TestProvider_GetTickers(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, "USD", assets.Init("assets.api"))
	data, err := provider.GetTickers()
	assert.Nil(t, err)
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	assert.Nil(t, err)

	wantedTicker := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BTC",
		TokenId:  "",
		CoinType: "coin",
		Price: watchmarket.Price{
			Change24h: 5.47477,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     9862.53985763,
		},
		LastUpdate: time.Time{},
		Error:      "",
		Volume:     0,
		MarketCap:  0,
		ShowOption: 0,
	}
	var isOk bool
	for _, d := range data {
		if d.Coin == wantedTicker.Coin && d.CoinName == wantedTicker.CoinName && d.Price == wantedTicker.Price {
			isOk = true
		}
	}
	assert.True(t, isOk)
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
						Provider:  watchmarket.CoinMarketCap,
					},
				},
				watchmarket.Ticker{Coin: watchmarket.UnknownCoinID, CoinName: "ETH", CoinType: watchmarket.Coin, LastUpdate: time.Unix(333, 0),
					Price: watchmarket.Price{
						Value:     11.11,
						Change24h: 20,
						Currency:  "USD",
						Provider:  watchmarket.CoinMarketCap,
					},
				},
				watchmarket.Ticker{Coin: watchmarket.UnknownCoinID, CoinName: "ETH", TokenId: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9", CoinType: watchmarket.Token, LastUpdate: time.Unix(444, 0),
					Price: watchmarket.Price{
						Value:     463.22,
						Change24h: -3,
						Currency:  "USD",
						Provider:  watchmarket.CoinMarketCap,
					},
				},
				watchmarket.Ticker{Coin: 1023, CoinName: "ONE", CoinType: watchmarket.Coin, LastUpdate: time.Unix(555, 0),
					Price: watchmarket.Price{
						Value:     123.09,
						Change24h: -1.4,
						Currency:  "USD",
						Provider:  watchmarket.CoinMarketCap,
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
