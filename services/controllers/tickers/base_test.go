package tickerscontroller

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"sort"
	"testing"
	"time"
)

func TestController_HandleTickersRequest(t *testing.T) {
	timeUPD := time.Now()
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		LastUpdated:      timeUPD,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         "coingecko",
		Rate:             2,
		LastUpdated:      timeUPD,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         "fixer",
		Rate:             6,
		LastUpdated:      timeUPD,
	}

	ticker60ACMC := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coinmarketcap",
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker60ACG := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coingecko",
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker714ACG := models.Ticker{
		Coin:        714,
		CoinName:    "BNB",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coingecko",
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
	c := setupController(t, db)
	assert.NotNil(t, c)

	response, err := c.HandleTickersRequest(controllers.TickerRequest{Currency: "USD", Assets: []controllers.Coin{{Coin: 60, TokenId: "a"}, {Coin: 714, TokenId: "a"}}}, context.Background())
	assert.Nil(t, err)

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
			Provider:  "coingecko",
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
	sort.Slice(response.Tickers, func(i, j int) bool {
		return response.Tickers[i].Coin < response.Tickers[j].Coin
	})

	assert.Equal(t, wantedResp, response)
}

func TestController_HandleTickersRequest_Negative(t *testing.T) {
	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedRatesError = errors.New("not found")
	c := setupController(t, db)
	assert.NotNil(t, c)

	_, err := c.HandleTickersRequest(controllers.TickerRequest{}, context.Background())
	assert.Equal(t, err, errors.New(watchmarket.ErrBadRequest))
}

func TestController_HandleTickersRequestV2(t *testing.T) {
	timeUPD := time.Now()
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         "coinmarketcap",
		Rate:             1,
		LastUpdated:      timeUPD,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         "coingecko",
		Rate:             2,
		LastUpdated:      timeUPD,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         "fixer",
		Rate:             6,
		LastUpdated:      timeUPD,
	}

	ticker60ACMC := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coinmarketcap",
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker60ACG := models.Ticker{
		Coin:        60,
		CoinName:    "ETH",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coingecko",
		Value:       100,
		LastUpdated: timeUPD,
	}

	ticker714ACG := models.Ticker{
		Coin:        714,
		CoinName:    "BNB",
		TokenId:     "a",
		Change24h:   10,
		Currency:    "USD",
		Provider:    "coingecko",
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
	c := setupController(t, db)
	assert.NotNil(t, c)

	response, err := c.HandleTickersRequestV2(controllers.TickerRequestV2{Currency: "USD", Ids: []string{"60_a", "714_a"}}, context.Background())
	assert.Nil(t, err)

	wantedTicker1 := controllers.TickerPrice{
		ID:        "60_a",
		Change24h: 10,
		Provider:  "coinmarketcap",
		Price:     100,
	}
	wantedTicker2 := controllers.TickerPrice{
		ID:        "714_a",
		Change24h: 10,
		Provider:  "coingecko",
		Price:     100,
	}

	wantedResp := controllers.TickerResponseV2{
		Currency: "USD",
		Tickers:  []controllers.TickerPrice{wantedTicker2, wantedTicker1},
	}

	sort.Slice(wantedResp.Tickers, func(i, j int) bool {
		return wantedResp.Tickers[i].Change24h < wantedResp.Tickers[j].Change24h
	})
	sort.Slice(response.Tickers, func(i, j int) bool {
		return response.Tickers[i].Change24h < response.Tickers[j].Change24h
	})

	assert.Equal(t, wantedResp, response)
}

func TestNewController(t *testing.T) {
	assert.NotNil(t, setupController(t, getDbMock()))
}

func setupController(t *testing.T, d dbMock) Controller {
	c := config.Init("../../../config/test.yml")
	assert.NotNil(t, c)

	ratesPriority := c.Markets.Priority.Rates
	tickerPriority := c.Markets.Priority.Tickers

	controller := NewController(d, ratesPriority, tickerPriority, c)
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

func (d dbMock) GetRates(currency string, ctx context.Context) ([]models.Rate, error) {
	res := make([]models.Rate, 0)
	for _, r := range d.WantedRates {
		if r.Currency == currency {
			res = append(res, r)
		}
	}
	return res, d.WantedRatesError
}

func (d dbMock) AddRates(rates []models.Rate, batchLimit uint, ctx context.Context) error {
	return nil
}

func (d dbMock) AddTickers(tickers []models.Ticker, batchLimit uint, ctx context.Context) error {
	return nil
}

func (d dbMock) GetTickers(coin uint, tokenId string, ctx context.Context) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func (d dbMock) GetTickersByQueries(tickerQueries []models.TickerQuery, ctx context.Context) ([]models.Ticker, error) {
	return d.WantedTickers, d.WantedTickersError
}

func getCacheMock() cache.Provider {
	i := cacheMock{}
	return i
}

type cacheMock struct {
}

func (c cacheMock) GetID() string {
	return ""
}

func (c cacheMock) GenerateKey(data string) string {
	return ""
}

func (c cacheMock) Get(key string, ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) Set(key string, data []byte, ctx context.Context) error {
	return nil
}

func (c cacheMock) GetWithTime(key string, time int64, ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (c cacheMock) SetWithTime(key string, data []byte, time int64, ctx context.Context) error {
	return nil
}
