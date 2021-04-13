package config

import (
	"testing"
	"time"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	c, _ := Init("../config.yml")

	assert.Equal(t, []string{watchmarket.CoinMarketCap, watchmarket.CoinGecko}, c.Markets.Priority.Charts)
	assert.Equal(t, []string{watchmarket.Fixer, watchmarket.CoinMarketCap, watchmarket.CoinGecko}, c.Markets.Priority.Rates)
	assert.Equal(t, []string{watchmarket.CoinMarketCap, watchmarket.CoinGecko}, c.Markets.Priority.CoinInfo)

	assert.Equal(t, "USD", c.Markets.Coinmarketcap.Currency)
	assert.Equal(t, "https://pro-api.coinmarketcap.com", c.Markets.Coinmarketcap.API)
	assert.Equal(t, "", c.Markets.Coinmarketcap.Key)
	assert.Equal(t, "https://3rdparty-apis.coinmarketcap.com", c.Markets.Coinmarketcap.WidgetAPI)
	assert.Equal(t, "https://web-api.coinmarketcap.com", c.Markets.Coinmarketcap.WebAPI)

	assert.Equal(t, "https://api.coingecko.com/api", c.Markets.Coingecko.API)
	assert.Equal(t, "USD", c.Markets.Coingecko.Currency)

	assert.Equal(t, "https://data.fixer.io/api", c.Markets.Fixer.API)
	assert.Equal(t, "", c.Markets.Fixer.Key)
	assert.Equal(t, "USD", c.Markets.Fixer.Currency)

	assert.Equal(t, "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains", c.Markets.Assets)

	assert.Equal(t, "redis://localhost:6379", c.Storage.Redis.Url)
	assert.Equal(t, "postgresql://user:pass@localhost/watchmarket?sslmode=disable", c.Storage.Postgres.Url)
	assert.Equal(t, false, c.Storage.Postgres.Logs)

	assert.Equal(t, "5m", c.Worker.Tickers)
	assert.Equal(t, "5m", c.Worker.Rates)
	assert.Equal(t, time.Hour*72, c.RestAPI.Tickers.RespectableUpdateTime)

	assert.Equal(t, time.Minute, c.RestAPI.Tickers.CacheControl)
	assert.Equal(t, time.Minute*10, c.RestAPI.Charts.CacheControl)
	assert.Equal(t, time.Minute*10, c.RestAPI.Info.CacheControl)

	assert.Equal(t, "release", c.RestAPI.Mode)
	assert.Equal(t, "8421", c.RestAPI.Port)

	assert.Equal(t, time.Minute*15, c.RestAPI.Cache)
	assert.Equal(t, true, c.RestAPI.UseMemoryCache)
	assert.Equal(t, "5m", c.RestAPI.UpdateTime.Tickers)
	assert.Equal(t, "5m", c.RestAPI.UpdateTime.Rates)
}
