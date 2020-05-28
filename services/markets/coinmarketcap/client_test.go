package coinmarketcap

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("pro.api", "assets.api", "web.api", "widget.api", "key")
	assert.NotNil(t, c)
	assert.Equal(t, "pro.api", c.api.BaseUrl)
	assert.Equal(t, "web.api", c.web.BaseUrl)
	assert.Equal(t, "widget.api", c.widget.BaseUrl)
	assert.Equal(t, "assets.api", c.assets.BaseUrl)
}

func Test_fetchCoinMap(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	client := NewClient(server.URL, server.URL, server.URL, server.URL, "demo.key")
	data, err := client.fetchCoinMap(context.Background())
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedCoinMap, string(rawData))
}
