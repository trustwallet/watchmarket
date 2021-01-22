//nolint:unparam
package assets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/golibs/mock"
)

var (
	wantedInfo, _         = mock.JsonStringFromFilePath("mocks/info.json")
	mockedInfoResponse, _ = mock.JsonStringFromFilePath("mocks/info_response.json")
)

func TestInit(t *testing.T) {
	c := Init("url")
	assert.NotNil(t, c)
	assert.Equal(t, c.api, "url")
}

func TestClient_GetCoinInfo(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	c := Init(server.URL)
	assert.NotNil(t, c)

	data, err := c.GetCoinInfo(60, "")
	assert.Nil(t, err)

	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedInfo, string(rawData))
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/ethereum/info/info.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedInfoResponse); err != nil {
			panic(err)
		}
	})

	return r
}
