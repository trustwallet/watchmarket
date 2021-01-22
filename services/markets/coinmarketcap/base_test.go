package coinmarketcap

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/golibs/mock"
	"github.com/trustwallet/watchmarket/services/assets"
)

var (
	testMapping, _             = mock.JsonStringFromFilePath("mocks/mapping.json")
	mockedProApiResponse, _    = mock.JsonStringFromFilePath("mocks/pro_api_response.json")
	mockedAssetsApiResponse, _ = mock.JsonStringFromFilePath("mocks/assets_api_response.json")
	wantedRates, _             = mock.JsonStringFromFilePath("mocks/rates.json")

	wantedChartsSorted, _      = mock.JsonStringFromFilePath("mocks/charts_sorted.json")
	wantedCoinInfo, _          = mock.JsonStringFromFilePath("mocks/coin_info.json")
	mockedCmcResponse, _       = mock.JsonStringFromFilePath("mocks/cmc_response.json")
	mockedInfoResponse, _      = mock.JsonStringFromFilePath("mocks/info_response.json")
	mockedChartsCmcResponse, _ = mock.JsonStringFromFilePath("mocks/charts_cmc_response.json")
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("pro.api", "web.api", "widget.api", "key", "USD", assets.Init("assets.api"))
	assert.NotNil(t, provider)
	cm, err := setupCoinMap(testMapping)
	assert.Nil(t, err)
	provider.Cm = cm
	assert.Equal(t, "pro.api", provider.client.proApiURL)
	assert.Equal(t, "web.api", provider.client.webApiURL)
	assert.Equal(t, "widget.api", provider.client.widgetApiURL)
	assert.Equal(t, "USD", provider.currency)
	assert.Equal(t, "coinmarketcap", provider.id)
	assert.Less(t, 1, len(provider.Cm))
}

func TestProvider_GetProvider(t *testing.T) {
	provider := InitProvider("pro.api", "web.api", "widget.api", "key", "USD",
		assets.Init("assets.api"))
	assert.Equal(t, "coinmarketcap", provider.GetProvider())
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/v1/cryptocurrency/listings/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedProApiResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/Mapping.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedAssetsApiResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v2/ticker/1027", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedCmcResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/ethereum/info/info.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedInfoResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v1/cryptocurrency/quotes/historical", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedChartsCmcResponse); err != nil {
			panic(err)
		}
	})

	return r
}
