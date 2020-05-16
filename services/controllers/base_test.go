package controllers

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {

	assert.NotNil(t, setupController(t))
	//data, err := controller.HandleChartsRequest(ChartRequest{
	//	coinQuery:    "60",
	//	token:        "",
	//	currency:     "USD",
	//	timeStartRaw: "1577871126",
	//	maxItems:     "64",
	//})
	//
	//assert.Nil(t, err)
	//assert.NotNil(t, data)
	//
	//controller.HandleChartsRequest(ChartRequest{
	//	coinQuery:    "60",
	//	token:        "",
	//	currency:     "USD",
	//	timeStartRaw: "1577871126",
	//	maxItems:     "64",
	//})
}

func setupController(t *testing.T) Controller {
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

	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	cacheInstance := cache.Init(r, time.Minute, time.Minute, time.Minute, time.Minute)

	db := setupDb(t)

	controller := NewController(cacheInstance, db, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m)
	assert.NotNil(t, controller)
	return controller
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func setupDb(t *testing.T) dbMock {
	return dbMock("f")
}

type dbMock string

func (d dbMock) GetRates(currency string) ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) AddRates(rates []models.Rate) error {
	return nil
}

func (d dbMock) AddTickers(tickers []models.Ticker) error {
	return nil
}
func (d dbMock) GetTickers(coin uint, tokenId string) ([]models.Ticker, error) {
	ticker60ACMC := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coinmarketcap",
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ABNB := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "binancedex",
		Value:     100,
	}

	return []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}, nil
}
func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	ticker60ACMC := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coinmarketcap",
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ABNB := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "A",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "binancedex",
		Value:     100,
	}

	return []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}, nil
}
