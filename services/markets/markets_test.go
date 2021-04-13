package markets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/assets"
)

func TestInit(t *testing.T) {
	c, _ := config.Init("../../config.yml")
	assert.NotNil(t, c)

	a := assets.Init(c.Markets.Assets)

	apis, err := Init(c, a)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(apis.ChartsAPIs))
	assert.Equal(t, 3, len(apis.RatesAPIs))
	assert.Equal(t, 2, len(apis.TickersAPIs))

	assert.Equal(t, watchmarket.CoinGecko, apis.ChartsAPIs[watchmarket.CoinGecko].GetProvider())
	assert.Equal(t, watchmarket.CoinMarketCap, apis.ChartsAPIs[watchmarket.CoinMarketCap].GetProvider())

	assert.Equal(t, watchmarket.Fixer, apis.RatesAPIs[watchmarket.Fixer].GetProvider())
	assert.Equal(t, watchmarket.CoinMarketCap, apis.RatesAPIs[watchmarket.CoinMarketCap].GetProvider())
	assert.Equal(t, watchmarket.CoinGecko, apis.RatesAPIs[watchmarket.CoinGecko].GetProvider())

	assert.Equal(t, watchmarket.CoinMarketCap, apis.TickersAPIs[watchmarket.CoinMarketCap].GetProvider())
	assert.Equal(t, watchmarket.CoinGecko, apis.TickersAPIs[watchmarket.CoinGecko].GetProvider())
}
