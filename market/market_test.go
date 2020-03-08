package market

import (
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/watchmarket/internal"
	"github.com/trustwallet/watchmarket/market/market"
	"github.com/trustwallet/watchmarket/market/rate"
	rateprovider "github.com/trustwallet/watchmarket/mocks/market/rate"
	marketprovider "github.com/trustwallet/watchmarket/mocks/market/market"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"testing"
	"time"
)

func TestMarketObserver(t *testing.T) {
	// Setup
	s := setupRedis(t)
	defer s.Close()

	cache := internal.InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	mockRateProvider := &rateprovider.RateProvider{}
	mockRateProvider.On("Init", mock.AnythingOfType("*storage.Storage")).Return(nil)
	mockRateProvider.On("GetUpdateTime").Return("5m")
	mockRateProvider.On("GetLogType").Return("market-rate")
	mockRateProvider.On("GetId").Return("coingecko")
	testRate := watchmarket.Rate{
		Currency:         "USD",
		Rate:             1.0,
		Timestamp:        time.Now().Unix(),
		PercentChange24h: nil,
		Provider:         "coingecko",
	}
	mockRateProvider.On("FetchLatestRates").Return(watchmarket.Rates{testRate}, nil)
	rateProviders := &rate.Providers{
		0: mockRateProvider,
	}

	mockMarketProvider := &marketprovider.MarketProvider{}
	mockMarketProvider.On("Init", mock.AnythingOfType("*storage.Storage")).Return(nil)
	mockMarketProvider.On("GetUpdateTime").Return("5m")
	mockMarketProvider.On("GetLogType").Return("market-data")
	mockMarketProvider.On("GetId").Return("coingecko")
	coinObj, ok := coin.Coins[60] // ETH
	if !ok {
		t.Fatal(errors.New("coin does not exist"))
	}
	testTicker := watchmarket.Ticker{
		Coin:       coinObj.ID,
		CoinName:   coinObj.Symbol,
		TokenId:    "",
		CoinType:   "tbd",
		Price:      watchmarket.TickerPrice{
			Value:     50,
			Change24h: 0,
			Currency:  "USD",
			Provider:  "coingecko",
		},
		LastUpdate: time.Now(),
	}
	mockMarketProvider.On("GetData").Return(watchmarket.Tickers{&testTicker}, nil)
	marketProviders := &market.Providers{
		0: mockMarketProvider,
	}

	// Act
	rateCron := InitRates(cache, rateProviders)
	defer rateCron.Stop()
	rateCron.Start()
	marketCron := InitMarkets(cache, marketProviders)
	defer marketCron.Stop()
	marketCron.Start()

	// Verify
	resultRate, err := cache.GetRate("USD")
	if err != nil {
		t.Fatal(err)
	}
	resultTicker, err := cache.GetTicker("ETH", "")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resultRate.Currency, testRate.Currency)
	assert.Equal(t, resultRate.Provider, testRate.Provider)
	assert.Equal(t, resultRate.Rate, testRate.Rate)
	assert.Equal(t, resultRate.PercentChange24h, testRate.PercentChange24h)
	assert.Equal(t, resultTicker.TokenId, testTicker.TokenId)
	assert.Equal(t, resultTicker.Coin, testTicker.Coin)
	assert.Equal(t, resultTicker.CoinName, testTicker.CoinName)
	assert.Equal(t, resultTicker.CoinType, testTicker.CoinType)
	assert.Equal(t, resultTicker.Price.Provider, testTicker.Price.Provider)
	assert.Equal(t, resultTicker.Price.Currency, testTicker.Price.Currency)
	assert.Equal(t, resultTicker.Price.Provider, testTicker.Price.Provider)
	assert.Equal(t, resultTicker.Price.Value, testTicker.Price.Value)
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}