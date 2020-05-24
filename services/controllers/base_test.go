package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
	"testing"
)

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), getCacheMock()))
}

func setupController(t *testing.T, d dbMock, ch cache.Provider) Controller {
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

	a := assets.Init(c.Markets.Assets)

	m, err := markets.Init(c, a)
	assert.Nil(t, err)

	controller := NewController(ch, d, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m, c)
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
func (d dbMock) GetTickers(coin uint, tokenId string) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}
func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func getCacheMock() cache.Provider {
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
