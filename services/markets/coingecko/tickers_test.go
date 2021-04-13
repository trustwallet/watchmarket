package coingecko

import (
	"encoding/json"
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

	provider := InitProvider(server.URL, "", "USD", assets.Init("assets.api"))
	data, err := provider.GetTickers()
	assert.Nil(t, err)
	res, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedTickers, string(res))
}

func Test_normalizeTickers(t *testing.T) {
	coinsList := Coins{
		Coin{
			Id:        "ethereum",
			Symbol:    "eth",
			Name:      "eth",
			Platforms: nil,
		},
		Coin{
			Id:        "bitcoin",
			Symbol:    "btc",
			Name:      "btc",
			Platforms: nil,
		},
		Coin{
			Id:     "cREP",
			Symbol: "cREP",
			Name:   "cREP",
			Platforms: Platforms{
				"ethereum": "0x158079ee67fce2f58472a96584a73c7ab9ac95c1",
			},
		},
		Coin{
			Id:     "cUSDC",
			Symbol: "cUSDC",
			Name:   "cUSDC",
			Platforms: Platforms{
				"ethereum": "0x39aa39c021dfbae8fac545936693ac917d5e7563",
			},
		},
	}

	m := Provider{}

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
			"test normalize coingecko quote",
			args{prices: CoinPrices{
				{
					Id:           "cUSDC",
					Symbol:       "cUSDC",
					CurrentPrice: 0.0021,
					MarketCap:    2,
					TotalVolume:  5000,
				},
				{
					Id:           "cREP",
					Symbol:       "cREP",
					CurrentPrice: 0.02,
					MarketCap:    1,
					TotalVolume:  5000,
				},
			}, provider: id},
			watchmarket.Tickers{
				watchmarket.Ticker{Coin: 60, CoinName: "ETH", Volume: 5000, MarketCap: 2, TokenId: "0x39aa39c021dfbae8fac545936693ac917d5e7563", CoinType: watchmarket.Token, LastUpdate: time.Unix(222, 0),
					Price: watchmarket.Price{
						Value:    0.0021,
						Currency: "USD",
						Provider: id,
					},
				},
				watchmarket.Ticker{Coin: 60, CoinName: "ETH", Volume: 5000, MarketCap: 1, TokenId: "0x158079ee67fce2f58472a96584a73c7ab9ac95c1", CoinType: watchmarket.Token, LastUpdate: time.Unix(444, 0),
					Price: watchmarket.Price{
						Value:    0.02,
						Currency: "USD",
						Provider: id,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTickers := m.normalizeTickers(tt.args.prices, coinsList, id, "USD")
			now := time.Now()
			sort.Slice(gotTickers, func(i, j int) bool {
				gotTickers[i].LastUpdate = now
				gotTickers[j].LastUpdate = now
				return gotTickers[i].Coin > gotTickers[j].Coin
			})
			sort.Slice(tt.wantTickers, func(i, j int) bool {
				tt.wantTickers[i].LastUpdate = now
				tt.wantTickers[j].LastUpdate = now
				return tt.wantTickers[i].Coin > tt.wantTickers[j].Coin
			})
			assert.Equal(t, tt.wantTickers, gotTickers)
		})
	}
}

func Test_createTicker(t *testing.T) {
	prices := make(CoinPrices, 0)
	prices = append(prices, CoinPrice{
		Id:                       "SH",
		Symbol:                   "shitcoin",
		CurrentPrice:             0.00000001,
		MarketCap:                1,
		PriceChangePercentage24h: 1,
		TotalVolume:              5000,
	})
	prices = append(prices, CoinPrice{
		Id:                       "SH",
		Symbol:                   "shitcoin",
		CurrentPrice:             0.00000001,
		MarketCap:                5000,
		PriceChangePercentage24h: 1,
		TotalVolume:              5000,
	})

	emptyTicker := watchmarket.Ticker{
		Price: watchmarket.Price{
			Value:     0.00000001,
			Change24h: 1,
		},
	}

	normalTicker := watchmarket.Ticker{
		Price: watchmarket.Price{
			Value:     0.00000001,
			Change24h: 1,
		},
	}

	wantedTickers := make(watchmarket.Tickers, 0)
	wantedTickers = append(wantedTickers, emptyTicker)
	wantedTickers = append(wantedTickers, normalTicker)

	for i, price := range prices {
		ticker := createTicker(price, watchmarket.Token, watchmarket.UnknownCoinID, "shitcoin", "shitcoinID", watchmarket.CoinGecko, "USD")

		assert.Equal(t, ticker.Price.Value, wantedTickers[i].Price.Value)
		assert.Equal(t, ticker.Price.Change24h, wantedTickers[i].Price.Change24h)
	}

}
