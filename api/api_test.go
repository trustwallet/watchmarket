package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/mocks/storage"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/storage"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const(
	USDRate = 1
	ETHToUSDRate = 10
	ETHPrice = 10
)

func TestTickers(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	engine := setupEngine()
	db := internal.InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)

	Bootstrap(engine, db)

	server := httptest.NewServer(engine)
	defer server.Close()

	tests := []struct {
		name           string
		requestMethod  string
		requestUrl     string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "test bad payload",
			requestUrl: fmt.Sprintf("%s/v1/market/ticker", server.URL),
			requestMethod: "POST",
			requestBody: "bad payload",
			expectedStatus: 400,
			expectedBody: "{\"code\":400,\"error\":\"Invalid request payload\"}",
		},
		{
			name: "test unknown currency",
			requestUrl: fmt.Sprintf("%s/v1/market/ticker", server.URL),
			requestMethod: "POST",
			requestBody: "{\"currency\":\"i-do-not-exist\",\"assets\":[{\"type\":\"coin\",\"coin\":60}]}",
			expectedStatus: 404,
			expectedBody: "{\"code\":404,\"error\":\"Currency i-do-not-exist not found\"}",
		},
		{
			name: "without conversion",
			requestUrl: fmt.Sprintf("%s/v1/market/ticker", server.URL),
			requestMethod: "POST",
			requestBody: "{\"currency\":\"ETH\",\"assets\":[{\"type\":\"coin\",\"coin\":60}]}",
			expectedStatus: 200,
			expectedBody: fmt.Sprintf("{\"currency\":\"ETH\",\"docs\":[{\"coin\":60,\"type\":\"tbd\",\"price\":{\"value\":%d,\"change_24h\":0},\"last_update\":\"0001-01-01T00:00:00Z\"}]}",
				ETHPrice),
		},
		{
			name: "with conversion",
			requestUrl: fmt.Sprintf("%s/v1/market/ticker", server.URL),
			requestMethod: "POST",
			requestBody: "{\"currency\":\"USD\",\"assets\":[{\"type\":\"coin\",\"coin\":60}]}",
			expectedStatus: 200,
			expectedBody: fmt.Sprintf("{\"currency\":\"USD\",\"docs\":[{\"coin\":60,\"type\":\"tbd\",\"price\":{\"value\":%d,\"change_24h\":0},\"last_update\":\"0001-01-01T00:00:00Z\"}]}",
				ETHToUSDRate * ETHPrice),
		},
		{
			name: "with conversion when there is no rate for the currency of the coin price but a ticker for the currency of the coin price",
			requestUrl: fmt.Sprintf("%s/v1/market/ticker", server.URL),
			requestMethod: "POST",
			requestBody: "{\"currency\":\"USD\",\"assets\":[{\"type\":\"coin\",\"coin\":714}]}",
			expectedStatus: 200,
			expectedBody: fmt.Sprintf("{\"currency\":\"USD\",\"docs\":[{\"coin\":714,\"type\":\"tbd\",\"price\":{\"value\":%d,\"change_24h\":0},\"last_update\":\"0001-01-01T00:00:00Z\"}]}",
				ETHToUSDRate * ETHPrice),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(makeRequest(t, tt.requestMethod, tt.requestUrl, strings.NewReader(tt.requestBody)))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			assert.Equal(t, resp.StatusCode, tt.expectedStatus)
			responseBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Responded with %d: \n%s\n", resp.StatusCode, string(responseBytes))

			assert.Equal(t, parseJson(t, responseBytes), parseJson(t, []byte(tt.expectedBody)))
		})
	}
}

func parseJson(t *testing.T, data []byte) interface{} {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func setupEngine() *gin.Engine {
	internal.InitConfig("../../config.yml")
	tmp := sentrygin.New(sentrygin.Options{}); sg := &tmp
	return internal.InitEngine(sg, viper.GetString("gin.mode"))
}

func seedDb(t *testing.T, db *storage.Storage) {
	mockProviderList := &mocks.ProviderList{}
	mockProviderList.On("GetPriority", "myMockProvider").Return(0)
	rates := watchmarket.Rates{
		watchmarket.Rate{Currency: "USD", Rate: USDRate, Timestamp: time.Now().Unix(), Provider: "myMockProvider", PercentChange24h: big.NewFloat(0)},
		watchmarket.Rate{Currency: "ETH", Rate: ETHToUSDRate, Timestamp: time.Now().Unix(), Provider: "myMockProvider", PercentChange24h: big.NewFloat(0)},
	}

	db.SaveRates(rates, mockProviderList)
	saveTicker(t, db, mockProviderList, 60, "USD", ETHPrice)
	saveTicker(t, db, mockProviderList, 60, "ETH", ETHPrice)
	saveTicker(t, db, mockProviderList, 714, "BNB", ETHPrice)
}

func saveTicker(t *testing.T, db *storage.Storage, pl storage.ProviderList, coinId uint, coinCurrency string, coinPrice float64) {
	coinObj, ok := coin.Coins[coinId]
	if !ok {
		t.Fatal(errors.New("coin does not exist"))
	}
	_, err := db.SaveTicker(&watchmarket.Ticker{
		Coin:       coinObj.ID,
		CoinName:   coinObj.Symbol,
		TokenId:    "",
		CoinType:   "tbd",
		Price:      watchmarket.TickerPrice{
			Value:     coinPrice,
			Change24h: 0,
			Currency:  coinCurrency,
			Provider:  "myMockProvider",
		},
		LastUpdate: time.Time{},
	}, pl)
	if err != nil {
		t.Fatal(err)
	}
}

func makeRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	return r
}