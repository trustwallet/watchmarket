// +build integration

package db_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/services/controllers"
	"github.com/trustwallet/watchmarket/tests/integration/setup"

	"testing"
)

func TestAddTickers(t *testing.T) {
	setup.CleanupPgContainer(databaseInstance.Gorm)
	tickers := []models.Ticker{{
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "60",
		Value:     60,
	}, {
		Coin:      70,
		CoinName:  "70",
		CoinType:  "70",
		TokenId:   "70",
		Change24h: 70,
		Currency:  "70",
		Provider:  "70",
		Value:     70,
	}}

	d := db.Instance(databaseInstance)
	err := d.AddTickers(tickers)
	assert.Nil(t, err)

	result1, err := d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 60, TokenId: "60"}})
	assert.Nil(t, err)
	assert.Len(t, result1, 1)
	assert.Equal(t, uint(60), result1[0].Coin)

	result2, err := d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 70, TokenId: "70"}})
	assert.Nil(t, err)
	assert.Len(t, result2, 1)
	assert.Equal(t, uint(70), result2[0].Coin)

	tickers = append(tickers, models.Ticker{
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "61",
		Value:     60,
	})

	err = d.AddTickers(tickers)
	assert.Nil(t, err)

	result1, err = d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 60, TokenId: "60"}})
	assert.Nil(t, err)
	assert.Len(t, result1, 2)

	tickers[1].Value = 100500
	tickers[1].Change24h = 666

	err = d.AddTickers(tickers)
	assert.Nil(t, err)
	result2, err = d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 70, TokenId: "70"}})
	assert.Nil(t, err)
	assert.Len(t, result2, 1)
	assert.Equal(t, float64(100500), result2[0].Value)
	assert.Equal(t, float64(666), result2[0].Change24h)
}

func TestAddTickersMult(t *testing.T) {
	setup.CleanupPgContainer(databaseInstance.Gorm)
	tickers := []models.Ticker{{
		Coin:      70,
		CoinName:  "70",
		CoinType:  "70",
		TokenId:   "70",
		Change24h: 70,
		Currency:  "70",
		Provider:  "70",
		Value:     70,
		MarketCap: 1,
	}, {
		Coin:      70,
		CoinName:  "70",
		CoinType:  "70",
		TokenId:   "70",
		Change24h: 70,
		Currency:  "70",
		Provider:  "70",
		Value:     70,
	}, {
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "61",
		Value:     60,
	}}

	d := db.Instance(databaseInstance)
	err := d.AddTickers(tickers)
	assert.Nil(t, err)

	err = d.AddTickers(tickers)
	assert.Nil(t, err)

	result1, err := d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 60, TokenId: "60"}})
	assert.Nil(t, err)
	assert.Len(t, result1, 1)
	assert.Equal(t, uint(60), result1[0].Coin)

	tickers = append(tickers, models.Ticker{
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "61",
		Value:     60,
	})

	err = d.AddTickers(tickers)
	assert.Nil(t, err)

	result1, err = d.GetTickers([]controllers.Asset{controllers.Asset{CoinId: 60, TokenId: "60"}})
	assert.Nil(t, err)
	assert.Len(t, result1, 1)
}
