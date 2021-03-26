package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/golibs/mock"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

var (
	wantedRates, _           = mock.JsonStringFromFilePath("mocks/rates.json")
	mockedMarketsResponse, _ = mock.JsonStringFromFilePath("mocks/markets_response.json")
	wantedInfo, _            = mock.JsonStringFromFilePath("mocks/info.json")
	mockedInfoResponse, _    = mock.JsonStringFromFilePath("mocks/info_response.json")
	wantedCharts, _          = mock.JsonStringFromFilePath("mocks/charts.json")
	mockedChartResponse, _   = mock.JsonStringFromFilePath("mocks/chart_response.json")
	wantedTickers, _         = mock.JsonStringFromFilePath("mocks/tickers.json")
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("web.api", "", "USD", assets.Init("assets.api"))
	assert.NotNil(t, provider)
	assert.Equal(t, "web.api", provider.client.client.BaseUrl)
	assert.Equal(t, "USD", provider.currency)
	assert.Equal(t, watchmarket.CoinGecko, provider.id)
}

func TestProvider_GetProvider(t *testing.T) {
	provider := InitProvider("web.api", "", "USD", assets.Init("assets.api"))
	assert.Equal(t, watchmarket.CoinGecko, provider.GetProvider())
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
		if _, err := fmt.Fprint(w, mockedMarketsResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v3/coins/ethereum/market_chart/range", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedChartResponse); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/ethereum/info/info.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedInfoResponse); err != nil {
			panic(err)
		}
	})

	return r
}
