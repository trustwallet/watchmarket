package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetupBasicAPI(t *testing.T) {
	e := setupEngine()
	server := httptest.NewServer(e)
	defer server.Close()
	SetupBasicAPI(e)

	go func() {
		if err := e.Run(":8080"); err != nil {
			logger.Error(err)
		}
	}()

	resp3, err := http.Get(server.URL)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp3.Body)
	assert.Nil(t, err)
	assert.Equal(t, `"Watchmarket API"`, string(body))
}

func TestSetupTickersAPI(t *testing.T) {
	e := setupEngine()
	server := httptest.NewServer(e)
	defer server.Close()

	wantedTickers := controllers.TickerResponse{
		Currency: "USD",
		Tickers: []watchmarket.Ticker{
			{
				Coin: 60, TokenId: "a",
				Price: watchmarket.Price{
					Change24h: 2,
					Currency:  "",
					Provider:  "coinmarketcap",
					Value:     1,
				},
			},
		},
	}
	wantedTickersV2 := controllers.TickerResponseV2{
		Currency: "USD",
		Tickers:  []controllers.TickerPrice{{Change24h: 2, Provider: "coinmarketcap", Price: 1, ID: "c60_ta"}},
	}

	SetupTickersAPI(e, getTickersMock(wantedTickers, wantedTickersV2, nil), time.Minute)

	go func() {
		if err := e.Run(":8083"); err != nil {
			logger.Error(err)
		}
	}()

	givenV1Resp := controllers.TickerResponse{}

	cr1 := controllers.TickerRequest{
		Currency: "USD",
		Assets:   []controllers.Coin{{Coin: 60, TokenId: "a"}},
	}

	rawcr1, err := json.Marshal(&cr1)
	assert.Nil(t, err)

	resp, err := http.Post(server.URL+"/v1/market/ticker", "application/json", bytes.NewBuffer(rawcr1))
	assert.Nil(t, err)

	err = json.NewDecoder(resp.Body).Decode(&givenV1Resp)
	assert.Nil(t, err)
	assert.Equal(t, uint(60), givenV1Resp.Tickers[0].Coin)
	assert.Equal(t, "a", givenV1Resp.Tickers[0].TokenId)
	assert.Equal(t, float64(1), givenV1Resp.Tickers[0].Price.Value)
	assert.Equal(t, float64(2), givenV1Resp.Tickers[0].Price.Change24h)
	assert.Equal(t, "coinmarketcap", givenV1Resp.Tickers[0].Price.Provider)

	givenV2Resp := controllers.TickerResponseV2{}

	cr2 := controllers.TickerRequestV2{
		Currency: "USD",
		Ids:      []string{"60_a"},
	}

	rawcr2, err := json.Marshal(&cr2)
	assert.Nil(t, err)

	resp2, err := http.Post(server.URL+"/v2/market/tickers", "application/json", bytes.NewBuffer(rawcr2))
	assert.Nil(t, err)

	err = json.NewDecoder(resp2.Body).Decode(&givenV2Resp)
	assert.Nil(t, err)
	assert.Equal(t, "c60_ta", givenV2Resp.Tickers[0].ID)
	assert.Equal(t, float64(1), givenV2Resp.Tickers[0].Price)
	assert.Equal(t, float64(2), givenV2Resp.Tickers[0].Change24h)
	assert.Equal(t, "coinmarketcap", givenV2Resp.Tickers[0].Provider)

	resp3, err := http.Get(server.URL + "/v2/market/ticker/c60_ta")
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp3.Body)
	assert.Nil(t, err)

	givenV2Resp2 := controllers.TickerResponseV2{}

	err = json.Unmarshal(body, &givenV2Resp2)
	assert.Nil(t, err)

	assert.Equal(t, "c60_ta", givenV2Resp2.Tickers[0].ID)
	assert.Equal(t, float64(1), givenV2Resp2.Tickers[0].Price)
	assert.Equal(t, float64(2), givenV2Resp2.Tickers[0].Change24h)
	assert.Equal(t, "coinmarketcap", givenV2Resp2.Tickers[0].Provider)
}

func TestSetupChartsAPI(t *testing.T) {
	e := setupEngine()
	server := httptest.NewServer(e)
	defer server.Close()
	wantedCharts := watchmarket.Chart{
		Provider: "coinmarketcap",
		Prices:   []watchmarket.ChartPrice{{Price: 10, Date: 10}},
	}
	SetupChartsAPI(e, getChartsMock(wantedCharts, nil), time.Minute)

	go func() {
		if err := e.Run(":8082"); err != nil {
			logger.Error(err)
		}
	}()

	resp, err := http.Get(server.URL + "/v2/market/charts/c60_ta")
	assert.Nil(t, err)

	givenResp := watchmarket.Chart{}

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &givenResp)
	assert.Nil(t, err)

	assert.Equal(t, wantedCharts, givenResp)

	resp2, err := http.Get(server.URL + "/v1/market/charts?coin=60&token=a&time_start=1000000000")
	assert.Nil(t, err)

	givenResp2 := watchmarket.Chart{}

	body2, err := ioutil.ReadAll(resp2.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body2, &givenResp2)
	assert.Nil(t, err)

	assert.Equal(t, wantedCharts, givenResp2)
	assert.Nil(t, err)
}

func TestSetupInfoAPI(t *testing.T) {
	e := setupEngine()
	server := httptest.NewServer(e)
	defer server.Close()
	wantedInfo := watchmarket.CoinDetails{
		Provider:          "coinmarketcap",
		Vol24:             1,
		MarketCap:         2,
		CirculatingSupply: 3,
		TotalSupply:       4,
		Info: &watchmarket.Info{
			Name:             "a",
			Website:          "b",
			SourceCode:       "c",
			WhitePaper:       "d",
			Description:      "s",
			ShortDescription: "a",
			Explorer:         "q",
			Socials: []watchmarket.SocialLink{
				{
					Name:   "a",
					Url:    "",
					Handle: "",
				},
			},
		},
	}
	SetupInfoAPI(e, getInfoMock(wantedInfo, nil), time.Minute)

	go func() {
		if err := e.Run(":8081"); err != nil {
			logger.Error(err)
		}
	}()

	resp, err := http.Get(server.URL + "/v2/market/info/c60_ta")
	assert.Nil(t, err)

	givenResp := watchmarket.CoinDetails{}

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &givenResp)
	assert.Nil(t, err)

	assert.Equal(t, wantedInfo, givenResp)

	resp2, err := http.Get(server.URL + "/v1/market/info?coin=60&token=a")
	assert.Nil(t, err)

	givenResp2 := watchmarket.CoinDetails{}

	body2, err := ioutil.ReadAll(resp2.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body2, &givenResp2)

	assert.Equal(t, wantedInfo, givenResp2)
	assert.Nil(t, err)
}

func TestSetupSwaggerAPI(t *testing.T) {
	e := setupEngine()
	server := httptest.NewServer(e)
	defer server.Close()
	SetupSwaggerAPI(e)
	go func() {
		if err := e.Run(":8084"); err != nil {
			logger.Error(err)
		}
	}()

	resp, err := http.Get(server.URL + "/swagger/index.html")
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

type (
	chartsControllerMock struct {
		wantedCharts watchmarket.Chart
		wantedError  error
	}

	tickersControllerMock struct {
		wantedTickersV1 controllers.TickerResponse
		wantedTickersV2 controllers.TickerResponseV2
		wantedError     error
	}

	infoControllerMock struct {
		wantedInfo  watchmarket.CoinDetails
		wantedError error
	}
)

func getChartsMock(wantedCharts watchmarket.Chart, wantedError error) controllers.ChartsController {
	return chartsControllerMock{
		wantedCharts,
		wantedError,
	}
}

func getInfoMock(wantedInfo watchmarket.CoinDetails, wantedError error) controllers.InfoController {
	return infoControllerMock{
		wantedInfo,
		wantedError,
	}
}

func getTickersMock(wantedTickersV1 controllers.TickerResponse, wantedTickersV2 controllers.TickerResponseV2, wantedError error) controllers.TickersController {
	return tickersControllerMock{
		wantedTickersV1,
		wantedTickersV2,
		wantedError,
	}
}

func (c chartsControllerMock) HandleChartsRequest(cr controllers.ChartRequest, ctx context.Context) (watchmarket.Chart, error) {
	return c.wantedCharts, c.wantedError
}

func (c infoControllerMock) HandleInfoRequest(dr controllers.DetailsRequest, ctx context.Context) (watchmarket.CoinDetails, error) {
	return c.wantedInfo, c.wantedError
}

func (c tickersControllerMock) HandleTickersRequestV2(tr controllers.TickerRequestV2, ctx context.Context) (controllers.TickerResponseV2, error) {
	return c.wantedTickersV2, c.wantedError
}

func (c tickersControllerMock) HandleTickersRequest(tr controllers.TickerRequest, ctx context.Context) (controllers.TickerResponse, error) {
	return c.wantedTickersV1, c.wantedError
}

func setupEngine() *gin.Engine {
	return gin.New()
}
