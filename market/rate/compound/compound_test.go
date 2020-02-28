package compound

import (
	"github.com/stretchr/testify/assert"
	c "github.com/trustwallet/watchmarket/market/clients/compound"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sort"
	"testing"
	"time"
)

func Test_normalizeRates(t *testing.T) {
	provider := "compound"
	tests := []struct {
		name      string
		prices    c.CoinPrices
		wantRates watchmarket.Rates
	}{
		{
			"test normalize compound rate 1",
			c.CoinPrices{
				Data: []c.CToken{
					{
						Symbol:          "cUSDC",
						UnderlyingPrice: c.Amount{Value: 0.0021},
					},
					{
						Symbol:          "cREP",
						UnderlyingPrice: c.Amount{Value: 0.02},
					},
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: 1 / 0.0021, Timestamp: 333, Provider: provider},
				watchmarket.Rate{Currency: "CREP", Rate: 1 / 0.02, Timestamp: 333, Provider: provider},
			},
		},
		{
			"test normalize compound rate 2",
			c.CoinPrices{
				Data: []c.CToken{
					{
						Symbol:          "cUSDC",
						UnderlyingPrice: c.Amount{Value: 110.0021},
					},
					{
						Symbol:          "cREP",
						UnderlyingPrice: c.Amount{Value: 110.02},
					},
				},
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "CUSDC", Rate: 1 / 110.0021, Timestamp: 123, Provider: provider},
				watchmarket.Rate{Currency: "CREP", Rate: 1 / 110.02, Timestamp: 123, Provider: provider},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRates := normalizeRates(tt.prices, provider)
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
