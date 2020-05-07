package cache

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
	seedDbCharts(t, db)
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

	err = db.Set("testKEY", r)
	assert.Nil(t, err)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 10000)
	assert.Equal(t, "record does not exist", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_GetChartsCache_notExistingKey(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbCharts(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY+1", 1)
	assert.Equal(t, "record does not exist", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_GetChartsCache_Outdated(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbCharts(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100000)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)
}

func TestProvider_GetChartsCache_OutdatedCacheIsNotReturned(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbCharts(t, db)
	provider := InitCaching(db)

	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100000)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res)
}

func TestProvider_GetChartsCache_ValidCacheIsReturned(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbCharts(t, db)

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	res, err := provider.DB.Get("data_key")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestProvider_GetChartsCache_StartTimeIsEarlierThatWasCached(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbCharts(t, db)

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)

	data, err := provider.GetChartsCache("testKEY", -1)
	assert.Equal(t, watchmarket.ChartData{}, data)
	assert.NotNil(t, err)

	res, err := provider.DB.Get("testKEY")
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Nil(t, res)

	// emulate that cache was created
	seedDbCharts(t, db)

	dataTwo, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), dataTwo)
	assert.Nil(t, err)

	resTwo, err := provider.DB.Get("data_key")
	assert.Nil(t, err)
	assert.NotNil(t, resTwo)
}

func TestProvider_GetChartsCache_BadCachingDataWasDeletedAndHandledRight(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	err := db.Set("data_key", []byte{0, 1, 2})
	assert.Nil(t, err)
	err = db.UpdateInterval("testKEY", storage.CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
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
	assert.Nil(t, res)
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

	res, err := provider.DB.Get("xQNa0B7ITYf1gJY0dGG3fabGPic=")
	mocked, _ := makeRawDataMockCharts()
	assert.Equal(t, mocked, res)
	assert.Nil(t, err)
}

func TestProvider_Mixed(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingDuration(testedCachingDuration)
	err := provider.SaveChartsCache("testKEY", makeChartDataMock(), 0)
	assert.Nil(t, err)

	data, err := provider.GetChartsCache("testKEY", 100)
	assert.Equal(t, makeChartDataMock(), data)
	assert.Nil(t, err)

	dataTwo, err := provider.GetChartsCache("testKEY", 10001)
	assert.NotNil(t, err)
	assert.Equal(t, "no suitable intervals", err.Error())
	assert.Equal(t, watchmarket.ChartData{}, dataTwo)
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
	res, err := provider.GetChartsCache("testKEY", 0)
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Equal(t, watchmarket.ChartData{}, res)
}

func TestProvider_GetChartsCache_FailedToDBGet(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")

	mockDb.On("GetHMValue", storage.EntityInterval, "testKEY", mock.Anything).Return(addHMErr)

	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)

	data, err := provider.GetChartsCache("testKEY", 0)
	assert.Equal(t, addHMErr, err)
	assert.Equal(t, watchmarket.ChartData{}, data)
}

func TestProvider_SaveChartsCache_FailedToDBSet(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")

	mockDb.On("AddHM", storage.EntityInterval, "testKEY", mock.Anything).Return(addHMErr)
	mockDb.On("GetHMValue", storage.EntityInterval, "testKEY", mock.Anything).Return(nil)

	SetChartsCachingDuration(testedCachingDuration)
	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)

	err := provider.SaveChartsCache("testKEY", makeChartDataMock(), 0)
	assert.Equal(t, addHMErr, err)
}

func TestProvider_GetCoinInfoCache(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbChartsInfo(t, db)
	provider := InitCaching(db)
	assert.NotNil(t, provider)

	SetChartsCachingInfoDuration(testedCachingDuration)

	data, err := provider.GetCoinInfoCache("testKEY", 1)
	assert.NotNil(t, data)
	assert.Nil(t, err)
	assert.Equal(t, makeChartInfoMock(), data)
}

func TestProvider_GetCoinInfoCache_Expired(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbChartsInfo(t, db)
	provider := InitCaching(db)
	assert.NotNil(t, provider)

	SetChartsCachingInfoDuration(testedCachingDuration)

	data, err := provider.GetCoinInfoCache("testKEY", 1001)
	assert.NotNil(t, data)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.ChartCoinInfo{}, data)
}

func TestProvider_GetCoinInfoCache_Mixed(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))
	seedDbChartsInfo(t, db)
	provider := InitCaching(db)
	assert.NotNil(t, provider)

	SetChartsCachingInfoDuration(testedCachingDuration)

	data, err := provider.GetCoinInfoCache("testKEY", 1001)
	assert.NotNil(t, data)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.ChartCoinInfo{}, data)

	dataTwo, err := provider.GetCoinInfoCache("testKEY", 101)
	assert.NotNil(t, dataTwo)
	assert.Nil(t, err)
	assert.Equal(t, makeChartInfoMock(), dataTwo)
}

func TestProvider_SaveCoinInfoCache(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	db := InitRedis(fmt.Sprintf("redis://%s", s.Addr()))

	provider := InitCaching(db)
	assert.NotNil(t, provider)
	SetChartsCachingInfoDuration(testedCachingDuration)

	err := provider.SaveCoinInfoCache("testKEY", watchmarket.ChartCoinInfo{}, 0)
	assert.Equal(t, "data is empty", err.Error())
	res, err := provider.GetCoinInfoCache("testKEY", 0)
	assert.NotNil(t, err)
	assert.Equal(t, storage.ErrNotExist, err)
	assert.Equal(t, watchmarket.ChartCoinInfo{}, res)
}

func TestProvider_SaveCoinInfoCache_DbFails(t *testing.T) {
	mockDb := &mocks.DB{}

	addHMErr := errors.New("boom")

	mockDb.On("AddHM", storage.EntityInterval, "testKEY", mock.Anything).Return(addHMErr)
	mockDb.On("GetHMValue", storage.EntityInterval, "testKEY", mock.Anything).Return(nil)

	SetChartsCachingInfoDuration(testedCachingDuration)
	provider := InitCaching(&storage.Storage{DB: mockDb})
	assert.NotNil(t, provider)

	err := provider.SaveCoinInfoCache("testKEY", makeChartInfoMock(), 0)
	assert.Equal(t, addHMErr, err)
}

func setupRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func seedDbCharts(t *testing.T, db storage.Caching) {
	rawData, err := makeRawDataMockCharts()
	assert.NotNil(t, rawData)
	assert.Nil(t, err)
	_ = db.UpdateInterval("testKEY", storage.CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
	_ = db.Set("data_key", rawData)

}

func seedDbChartsInfo(t *testing.T, db storage.Caching) {

	rawData, err := makeRawDataMockChartsInfo()
	assert.NotNil(t, rawData)
	assert.Nil(t, err)
	_ = db.UpdateInterval("testKEY", storage.CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
	_ = db.Set("data_key", rawData)

}

func makeRawDataMockCharts() ([]byte, error) {
	rawData, err := json.Marshal(makeChartDataMock())
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func makeRawDataMockChartsInfo() ([]byte, error) {
	rawData, err := json.Marshal(makeChartInfoMock())
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func makeChartInfoMock() watchmarket.ChartCoinInfo {
	info := watchmarket.CoinInfo{
		Name:             "name test",
		Website:          "website test",
		SourceCode:       "source code",
		WhitePaper:       "paper",
		Description:      "desc",
		ShortDescription: "short",
		Explorer:         "explorer",
		Socials:          nil,
	}
	return watchmarket.ChartCoinInfo{
		Vol24:             10,
		MarketCap:         10,
		CirculatingSupply: 11,
		TotalSupply:       33,
		Info:              &info,
	}
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
