package coinmarketcap

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/rates"
	"github.com/trustwallet/watchmarket/services/tickers/providers/coinmarketcap"
	"math/big"
	"sort"
	"testing"
	"time"
)

func Test_normalizeRates(t *testing.T) {
	provider := "coinmarketcap"
	tests := []struct {
		name      string
		prices    coinmarketcap.CoinPrices
		wantRates rates.Rates
	}{
		{
			"test normalize coinmarketcap rate 1",
			coinmarketcap.CoinPrices{
				Data: []coinmarketcap.Data{
					{
						Coin: coinmarketcap.Coin{
							Symbol: "BTC",
						},
						Quote: coinmarketcap.Quote{
							USD: coinmarketcap.USD{
								Price:            223.5,
								PercentChange24h: 0.33,
							},
						},
						LastUpdated: time.Unix(333, 0),
					},
					{
						Coin: coinmarketcap.Coin{
							Symbol: "ETH",
						},
						Quote: coinmarketcap.Quote{
							USD: coinmarketcap.USD{
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
			coinmarketcap.CoinPrices{
				Data: []coinmarketcap.Data{
					{
						Coin: coinmarketcap.Coin{
							Symbol: "BNB",
						},
						Quote: coinmarketcap.Quote{
							USD: coinmarketcap.USD{
								Price:            30.333,
								PercentChange24h: 2.1,
							},
						},
						LastUpdated: time.Unix(123, 0),
					},
					{
						Coin: coinmarketcap.Coin{
							Symbol: "XRP",
						},
						Quote: coinmarketcap.Quote{
							USD: coinmarketcap.USD{
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
