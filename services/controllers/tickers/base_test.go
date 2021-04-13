package tickerscontroller

import (
	"encoding/json"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/cache/memory"
	"github.com/trustwallet/watchmarket/services/controllers"
)

func TestController_HandleTickersRequest(t *testing.T) {
	timeUPD := time.Now()
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         watchmarket.CoinMarketCap,
		Rate:             1,
		LastUpdated:      timeUPD,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         watchmarket.CoinGecko,
		Rate:             2,
		LastUpdated:      timeUPD,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         watchmarket.Fixer,
		Rate:             6,
		LastUpdated:      timeUPD,
	}

	ticker60ACMC := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    watchmarket.CoinMarketCap,
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker60ACG := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    watchmarket.CoinGecko,
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker714ACG := models.Ticker{
		Coin:        714,
		CoinName:    "BNB",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    watchmarket.CoinGecko,
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker714ABNB := models.Ticker{
		Coin:        714,
		CoinName:    "BNB",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "binancedex",
		Value:       100,
		LastUpdated: timeUPD,
	}

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG, ticker714ABNB}
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}
	c := setupController(t, db, false)
	assert.NotNil(t, c)

	response, err := c.HandleTickersRequest(controllers.TickerRequest{Currency: "USD", Assets: []controllers.Asset{{CoinId: 60, TokenId: "a"}, {CoinId: 714, TokenId: "a"}}})
	assert.Nil(t, err)

	wantedTicker1 := watchmarket.Ticker{
		Coin:     60,
		CoinName: "ETH",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     100,
		},
		TokenId:    "a",
		LastUpdate: timeUPD,
	}
	wantedTicker2 := watchmarket.Ticker{
		Coin:     714,
		CoinName: "BNB",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  watchmarket.CoinGecko,
			Value:     100,
		},
		TokenId:    "a",
		LastUpdate: timeUPD,
	}

	wantedResp := controllers.TickerResponse{
		Currency: "USD",
		Tickers:  []watchmarket.Ticker{wantedTicker2, wantedTicker1},
	}

	sort.Slice(wantedResp.Tickers, func(i, j int) bool {
		return wantedResp.Tickers[i].Coin < wantedResp.Tickers[j].Coin
	})
	sort.Slice(response, func(i, j int) bool {
		return response[i].Coin < response[j].Coin
	})

	assert.Equal(t, wantedResp.Tickers, response)

	controllerWithCache := setupController(t, db, true)
	assert.NotNil(t, controllerWithCache)
	wantedTicker1Raw, err := json.Marshal(&wantedTicker1)
	assert.Nil(t, err)
	wantedTicker2Raw, err := json.Marshal(&wantedTicker2)
	assert.Nil(t, err)
	rateRaw, err := json.Marshal(&watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         watchmarket.Fixer,
		Rate:             6,
		Timestamp:        timeUPD.Unix(),
	})
	assert.Nil(t, err)
	err = controllerWithCache.cache.Set("c60_ta", wantedTicker1Raw)
	assert.Nil(t, err)
	err = controllerWithCache.cache.Set("c714_ta", wantedTicker2Raw)
	assert.Nil(t, err)
	err = controllerWithCache.cache.Set("USD", rateRaw)
	assert.Nil(t, err)

	response2, err := controllerWithCache.HandleTickersRequest(controllers.TickerRequest{Currency: "USD", Assets: []controllers.Asset{{CoinId: 60, TokenId: "a"}, {CoinId: 714, TokenId: "a"}}})
	assert.Nil(t, err)

	sort.Slice(response2, func(i, j int) bool {
		return response2[i].Coin < response2[j].Coin
	})
	for i := range wantedResp.Tickers {
		assert.True(t, wantedResp.Tickers[i].LastUpdate.Equal(response2[i].LastUpdate))
		wantedResp.Tickers[i].LastUpdate = response2[i].LastUpdate
	}
	assert.Equal(t, wantedResp.Tickers, response2)
}

func TestController_HandleTickersRequest_Negative(t *testing.T) {
	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedRatesError = errors.New("not found")
	c := setupController(t, db, false)
	assert.NotNil(t, c)

	_, err := c.HandleTickersRequest(controllers.TickerRequest{})
	assert.Equal(t, err, errors.New(watchmarket.ErrNotFound))
}

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock(), false))
}

func setupController(t *testing.T, d dbMock, useMemoryCache bool) Controller {
	c, _ := config.Init("../../../config.yml")
	assert.NotNil(t, c)
	c.RestAPI.UseMemoryCache = useMemoryCache

	ratesPriority := c.Markets.Priority.Rates
	tickerPriority := c.Markets.Priority.Tickers
	var ch cache.Provider
	if useMemoryCache {
		ch = memory.Init()
	}
	controller := NewController(d, ch, ratesPriority, tickerPriority, c)
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

func (d dbMock) GetRatesByProvider(provider string) ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) AddTickers(tickers []models.Ticker) error {
	return nil
}

func (d dbMock) GetAllTickers() ([]models.Ticker, error) {
	return nil, nil
}

func (d dbMock) GetAllRates() ([]models.Rate, error) {
	return nil, nil
}

func (d dbMock) GetTickers(assets []controllers.Asset) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}
