package coingecko

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/tickers"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
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
	assert.NotNil(t, data)
	assert.True(t, verifyTickers(t, wantedTickers, data))
}

func verifyTickers(t *testing.T, wantedTickers, givenTickers tickers.Tickers) bool {
	assert.Equal(t, len(givenTickers), len(wantedTickers))
	var counter = 0
	for _, w := range wantedTickers {
		for _, g := range givenTickers {
			if w.CoinName == g.CoinName && w.TokenId == g.TokenId && w.Price == g.Price {
				assert.Equal(t, w.Coin, g.Coin)
				assert.Equal(t, w.TokenId, g.TokenId)
				assert.Equal(t, w.Price, g.Price)
				assert.Equal(t, w.CoinName, g.CoinName)
				assert.Equal(t, w.Error, g.Error)
				assert.Equal(t, w.CoinType, g.CoinType)
				counter++
			}
		}
	}
	if counter == len(givenTickers)-1 {
		return true
	}
	return false
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

func Test_normalizeTickers(t *testing.T) {
	coinsList := Coins{
		Coin{
			Id:        "ethtereum",
			Symbol:    "eth",
			Name:      "eth",
			Platforms: nil,
		},
		Coin{
			Id:        "bitcoin",
			Symbol:    "btc",
			Name:      "btc",
			Platforms: nil,
		},
		Coin{
			Id:     "cREP",
			Symbol: "cREP",
			Name:   "cREP",
			Platforms: Platforms{
				"ethtereum": "0x158079ee67fce2f58472a96584a73c7ab9ac95c1",
			},
		},
		Coin{
			Id:     "cUSDC",
			Symbol: "cUSDC",
			Name:   "cUSDC",
			Platforms: Platforms{
				"ethtereum": "0x39aa39c021dfbae8fac545936693ac917d5e7563",
			},
		},
	}

	m := Provider{}

	type args struct {
		prices   CoinPrices
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantTickers tickers.Tickers
	}{
		{
			"test normalize coingecko quote",
			args{prices: CoinPrices{
				{
					Id:           "cUSDC",
					Symbol:       "cUSDC",
					CurrentPrice: 0.0021,
					MarketCap:    2,
					TotalVolume:  5000,
				},
				{
					Id:           "cREP",
					Symbol:       "cREP",
					CurrentPrice: 0.02,
					MarketCap:    1,
					TotalVolume:  5000,
				},
			}, provider: id},
			tickers.Tickers{
				tickers.Ticker{CoinName: "ETH", TokenId: "0x39aa39c021dfbae8fac545936693ac917d5e7563", CoinType: tickers.Token, LastUpdate: time.Unix(222, 0),
					Price: tickers.Price{
						Value:    0.0021,
						Currency: "USD",
						Provider: id,
					},
				},
				tickers.Ticker{CoinName: "ETH", TokenId: "0x158079ee67fce2f58472a96584a73c7ab9ac95c1", CoinType: tickers.Token, LastUpdate: time.Unix(444, 0),
					Price: tickers.Price{
						Value:    0.02,
						Currency: "USD",
						Provider: id,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTickers := m.normalizeTickers(tt.args.prices, coinsList, id, "USD")
			now := time.Now()
			sort.Slice(gotTickers, func(i, j int) bool {
				gotTickers[i].LastUpdate = now
				gotTickers[j].LastUpdate = now
				return gotTickers[i].Coin > gotTickers[j].Coin
			})
			sort.Slice(tt.wantTickers, func(i, j int) bool {
				tt.wantTickers[i].LastUpdate = now
				tt.wantTickers[j].LastUpdate = now
				return tt.wantTickers[i].Coin > tt.wantTickers[j].Coin
			})
			assert.Equal(t, tt.wantTickers, gotTickers)
		})
	}
}

func Test_createTicker(t *testing.T) {
	prices := make(CoinPrices, 0)
	prices = append(prices, CoinPrice{
		Id:                       "SH",
		Symbol:                   "shitcoin",
		CurrentPrice:             0.00000001,
		MarketCap:                1,
		PriceChangePercentage24h: 1,
		TotalVolume:              5000,
	})
	prices = append(prices, CoinPrice{
		Id:                       "SH",
		Symbol:                   "shitcoin",
		CurrentPrice:             0.00000001,
		MarketCap:                5000,
		PriceChangePercentage24h: 1,
		TotalVolume:              5000,
	})

	emptyTicker := tickers.Ticker{
		Price: tickers.Price{
			Value:     0.00000001,
			Change24h: 1,
		},
	}

	normalTicker := tickers.Ticker{
		Price: tickers.Price{
			Value:     0.00000001,
			Change24h: 1,
		},
	}

	wantedTickers := make(tickers.Tickers, 0)
	wantedTickers = append(wantedTickers, emptyTicker)
	wantedTickers = append(wantedTickers, normalTicker)

	for i, price := range prices {
		ticker := createTicker(price, tickers.Token, unknownCoinID, "shitcoin", "shitcoinID", "coingecko", "USD")

		assert.Equal(t, ticker.Price.Value, wantedTickers[i].Price.Value)
		assert.Equal(t, ticker.Price.Change24h, wantedTickers[i].Price.Change24h)
	}

}

var (
	mockedMarketsResponse = `[ { "id": "bitcoin", "symbol": "btc", "name": "Bitcoin", "image": "https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579", "current_price": 9696.96, "market_cap": 177446468003, "market_cap_rank": 1, "total_volume": 51778003346, "high_24h": 9661.04, "low_24h": 9099.5, "price_change_24h": 459.99, "price_change_percentage_24h": 4.97984, "market_cap_change_24h": 7519669329, "market_cap_change_percentage_24h": 4.42524, "circulating_supply": 18367225.0, "total_supply": 21000000.0, "ath": 19665.39, "ath_change_percentage": -51.20717, "ath_date": "2017-12-16T00:00:00.000Z", "atl": 67.81, "atl_change_percentage": 14050.4866, "atl_date": "2013-07-06T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:12:34.220Z" }, { "id": "ethereum", "symbol": "eth", "name": "Ethereum", "image": "https://assets.coingecko.com/coins/images/279/large/ethereum.png?1547034048", "current_price": 206.55, "market_cap": 22851909019, "market_cap_rank": 2, "total_volume": 17356592769, "high_24h": 207.93, "low_24h": 199.21, "price_change_24h": -1.32300849, "price_change_percentage_24h": -0.63646, "market_cap_change_24h": -218200222.339359, "market_cap_change_percentage_24h": -0.94581, "circulating_supply": 110830501.624, "total_supply": null, "ath": 1448.18, "ath_change_percentage": -85.76061, "ath_date": "2018-01-13T00:00:00.000Z", "atl": 0.432979, "atl_change_percentage": 47526.36618, "atl_date": "2015-10-20T00:00:00.000Z", "roi": { "times": 27.445825905402025, "currency": "btc", "percentage": 2744.5825905402025 }, "last_updated": "2020-05-07T17:12:38.629Z" }, { "id": "binancecoin", "symbol": "bnb", "name": "Binance Coin", "image": "https://assets.coingecko.com/coins/images/825/large/binance-coin-logo.png?1547034615", "current_price": 16.76, "market_cap": 2474368090, "market_cap_rank": 9, "total_volume": 386293713, "high_24h": 16.96, "low_24h": 16.34, "price_change_24h": -0.15954337, "price_change_percentage_24h": -0.9431, "market_cap_change_24h": -28933861.2562523, "market_cap_change_percentage_24h": -1.15583, "circulating_supply": 147883948.0, "total_supply": 179883948.0, "ath": 39.68, "ath_change_percentage": -57.95222, "ath_date": "2019-06-22T12:20:21.894Z", "atl": 0.0398177, "atl_change_percentage": 41801.07312, "atl_date": "2017-10-19T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:10:27.413Z" }, { "id": "tron", "symbol": "trx", "name": "TRON", "image": "https://assets.coingecko.com/coins/images/1094/large/tron-logo.png?1547035066", "current_price": 0.01594768, "market_cap": 1057389936, "market_cap_rank": 17, "total_volume": 1697721558, "high_24h": 0.01620348, "low_24h": 0.0155473, "price_change_24h": -8.271e-05, "price_change_percentage_24h": -0.51593, "market_cap_change_24h": -4768379.66360569, "market_cap_change_percentage_24h": -0.44893, "circulating_supply": 66140232427.0, "total_supply": 99281283754.0, "ath": 0.231673, "ath_change_percentage": -93.12752, "ath_date": "2018-01-05T00:00:00.000Z", "atl": 0.00180434, "atl_change_percentage": 782.40789, "atl_date": "2017-11-12T00:00:00.000Z", "roi": { "times": 7.39351342638589, "currency": "usd", "percentage": 739.3513426385889 }, "last_updated": "2020-05-07T17:10:27.759Z" }, { "id": "01coin", "symbol": "zoc", "name": "01coin", "image": "https://assets.coingecko.com/coins/images/5720/large/F1nTlw9I_400x400.jpg?1547041588", "current_price": 0.00135115, "market_cap": 14384.83, "market_cap_rank": 1587, "total_volume": 888.82, "high_24h": 0.0013974, "low_24h": 0.00120559, "price_change_24h": 3.876e-05, "price_change_percentage_24h": 2.95316, "market_cap_change_24h": 480.77, "market_cap_change_percentage_24h": 3.45776, "circulating_supply": 10646360.834599, "total_supply": 65658824.0, "ath": 0.03418169, "ath_change_percentage": -96.04715, "ath_date": "2018-10-10T17:27:38.034Z", "atl": 0.00070641, "atl_change_percentage": 91.26875, "atl_date": "2020-03-16T10:22:30.944Z", "roi": null, "last_updated": "2020-05-07T16:57:12.616Z" }, { "id": "02-token", "symbol": "o2t", "name": "O2 Token", "image": "https://assets.coingecko.com/coins/images/6925/large/44429612.jpeg?1547043298", "current_price": 0.00083971, "market_cap": 0.0, "market_cap_rank": 7111, "total_volume": 69.52, "high_24h": null, "low_24h": null, "price_change_24h": null, "price_change_percentage_24h": null, "market_cap_change_24h": null, "market_cap_change_percentage_24h": null, "circulating_supply": 0.0, "total_supply": 28520100.0, "ath": 0.00439107, "ath_change_percentage": -80.87694, "ath_date": "2018-11-20T05:12:22.611Z", "atl": 0.00057411, "atl_change_percentage": 46.26319, "atl_date": "2018-11-26T00:00:00.000Z", "roi": null, "last_updated": "2019-12-26T04:00:21.046Z" }, { "id": "xrp-bep2", "symbol": "xrp-bf2", "name": "XRP BEP2", "image": "https://assets.coingecko.com/coins/images/9686/large/12-122739_xrp-logo-png-clipart.png?1570790408", "current_price": 0.21726, "market_cap": 0.0, "market_cap_rank": 5069, "total_volume": 301.34, "high_24h": 0.219438, "low_24h": 0.212267, "price_change_24h": -0.00202035, "price_change_percentage_24h": -0.92136, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 10000000.0, "ath": 0.360995, "ath_change_percentage": -40.44337, "ath_date": "2019-10-21T13:44:21.822Z", "atl": 0.115982, "atl_change_percentage": 85.36984, "atl_date": "2020-03-13T02:02:33.103Z", "roi": null, "last_updated": "2020-05-07T17:14:13.364Z" }, { "id": "lovehearts", "symbol": "lvh", "name": "LoveHearts", "image": "https://assets.coingecko.com/coins/images/9360/large/1_d3hJ7JQeQ84goeTVWLI9Qw.png?1566528108", "current_price": 8.08e-06, "market_cap": 0.0, "market_cap_rank": 5528, "total_volume": 7.87, "high_24h": 8.73e-06, "low_24h": 7.85e-06, "price_change_24h": 1.9e-07, "price_change_percentage_24h": 2.39645, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 100000000000.0, "ath": 8.596e-05, "ath_change_percentage": -90.70838, "ath_date": "2019-08-23T03:49:38.791Z", "atl": 3.13e-06, "atl_change_percentage": 155.38143, "atl_date": "2020-02-21T21:25:35.813Z", "roi": null, "last_updated": "2020-05-07T17:10:13.895Z" } ]`
	wantedTickers         = tickers.Tickers([]tickers.Ticker{
		{
			Coin:     0,
			CoinName: "BTC",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     9696.96,
				Change24h: 4.97984,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     60,
			CoinName: "ETH",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     206.55,
				Change24h: -0.63646,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     714,
			CoinName: "BNB",
			TokenId:  "bnb",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     16.76,
				Change24h: -0.9431,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     195,
			CoinName: "TRX",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     0.01594768,
				Change24h: -0.51593,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     60,
			CoinName: "ETH",
			TokenId:  "0xb1bafca3737268a96673a250173b6ed8f1b5b65f",
			CoinType: tickers.Token,
			Price: tickers.Price{
				Value:     0.00083971,
				Change24h: 0,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     714,
			CoinName: "BNB",
			TokenId:  "xrp-bf2",
			CoinType: tickers.Token,
			Price: tickers.Price{
				Value:     0.21726,
				Change24h: -0.92136,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     195,
			CoinName: "TRX",
			TokenId:  "1000451",
			CoinType: tickers.Token,
			Price: tickers.Price{
				Value:     0.00000808,
				Change24h: 2.39645,
				Currency:  "USD",
				Provider:  "coingecko",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
	})
)
