package coingecko

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestProvider_GetRates(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, "USD", assets.Init("assets.api"))
	data, err := provider.GetRates()
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedRates, string(rawData))
}

var (
	wantedRates = `[{"currency":"BTC","provider":"coingecko","rate":9696.96,"timestamp":1588871554},{"currency":"ETH","provider":"coingecko","rate":206.55,"timestamp":1588871558},{"currency":"BNB","provider":"coingecko","rate":16.76,"timestamp":1588871427},{"currency":"TRX","provider":"coingecko","rate":0.01594768,"timestamp":1588871427},{"currency":"ZOC","provider":"coingecko","rate":0.00135115,"timestamp":1588870632},{"currency":"O2T","provider":"coingecko","rate":0.00083971,"timestamp":1577332821},{"currency":"XRP-BF2","provider":"coingecko","rate":0.21726,"timestamp":1588871653},{"currency":"LVH","provider":"coingecko","rate":0.00000808,"timestamp":1588871413}]`
)

func Test_normalizeRates(t *testing.T) {
	tests := []struct {
		name      string
		prices    CoinPrices
		wantRates watchmarket.Rates
	}{
		{
			"test normalize coingecko rate 1",
			CoinPrices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 0.0021,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 0.02,
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: 0.0021, Timestamp: 333, Provider: id},
				watchmarket.Rate{Currency: "CREP", Rate: 0.02, Timestamp: 333, Provider: id},
			},
		},
		{
			"test normalize coingecko rate 2",
			CoinPrices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 110.0021,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: 110.0021, Timestamp: 123, Provider: id},
				watchmarket.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
			},
		},
		{
			"test normalize 0 rates",
			CoinPrices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 0.0,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: 0.0, Timestamp: 123, Provider: id},
				watchmarket.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
			},
		},
		{
			"test normalize negative rates (you never know...)",
			CoinPrices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: -5.0,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: -5.0, Timestamp: 123, Provider: id},
				watchmarket.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRates := normalizeRates(tt.prices, id)
			now := time.Now().Unix()
			sort.Slice(gotRates, func(i, j int) bool {
				gotRates[i].Timestamp = now
				gotRates[j].Timestamp = now
				return gotRates[i].Rate < gotRates[j].Rate
			})
			sort.Slice(tt.wantRates, func(i, j int) bool {
				tt.wantRates[i].Timestamp = now
				tt.wantRates[j].Timestamp = now
				return tt.wantRates[i].Rate < tt.wantRates[j].Rate
			})
			if !assert.ObjectsAreEqualValues(gotRates, tt.wantRates) {
				t.Errorf("normalizeRates() = %v, want %v", gotRates, tt.wantRates)
			}
		})
	}
}
