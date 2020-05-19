package cache

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
	p := Init(r, time.Minute, time.Minute, time.Minute, time.Minute)
	assert.NotNil(t, p)
	assert.Equal(t, "redis", p.RatesCache["redis"].GetID())
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}
