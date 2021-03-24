package worker

import (
	"reflect"
	"testing"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func TestFilterRates(t *testing.T) {
	type args struct {
		rates          []watchmarket.Rate
		cryptoCurrency map[string]bool
	}
	tests := []struct {
		name string
		args args
		want []watchmarket.Rate
	}{
		{
			"Test only include allow rates",
			args{
				rates: []watchmarket.Rate{
					{
						Currency: "USD",
						Provider: watchmarket.Fixer,
					},
					{
						Currency: "BTC",
						Provider: watchmarket.CoinMarketCap,
					},
					{
						Currency: "TRX",
						Provider: watchmarket.CoinGecko,
					},
					{
						Currency: "RU",
						Provider: watchmarket.CoinGecko,
					},
					{
						Currency: "ETH",
						Provider: watchmarket.CoinGecko,
					},
				},
				cryptoCurrency: map[string]bool{"BTC": true, "ETH": true},
			},
			[]watchmarket.Rate{
				{
					Currency: "USD",
					Provider: watchmarket.Fixer,
				},
				{
					Currency: "BTC",
					Provider: watchmarket.CoinMarketCap,
				},
				{
					Currency: "ETH",
					Provider: watchmarket.CoinGecko,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterRates(tt.args.rates, tt.args.cryptoCurrency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterRates() = %v, want %v", got, tt.want)
			}
		})
	}
}
