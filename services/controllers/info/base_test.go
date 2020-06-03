package infocontroller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/services/markets"
	"testing"
)

func TestController_HandleDetailsRequest(t *testing.T) {
	cm := getChartsMock()
	wantedD := watchmarket.CoinDetails{
		Provider:          "coinmarketcap",
		Vol24:             1,
		MarketCap:         2,
		CirculatingSupply: 3,
		TotalSupply:       4,
		Info: &watchmarket.Info{
			Name:             "2",
			Website:          "2",
			SourceCode:       "2",
			WhitePaper:       "2",
			Description:      "2",
			ShortDescription: "2",
			Explorer:         "2",
			Socials:          nil,
		},
	}
	cm.wantedDetails = wantedD
	c := setupController(t, getCacheMock(), cm)
	assert.NotNil(t, c)
	details, err := c.HandleInfoRequest(controllers.DetailsRequest{
		CoinQuery: "0",
		Token:     "2",
		Currency:  "3",
	}, context.Background())
	assert.Nil(t, err)
	assert.Equal(t, wantedD, details)
}

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getCacheMock(), getChartsMock()))
}

func setupController(t *testing.T, ch cache.Provider, cm chartsMock) Controller {
	c := config.Init("../../../config/test.yml")
	assert.NotNil(t, c)

	chartsPriority := []string{"coinmarketcap"}
	ratesPriority := c.Markets.Priority.Rates
	tickerPriority := c.Markets.Priority.Tickers
	coinInfoPriority := c.Markets.Priority.CoinInfo

	chartsAPIs := make(markets.ChartsAPIs, 1)
	chartsAPIs[cm.GetProvider()] = cm

	controller := NewController(ch, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, chartsAPIs, c)
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

func (c cacheMock) Get(key string, ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) Set(key string, data []byte, ctx context.Context) error {
	return nil
}

func (c cacheMock) GetWithTime(key string, time int64, ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) SetWithTime(key string, data []byte, time int64, ctx context.Context) error {
	return nil
}

func getChartsMock() chartsMock {
	cm := chartsMock{}
	return cm
}

type chartsMock struct {
	wantedCharts  watchmarket.Chart
	wantedDetails watchmarket.CoinDetails
}

func (cm chartsMock) GetChartData(coinID uint, token, currency string, timeStart int64, ctx context.Context) (watchmarket.Chart, error) {
	return cm.wantedCharts, nil
}

func (cm chartsMock) GetCoinData(coinID uint, token, currency string, ctx context.Context) (watchmarket.CoinDetails, error) {
	return cm.wantedDetails, nil
}

func (cm chartsMock) GetProvider() string {
	return "coinmarketcap"
}
