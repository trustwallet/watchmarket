package postgres

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"testing"
)

func Test_normalizeTickers(t *testing.T) {
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
	}, {
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "60",
		Value:     60,
	},
		{
			Coin:      60,
			CoinName:  "60",
			CoinType:  "60",
			TokenId:   "60",
			Change24h: 60,
			Currency:  "60",
			Provider:  "60",
			Value:     60,
		},
	}
	result := normalizeTickers(tickers)
	assert.Len(t, result, 2)
	assert.NotEqual(t, result[0], result[1])

	tickers = []models.Ticker{{
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 60,
		Currency:  "60",
		Provider:  "60",
		Value:     60,
	}, {
		Coin:      60,
		CoinName:  "60",
		CoinType:  "60",
		TokenId:   "60",
		Change24h: 100,
		Currency:  "60",
		Provider:  "60",
		Value:     100,
	}}
	result = normalizeTickers(tickers)
	assert.Len(t, result, 1)
	assert.Equal(t, uint(60), result[0].Coin)
	assert.Equal(t, float64(60), result[0].Change24h)
	assert.Equal(t, float64(60), result[0].Value)
}
