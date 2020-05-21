package controllers

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strconv"
	"sync"
	"testing"
)

func TestController_getRateByPriority(t *testing.T) {
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		LastUpdated:      12,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         "coingecko",
		Rate:             2,
		LastUpdated:      12,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         "fixer",
		Rate:             6,
		LastUpdated:      12,
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
		Timestamp:        12,
	}, r)
}

func TestController_getTickersByPriority(t *testing.T) {

	ticker60ACMC := models.Ticker{
		Coin:      "60",
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coinmarketcap",
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      "60",
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      "714",
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coingecko",
		Value:     100,
	}

	ticker714ABNB := models.Ticker{
		Coin:      "714",
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
	modelRate := models.Rate{
		Currency:         "BTC",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             21,
		LastUpdated:      12,
	}
	modelRate2 := models.Rate{
		Currency:         "EUR",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             12,
		LastUpdated:      12,
	}

	rate := watchmarket.Rate{
		Currency:         "EUR",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             12,
		Timestamp:        12,
	}

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

	wantedTicker3 := watchmarket.Ticker{
		Coin:     118,
		CoinName: "ATOM",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "BTC",
			Provider:  "coingecko",
			Value:     0.001,
		},
		TokenId: "",
	}

	db := getDbMock()
	db.WantedRates = []models.Rate{modelRate, modelRate2}

	c := setupController(t, db, getCacheMock())
	assert.NotNil(t, c)

	result := c.normalizeTickers([]watchmarket.Ticker{wantedTicker1, wantedTicker2, wantedTicker3}, rate)
	assert.NotNil(t, result)
}

func Test_findBestProviderForQuery(t *testing.T) {
	tickerQueries := []Coin{{Coin: 60, TokenId: "A"}}

	ticker60ACMC := models.Ticker{
		Coin:      "60",
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  "coinmarketcap",
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      "60",
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
		t.Coin = strconv.Itoa(i)
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
