package cache

import (
	"encoding/json"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
	"time"
)

func TestUnixToDuration(t *testing.T) {
	wantedDuration := time.Second * 10

	assert.Equal(t, wantedDuration, UnixToDuration(10))
	assert.Equal(t, time.Second*0, UnixToDuration(0))
	assert.Equal(t, time.Minute, UnixToDuration(60))
}

func TestDurationToUnix(t *testing.T) {
	wantedUnixTime := 10
	assert.Equal(t, uint(wantedUnixTime), DurationToUnix(time.Second*10))
	assert.Equal(t, uint(0), DurationToUnix(time.Second*0))
}

func TestInstance_GetCharts_notOutdated(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)
	seedDbCharts(t, i)

	data, err := i.getIntervalKey("testKEY", 1)
	assert.NotNil(t, data)
	assert.Nil(t, err)

	charts, err := i.GetCharts("testKEY", 0)
	assert.Nil(t, err)
	assert.Equal(t, makeChartDataMock(), charts)
}

func TestInstance_GetCharts_CachingDataWasEmpty(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	res, err := json.Marshal([]CachedInterval{{Timestamp: 0, Duration: 100000, Key: "A"}})
	assert.Nil(t, err)
	assert.NotNil(t, res)

	err = i.redis.Set("testKEY", res, UnixToDuration(1000))
	assert.Nil(t, err)

	data, err := i.GetCharts("testKEY", 10000)
	assert.Equal(t, "Not found", err.Error())
	assert.Equal(t, watchmarket.Chart{}, data)
}

func TestInstance_GetCharts_notExistingKey(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()
	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	seedDbCharts(t, i)

	data, err := i.GetCharts("testKEY+1", 1)
	assert.Equal(t, "Not found", err.Error())
	assert.Equal(t, watchmarket.Chart{}, data)
}

func TestInstance_GetCharts_Outdated(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()
	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	data, err := i.GetCharts("testKEY", 100000)
	assert.Equal(t, watchmarket.Chart{}, data)
	assert.NotNil(t, err)
}

func TestInstance_GetCharts_OutdatedCacheIsNotReturned(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	data, err := i.GetCharts("testKEY", 100000)
	assert.Equal(t, watchmarket.Chart{}, data)
	assert.NotNil(t, err)

	res, err := i.redis.Get("testKEY")
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestInstance_GetCharts_ValidCacheIsReturned(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	seedDbCharts(t, i)

	data, err := i.GetCharts("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	res, err := i.redis.Get("data_key")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestInstance_GetCharts_StartTimeIsEarlierThatWasCached(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	data, err := i.GetCharts("testKEY", -1)
	assert.Equal(t, watchmarket.Chart{}, data)
	assert.NotNil(t, err)

	res, err := i.redis.Get("testKEY")
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// emulate that cache was created
	seedDbCharts(t, i)

	dataTwo, err := i.GetCharts("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), dataTwo)
	assert.Nil(t, err)

	resTwo, err := i.redis.Get("data_key")
	assert.Nil(t, err)
	assert.NotNil(t, resTwo)
}

func TestInstance_GetCharts(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	err = r.Set("data_key", []byte{0, 1, 2}, time.Minute)
	assert.Nil(t, err)

	err = i.updateInterval("testKEY", CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
	assert.Nil(t, err)

	data, err := i.GetCharts("testKEY", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "cache is not valid", err.Error())
	assert.Equal(t, watchmarket.Chart{}, data)

	res, err := i.redis.Get("testKEY")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestInstance_SaveCharts(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	err = i.SaveCharts("testKEY", makeChartDataMock(), 0)
	assert.Nil(t, err)

	res, err := i.redis.Get("xQNa0B7ITYf1gJY0dGG3fabGPic=")
	mocked, _ := makeRawDataMockCharts()
	assert.Equal(t, mocked, res)
	assert.Nil(t, err)
}

func TestProvider_Mixed(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)
	err = i.SaveCharts("testKEY", makeChartDataMock(), 0)
	assert.Nil(t, err)

	data, err := i.GetCharts("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	dataTwo, err := i.GetCharts("testKEY", 10001)
	assert.NotNil(t, err)
	assert.Equal(t, "no suitable intervals", err.Error())
	assert.Equal(t, watchmarket.Chart{}, dataTwo)
}

func TestInstance_SaveCharts_DataIsEmpty(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	err = i.SaveCharts("testKEY", watchmarket.Chart{Prices: nil, Error: ""}, 0)
	assert.Equal(t, "data is empty", err.Error())
	res, err := i.GetCharts("testKEY", 0)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.Chart{}, res)
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func seedDbCharts(t *testing.T, instance Instance) {
	rawData, err := makeRawDataMockCharts()
	assert.NotNil(t, rawData)
	assert.Nil(t, err)
	_ = instance.updateInterval("testKEY", CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
	_ = instance.redis.Set("data_key", rawData, UnixToDuration(1000))

}

func makeRawDataMockCharts() ([]byte, error) {
	rawData, err := json.Marshal(makeChartDataMock())
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func makeChartDataMock() watchmarket.Chart {
	price := watchmarket.ChartPrice{
		Price: 100000,
		Date:  0,
	}

	prices := make([]watchmarket.ChartPrice, 0)
	prices = append(prices, price)
	prices = append(prices, price)

	return watchmarket.Chart{
		Prices: prices,
		Error:  "",
	}
}
