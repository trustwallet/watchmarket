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

func TestInstance_GetRates(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000, time.Second*1000, time.Second*1000, time.Second*1000)
	assert.NotNil(t, i)

	rate := watchmarket.Rate{
		Currency:  "USD",
		Provider:  "",
		Rate:      0,
		Timestamp: 0,
	}
	rates := watchmarket.Rates{rate}

	d, err := json.Marshal(rates)
	assert.NotNil(t, d)
	assert.Nil(t, err)
	assert.Nil(t, i.redis.Set("test", d, i.ratesCaching))

	newRates, err := i.GetRates("test")
	assert.Nil(t, err)
	assert.Equal(t, rates[0].Rate, newRates[0].Rate)
	assert.Equal(t, rates[0].Timestamp, newRates[0].Timestamp)
	assert.Equal(t, rates[0].Provider, newRates[0].Provider)
	assert.Equal(t, rates[0].Currency, newRates[0].Currency)
}

func TestInstance_SaveRates(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)
	rate := watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: float64(0),
		Provider:         "",
		Rate:             0,
		Timestamp:        0,
	}
	rates := watchmarket.Rates{rate}
	i := Init(r, time.Second*1000, time.Second*1000, time.Second*1000, time.Second*1000)
	assert.NotNil(t, i)

	d, err := json.Marshal(rates)
	assert.NotNil(t, d)
	assert.Nil(t, err)
	assert.Nil(t, i.redis.Set("test", d, i.ratesCaching))

	err = i.SaveRates("test", rates)
	assert.Nil(t, err)

	newRates, err := i.GetRates("test")
	assert.Nil(t, err)
	assert.Equal(t, rates[0].Rate, newRates[0].Rate)
	assert.Equal(t, rates[0].Timestamp, newRates[0].Timestamp)
	assert.Equal(t, rates[0].Provider, newRates[0].Provider)
	assert.Equal(t, rates[0].Currency, newRates[0].Currency)
}
