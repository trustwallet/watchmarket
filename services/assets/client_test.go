//nolint:unparam
package assets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"

	"github.com/stretchr/testify/assert"
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

func Test_normalize(t *testing.T) {
	type args struct {
		info watchmarket.Info
	}
	tests := []struct {
		name string
		args args
		want watchmarket.Info
	}{
		{
			"Url properly formatted",
			args{info: watchmarket.Info{
				Website: "https://www.google.com",
			}},
			watchmarket.Info{
				Website: "https://google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
