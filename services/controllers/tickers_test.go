package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
	"testing"
	"time"
)

func TestController_getRateByPriority(t *testing.T) {
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

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}

	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	r, err := c.getRateByPriority("USD")
	assert.Nil(t, err)

	assert.Equal(t, watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		Timestamp:        time.Now().Unix(),
	}, r)
}

func TestController_getTickersByPriority(t *testing.T) {
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
	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	tickers, err := c.getTickersByPriority(makeTickerQueries(
		[]Coin{{Coin: 60, TokenId: "A"}, {Coin: 714, TokenId: "A"}},
	))
	assert.Nil(t, err)
	assert.NotNil(t, tickers)
	assert.Equal(t, 2, len(tickers))

	wantedTicker1 := watchmarket.Ticker{
		Coin:     60,
		CoinName: "ETH",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  "coinmarketcap",
			Value:     100,
		},
		TokenId: "a",
	}
	wantedTicker2 := watchmarket.Ticker{
		Coin:     714,
		CoinName: "BNB",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  "coingecko",
			Value:     100,
		},
		TokenId: "a",
	}
	var counter int
	for _, t := range tickers {
		if t == wantedTicker1 || t == wantedTicker2 {
			counter++
		}
	}
	assert.Equal(t, 2, counter)
	db2 := getDbMock()
	db2.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG}
	c2 := setupController(t, db2, getCacheMock())
	tickers2, err := c2.getTickersByPriority(makeTickerQueries([]Coin{{Coin: 60, TokenId: "A"}}))
	assert.Nil(t, err)
	assert.NotNil(t, tickers2)
	assert.Equal(t, 1, len(tickers2))
	assert.Equal(t, wantedTicker1, tickers2[0])
}

func TestController_HandleTickersRequest_Negative(t *testing.T) {
	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedRatesError = errors.E("Not found")
	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	_, err := c.HandleTickersRequest(TickerRequest{})
	assert.Equal(t, err, errors.E("Not found"))
}

func TestController_normalizeTickers(t *testing.T) {
	modelRate2 := models.Rate{
		Currency:         "EUR",
		PercentChange24h: 0,
		Provider:         "fixer",
		Rate:             1.0992876616,
		LastUpdated:      time.Now(),
	}

	rate := watchmarket.Rate{
		Currency:         "EUR",
		PercentChange24h: 0,
		Provider:         "fixer",
		Rate:             1.0992876616,
		Timestamp:        12,
	}

	gotTicker1 := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BTC",
		CoinType: "coin",
		Price: watchmarket.Price{
			Change24h: -4.03168,
			Currency:  "USD",
			Provider:  "coinmarketcap",
			Value:     9360.20314131,
		},
		TokenId: "",
	}
	db := getDbMock()
	db.WantedRates = []models.Rate{modelRate2}

	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	result := c.normalizeTickers([]watchmarket.Ticker{gotTicker1}, rate)
	wanted := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BTC",
		CoinType: "coin",
		Error:    "",
		Price: watchmarket.Price{
			Change24h: -4.03168,
			Currency:  "EUR",
			Provider:  "coinmarketcap",
			Value:     8514.78959355037,
		},
		TokenId: "",
	}
	assert.Equal(t, wanted, result[0])
}

func TestController_normalizeTickers_advanced(t *testing.T) {
	modelRate := models.Rate{
		Currency:         "BNB",
		PercentChange24h: 0,
		Provider:         "coingecko",
		Rate:             16.16,
		LastUpdated:      time.Now(),
	}

	modelRate2 := models.Rate{
		Currency:         "EUR",
		PercentChange24h: 0,
		Provider:         "fixer",
		Rate:             1.0992876616,
		LastUpdated:      time.Now(),
	}

	rate := watchmarket.Rate{
		Currency:         "EUR",
		PercentChange24h: 0,
		Provider:         "fixer",
		Rate:             1.0992876616,
		Timestamp:        12,
	}

	gotTicker1 := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BNB",
		CoinType: "token",
		Price: watchmarket.Price{
			Change24h: -10.24,
			Currency:  "BNB",
			Provider:  "binancedex",
			Value:     1,
		},
		TokenId:   "raven-f66",
		Volume:    10,
		MarketCap: 10,
	}
	db := getDbMock()
	db.WantedRates = []models.Rate{modelRate, modelRate2}

	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	result := c.normalizeTickers([]watchmarket.Ticker{gotTicker1}, rate)
	wanted := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BNB",
		CoinType: "token",
		Error:    "",
		Price: watchmarket.Price{
			Change24h: -10.24,
			Currency:  "EUR",
			Provider:  "binancedex",
			Value:     14.700428799936965,
		},
		TokenId:   "raven-f66",
		Volume:    147.00428799936964,
		MarketCap: 147.00428799936964,
	}
	assert.Equal(t, wanted, result[0])
}

func TestController_createResponse(t *testing.T) {
	ticker := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BNB",
		CoinType: "token",
		Error:    "",
		Price: watchmarket.Price{
			Change24h: -10.24,
			Currency:  "EUR",
			Provider:  "binancedex",
			Value:     14.700428799936965,
		},
		TokenId:   "raven-f66",
		Volume:    147.00428799936964,
		MarketCap: 147.00428799936964,
	}

	tr := TickerRequest{
		Currency: "EUR",
		Assets:   []Coin{{Coin: 0, CoinType: "token", TokenId: "RAVEN-F66"}},
	}

	response := createResponse(tr, watchmarket.Tickers{ticker})
	wantedResponse := TickerResponse{
		Currency: "EUR",
		Tickers: watchmarket.Tickers{
			{
				Coin:     0,
				CoinName: "BNB",
				CoinType: "token",
				Error:    "",
				Price: watchmarket.Price{
					Change24h: -10.24,
					Currency:  "EUR",
					Provider:  "binancedex",
					Value:     14.700428799936965,
				},
				TokenId:   "RAVEN-F66",
				Volume:    147.00428799936964,
				MarketCap: 147.00428799936964,
			},
		},
	}
	assert.Equal(t, wantedResponse, response)
}

func Test_findBestProviderForQuery(t *testing.T) {
	tickerQueries := []Coin{{Coin: 60, TokenId: "A"}}

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

	providers := []string{"coingecko", "coinmarketcap"}
	dbTickers := []models.Ticker{ticker60ACMC, ticker60ACG}
	for i := 0; i < 10000; i++ {
		t := ticker60ACG
		t.Value = t.Value + float64(i)
		t.Coin = uint(i)
		dbTickers = append(dbTickers, t)
	}

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res)
	}

	wg.Wait()

	assert.NotNil(t, res.tickers)
}
