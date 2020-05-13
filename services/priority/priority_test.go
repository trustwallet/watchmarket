package priority

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"testing"
)

func TestInit(t *testing.T) {
	var providersList = []string{"coinmarketcap", "fixer", "coingecko"}
	c, err := Init(providersList)
	assert.Nil(t, err)
	assert.Equal(t, "coinmarketcap", c.GetCurrentProvider())
	assert.Equal(t, providersList, c.GetAllProviders())
}

func TestController_GetNextProvider(t *testing.T) {
	var providersList = []string{"coinmarketcap", "fixer", "coingecko"}
	c, err := Init(providersList)
	assert.Nil(t, err)

	nextProvider, err := c.GetNextProvider()
	assert.Nil(t, err)
	assert.Equal(t, "fixer", nextProvider)

	anotherProvider, err := c.GetNextProvider()
	assert.Nil(t, err)
	assert.Equal(t, "coingecko", anotherProvider)

	empty, err := c.GetNextProvider()
	assert.Equal(t, errors.E("There is no next provider"), err)
	assert.Equal(t, "", empty)
}

func TestController_GetCurrentProvider(t *testing.T) {
	var providersList = []string{"coinmarketcap", "fixer", "coingecko"}
	c, err := Init(providersList)
	assert.Nil(t, err)
	assert.Equal(t, "coinmarketcap", c.GetCurrentProvider())
	_, err = c.GetNextProvider()
	assert.Nil(t, err)
	assert.Equal(t, "fixer", c.GetCurrentProvider())
	_, err = c.GetNextProvider()
	assert.Nil(t, err)
	assert.Equal(t, "coingecko", c.GetCurrentProvider())
	_, err = c.GetNextProvider()
	assert.NotNil(t, err)
	assert.Equal(t, "coingecko", c.GetCurrentProvider())
}

func TestController_GetAllProviders(t *testing.T) {
	var providersList = []string{"coinmarketcap", "fixer", "coingecko"}
	c, err := Init(providersList)
	assert.Nil(t, err)
	assert.Equal(t, "coinmarketcap", c.GetCurrentProvider())
	assert.Equal(t, providersList, c.GetAllProviders())
}
