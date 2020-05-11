package coinmarketcap

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/rates"
	"math/big"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestProvider_GetRates(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, server.URL, "", "USD")
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
		wantRates rates.Rates
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
			rates.Rates{
				rates.Rate{Currency: "ETH", Rate: 11.11, Timestamp: 333, Provider: provider, PercentChange24h: *big.NewFloat(-1.22)},
				rates.Rate{Currency: "BTC", Rate: 223.5, Timestamp: 333, Provider: provider, PercentChange24h: *big.NewFloat(0.33)},
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
			rates.Rates{
				rates.Rate{Currency: "XRP", Rate: 0.4687, Timestamp: 123, Provider: provider, PercentChange24h: *big.NewFloat(0)},
				rates.Rate{Currency: "BNB", Rate: 30.333, Timestamp: 123, Provider: provider, PercentChange24h: *big.NewFloat(2.1)},
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

var (
	wantedRates = `[{"currency":"BTC","percent_change_24h":"5.47477","provider":"coinmarketcap","rate":9862.53985763,"timestamp":1588890514},{"currency":"ETH","percent_change_24h":"2.48595","provider":"coinmarketcap","rate":213.544073721,"timestamp":1588890505},{"currency":"XRP","percent_change_24h":"0.191592","provider":"coinmarketcap","rate":0.219026524813,"timestamp":1588890545},{"currency":"BCH","percent_change_24h":"1.42361","provider":"coinmarketcap","rate":253.196839314,"timestamp":1588890547},{"currency":"BSV","percent_change_24h":"0.471004","provider":"coinmarketcap","rate":210.045381795,"timestamp":1588890553},{"currency":"LTC","percent_change_24h":"1.09024","provider":"coinmarketcap","rate":47.5118592612,"timestamp":1588890544},{"currency":"BNB","percent_change_24h":"0.906496","provider":"coinmarketcap","rate":17.1140649277,"timestamp":1588890547},{"currency":"EOS","percent_change_24h":"-0.562061","provider":"coinmarketcap","rate":2.76250457645,"timestamp":1588890546},{"currency":"XTZ","percent_change_24h":"1.86731","provider":"coinmarketcap","rate":2.78977181773,"timestamp":1588890547},{"currency":"XLM","percent_change_24h":"0.655066","provider":"coinmarketcap","rate":0.0728177647038,"timestamp":1588890545},{"currency":"ADA","percent_change_24h":"0.826166","provider":"coinmarketcap","rate":0.0510633352483,"timestamp":1588890547},{"currency":"XMR","percent_change_24h":"6.70904","provider":"coinmarketcap","rate":63.6958970684,"timestamp":1588890543},{"currency":"TRX","percent_change_24h":"0.328149","provider":"coinmarketcap","rate":0.0161019012186,"timestamp":1588890547},{"currency":"ETC","percent_change_24h":"-0.447855","provider":"coinmarketcap","rate":7.15171701655,"timestamp":1588890544},{"currency":"DASH","percent_change_24h":"-0.531529","provider":"coinmarketcap","rate":79.2516409705,"timestamp":1588890544}]`
)
