package infocontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
)

func TestController_HandleDetailsRequest(t *testing.T) {
	cm := getChartsMock()
	wantedD := watchmarket.CoinDetails{
		Provider: watchmarket.CoinMarketCap,
		Info: &watchmarket.Info{
			Name:             "2",
			Website:          "2",
			SourceCode:       "2",
			WhitePaper:       "2",
			Description:      "2",
			ShortDescription: "2",
			Research:         "2",
			Explorer:         "2",
			Socials:          nil,
		},
	}
	cm.wantedDetails = wantedD

	db := getDbMock()
	db.WantedTickers = []models.Ticker{{CirculatingSupply: 1, TotalSupply: 1, MarketCap: 1, Volume: 1, Provider: "coinmarketcap"}}
	db.WantedRates = []models.Rate{{Currency: "RUB", Rate: 10, Provider: "coinmarketcap"}}
	c := setupController(t, db, getCacheMock(), cm)
	assert.NotNil(t, c)
	details, err := c.HandleInfoRequest(controllers.DetailsRequest{
		Asset: controllers.Asset{
			CoinId:  0,
			TokenId: "2",
		},
		Currency: "RUB",
	})
	assert.Nil(t, err)
	assert.Equal(t, controllers.InfoResponse{
		Provider:          wantedD.Provider,
		ProviderURL:       wantedD.ProviderURL,
		Vol24:             0.1,
		MarketCap:         0.1,
		CirculatingSupply: 1,
		TotalSupply:       1,
		Info:              wantedD.Info,
	}, details)
}

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), getCacheMock(), getChartsMock()))
}

func setupController(t *testing.T, db dbMock, ch cache.Provider, cm chartsMock) Controller {
	c, _ := config.Init("../../../config.yml")
	assert.NotNil(t, c)

	ratesPriority := []string{"coinmarketcap"}
	coinInfoPriority := []string{"coinmarketcap"}

	chartsAPIs := make(markets.ChartsAPIs, 1)
	chartsAPIs[cm.GetProvider()] = cm

	controller := NewController(db, ch, coinInfoPriority, ratesPriority, chartsAPIs)
	assert.NotNil(t, controller)
	return controller

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
	return "coinmarketcap"
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

func (d dbMock) GetRatesByProvider(provider string) ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) GetAllTickers() ([]models.Ticker, error) {
	return nil, nil
}

func (d dbMock) GetAllRates() ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) GetTickers(asset []controllers.Asset) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}
