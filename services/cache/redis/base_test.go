package rediscache

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute, time.Minute, time.Minute, time.Minute)

	assert.NotNil(t, i)
	assert.True(t, i.redis.IsAvailable())
	assert.Equal(t, i.chartsCaching, time.Minute)
	assert.Equal(t, i.tickersCaching, time.Minute)
	assert.Equal(t, i.ratesCaching, time.Minute)
	assert.Equal(t, i.detailsCaching, time.Minute)
}

func TestInstance_GenerateKey(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute, time.Minute, time.Minute, time.Minute)

	expected := "bc1M4j2I4u6VaLpUbAB8Y9kTHBs="

	assert.Equal(t, expected, i.GenerateKey("A"))
	assert.NotEqual(t, expected, i.GenerateKey("a"))
}

func TestInstance_GetID(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	i := Init(r, time.Minute, time.Minute, time.Minute, time.Minute)

	assert.NotNil(t, i)

	assert.Equal(t, "redis", i.GetID())
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}
