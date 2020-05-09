package coingecko


import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/market/clients/coingecko"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)


func TestNewClient(t *testing.T) {
	client := NewClient("web.api")
	assert.NotNil(t, client)
	assert.Equal(t, "web.api", client.BaseUrl)
}

func TestInitProvider(t *testing.T) {
	provider := InitProvider("web.api", "info.api")
	assert.NotNil(t, provider)
	assert.Equal(t, "web.api", provider.client.BaseUrl)
	assert.Equal(t, "info.api", provider.info.BaseUrl)
	assert.Equal(t, "coingecko", provider.ID)
}

func Test_normalizeInfo(t *testing.T) {
	type args struct {
		data CoinPrice
	}
	tests := []struct {
		name     string
		args     args
		wantInfo watchmarket.ChartCoinInfo
		wantInfo charts.CoinDetails}{
		{
			"test normalize coingecko chart info 1",
			args{
				data: coingecko.CoinPrice{
				data: CoinPrice{
					MarketCap:         555,
					TotalVolume:       444,
					CirculatingSupply: 111,
					TotalSupply:       222,
				},
			},
			watchmarket.ChartCoinInfo{
				Provider:          id,
			charts.CoinDetails{
				Vol24:             444,
				MarketCap:         555,
				CirculatingSupply: 111
				TotalSupply:       222,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo := normalizeInfo(tt.args.data)
			assert.True(t, reflect.DeepEqual(tt.wantInfo, gotInfo))
		})
	}
}

func Test_normalizeCharts(t *testing.T) {
	type args struct {

		charts coingecko.Charts
		chartsList Charts

	}

	timeStr1 := "2019-12-19T18:27:23.453Z"
	d1, _ := time.Parse(time.RFC3339, timeStr1)
	timeStr2 := "2019-11-19T18:27:23.453Z"
	d2, _ := time.Parse(time.RFC3339, timeStr2)
	tests := []struct {
		name     string
		args     args
		wantInfo watchmarket.ChartData

		wantInfo charts.Data

	}{
		{
			"test normalize coingecko chart 1",
			args{

				charts: coingecko.Charts{
					Prices: []coingecko.ChartVolume{

				chartsList: Charts{
					Prices: []Volume{

						[]float64{float64(d1.UnixNano() / int64(time.Millisecond)), 222},
						[]float64{float64(d2.UnixNano() / int64(time.Millisecond)), 333},
					},
				},
			},
			watchmarket.ChartData{
				Provider: id,
				Prices: []watchmarket.ChartPrice{
					{
						Price: 222,
						Date:  d1.Unix(),
					},

			charts.Data{
				Prices: []charts.Price{

					{
						Price: 333,
						Date:  d2.Unix(),
					},

					{
						Price: 222,
						Date:  d1.Unix(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotInfo := normalizeCharts(tt.args.charts)

			gotInfo := normalizeCharts(tt.args.chartsList)

			assert.True(t, reflect.DeepEqual(tt.wantInfo, gotInfo))
		})
	}
}


//func Test_normalizeInfo(t *testing.T) {
//	type args struct {
//		data coingecko.CoinPrice
//	}
//	tests := []struct {
//		name     string
//		args     args
//		wantInfo watchmarket.ChartCoinInfo
//	}{
//		{
//			"test normalize coingecko chart info 1",
//			args{
//				data: coingecko.CoinPrice{
//					MarketCap:         555,
//					TotalVolume:       444,
//					CirculatingSupply: 111,
//					TotalSupply:       222,
//				},
//			},
//			watchmarket.ChartCoinInfo{
//				Vol24:             444,
//				MarketCap:         555,
//				CirculatingSupply: 111,
//				TotalSupply:       222,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotInfo := normalizeInfo(tt.args.data)
//			assert.True(t, reflect.DeepEqual(tt.wantInfo, gotInfo))
//		})
//	}
//}
//
//func Test_normalizeCharts(t *testing.T) {
//	type args struct {
//		charts coingecko.Charts
//	}
//
//	timeStr1 := "2019-12-19T18:27:23.453Z"
//	d1, _ := time.Parse(time.RFC3339, timeStr1)
//	timeStr2 := "2019-11-19T18:27:23.453Z"
//	d2, _ := time.Parse(time.RFC3339, timeStr2)
//	tests := []struct {
//		name     string
//		args     args
//		wantInfo watchmarket.ChartData
//	}{
//		{
//			"test normalize coingecko chart 1",
//			args{
//				charts: coingecko.Charts{
//					Prices: []coingecko.ChartVolume{
//						[]float64{float64(d1.UnixNano() / int64(time.Millisecond)), 222},
//						[]float64{float64(d2.UnixNano() / int64(time.Millisecond)), 333},
//					},
//				},
//			},
//			watchmarket.ChartData{
//				Prices: []watchmarket.ChartPrice{
//					{
//						Price: 222,
//						Date:  d1.Unix(),
//					},
//					{
//						Price: 333,
//						Date:  d2.Unix(),
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotInfo := normalizeCharts(tt.args.charts)
//			assert.True(t, reflect.DeepEqual(tt.wantInfo, gotInfo))
//		})
//	}
//}

