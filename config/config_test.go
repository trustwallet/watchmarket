package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	c := Init("test.yml")

	assert.Equal(t, []string{"coinmarketcap", "coingecko"}, c.Markets.Priority.Charts)
	assert.Equal(t, []string{"coinmarketcap", "coingecko", "binancedex"}, c.Markets.Priority.Tickers)
	assert.Equal(t, []string{"fixer", "coinmarketcap", "coingecko"}, c.Markets.Priority.Rates)
	assert.Equal(t, []string{"coinmarketcap", "coingecko"}, c.Markets.Priority.CoinInfo)

	assert.Equal(t, "USD", c.Markets.Coinmarketcap.Currency)
	assert.Equal(t, "https://pro-api.coinmarketcap.com", c.Markets.Coinmarketcap.API)
	assert.Equal(t, "1", c.Markets.Coinmarketcap.Key)
	assert.Equal(t, "https://widgets.coinmarketcap.com", c.Markets.Coinmarketcap.WidgetAPI)
	assert.Equal(t, "https://web-api.coinmarketcap.com", c.Markets.Coinmarketcap.WebAPI)

	assert.Equal(t, "https://api.coingecko.com/api", c.Markets.Coingecko.API)
	assert.Equal(t, "USD", c.Markets.Coingecko.Currency)

	assert.Equal(t, "https://dex.binance.org/api", c.Markets.BinanceDex.API)

	assert.Equal(t, "https://data.fixer.io/api", c.Markets.Fixer.API)
	assert.Equal(t, "1", c.Markets.Fixer.Key)
	assert.Equal(t, "USD", c.Markets.Fixer.Currency)

	assert.Equal(t, "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains", c.Markets.Assets)

	assert.Equal(t, "redis://localhost:6379", c.Storage.Redis)
	assert.Equal(t, "postgresql://user:pass@localhost/watchmarket?sslmode=disable", c.Storage.Postgres.Uri)
	assert.Equal(t, false, c.Storage.Postgres.Logs)
	assert.Equal(t, true, c.Storage.Postgres.APM)

	assert.Equal(t, "5m", c.Worker.Tickers)
	assert.Equal(t, "5m", c.Worker.Rates)
	assert.Equal(t, uint(3000), c.Worker.BatchLimit)
	assert.Equal(t, time.Hour*72, c.RestAPI.Tickers.RespectableUpdateTime)

	assert.Equal(t, time.Minute, c.RestAPI.Tickers.CacheControl)
	assert.Equal(t, time.Minute, c.RestAPI.Charts.CacheControl)
	assert.Equal(t, time.Minute, c.RestAPI.Info.CacheControl)

	assert.Equal(t, "release", c.RestAPI.Mode)
	assert.Equal(t, "8420", c.RestAPI.Port)

	assert.Equal(t, time.Minute*15, c.RestAPI.Cache)
	assert.Equal(t, 300, c.RestAPI.RequestLimit)
	assert.Equal(t, true, c.RestAPI.UseMemoryCache)
	assert.Equal(t, "30s", c.RestAPI.UpdateTime.Tickers)
	assert.Equal(t, "30s", c.RestAPI.UpdateTime.Rates)
}
