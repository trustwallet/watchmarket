package coingecko

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/rates"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api", "USD", 5)
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
	assert.Equal(t, "USD", client.currency)
	assert.Equal(t, 5, client.bucketSize)
}

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api", "USD")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.BaseUrl)
	assert.Equal(t, "coingecko", provider.ID)
	assert.Equal(t, "USD", provider.currency)
}

func TestProvider_GetData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, "USD")
	data, err := provider.GetData()
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	fmt.Println(string(rawData))
	assert.Equal(t, wantedRates, string(rawData))
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/v3/coins/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		coin1 := Coin{
			Id:        "01coin",
			Symbol:    "zoc",
			Name:      "01coin",
			Platforms: nil,
		}
		coin2 := Coin{
			Id:        "02-token",
			Symbol:    "o2t",
			Name:      "O2 Token",
			Platforms: map[string]string{"ethereum": "0xb1bafca3737268a96673a250173b6ed8f1b5b65f"},
		}
		coin3 := Coin{
			Id:        "lovehearts",
			Symbol:    "lvh",
			Name:      "LoveHearts",
			Platforms: map[string]string{"tron": "1000451"},
		}
		coin4 := Coin{
			Id:        "xrp-bep2",
			Symbol:    "xrp-bf2",
			Name:      "XRP BEP2",
			Platforms: map[string]string{"binancecoin": "XRP-BF2"},
		}
		coin5 := Coin{
			Id:        "bitcoin",
			Symbol:    "btc",
			Name:      "Bitcoin",
			Platforms: nil,
		}
		coin6 := Coin{
			Id:        "binancecoin",
			Symbol:    "bnb",
			Name:      "Binance Coin",
			Platforms: map[string]string{"binancecoin": "BNB"},
		}
		coin7 := Coin{
			Id:        "tron",
			Symbol:    "trx",
			Name:      "TRON",
			Platforms: nil,
		}
		coin8 := Coin{
			Id:        "ethereum",
			Symbol:    "eth",
			Name:      "ethereum",
			Platforms: nil,
		}

		rawBytes, err := json.Marshal(Coins([]Coin{coin1, coin2, coin3, coin4, coin5, coin6, coin7, coin8}))
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(rawBytes); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v3/coins/markets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, mockedMarketsResponse); err != nil {
			panic(err)
		}
	})

	return r
}

var (
	mockedMarketsResponse = `[ { "id": "bitcoin", "symbol": "btc", "name": "Bitcoin", "image": "https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579", "current_price": 9696.96, "market_cap": 177446468003, "market_cap_rank": 1, "total_volume": 51778003346, "high_24h": 9661.04, "low_24h": 9099.5, "price_change_24h": 459.99, "price_change_percentage_24h": 4.97984, "market_cap_change_24h": 7519669329, "market_cap_change_percentage_24h": 4.42524, "circulating_supply": 18367225.0, "total_supply": 21000000.0, "ath": 19665.39, "ath_change_percentage": -51.20717, "ath_date": "2017-12-16T00:00:00.000Z", "atl": 67.81, "atl_change_percentage": 14050.4866, "atl_date": "2013-07-06T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:12:34.220Z" }, { "id": "ethereum", "symbol": "eth", "name": "Ethereum", "image": "https://assets.coingecko.com/coins/images/279/large/ethereum.png?1547034048", "current_price": 206.55, "market_cap": 22851909019, "market_cap_rank": 2, "total_volume": 17356592769, "high_24h": 207.93, "low_24h": 199.21, "price_change_24h": -1.32300849, "price_change_percentage_24h": -0.63646, "market_cap_change_24h": -218200222.339359, "market_cap_change_percentage_24h": -0.94581, "circulating_supply": 110830501.624, "total_supply": null, "ath": 1448.18, "ath_change_percentage": -85.76061, "ath_date": "2018-01-13T00:00:00.000Z", "atl": 0.432979, "atl_change_percentage": 47526.36618, "atl_date": "2015-10-20T00:00:00.000Z", "roi": { "times": 27.445825905402025, "currency": "btc", "percentage": 2744.5825905402025 }, "last_updated": "2020-05-07T17:12:38.629Z" }, { "id": "binancecoin", "symbol": "bnb", "name": "Binance Coin", "image": "https://assets.coingecko.com/coins/images/825/large/binance-coin-logo.png?1547034615", "current_price": 16.76, "market_cap": 2474368090, "market_cap_rank": 9, "total_volume": 386293713, "high_24h": 16.96, "low_24h": 16.34, "price_change_24h": -0.15954337, "price_change_percentage_24h": -0.9431, "market_cap_change_24h": -28933861.2562523, "market_cap_change_percentage_24h": -1.15583, "circulating_supply": 147883948.0, "total_supply": 179883948.0, "ath": 39.68, "ath_change_percentage": -57.95222, "ath_date": "2019-06-22T12:20:21.894Z", "atl": 0.0398177, "atl_change_percentage": 41801.07312, "atl_date": "2017-10-19T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:10:27.413Z" }, { "id": "tron", "symbol": "trx", "name": "TRON", "image": "https://assets.coingecko.com/coins/images/1094/large/tron-logo.png?1547035066", "current_price": 0.01594768, "market_cap": 1057389936, "market_cap_rank": 17, "total_volume": 1697721558, "high_24h": 0.01620348, "low_24h": 0.0155473, "price_change_24h": -8.271e-05, "price_change_percentage_24h": -0.51593, "market_cap_change_24h": -4768379.66360569, "market_cap_change_percentage_24h": -0.44893, "circulating_supply": 66140232427.0, "total_supply": 99281283754.0, "ath": 0.231673, "ath_change_percentage": -93.12752, "ath_date": "2018-01-05T00:00:00.000Z", "atl": 0.00180434, "atl_change_percentage": 782.40789, "atl_date": "2017-11-12T00:00:00.000Z", "roi": { "times": 7.39351342638589, "currency": "usd", "percentage": 739.3513426385889 }, "last_updated": "2020-05-07T17:10:27.759Z" }, { "id": "01coin", "symbol": "zoc", "name": "01coin", "image": "https://assets.coingecko.com/coins/images/5720/large/F1nTlw9I_400x400.jpg?1547041588", "current_price": 0.00135115, "market_cap": 14384.83, "market_cap_rank": 1587, "total_volume": 888.82, "high_24h": 0.0013974, "low_24h": 0.00120559, "price_change_24h": 3.876e-05, "price_change_percentage_24h": 2.95316, "market_cap_change_24h": 480.77, "market_cap_change_percentage_24h": 3.45776, "circulating_supply": 10646360.834599, "total_supply": 65658824.0, "ath": 0.03418169, "ath_change_percentage": -96.04715, "ath_date": "2018-10-10T17:27:38.034Z", "atl": 0.00070641, "atl_change_percentage": 91.26875, "atl_date": "2020-03-16T10:22:30.944Z", "roi": null, "last_updated": "2020-05-07T16:57:12.616Z" }, { "id": "02-token", "symbol": "o2t", "name": "O2 Token", "image": "https://assets.coingecko.com/coins/images/6925/large/44429612.jpeg?1547043298", "current_price": 0.00083971, "market_cap": 0.0, "market_cap_rank": 7111, "total_volume": 69.52, "high_24h": null, "low_24h": null, "price_change_24h": null, "price_change_percentage_24h": null, "market_cap_change_24h": null, "market_cap_change_percentage_24h": null, "circulating_supply": 0.0, "total_supply": 28520100.0, "ath": 0.00439107, "ath_change_percentage": -80.87694, "ath_date": "2018-11-20T05:12:22.611Z", "atl": 0.00057411, "atl_change_percentage": 46.26319, "atl_date": "2018-11-26T00:00:00.000Z", "roi": null, "last_updated": "2019-12-26T04:00:21.046Z" }, { "id": "xrp-bep2", "symbol": "xrp-bf2", "name": "XRP BEP2", "image": "https://assets.coingecko.com/coins/images/9686/large/12-122739_xrp-logo-png-clipart.png?1570790408", "current_price": 0.21726, "market_cap": 0.0, "market_cap_rank": 5069, "total_volume": 301.34, "high_24h": 0.219438, "low_24h": 0.212267, "price_change_24h": -0.00202035, "price_change_percentage_24h": -0.92136, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 10000000.0, "ath": 0.360995, "ath_change_percentage": -40.44337, "ath_date": "2019-10-21T13:44:21.822Z", "atl": 0.115982, "atl_change_percentage": 85.36984, "atl_date": "2020-03-13T02:02:33.103Z", "roi": null, "last_updated": "2020-05-07T17:14:13.364Z" }, { "id": "lovehearts", "symbol": "lvh", "name": "LoveHearts", "image": "https://assets.coingecko.com/coins/images/9360/large/1_d3hJ7JQeQ84goeTVWLI9Qw.png?1566528108", "current_price": 8.08e-06, "market_cap": 0.0, "market_cap_rank": 5528, "total_volume": 7.87, "high_24h": 8.73e-06, "low_24h": 7.85e-06, "price_change_24h": 1.9e-07, "price_change_percentage_24h": 2.39645, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 100000000000.0, "ath": 8.596e-05, "ath_change_percentage": -90.70838, "ath_date": "2019-08-23T03:49:38.791Z", "atl": 3.13e-06, "atl_change_percentage": 155.38143, "atl_date": "2020-02-21T21:25:35.813Z", "roi": null, "last_updated": "2020-05-07T17:10:13.895Z" } ]`
	wantedRates           = `[{"currency":"BTC","percent_change_24h":"0","provider":"coingecko","rate":9696.96,"timestamp":1588871554},{"currency":"ETH","percent_change_24h":"0","provider":"coingecko","rate":206.55,"timestamp":1588871558},{"currency":"BNB","percent_change_24h":"0","provider":"coingecko","rate":16.76,"timestamp":1588871427},{"currency":"TRX","percent_change_24h":"0","provider":"coingecko","rate":0.01594768,"timestamp":1588871427},{"currency":"ZOC","percent_change_24h":"0","provider":"coingecko","rate":0.00135115,"timestamp":1588870632},{"currency":"O2T","percent_change_24h":"0","provider":"coingecko","rate":0.00083971,"timestamp":1577332821},{"currency":"XRP-BF2","percent_change_24h":"0","provider":"coingecko","rate":0.21726,"timestamp":1588871653},{"currency":"LVH","percent_change_24h":"0","provider":"coingecko","rate":0.00000808,"timestamp":1588871413}]`
)

func Test_normalizeRates(t *testing.T) {
	tests := []struct {
		name      string
		prices    Prices
		wantRates rates.Rates
	}{
		{
			"test normalize coingecko rate 1",
			Prices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 0.0021,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 0.02,
				},
			},
			rates.Rates{
				rates.Rate{Currency: "CUSDC", Rate: 0.0021, Timestamp: 333, Provider: id},
				rates.Rate{Currency: "CREP", Rate: 0.02, Timestamp: 333, Provider: id},
			},
		},
		{
			"test normalize coingecko rate 2",
			Prices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 110.0021,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			rates.Rates{
				rates.Rate{Currency: "CUSDC", Rate: 110.0021, Timestamp: 123, Provider: id},
				rates.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
			},
		},
		{
			"test normalize 0 rates",
			Prices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: 0.0,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			rates.Rates{
				rates.Rate{Currency: "CUSDC", Rate: 0.0, Timestamp: 123, Provider: id},
				rates.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
			},
		},
		{
			"test normalize negative rates (you never know...)",
			Prices{
				{
					Symbol:       "cUSDC",
					CurrentPrice: -5.0,
				},
				{
					Symbol:       "cREP",
					CurrentPrice: 110.02,
				},
			},
			rates.Rates{
				rates.Rate{Currency: "CUSDC", Rate: -5.0, Timestamp: 123, Provider: id},
				rates.Rate{Currency: "CREP", Rate: 110.02, Timestamp: 123, Provider: id},
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
