package cache

import (
	"fmt"
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
	i := Init(r, time.Minute)

	assert.NotNil(t, i)
	assert.True(t, i.redis.IsAvailable())
	assert.Equal(t, i.chartsCaching, time.Minute)
}

func TestGenerateKey(t *testing.T) {
	expected := "bc1M4j2I4u6VaLpUbAB8Y9kTHBs="

	assert.Equal(t, expected, GenerateKey("A"))
	assert.NotEqual(t, expected, GenerateKey("a"))
}
