package coinmarketcap

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

func TestProvider_GetRates(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, watchmarket.DefaultCurrency, assets.Init("assets.api"))
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	data, err := provider.GetRates()
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedRates, string(rawData))
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
				watchmarket.Rate{Currency: "ETH", Rate: watchmarket.TruncateWithPrecision(11.11, watchmarket.DefaultPrecision), Timestamp: 333, Provider: provider, PercentChange24h: float64(-1.22)},
				watchmarket.Rate{Currency: "BTC", Rate: watchmarket.TruncateWithPrecision(223.5, watchmarket.DefaultPrecision), Timestamp: 333, Provider: provider, PercentChange24h: float64(0.33)},
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
				watchmarket.Rate{Currency: "XRP", Rate: watchmarket.TruncateWithPrecision(0.4687, watchmarket.DefaultPrecision), Timestamp: 123, Provider: provider, PercentChange24h: float64(0)},
				watchmarket.Rate{Currency: "BNB", Rate: watchmarket.TruncateWithPrecision(30.333, watchmarket.DefaultPrecision), Timestamp: 123, Provider: provider, PercentChange24h: float64(2.1)},
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
