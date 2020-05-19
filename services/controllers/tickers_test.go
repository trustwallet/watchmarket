package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"sync"
	"testing"
)

func TestController_HandleTickersRequest(t *testing.T) {
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		Timestamp:        12,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         "coingecko",
		Rate:             2,
		Timestamp:        12,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         "fixer",
		Rate:             6,
		Timestamp:        12,
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

	db := setupDb(t)

	db.WantedTickersError = nil
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}
	db.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}
	c := setupController(t, db)
	assert.NotNil(t, c)

	tickers, err := c.HandleTickersRequest(TickerRequest{
		Currency: "USD",
		Assets:   []Coin{{Coin: 60, TokenId: "A"}, {Coin: 714, TokenId: "A"}},
	})
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
	db2 := setupDb(t)
	db2.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG}
	db2.WantedRates = []models.Rate{rate, rate2, rate3}
	c2 := setupController(t, db2)
	tickers2, err := c2.HandleTickersRequest(TickerRequest{
		Currency: "USD",
		Assets:   []Coin{{Coin: 60, TokenId: "A"}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, tickers2)
	assert.Equal(t, 1, len(tickers2))
	assert.Equal(t, wantedTicker1, tickers2[0])
}

func TestController_HandleTickersRequest_Negative(t *testing.T) {
	db := setupDb(t)

	db.WantedTickersError = nil
	db.WantedRatesError = errors.E("Not found")
	c := setupController(t, db)
	assert.NotNil(t, c)

	_, err := c.HandleTickersRequest(TickerRequest{})
	assert.Equal(t, err, errors.E("Not found"))
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
		t.Coin = t.Coin + uint(i)
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
