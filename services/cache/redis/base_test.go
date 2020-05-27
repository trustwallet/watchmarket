package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute)

	assert.NotNil(t, i)
	assert.True(t, i.redis.IsAvailable())
	assert.Equal(t, i.cachingPeriod, time.Minute)
}

func TestInstance_GenerateKey(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute)

	expected := "bc1M4j2I4u6VaLpUbAB8Y9kTHBs="

	assert.Equal(t, expected, i.GenerateKey("A"))
	assert.NotEqual(t, expected, i.GenerateKey("a"))
}

func TestInstance_GetID(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute)

	assert.NotNil(t, i)

	assert.Equal(t, "redis", i.GetID())
}

func TestInstance_Get(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	ticker := watchmarket.Ticker{
		Coin:       60,
		CoinName:   "",
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
	assert.Nil(t, i.redis.Set("test", d, i.cachingPeriod, context.Background()))

	nd, err := i.Get("test", context.Background())
	assert.Nil(t, err)
	var ta watchmarket.Tickers
	assert.Nil(t, json.Unmarshal(nd, &ta))
	assert.Equal(t, tickers, ta)
}

func TestInstance_Set(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	ticker := watchmarket.Ticker{
		Coin:       60,
		CoinName:   "",
		CoinType:   "coin",
		Error:      "",
		LastUpdate: time.Time{},
		Price:      watchmarket.Price{},
		TokenId:    "",
	}
	tickers := watchmarket.Tickers{ticker}

	d, err := json.Marshal(tickers)
	assert.Nil(t, err)
	err = i.Set("test", d, context.Background())
	assert.Nil(t, err)

	nd, err := i.Get("test", context.Background())
	assert.Nil(t, err)
	var ta watchmarket.Tickers
	assert.Nil(t, json.Unmarshal(nd, &ta))
	assert.Equal(t, tickers, ta)
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}
