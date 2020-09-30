package tickerscontroller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestController_createResponse(t *testing.T) {
	ticker := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BNB",
		CoinType: "Token",
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

	tr := controllers.TickerRequest{
		Currency: "EUR",
		Assets:   []controllers.Coin{{Coin: 0, CoinType: "Token", TokenId: "RAVEN-F66"}},
	}

	response := createResponse(tr, watchmarket.Tickers{ticker})
	wantedResponse := controllers.TickerResponse{
		Currency: "EUR",
		Tickers: watchmarket.Tickers{
			{
				Coin:     0,
				CoinName: "BNB",
				CoinType: "Token",
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

func Test_makeTickerQueriesV2(t *testing.T) {
	ids := []string{"c60_ta", "c714", "c714_ta"}
	wantedRes := []models.TickerQuery{
		{
			Coin:    60,
			TokenId: "a",
		},
		{
			Coin:    714,
			TokenId: "",
		},
		{
			Coin:    714,
			TokenId: "a",
		},
	}
	res := makeTickerQueriesV2(ids)

	sort.Slice(res, func(i, j int) bool {
		return res[i].Coin < res[j].Coin
	})
	sort.Slice(wantedRes, func(i, j int) bool {
		return wantedRes[i].Coin < wantedRes[j].Coin
	})

	for i, r := range res {
		assert.Equal(t, wantedRes[i], r)
	}
}

func Test_createResponseV2(t *testing.T) {
	timeUPD := time.Now()
	given1 := watchmarket.Ticker{
		Coin:     60,
		CoinName: "ETH",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  "coinmarketcap",
			Value:     100,
		},
		TokenId:    "a",
		LastUpdate: timeUPD,
	}
	given2 := watchmarket.Ticker{
		Coin:     714,
		CoinName: "BNB",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  "coingecko",
			Value:     100,
		},
		TokenId:    "a",
		LastUpdate: timeUPD,
	}
	r := createResponseV2(controllers.TickerRequestV2{Currency: "USD", Ids: []string{"c60_ta", "c714_ta"}}, []watchmarket.Ticker{given1, given2})

	wantedTicker1 := controllers.TickerPrice{
		ID:        "c60_ta",
		Change24h: 10,
		Provider:  "coinmarketcap",
		Price:     100,
	}
	wantedTicker2 := controllers.TickerPrice{
		ID:        "c714_ta",
		Change24h: 10,
		Provider:  "coingecko",
		Price:     100,
	}

	wp := []controllers.TickerPrice{wantedTicker2, wantedTicker1}
	sort.Slice(wp, func(i, j int) bool {
		return wp[i].ID < wp[j].ID
	})
	wantedResp := controllers.TickerResponseV2{
		Currency: "USD",
		Tickers:  wp,
	}

	sort.Slice(r.Tickers, func(i, j int) bool {
		return r.Tickers[i].ID < r.Tickers[j].ID
	})

	assert.Equal(t, wantedResp, r)
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

	now := time.Now()

	gotTicker1 := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BTC",
		CoinType: "Coin",
		Price: watchmarket.Price{
			Change24h: -4.03168,
			Currency:  "USD",
			Provider:  "coinmarketcap",
			Value:     9360.20314131,
		},
		TokenId:    "",
		LastUpdate: now,
	}
	db := getDbMock()
	db.WantedRates = []models.Rate{modelRate2}

	c := setupController(t, db, false)
	assert.NotNil(t, c)

	result := c.normalizeTickers([]watchmarket.Ticker{gotTicker1}, rate, context.Background())
	wanted := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BTC",
		CoinType: "Coin",
		Error:    "",
		Price: watchmarket.Price{
			Change24h: -4.03168,
			Currency:  "EUR",
			Provider:  "coinmarketcap",
			Value:     8514.78959355037,
		},
		TokenId:    "",
		LastUpdate: now,
	}
	assert.Equal(t, wanted, result[0])
}

func TestController_normalizeTickers_advanced(t *testing.T) {
	timeUPD := time.Now()
	modelRate := models.Rate{
		Currency:         "BNB",
		PercentChange24h: 0,
		Provider:         "coingecko",
		Rate:             16.16,
		LastUpdated:      timeUPD,
	}

	modelRate2 := models.Rate{
		Currency:         "EUR",
		PercentChange24h: 0,
		Provider:         "fixer",
		Rate:             1.0992876616,
		LastUpdated:      timeUPD,
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
		CoinType: "Token",
		Price: watchmarket.Price{
			Change24h: -10.24,
			Currency:  "BNB",
			Provider:  "binancedex",
			Value:     1,
		},
		TokenId:    "raven-f66",
		Volume:     10,
		MarketCap:  10,
		LastUpdate: timeUPD,
	}

	db := getDbMock()
	db.WantedRates = []models.Rate{modelRate, modelRate2}

	c := setupController(t, db, false)
	assert.NotNil(t, c)

	result := c.normalizeTickers([]watchmarket.Ticker{gotTicker1}, rate, context.Background())
	wanted := watchmarket.Ticker{
		Coin:     0,
		CoinName: "BNB",
		CoinType: "Token",
		Error:    "",
		Price: watchmarket.Price{
			Change24h: -10.24,
			Currency:  "EUR",
			Provider:  "binancedex",
			Value:     14.700428799936965,
		},
		TokenId:    "raven-f66",
		Volume:     147.00428799936964,
		MarketCap:  147.00428799936964,
		LastUpdate: timeUPD,
	}
	assert.Equal(t, wanted, result[0])
}

func Test_findBestProviderForQuery(t *testing.T) {
	tickerQueries := []controllers.Coin{{Coin: 60, TokenId: "A"}}

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

	providers := []string{"coinmarketcap", "coingecko"}
	dbTickers := []models.Ticker{ticker60ACMC, ticker60ACG}
	for i := 0; i < 10000; i++ {
		t := ticker60ACG
		t.Value = t.Value + float64(i)
		t.Coin = uint(i)
		dbTickers = append(dbTickers, t)
	}

	c := config.Init("../../../config/test.yml")
	assert.NotNil(t, c)

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res, c)
	}

	wg.Wait()

	assert.Equal(t, ticker60ACMC, res.tickers[0])
}

func Test_findBestProviderForQuery_advanced(t *testing.T) {
	tickerQueries := []controllers.Coin{{Coin: 60, TokenId: "A"}}

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
		Coin:       60,
		CoinName:   "ETH",
		TokenId:    "a",
		Change24h:  10,
		Currency:   "USD",
		Provider:   "coingecko",
		Value:      100,
		ShowOption: models.NeverShow,
	}

	providers := []string{"coingecko", "coinmarketcap"}
	dbTickers := []models.Ticker{ticker60ACMC, ticker60ACG}
	for i := 0; i < 10000; i++ {
		t := ticker60ACG
		t.Value = t.Value + float64(i)
		t.Coin = uint(i)
		dbTickers = append(dbTickers, t)
	}

	c := config.Init("../../../config/test.yml")
	assert.NotNil(t, c)

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res, c)
	}

	wg.Wait()

	assert.Equal(t, ticker60ACMC, res.tickers[0])
}

func Test_findBestProviderForQuery_showOption(t *testing.T) {
	tickerQueries := []controllers.Coin{{Coin: 60, TokenId: "A"}}

	ticker60ACMC := models.Ticker{
		Coin:       60,
		CoinName:   "ETH",
		TokenId:    "a",
		Change24h:  10,
		Currency:   "USD",
		Provider:   "coinmarketcap",
		Value:      100,
		ShowOption: models.AlwaysShow,
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

	c := config.Init("../../../config/test.yml")
	assert.NotNil(t, c)

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res, c)
	}

	wg.Wait()

	assert.Equal(t, ticker60ACMC, res.tickers[0])
}
