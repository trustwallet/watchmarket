package fixer

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sort"
	"testing"
	"time"
)

func Test_normalizeRates(t *testing.T) {
	provider := "binance_dex"
	tests := []struct {
		name      string
		latest    Latest
		wantRates watchmarket.Rates
	}{
		{
			"test normalize fixer rate 1",
			Latest{
				Timestamp: 123,
				Rates:     map[string]float64{"USD": 22.111, "BRL": 33.2, "BTC": 44.99},
				UpdatedAt: time.Now(),
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "USD", Rate: 22.111, Timestamp: 123, Provider: provider},
				watchmarket.Rate{Currency: "BRL", Rate: 33.2, Timestamp: 123, Provider: provider},
				watchmarket.Rate{Currency: "BTC", Rate: 44.99, Timestamp: 123, Provider: provider},
			},
		},
		{
			"test normalize fixer rate 2",
			Latest{
				Timestamp: 333,
				Rates:     map[string]float64{"LSK": 123.321, "IFC": 34.973, "DUO": 998.3},
				UpdatedAt: time.Now(),
			},
			watchmarket.Rates{
				watchmarket.Rate{Currency: "IFC", Rate: 34.973, Timestamp: 333, Provider: provider},
				watchmarket.Rate{Currency: "LSK", Rate: 123.321, Timestamp: 333, Provider: provider},
				watchmarket.Rate{Currency: "DUO", Rate: 998.3, Timestamp: 333, Provider: provider},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRates := normalizeRates(tt.latest, provider)
			sort.SliceStable(gotRates, func(i, j int) bool {
				return gotRates[i].Rate < gotRates[j].Rate
			})
			if !assert.ObjectsAreEqualValues(gotRates, tt.wantRates) {
				t.Errorf("normalizeRates() = %v, want %v", gotRates, tt.wantRates)
			}
		})
	}
}
