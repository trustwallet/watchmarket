package binancedex

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.BaseUrl)
	assert.Equal(t, "binancedex", provider.id)
}

func TestProvider_GetProvider(t *testing.T) {
	provider := InitProvider("demo.api")
	assert.Equal(t, "binancedex", provider.GetProvider())
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/v1/ticker/24hr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		p := CoinPrice{BaseAssetName: "BaseName", QuoteAssetName: BNBAsset, PriceChangePercent: "10", LastPrice: "123"}
		rawBytes, err := json.Marshal([]CoinPrice{p})
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(rawBytes); err != nil {
			panic(err)
		}
	})

	return r
}
