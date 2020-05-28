package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"testing"
	"time"
)

func TestController_HandleChartsRequest(t *testing.T) {
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		LastUpdated:      time.Now(),
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         "coingecko",
		Rate:             2,
		LastUpdated:      time.Now(),
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         "fixer",
		Rate:             6,
		LastUpdated:      time.Now(),
	}

	ticker60ACMC := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coinmarketcap",
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ABNB := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "binancedex",
		Value:     100,
	}

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}

	wCharts := watchmarket.Chart{Provider: "coinmarketcap", Error: "", Prices: []watchmarket.ChartPrice{{Price: 1, Date: 1}, {Price: 3, Date: 3}}}
	cm := getChartsMock()
	cm.wantedCharts = wCharts

	c := setupController(t, db, getCacheMock(), cm)
	assert.NotNil(t, c)

	chart, err := c.HandleChartsRequest(ChartRequest{
		CoinQuery:    "60",
		Token:        "a",
		Currency:     "USD",
		TimeStartRaw: "1577871126",
		MaxItems:     "64",
	}, context.Background())
	assert.Nil(t, err)

	assert.Equal(t, wCharts, chart)
}
