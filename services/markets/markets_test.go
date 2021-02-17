package markets

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/services/assets"
	"testing"
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

	assert.Equal(t, "coingecko", apis.ChartsAPIs["coingecko"].GetProvider())
	assert.Equal(t, "coinmarketcap", apis.ChartsAPIs["coinmarketcap"].GetProvider())

	assert.Equal(t, "fixer", apis.RatesAPIs["fixer"].GetProvider())
	assert.Equal(t, "coinmarketcap", apis.RatesAPIs["coinmarketcap"].GetProvider())
	assert.Equal(t, "coingecko", apis.RatesAPIs["coingecko"].GetProvider())

	assert.Equal(t, "coinmarketcap", apis.TickersAPIs["coinmarketcap"].GetProvider())
	assert.Equal(t, "coingecko", apis.TickersAPIs["coingecko"].GetProvider())
}
