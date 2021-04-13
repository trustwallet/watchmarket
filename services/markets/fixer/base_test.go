package fixer

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/golibs/mock"
)

var (
	mockedResponse, _ = mock.JsonStringFromFilePath("mocks/fixer_response.json")
	wantedRates, _    = mock.JsonStringFromFilePath("mocks/rates.json")
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api", "key", "USD")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.client.BaseUrl)
	assert.Equal(t, "key", provider.client.key)
	assert.Equal(t, watchmarket.Fixer, provider.id)
	assert.Equal(t, "USD", provider.currency)
}

func TestProvider_GetProvider(t *testing.T) {
	provider := InitProvider("demo.api", "key", "USD")
	assert.Equal(t, watchmarket.Fixer, provider.GetProvider())
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedResponse); err != nil {
			panic(err)
		}
	})

	return r
}
