package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
	"testing"
)

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), getCacheMock()))
}

func setupController(t *testing.T, d dbMock, ch cache.Charts) Controller {
	c := config.Init("../../config/test.yml")
	assert.NotNil(t, c)

	chartsPriority, err := priority.Init(c.Markets.Priority.Charts)
	assert.Nil(t, err)

	ratesPriority, err := priority.Init(c.Markets.Priority.Rates)
	assert.Nil(t, err)

	tickerPriority, err := priority.Init(c.Markets.Priority.Tickers)
	assert.Nil(t, err)

	coinInfoPriority, err := priority.Init(c.Markets.Priority.CoinInfo)
	assert.Nil(t, err)

	a := assets.NewClient(c.Markets.Assets)

	m, err := markets.Init(c, a)
	assert.Nil(t, err)

	controller := NewController(ch, d, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m)
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
	return d.WantedRates, d.WantedRatesError
}

func (d dbMock) AddRates(rates []models.Rate) error {
	return nil
}

func (d dbMock) AddTickers(tickers []models.Ticker) error {
	return nil
}
func (d dbMock) GetTickers(coin uint, tokenId string) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}
func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func getCacheMock() cache.Charts {
	i := cacheMock{}
	return i
}

type cacheMock struct {
	res string
}

func (c cacheMock) GetID() string {
	return ""
}

func (c cacheMock) GenerateKey(data string) string {
	return ""
}

func (c cacheMock) GetCharts(key string, timeStart int64) (watchmarket.Chart, error) {
	return watchmarket.Chart{}, nil
}

func (c cacheMock) SaveCharts(key string, data watchmarket.Chart, timeStart int64) error {
	return nil
}

func (c cacheMock) SaveCoinDetails(key string, data watchmarket.CoinDetails, timeStart int64) error {
	return nil
}

func (c cacheMock) GetCoinDetails(key string, timeStart int64) (watchmarket.CoinDetails, error) {
	return watchmarket.CoinDetails{}, nil
}

func (c cacheMock) GetTickers(key string) (watchmarket.Tickers, error) {
	return watchmarket.Tickers{}, nil
}

func (c cacheMock) SaveTickers(key string, tickers watchmarket.Tickers) error {
	return nil
}

func (c cacheMock) GetRates(key string) (watchmarket.Rates, error) {
	return watchmarket.Rates{}, nil
}

func (c cacheMock) SaveRates(key string, tickers watchmarket.Rates) error {
	return nil
}
