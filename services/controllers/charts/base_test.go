package chartscontroller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

func TestController_HandleChartsRequest(t *testing.T) {
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         watchmarket.CoinMarketCap,
		Rate:             1,
		LastUpdated:      time.Now(),
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         watchmarket.CoinGecko,
		Rate:             2,
		LastUpdated:      time.Now(),
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         watchmarket.Fixer,
		Rate:             6,
		LastUpdated:      time.Now(),
	}

	ticker60ACMC := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinMarketCap,
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinGecko,
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinGecko,
		Value:     100,
	}

	ticker714ABNB := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinMarketCap,
		Value:     100,
	}

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}

	wCharts := watchmarket.Chart{Provider: watchmarket.CoinMarketCap, Error: "", Prices: []watchmarket.ChartPrice{{Price: 1, Date: 1}, {Price: 3, Date: 3}}}
	cm := getChartsMock()
	cm.wantedCharts = wCharts

	c := setupController(t, db, getCacheMock(), cm)
	assert.NotNil(t, c)

	chart, err := c.HandleChartsRequest(controllers.ChartRequest{
		Asset: controllers.Asset{
			CoinId:  60,
			TokenId: "a",
		},
		Currency:  "USD",
		TimeStart: 1577871126,
		MaxItems:  64,
	})
	assert.Nil(t, err)

	assert.Equal(t, wCharts, chart)
}

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), getCacheMock(), getChartsMock()))
}

func setupController(t *testing.T, d dbMock, ch cache.Provider, cm chartsMock) Controller {
	c, _ := config.Init("../../../config.yml")
	assert.NotNil(t, c)
	c.RestAPI.UseMemoryCache = false

	chartsPriority := []string{watchmarket.CoinMarketCap}

	chartsAPIs := make(markets.ChartsAPIs, 1)
	chartsAPIs[cm.GetProvider()] = cm

	controller := NewController(ch, ch, d, chartsPriority, chartsAPIs, c)
	assert.NotNil(t, controller)
	return controller

}
func getDbMock() dbMock {
	return dbMock{}
}

type dbMock struct {
	WantedRates        []models.Rate
	WantedTickers      []models.Ticker
	WantedTickersError error
	WantedRatesError   error
}

func (d dbMock) GetRates(currency string) ([]models.Rate, error) {
	res := make([]models.Rate, 0)
	for _, r := range d.WantedRates {
		if r.Currency == currency {
			res = append(res, r)
		}
	}
	return res, d.WantedRatesError
}

func (d dbMock) AddRates(rates []models.Rate) error {
	return nil
}

func (d dbMock) AddTickers(tickers []models.Ticker) error {
	return nil
}

func (d dbMock) GetAllTickers() ([]models.Ticker, error) {
	return nil, nil
}

func (d dbMock) GetRatesByProvider(provider string) ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) GetAllRates() ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) GetTickers(asset []controllers.Asset) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func getCacheMock() cache.Provider {
	i := cacheMock{}
	return i
}

type cacheMock struct {
}

func (c cacheMock) GetID() string {
	return ""
}

func (c cacheMock) GenerateKey(data string) string {
	return ""
}

func (c cacheMock) Get(key string) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) Set(key string, data []byte) error {
	return nil
}

func (c cacheMock) GetWithTime(key string, time int64) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) SetWithTime(key string, data []byte, time int64) error {
	return nil
}

func (c cacheMock) GetLenOfSavedItems() int {
	return 0
}

func getChartsMock() chartsMock {
	cm := chartsMock{}
	return cm
}

type chartsMock struct {
	wantedCharts  watchmarket.Chart
	wantedDetails watchmarket.CoinDetails
}

func (cm chartsMock) GetChartData(asset controllers.Asset, currency string, timeStart int64) (watchmarket.Chart, error) {
	return cm.wantedCharts, nil
}

func (cm chartsMock) GetCoinData(asset controllers.Asset, currency string) (watchmarket.CoinDetails, error) {
	return cm.wantedDetails, nil
}

func (cm chartsMock) GetProvider() string {
	return watchmarket.CoinMarketCap
}
