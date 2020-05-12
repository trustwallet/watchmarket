package cache

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/redis"
	"testing"
	"time"
)

func TestInstance_GetCoinDetails(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)
	seedDbChartsInfo(t, i)

	data, err := i.GetCoinDetails("testKEY", 1)
	assert.NotNil(t, data)
	assert.Nil(t, err)
	assert.Equal(t, makeChartInfoMock(), data)
}

func TestProvider_GetCoinInfoCache_Expired(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)
	seedDbChartsInfo(t, i)

	data, err := i.GetCoinDetails("testKEY", 1001)
	assert.NotNil(t, data)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.CoinDetails{}, data)
}

func TestProvider_GetCoinInfoCache_Mixed(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)
	seedDbChartsInfo(t, i)

	data, err := i.GetCoinDetails("testKEY", 1001)
	assert.NotNil(t, data)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.CoinDetails{}, data)

	dataTwo, err := i.GetCoinDetails("testKEY", 101)
	assert.NotNil(t, dataTwo)
	assert.Nil(t, err)
	assert.Equal(t, makeChartInfoMock(), dataTwo)
}

func TestInstance_SaveCoinDetails(t *testing.T) {
	s := setupRedis(t)
	defer s.Close()

	r, err := redis.Init(fmt.Sprintf("redis://%s", s.Addr()))
	assert.Nil(t, err)

	i := Init(r, time.Second*1000)
	assert.NotNil(t, i)

	err = i.SaveCoinDetails("testKEY", watchmarket.CoinDetails{}, 0)
	assert.Equal(t, "data is empty", err.Error())
	res, err := i.GetCoinDetails("testKEY", 0)
	assert.NotNil(t, err)
	assert.Equal(t, watchmarket.CoinDetails{}, res)
}

func makeRawDataMockChartsInfo() ([]byte, error) {
	rawData, err := json.Marshal(makeChartInfoMock())
	if err != nil {
		return nil, err
	}

	return rawData, nil
}
func makeChartInfoMock() watchmarket.CoinDetails {
	info := watchmarket.Info{
		Name:             "name test",
		Website:          "website test",
		SourceCode:       "source code",
		WhitePaper:       "paper",
		Description:      "desc",
		ShortDescription: "short",
		Explorer:         "explorer",
		Socials:          nil,
	}
	return watchmarket.CoinDetails{
		Vol24:             10,
		MarketCap:         10,
		CirculatingSupply: 11,
		TotalSupply:       33,
		Info:              info,
	}
}

func seedDbChartsInfo(t *testing.T, i Instance) {
	rawData, err := makeRawDataMockChartsInfo()
	assert.NotNil(t, rawData)
	assert.Nil(t, err)
	_ = i.updateInterval("testKEY", CachedInterval{
		Timestamp: 0,
		Duration:  1000,
		Key:       "data_key",
	})
	_ = i.redis.Set("data_key", rawData, i.chartsCaching)

}
