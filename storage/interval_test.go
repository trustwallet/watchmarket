package storage

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
)

func TestStorage_GetIntervalKey(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  2,
		Key:       "3",
	}

	err := s.DB.AddHM(EntityInterval, "1", []CachedInterval{excpected})
	assert.Nil(t, err)

	res, err := s.GetIntervalKey("1", 1)
	assert.Nil(t, err)
	assert.Equal(t, excpected.Key, res)
}

func TestStorage_GetIntervalKey_Mixed(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  2,
		Key:       "3",
	}
	ivOne := CachedInterval{
		Timestamp: 3,
		Duration:  2,
		Key:       "4",
	}
	err := s.DB.AddHM(EntityInterval, "1", []CachedInterval{excpected, ivOne})
	assert.Nil(t, err)

	res, err := s.GetIntervalKey("1", 2)
	assert.Nil(t, err)
	assert.Equal(t, excpected.Key, res)

	resTwo, err := s.GetIntervalKey("1", 3)
	assert.Nil(t, err)
	assert.Equal(t, ivOne.Key, resTwo)
}

func TestStorage_GetIntervalKey_MultIntervals(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  5,
		Key:       "3",
	}
	excpectedM := CachedInterval{
		Timestamp: 2,
		Duration:  6,
		Key:       "4",
	}
	err := s.DB.AddHM(EntityInterval, "1", []CachedInterval{excpected, excpectedM})
	assert.Nil(t, err)

	res, err := s.GetIntervalKey("1", 3)
	assert.Nil(t, err)
	assert.Equal(t, "3", res)

}

func TestStorage_GetIntervalKey_Empty(t *testing.T) {
	s := initRedis(t)
	res, err := s.GetIntervalKey("empty", 1)
	assert.Equal(t, ErrNotExist, err)
	assert.Equal(t, "", res)
}

func TestStorage_GetIntervalKey_ToOld(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  2,
		Key:       "3",
	}
	err := s.UpdateInterval("1", excpected)
	assert.Nil(t, err)
	res, err := s.GetIntervalKey("1", 4)
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestStorage_GetIntervalKey_ToEarly(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  2,
		Key:       "3",
	}
	err := s.UpdateInterval("1", excpected)
	assert.Nil(t, err)
	res, err := s.GetIntervalKey("1", 0)
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestStorage_UpdateInterval(t *testing.T) {
	s := initRedis(t)
	excpected := CachedInterval{
		Timestamp: 1,
		Duration:  2,
		Key:       "3",
	}
	err := s.UpdateInterval("1", excpected)
	assert.Nil(t, err)

	var d []CachedInterval
	err = s.DB.GetHMValue(EntityInterval, "1", &d)
	assert.Nil(t, err)

	assert.Equal(t, []CachedInterval{excpected}, d)
}

func initRedis(t *testing.T) *Storage {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	cache := &Storage{DB: &redis.Redis{}}
	err = cache.Init(fmt.Sprintf("redis://%s", s.Addr()))
	if err != nil {
		logger.Fatal(err)
	}
	return cache
}
