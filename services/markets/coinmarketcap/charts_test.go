package coinmarketcap

import (
	"encoding/json"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

func TestProvider_GetCoinData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, "USD", assets.Init(server.URL))
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	data, _ := provider.GetCoinData(controllers.Asset{CoinId: 60}, "USD")
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedCoinInfo, string(rawData))
}

func TestProvider_GetChartData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, server.URL, server.URL, "USD", assets.Init(server.URL))
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	data, _ := provider.GetChartData(controllers.Asset{CoinId: 60}, "USD", 1577871126)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	isSorted := sort.SliceIsSorted(data.Prices, func(i, j int) bool {
		return data.Prices[i].Date < data.Prices[j].Date
	})
	assert.True(t, isSorted)
	assert.JSONEq(t, wantedChartsSorted, string(rawData))
}

func Test_normalizeInfo(t *testing.T) {
	type args struct {
		currency string
		cmcCoin  uint
		data     ChartInfo
	}
	tests := []struct {
		name     string
		args     args
		wantInfo watchmarket.CoinDetails
	}{
		{
			"test normalize coinmarketcap chart assets 1",
			args{
				currency: "USD",
				cmcCoin:  1,
				data: ChartInfo{
					Data: map[int]ChartInfoData{
						1: {
							Rank:              1,
							Slug:              "test",
							CirculatingSupply: 111,
							TotalSupply:       222,
							Quotes: map[string]ChartInfoQuote{
								"USD": {Price: 333, Volume24: 444, MarketCap: 555},
							},
						},
					},
				},
			},
			watchmarket.CoinDetails{
				Provider:    id,
				ProviderURL: "https://coinmarketcap.com/currencies/test/",
			},
		},
		{
			"test normalize coinmarketcap chart assets 2",
			args{
				currency: "EUR",
				cmcCoin:  2,
				data: ChartInfo{
					Data: map[int]ChartInfoData{
						2: {
							Rank:              2,
							Slug:              "test",
							CirculatingSupply: 111,
							TotalSupply:       222,
							Quotes: map[string]ChartInfoQuote{
								"EUR": {Price: 333, Volume24: 444, MarketCap: 555},
							},
						},
					},
				},
			},
			watchmarket.CoinDetails{
				Provider:    id,
				ProviderURL: "https://coinmarketcap.com/currencies/test/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := normalizeInfo(tt.args.data, nil)
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(tt.wantInfo, gotInfo))
		})
	}
}

func Test_normalizeCharts(t *testing.T) {
	type args struct {
		currency string
		symbol   string
		charts   Charts
	}

	timeStr1 := "2019-12-19T18:27:23.453Z"
	d1, _ := time.Parse(time.RFC3339, timeStr1)
	timeStr2 := "2019-11-19T18:27:23.453Z"
	d2, _ := time.Parse(time.RFC3339, timeStr2)
	tests := []struct {
		name     string
		args     args
		wantInfo watchmarket.Chart
	}{
		{
			"test normalize coinmarketcap chart 1",
			args{
				currency: "USD",
				symbol:   "BTC",
				charts: Charts{
					Data: ChartQuotes{
						timeStr1: ChartQuoteValues{
							"USD": []float64{111, 222, 333},
						},
					},
				},
			},
			watchmarket.Chart{
				Provider: id,
				Prices: []watchmarket.ChartPrice{
					{
						Price: 111,
						Date:  d1.Unix(),
					},
				},
			},
		},
		{
			"test normalize coinmarketcap chart 2",
			args{
				currency: "EUR",
				symbol:   "BTC",
				charts: Charts{
					Data: ChartQuotes{
						timeStr1: ChartQuoteValues{
							"EUR": []float64{333, 444, 555},
						},
						timeStr2: ChartQuoteValues{
							"EUR": []float64{555, 666, 777},
						},
					},
				},
			},
			watchmarket.Chart{
				Provider: id,
				Prices: []watchmarket.ChartPrice{
					{
						Price: 333,
						Date:  d1.Unix(),
					},
					{
						Price: 555,
						Date:  d2.Unix(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo := normalizeCharts(tt.args.currency, tt.args.charts)
			sort.SliceStable(gotInfo.Prices, func(i, j int) bool {
				return gotInfo.Prices[i].Price < gotInfo.Prices[j].Price
			})
			if !assert.ObjectsAreEqualValues(tt.wantInfo, gotInfo) {
				t.Errorf("normalizeCharts() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func Test_getInterval(t *testing.T) {
	tests := []struct {
		name     string
		days     int
		wantInfo string
	}{
		{
			"test getInterval 1",
			1,
			"5m",
		},
		{
			"test getInterval 2",
			5,
			"5m",
		},
		{
			"test getInterval 3",
			7,
			"15m",
		},
		{
			"test getInterval 4",
			8,
			"15m",
		},
		{
			"test getInterval 5",
			30,
			"1h",
		},
		{
			"test getInterval 6",
			40,
			"1h",
		},
		{
			"test getInterval 7",
			90,
			"2h",
		},
		{
			"test getInterval 8",
			120,
			"2h",
		},
		{
			"test getInterval 9",
			360,
			"1d",
		},
		{
			"test getInterval 10",
			800,
			"1d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo := getInterval(tt.days)
			assert.Equal(t, tt.wantInfo, gotInfo)
		})
	}
}
