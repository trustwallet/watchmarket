package rediscache

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
	"time"
)

func TestInstance_GetTickers(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000, time.Second*1000, time.Second*1000, time.Second*1000)
	assert.NotNil(t, i)

	ticker := watchmarket.Ticker{
		Coin:       60,
		CoinName:   "ETH",
		CoinType:   "coin",
		Error:      "",
		LastUpdate: time.Time{},
		Price:      watchmarket.Price{},
		TokenId:    "",
	}
	tickers := watchmarket.Tickers{ticker}

	d, err := json.Marshal(tickers)
	assert.NotNil(t, d)
	assert.Nil(t, err)
	assert.Nil(t, i.redis.Set("test", d, i.tickersCaching))

	newTickers, err := i.GetTickers("test")
	assert.Nil(t, err)
	assert.Equal(t, tickers, newTickers)
}

func TestInstance_SaveTickers(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000, time.Second*1000, time.Second*1000, time.Second*1000)
	assert.NotNil(t, i)

	ticker := watchmarket.Ticker{
		Coin:       60,
		CoinName:   "ETH",
		CoinType:   "coin",
		Error:      "",
		LastUpdate: time.Time{},
		Price:      watchmarket.Price{},
		TokenId:    "",
	}
	tickers := watchmarket.Tickers{ticker}

	err = i.SaveTickers("test", tickers)
	assert.Nil(t, err)

	newTickers, err := i.GetTickers("test")
	assert.Nil(t, err)
	assert.Equal(t, tickers, newTickers)
}
