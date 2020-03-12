package caching

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trustwallet/blockatlas/pkg/logger"
	mocks "github.com/trustwallet/watchmarket/mocks/storage"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"github.com/trustwallet/watchmarket/storage"
	"testing"
)

const testedCachingDuration int64 = 60 * 5

func TestProvider_GetChartsCache_notOutdated(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)
	provider := InitCaching(db)
	assert.NotNil(t, provider)

	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 1)
	assert.NotNil(t, data)
	assert.Nil(t, err)
	assert.Equal(t, makeChartDataMock(), data)
}

func TestProvider_GetChartsCache_CachingDataWasEmpty(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	r, err := json.Marshal(watchmarket.ChartData{Prices: nil, Error: ""})
	assert.Nil(t, err)
	assert.NotNil(t, r)

	_, err = db.Set("testKEY", storage.CacheData{RawData: r, WasSavedTime: 0})
	assert.Nil(t, err)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 1)
	assert.Equal(t, "cache is not valid", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)
}

func TestProvider_GetChartsCache_notExistingKey(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY+1", 1)
	assert.Equal(t, "record does not exist", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_GetChartsCache_Delete_Error(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")
	mockDb.On("GetHMValue", storage.EntityCache, "testKEY", mock.AnythingOfType("*storage.CacheData")).Return(nil)
	mockDb.On("DeleteHM", storage.EntityCache, "testKEY").Return(addHMErr)

	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 0)
	assert.Equal(t, "invalid cache is not deleted", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_GetChartsCache_Outdated(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100000)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)
}

func TestProvider_GetChartsCache_OutdatedCacheIsDeleted(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100000)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)
}

func TestProvider_GetChartsCache_ValidCacheIsNotDeleted(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	res, err := provider.DB.Get("testKEY")
	assert.Nil(t, err)
	assert.NotNil(t, res.RawData)
}

func TestProvider_GetChartsCache_StartTimeIsEarlierThatWasCached(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDb(t, db)

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", -1)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)

	// emulate that cache was created
	seedDb(t, db)

	dataTwo, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), dataTwo)
	assert.Nil(t, err)

	resTwo, err := provider.DB.Get("testKEY")
	assert.Nil(t, err)
	assert.NotNil(t, resTwo.RawData)
}

func TestProvider_GetChartsCache_BadCachingDataWasDeletedAndHandledRight(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	_, err := db.Set("testKEY", storage.CacheData{RawData: []byte{0, 1, 2}, WasSavedTime: 0})
	assert.Nil(t, err)

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "cache is not valid", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)
}

func TestProvider_SaveChartsCache_Success(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	err := provider.SaveChartsCache("testKEY", makeChartDataMock(), 0)
	assert.Nil(t, err)

	res, err := provider.DB.Get("testKEY")
	mocked, _ := makeRawDataMock()
	assert.Equal(t, mocked, res.RawData)
	assert.Nil(t, err)
}

func TestProvider_Mixed(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	provider := InitCaching(db)
	assert.NotNil(t, provider)

	err := provider.SaveChartsCache("testKEY", makeChartDataMock(), 0)
	assert.Nil(t, err)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	dataTwo, err := provider.GetChartsCache("testKEY", 10001)
	assert.NotNil(t, err)
	assert.Equal(t, "cache is not valid", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, dataTwo)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)
}

func TestProvider_SaveChartsCache_DataIsEmpty(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	err := provider.SaveChartsCache("testKEY", watchmarket.ChartData{Prices: nil, Error: ""}, 0)
	assert.Equal(t, "data is empty", err.Error())
	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res.RawData)
}

func TestProvider_GetChartsCache_FailedToDBGet(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")

	mockDb.On("GetHMValue", storage.EntityCache, "testKEY", mock.AnythingOfType("*storage.CacheData")).Return(addHMErr)

	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)

	data, err := provider.GetChartsCache("testKEY", 0)
	assert.Equal(t, addHMErr, err)
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_SaveChartsCache_FailedToDBSet(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")

	mockDb.On("AddHM", storage.EntityCache, "testKEY", mock.AnythingOfType("*storage.CacheData")).Return(addHMErr)

	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	err := provider.SaveChartsCache("testKEY", makeChartDataMock(), 0)
	assert.Equal(t, addHMErr, err)
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func seedDb(t *testing.T, db storage.Caching) {
	rawData, err := makeRawDataMock()
	assert.NotNil(t, rawData)
	assert.Nil(t, err)
	db.Set("testKEY", storage.CacheData{RawData: rawData, WasSavedTime: 0})
}

func makeRawDataMock() ([]byte, error) {
	rawData, err := json.Marshal(makeChartDataMock())
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func makeChartDataMock() watchmarket.ChartData {
	price := watchmarket.ChartPrice{
		Price: 100000,
		Date:  0,
	}

	prices := make([]watchmarket.ChartPrice, 0)
	prices = append(prices, price)
	prices = append(prices, price)

	return watchmarket.ChartData{
		Prices: prices,
		Error:  "",
	}
}

func InitRedis(host string) *storage.Storage {
	cache := &storage.Storage{DB: &redis.Redis{}}
	err := cache.Init(host)
	if err != nil {
		logger.Fatal(err)
	}
	return cache
}
