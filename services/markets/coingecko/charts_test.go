package coingecko

import (
	"encoding/json"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/assets"
)

func TestProvider_GetCoinData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, "USD", assets.Init(server.URL))
	data, _ := provider.GetCoinData(60, "", "USD")
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedInfo, string(rawData))
}

func TestProvider_GetChartData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, "USD", assets.Init("assets.api"))
	data, _ := provider.GetChartData(60, "", "USD", 1577871126)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	isSorted := sort.SliceIsSorted(data.Prices, func(i, j int) bool {
		return data.Prices[i].Date < data.Prices[j].Date
	})
	assert.True(t, isSorted)
	assert.JSONEq(t, wantedCharts, string(rawData))
}
