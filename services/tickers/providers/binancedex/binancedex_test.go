package binancedex

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/tickers"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
}

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.BaseUrl)
	assert.Equal(t, "binancedex", provider.ID)
}

func TestProvider_GetData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	provider := InitProvider(server.URL)
	data, err := provider.GetData()
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "BNB", data[0].CoinName)
	assert.Equal(t, uint(714), data[0].Coin)
	assert.Equal(t, tickers.Price{Value: 123, Change24h: 10, Currency: "BNB", Provider: "binancedex"}, data[0].Price)
	assert.Equal(t, tickers.CoinType("token"), data[0].CoinType)
	assert.Equal(t, "", data[0].Error)
	assert.LessOrEqual(t, data[0].LastUpdate.Unix(), time.Now().Unix())
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/v1/ticker/24hr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		p := CoinPrice{BaseAssetName: "BaseName", QuoteAssetName: BNBAsset, PriceChangePercent: "10", LastPrice: "123"}
		rawBytes, err := json.Marshal([]CoinPrice{p})
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(rawBytes); err != nil {
			panic(err)
		}
	})

	return r
}

func Test_normalizeTickers(t *testing.T) {
	type args struct {
		prices   []CoinPrice
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantTickers tickers.Tickers
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
			tickers.Tickers{
				tickers.Ticker{Coin: uint(714), CoinName: "BNB", TokenId: "RAVEN-F66", CoinType: tickers.Token, LastUpdate: time.Now(),
					Price: tickers.Price{
						Value:     0.00001082,
						Change24h: -2.2500,
						Currency:  "BNB",
						Provider:  "binancedex",
					},
				},
				tickers.Ticker{Coin: uint(714), CoinName: "BNB", TokenId: "SLV-986", CoinType: tickers.Token, LastUpdate: time.Now(),
					Price: tickers.Price{
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
