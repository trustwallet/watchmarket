package worker

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/trustwallet/watchmarket/services/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache/memory"
)

func TestWorker_SaveRatesToMemory(t *testing.T) {
	c, _ := config.Init("../../config.yml")
	assert.NotNil(t, c)

	testRatesBasic(t, c)
	testRatesShowOptionAlways(t, c)
	testRatesShowOptionNever(t, c)
}

func testRatesBasic(t *testing.T, c config.Configuration) {
	now := time.Now()
	dbMock := getDbMock()
	dbMock.WantedRates = []models.Rate{
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinMarketCap,
			PercentChange24h: 1,
			Rate:             1,
			ShowOption:       0,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinGecko,
			PercentChange24h: 2,
			Rate:             2,
			ShowOption:       0,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.Fixer,
			PercentChange24h: 11,
			Rate:             1.5,
			ShowOption:       0,
			LastUpdated:      now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveRatesToMemory()
	resRaw, err := w.cache.Get("USD")
	assert.Nil(t, err)

	var res watchmarket.Rate
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 11,
		Provider:         watchmarket.Fixer,
		Rate:             1.5,
		Timestamp:        now.Unix(),
	}, res)
}

func testRatesShowOptionAlways(t *testing.T, c config.Configuration) {
	now := time.Now()
	dbMock := getDbMock()
	dbMock.WantedRates = []models.Rate{
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinMarketCap,
			PercentChange24h: 1,
			Rate:             1,
			ShowOption:       0,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinGecko,
			PercentChange24h: 2,
			Rate:             2,
			ShowOption:       0,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.Fixer,
			PercentChange24h: 11,
			Rate:             1.5,
			ShowOption:       models.NeverShow,
			LastUpdated:      now,
		},
	}

	w2 := Init(nil, nil, dbMock, memory.Init(), c)
	w2.SaveRatesToMemory()
	resRaw2, err := w2.cache.Get("USD")
	assert.Nil(t, err)

	var res2 watchmarket.Rate
	assert.Nil(t, json.Unmarshal(resRaw2, &res2))
	assert.Equal(t, watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         watchmarket.CoinMarketCap,
		Rate:             1,
		Timestamp:        now.Unix(),
	}, res2)
}

func testRatesShowOptionNever(t *testing.T, c config.Configuration) {
	now := time.Now()
	dbMock := getDbMock()
	dbMock.WantedRates = []models.Rate{
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinMarketCap,
			PercentChange24h: 1,
			Rate:             1,
			ShowOption:       0,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.CoinGecko,
			PercentChange24h: 2,
			Rate:             2,
			ShowOption:       models.AlwaysShow,
			LastUpdated:      now,
		},
		{
			Currency:         "USD",
			Provider:         watchmarket.Fixer,
			PercentChange24h: 11,
			Rate:             1.5,
			ShowOption:       0,
			LastUpdated:      now,
		},
	}

	w3 := Init(nil, nil, dbMock, memory.Init(), c)
	w3.SaveRatesToMemory()
	resRaw3, err := w3.cache.Get("USD")
	assert.Nil(t, err)

	var res3 watchmarket.Rate
	assert.Nil(t, json.Unmarshal(resRaw3, &res3))
	assert.Equal(t, watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         watchmarket.CoinGecko,
		Rate:             2,
		Timestamp:        now.Unix(),
	}, res3)

}

func TestWorker_SaveTickersToMemory(t *testing.T) {
	c, _ := config.Init("../../config.yml")
	assert.NotNil(t, c)

	testTickersBasic(t, c)
	testTickersShowOptionNever(t, c)
	testTickersShowOptionAlways(t, c)
	testTickersOutdated(t, c)
	testTickersVolume(t, c)
}

func testTickersBasic(t *testing.T, c config.Configuration) {
	dbMock := getDbMock()
	now := time.Now()
	dbMock.WantedTickers = []models.Ticker{
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       11,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinGecko,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   2,
			Value:       12,
			LastUpdated: now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveTickersToMemory()
	resRaw, err := w.cache.Get("c1")
	assert.Nil(t, err)

	var res watchmarket.Ticker
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Ticker{
		Coin:    1,
		TokenId: "",
		Price: watchmarket.Price{
			Change24h: 1,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     11,
		},
		LastUpdate: res.LastUpdate,
	}, res)
	assert.Equal(t, res.LastUpdate.Unix(), now.Unix())
}

func testTickersShowOptionNever(t *testing.T, c config.Configuration) {
	dbMock := getDbMock()
	now := time.Now()
	dbMock.WantedTickers = []models.Ticker{
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			ShowOption:  2,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       11,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinGecko,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   2,
			Value:       12,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       14,
			ShowOption:  0,
			LastUpdated: now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveTickersToMemory()
	resRaw, err := w.cache.Get("c1")
	assert.Nil(t, err)

	var res watchmarket.Ticker
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Ticker{
		Coin:    1,
		TokenId: "",
		Price: watchmarket.Price{
			Change24h: 1,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     14,
		},
		LastUpdate: res.LastUpdate,
	}, res)
	assert.Equal(t, res.LastUpdate.Unix(), now.Unix())
}

func testTickersShowOptionAlways(t *testing.T, c config.Configuration) {
	dbMock := getDbMock()
	now := time.Now()
	dbMock.WantedTickers = []models.Ticker{
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			ShowOption:  2,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       11,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinGecko,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   2,
			Value:       12,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       14,
			ShowOption:  1,
			LastUpdated: now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveTickersToMemory()
	resRaw, err := w.cache.Get("c1")
	assert.Nil(t, err)

	var res watchmarket.Ticker
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Ticker{
		Coin:    1,
		TokenId: "",
		Price: watchmarket.Price{
			Change24h: 1,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     14,
		},
		LastUpdate: res.LastUpdate,
	}, res)
	assert.Equal(t, res.LastUpdate.Unix(), now.Unix())
}

func testTickersOutdated(t *testing.T, c config.Configuration) {
	dbMock := getDbMock()
	now := time.Now()
	dbMock.WantedTickers = []models.Ticker{
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			ShowOption:  2,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       11,
			LastUpdated: time.Date(1999, 1, 1, 1, 1, 1, 1, time.Local),
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinGecko,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   2,
			Value:       12,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       14,
			ShowOption:  0,
			LastUpdated: now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveTickersToMemory()
	resRaw, err := w.cache.Get("c1")
	assert.Nil(t, err)

	var res watchmarket.Ticker
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Ticker{
		Coin:    1,
		TokenId: "",
		Price: watchmarket.Price{
			Change24h: 1,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     14,
		},
		LastUpdate: res.LastUpdate,
	}, res)
	assert.Equal(t, res.LastUpdate.Unix(), now.Unix())
}

func testTickersVolume(t *testing.T, c config.Configuration) {
	dbMock := getDbMock()
	now := time.Now()
	c.RestAPI.Tickers.RespsectableVolume = 999
	dbMock.WantedTickers = []models.Ticker{
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			ShowOption:  2,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       11,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinGecko,
			ShowOption:  0,
			Coin:        1,
			TokenId:     "",
			Change24h:   2,
			Value:       12,
			LastUpdated: now,
		},
		{
			ID:          "c1",
			Currency:    "USD",
			Provider:    watchmarket.CoinMarketCap,
			Coin:        1,
			TokenId:     "",
			Change24h:   1,
			Value:       14,
			Volume:      1000,
			ShowOption:  0,
			LastUpdated: now,
		},
	}

	w := Init(nil, nil, dbMock, memory.Init(), c)
	w.SaveTickersToMemory()
	resRaw, err := w.cache.Get("c1")
	assert.Nil(t, err)

	var res watchmarket.Ticker
	assert.Nil(t, json.Unmarshal(resRaw, &res))
	assert.Equal(t, watchmarket.Ticker{
		Coin:    1,
		TokenId: "",
		Price: watchmarket.Price{
			Change24h: 1,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     14,
		},
		LastUpdate: res.LastUpdate,
	}, res)
	assert.Equal(t, res.LastUpdate.Unix(), now.Unix())
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
	return d.WantedTickers, nil
}

func (d dbMock) GetAllRates() ([]models.Rate, error) {
	return d.WantedRates, nil
}

func (d dbMock) GetTickers(asset []controllers.Asset) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}
