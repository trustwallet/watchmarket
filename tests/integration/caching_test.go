// +build integration

package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/watchmarket/api"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/caching"
	"github.com/trustwallet/watchmarket/storage"
	"github.com/trustwallet/watchmarket/tests/integration/setup"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	cache *caching.Provider
)

func TestChartsCachingInit(t *testing.T) {
	var mockedCharts market.Charts
	mockedCharts.ChartProviders = *setup.InitChartProviders()
	assert.NotNil(t, mockedCharts.ChartProviders)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	assert.NotNil(t, engine)

	caching.SetChartsCachingDuration(300)
	cache = caching.InitCaching(setup.Cache)
	assert.NotNil(t, cache)

	saveTicker(t, setup.Cache, nil, 60, "ETHToken", "ETH", 10)

	api.Bootstrap(api.BootstrapProviders{
		Engine: engine,
		Market: setup.Cache,
		Charts: &mockedCharts,
		Ac:     &assets.HttpAssetClient{HttpClient: resty.New()},
		Cache:  cache,
	})

	go internal.SetupGracefulShutdown("8080", engine)
}

func TestWithAlreadySetupedCache(t *testing.T) {
	cleanupCache(*setup.Cache)
	rawData, err := makeRawDataMock()
	if err != nil {
		t.Fatal(err)
	}

	timeFirst := 1574483028
	url, key := buildUrlAndKey(timeFirst)
	SetCachedData(*cache, key, rawData, "oZGj-pQpMkoKoxDLy07SEb1XwH4=", int64(timeFirst))
	makeRequestAndTestIt(t, url, `{"prices":[{"price":100000,"date":0},{"price":100000,"date":0}],"provider":""}`)
}

func TestWithThatCacheResetsWithTimeBefore(t *testing.T) {
	cleanupCache(*setup.Cache)
	rawData, err := makeRawDataMock()
	if err != nil {
		t.Fatal(err)
	}

	timeBeforeFirst := 1574483026
	url, key := buildUrlAndKey(timeBeforeFirst)
	SetCachedData(*cache, key, rawData, "oZGj-pQpMkoKoxDLy07SEb1XwH4=", int64(timeBeforeFirst+1))

	makeRequestAndTestIt(t, url, `{"prices":[{"price":10,"date":0},{"price":10,"date":0}],"provider":""}`)
}

func TestWithThatCacheIsNotDisplayedIfOutdated(t *testing.T) {
	cleanupCache(*setup.Cache)
	rawData, err := makeRawDataMock()
	if err != nil {
		t.Fatal(err)
	}

	timeWithInvalidPeriod := 1574483128
	url, key := buildUrlAndKey(timeWithInvalidPeriod)
	SetCachedData(*cache, key, rawData, "oZGj-pQpMkoKoxDLy07SEb1XwH4=", int64(timeWithInvalidPeriod+100000))

	makeRequestAndTestIt(t, url, `{"prices":[{"price":10,"date":0},{"price":10,"date":0}],"provider":""}`)

}

func makeRawDataMock() ([]byte, error) {
	price := watchmarket.ChartPrice{
		Price: 100000,
		Date:  0,
	}

	prices := make([]watchmarket.ChartPrice, 0)
	prices = append(prices, price)
	prices = append(prices, price)

	data := watchmarket.ChartData{
		Prices: prices,
		Error:  "",
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func cleanupCache(db storage.Storage) {
	db.DeleteHM(storage.EntityInterval, "Zfb7OX4NBbtWQn_Wtz8e3YfmWiM=")
	db.DeleteHM(storage.EntityCache, "oZGj-pQpMkoKoxDLy07SEb1XwH4=")
}

func buildUrlAndKey(timeStart int) (string, string) {
	url := fmt.Sprintf("http://localhost:8080/v1/market/charts?coin=60&time_start=%s&token=ETHToken", strconv.Itoa(timeStart))
	return url, "Zfb7OX4NBbtWQn_Wtz8e3YfmWiM="
}

func makeRequestAndTestIt(t *testing.T, url, wantRes string) {
	resp, err := http.DefaultClient.Do(makeRequest(t, "GET", url, strings.NewReader("")))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, parseJson(t, []byte(wantRes)), parseJson(t, responseBytes))
}

func SetCachedData(cache caching.Provider, key string, rawData []byte, keyTwo string, timestamp int64) {
	cache.DB.Set(keyTwo, rawData)
	cache.DB.UpdateInterval(key, storage.CachedInterval{
		Timestamp: timestamp,
		Duration:  300,
		Key:       keyTwo,
	})
}

func parseJson(t *testing.T, data []byte) interface{} {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		t.Fatal(err)
	}
	return value
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

func saveTicker(t *testing.T, db *storage.Storage, pl storage.ProviderList, coinId uint, tokenId, coinCurrency string, coinPrice float64) {
	coinObj, ok := coin.Coins[coinId]
	if !ok {
		t.Fatal(errors.New("coin does not exist"))
	}
	_, err := db.SaveTicker(&watchmarket.Ticker{
		Coin:     coinObj.ID,
		CoinName: coinObj.Symbol,
		TokenId:  tokenId,
		CoinType: "tbd",
		Price: watchmarket.TickerPrice{
			Value:     coinPrice,
			Change24h: 0,
			Currency:  coinCurrency,
			Provider:  "myMockProvider",
		},
		LastUpdate: time.Time{},
	}, A)
	if err != nil {
		t.Fatal(err)
	}
}

type mockProviderList string

var A mockProviderList

func (a mockProviderList) GetPriority(providerId string) int {
	return 0
}
