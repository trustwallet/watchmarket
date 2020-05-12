package coingecko

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/assets"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestProvider_GetCoinData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, "USD", assets.NewClient(server.URL))
	data, _ := provider.GetCoinData(60, "", "USD")
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedInfo, string(rawData))
}

func TestProvider_GetChartData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, server.URL, "USD", assets.NewClient("assets.api"))
	data, _ := provider.GetChartData(60, "", "USD", 1577871126)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	isSorted := sort.SliceIsSorted(data.Prices, func(i, j int) bool {
		return data.Prices[i].Date < data.Prices[j].Date
	})
	assert.True(t, isSorted)
	assert.Equal(t, wantedCharts, string(rawData))
}
