package coinmarketcap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("pro.api", "assets.api", "web.api", "widget.api", "info.api", "key", "USD")
	assert.NotNil(t, provider)
	assert.Equal(t, "pro.api", provider.client.api.BaseUrl)
	assert.Equal(t, "web.api", provider.client.web.BaseUrl)
	assert.Equal(t, "widget.api", provider.client.widget.BaseUrl)
	assert.Equal(t, "assets.api", provider.client.assets.BaseUrl)
	assert.Equal(t, "info.api", provider.info.BaseUrl)
	assert.Equal(t, "USD", provider.currency)
	assert.Equal(t, "coinmarketcap", provider.ID)
}
