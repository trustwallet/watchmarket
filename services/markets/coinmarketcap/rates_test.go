package coinmarketcap

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"math/big"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestProvider_GetRates(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, server.URL, "", assets.NewClient("assets.api"))
	data, err := provider.GetRates()
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedRates, string(rawData))
}

func Test_normalizeRates(t *testing.T) {
	provider := "coinmarketcap"
	tests := []struct {
		name      string
		prices    CoinPrices
		wantRates watchmarket.Rates
	}{
		{
			"test normalize coinmarketcap rate 1",
			CoinPrices{
				Data: []Data{
					{
						Coin: Coin{
							Symbol: "BTC",
						},
						Quote: Quote{
							USD: USD{
								Price:            223.5,
								PercentChange24h: 0.33,
							},
						},
						LastUpdated: time.Unix(333, 0),
					},
					{
						Coin: Coin{
							Symbol: "ETH",
						},
						Quote: Quote{
							USD: USD{
								Price:            11.11,
								PercentChange24h: -1.22,
							},
						},
						LastUpdated: time.Unix(333, 0),
					},
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "ETH", Rate: 11.11, Timestamp: 333, Provider: provider, PercentChange24h: *big.NewFloat(-1.22)},
				watchmarket.Rate{Currency: "BTC", Rate: 223.5, Timestamp: 333, Provider: provider, PercentChange24h: *big.NewFloat(0.33)},
			},
		},
		{
			"test normalize coinmarketcap rate 2",
			CoinPrices{
				Data: []Data{
					{
						Coin: Coin{
							Symbol: "BNB",
						},
						Quote: Quote{
							USD: USD{
								Price:            30.333,
								PercentChange24h: 2.1,
							},
						},
						LastUpdated: time.Unix(123, 0),
					},
					{
						Coin: Coin{
							Symbol: "XRP",
						},
						Quote: Quote{
							USD: USD{
								Price: 0.4687,
							},
						},
						LastUpdated: time.Unix(123, 0),
					},
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "XRP", Rate: 0.4687, Timestamp: 123, Provider: provider, PercentChange24h: *big.NewFloat(0)},
				watchmarket.Rate{Currency: "BNB", Rate: 30.333, Timestamp: 123, Provider: provider, PercentChange24h: *big.NewFloat(2.1)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRates := normalizeRates(tt.prices, provider)
			sort.SliceStable(gotRates, func(i, j int) bool {
				return gotRates[i].Rate < gotRates[j].Rate
			})
			if !assert.ObjectsAreEqualValues(gotRates, tt.wantRates) {
				t.Errorf("normalizeRates() = %v, want %v", gotRates, tt.wantRates)
			}
		})
	}
}
