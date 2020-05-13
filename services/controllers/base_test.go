package controllers

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/services/assets"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/markets"
	"github.com/trustwallet/watchmarket/services/priority"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	c := config.Init("../../config/test.yml")
	assert.NotNil(t, c)

	chartsPriority, err := priority.Init(c.Markets.Priority.Charts)
	assert.Nil(t, err)

	ratesPriority, err := priority.Init(c.Markets.Priority.Rates)
	assert.Nil(t, err)

	tickerPriority, err := priority.Init(c.Markets.Priority.Tickers)
	assert.Nil(t, err)

	coinInfoPriority, err := priority.Init(c.Markets.Priority.CoinInfo)
	assert.Nil(t, err)

	a := assets.NewClient(c.Markets.Assets)

	m, err := markets.Init(c, a)
	assert.Nil(t, err)

	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	cacheInstance := cache.Init(r, time.Minute, time.Minute, time.Minute, time.Minute)

	controller := NewController(cacheInstance, chartsPriority, coinInfoPriority, ratesPriority, tickerPriority, m)
	assert.NotNil(t, controller)

	data, err := controller.HandleChartsRequest(ChartRequest{
		coinQuery:    "60",
		token:        "",
		currency:     "USD",
		timeStartRaw: "1577871126",
		maxItems:     "64",
	})

	assert.Nil(t, err)
	assert.NotNil(t, data)

	controller.HandleChartsRequest(ChartRequest{
		coinQuery:    "60",
		token:        "",
		currency:     "USD",
		timeStartRaw: "1577871126",
		maxItems:     "64",
	})
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}
