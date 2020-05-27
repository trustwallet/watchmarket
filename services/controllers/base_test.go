package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"testing"
)

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), getCacheMock(), getChartsMock()))
}

func setupController(t *testing.T, d dbMock, ch cache.Provider, cm chartsMock) Controller {
	c := config.Init("../../config/test.yml")
	assert.NotNil(t, c)

	chartsPriority := []string{"coinmarketcap"}
	ratesPriority := c.Markets.Priority.Rates
	tickerPriority := c.Markets.Priority.Tickers
	coinInfoPriority := c.Markets.Priority.CoinInfo

	chartsAPIs := make(markets.ChartsAPIs, 1)
	chartsAPIs[cm.GetProvider()] = cm

	controller := NewController(ch, d, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, chartsAPIs, c)
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

func (d dbMock) GetRates(currency string, ctx context.Context) ([]models.Rate, error) {
	res := make([]models.Rate, 0)
	for _, r := range d.WantedRates {
		if r.Currency == currency {
			res = append(res, r)
		}
	}
	return res, d.WantedRatesError
}

func (d dbMock) AddRates(rates []models.Rate, ctx context.Context) error {
	return nil
}

func (d dbMock) AddTickers(tickers []models.Ticker, ctx context.Context) error {
	return nil
}

func (d dbMock) GetTickers(coin uint, tokenId string, ctx context.Context) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery, ctx context.Context) ([]models.Ticker, error) {
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
