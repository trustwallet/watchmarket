//nolint:unparam
package assets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit(t *testing.T) {
	c := Init("url")
	assert.NotNil(t, c)
	assert.Equal(t, c.BaseUrl, "url")
}

func TestClient_GetCoinInfo(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	c := Init(server.URL)
	assert.NotNil(t, c)

	data, err := c.GetCoinInfo(60, "", context.Background())
	assert.Nil(t, err)

	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedInfo, string(rawData))
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

var (
	wantedInfo         = `{"name":"Ethereum","website":"https://ethereum.org/","source_code":"https://github.com/ethereum","white_paper":"https://github.com/ethereum/wiki/wiki/White-Paper","description":"Open source platform to write and distribute decentralized applications.","short_description":"Open source platform to write and distribute decentralized applications.","explorer":"https://etherscan.io/","socials":[{"name":"Twitter","url":"https://twitter.com/ethereum","handle":"ethereum"},{"name":"Reddit","url":"https://www.reddit.com/r/ethereum","handle":"ethereum"}]}`
	mockedInfoResponse = `{ "name": "Ethereum", "website": "https://ethereum.org/", "source_code": "https://github.com/ethereum", "white_paper": "https://github.com/ethereum/wiki/wiki/White-Paper", "short_description": "Open source platform to write and distribute decentralized applications.", "description": "Open source platform to write and distribute decentralized applications.", "socials": [ { "name": "Twitter", "url": "https://twitter.com/ethereum", "handle": "ethereum" }, { "name": "Reddit", "url": "https://www.reddit.com/r/ethereum", "handle": "ethereum" } ], "explorer": "https://etherscan.io/" }`
)
