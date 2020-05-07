package coinnmarketcap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/tickers"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api", "demo.key")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
	assert.Equal(t, "demo.key", client.Headers["X-CMC_PRO_API_KEY"])
}

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api", "demo.key", "USD")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.BaseUrl)
	assert.Equal(t, "demo.key", provider.client.Headers["X-CMC_PRO_API_KEY"])
	assert.Equal(t, "coinmarketcap", provider.ID)
	assert.Equal(t, "USD", provider.currency)
}

func TestProvider_GetData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	provider := InitProvider(server.URL, "demo.key", "USD")
	data, _ := provider.GetData()
	assert.True(t, verifyTickers(t, data, data))
}

func verifyTickers(t *testing.T, wantedTickers, givenTickers tickers.Tickers) bool {
	assert.Equal(t, len(givenTickers), len(wantedTickers))
	var counter = 0
	for _, w := range wantedTickers {
		for _, g := range givenTickers {
			if w.CoinName == g.CoinName && w.TokenId == g.TokenId && w.Price == g.Price {
				assert.Equal(t, w.Coin, g.Coin)
				assert.Equal(t, w.TokenId, g.TokenId)
				assert.Equal(t, w.Price, g.Price)
				assert.Equal(t, w.CoinName, g.CoinName)
				assert.Equal(t, w.Error, g.Error)
				assert.Equal(t, w.CoinType, g.CoinType)
				counter++
			}
		}
	}
	if counter == len(givenTickers) {
		return true
	}
	return false
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/v1/cryptocurrency/listings/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, mockedResponse); err != nil {
			panic(err)
		}
	})

	return r
}

func Test_normalizeTickers(t *testing.T) {
	type args struct {
		prices   CoinPrices
		provider string
	}
	tests := []struct {
		name        string
		args        args
		wantTickers tickers.Tickers
	}{
		{
			"test normalize coinnmarketcap quote",
			args{prices: CoinPrices{Data: []Data{
				{Coin: Coin{Symbol: "BTC", Id: 0}, LastUpdated: time.Unix(111, 0), Quote: Quote{
					USD: USD{Price: 223.55, PercentChange24h: 10}}},
				{Coin: Coin{Symbol: "ETH", Id: 60}, LastUpdated: time.Unix(333, 0), Quote: Quote{
					USD: USD{Price: 11.11, PercentChange24h: 20}}},
				{Coin: Coin{Symbol: "SWP", Id: 6969}, LastUpdated: time.Unix(444, 0), Quote: Quote{
					USD: USD{Price: 463.22, PercentChange24h: -3}},
					Platform: Platform{Coin: Coin{Symbol: "ETH"}, TokenAddress: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9"}},
			}}, provider: "coinnmarketcap"},
			tickers.Tickers{
				tickers.Ticker{CoinName: "BTC", CoinType: tickers.Coin, LastUpdate: time.Unix(111, 0),
					Price: tickers.Price{
						Value:     223.55,
						Change24h: 10,
						Currency:  "USD",
						Provider:  "coinnmarketcap",
					},
				},
				tickers.Ticker{CoinName: "ETH", CoinType: tickers.Coin, LastUpdate: time.Unix(333, 0),
					Price: tickers.Price{
						Value:     11.11,
						Change24h: 20,
						Currency:  "USD",
						Provider:  "coinnmarketcap",
					},
				},
				tickers.Ticker{CoinName: "ETH", TokenId: "0x8ce9137d39326ad0cd6491fb5cc0cba0e089b6a9", CoinType: tickers.Token, LastUpdate: time.Unix(444, 0),
					Price: tickers.Price{
						Value:     463.22,
						Change24h: -3,
						Currency:  "USD",
						Provider:  "coinnmarketcap",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTickers := normalizeTickers(tt.args.prices, tt.args.provider, "USD")
			sort.SliceStable(gotTickers, func(i, j int) bool {
				return gotTickers[i].LastUpdate.Unix() < gotTickers[j].LastUpdate.Unix()
			})
			if !assert.Equal(t, len(tt.wantTickers), len(gotTickers)) {
				t.Fatal("invalid tickers length")
			}
			for i, obj := range tt.wantTickers {
				assert.Equal(t, obj, gotTickers[i])
			}
		})
	}
}

var (
	mockedResponse = `{"status":{"timestamp":"2020-05-07T22:29:57.541Z","error_code":0,"error_message":null,"elapsed":11,"credit_count":1,"notice":null},"data":[{"id":1,"name":"Bitcoin","symbol":"BTC","slug":"bitcoin","num_market_pairs":8013,"date_added":"2013-04-28T00:00:00.000Z","tags":["mineable"],"max_supply":21000000,"circulating_supply":18367687,"total_supply":18367687,"platform":null,"cmc_rank":1,"last_updated":"2020-05-07T22:28:34.000Z","quote":{"USD":{"price":9862.53985763,"volume_24h":59939626493.0469,"percent_change_1h":0.506016,"percent_change_24h":5.47477,"percent_change_7d":12.7368,"market_cap":181152045129.97238,"last_updated":"2020-05-07T22:28:34.000Z"}}},{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","num_market_pairs":5141,"date_added":"2015-08-07T00:00:00.000Z","tags":["mineable"],"max_supply":null,"circulating_supply":110833575.5615,"total_supply":110833575.5615,"platform":null,"cmc_rank":2,"last_updated":"2020-05-07T22:28:25.000Z","quote":{"USD":{"price":213.544073721,"volume_24h":23877577618.5091,"percent_change_1h":0.272421,"percent_change_24h":2.48595,"percent_change_7d":1.42623,"market_cap":23667853230.46698,"last_updated":"2020-05-07T22:28:25.000Z"}}},{"id":52,"name":"XRP","symbol":"XRP","slug":"xrp","num_market_pairs":538,"date_added":"2013-08-04T00:00:00.000Z","tags":[],"max_supply":100000000000,"circulating_supply":44112853111,"total_supply":99990976125,"platform":null,"cmc_rank":3,"last_updated":"2020-05-07T22:29:05.000Z","quote":{"USD":{"price":0.219026524813,"volume_24h":2583495738.19139,"percent_change_1h":0.289071,"percent_change_24h":0.191592,"percent_change_7d":1.82831,"market_cap":9661884916.488667,"last_updated":"2020-05-07T22:29:05.000Z"}}},{"id":825,"name":"Tether","symbol":"USDT","slug":"tether","num_market_pairs":4777,"date_added":"2015-02-25T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":6361032509,"total_supply":6998318752.17,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"31"},"cmc_rank":4,"last_updated":"2020-05-07T22:28:22.000Z","quote":{"USD":{"price":1.00583313202,"volume_24h":71816766577.5082,"percent_change_1h":-0.267651,"percent_change_24h":-0.498649,"percent_change_7d":0.0474062,"market_cap":6398137251.408509,"last_updated":"2020-05-07T22:28:22.000Z"}}},{"id":1831,"name":"Bitcoin Cash","symbol":"BCH","slug":"bitcoin-cash","num_market_pairs":476,"date_added":"2017-07-23T00:00:00.000Z","tags":["mineable"],"max_supply":21000000,"circulating_supply":18400675,"total_supply":18400675,"platform":null,"cmc_rank":5,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":253.196839314,"volume_24h":3675705558.16027,"percent_change_1h":0.273533,"percent_change_24h":1.42361,"percent_change_7d":0.434299,"market_cap":4658992751.244137,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":3602,"name":"Bitcoin SV","symbol":"BSV","slug":"bitcoin-sv","num_market_pairs":182,"date_added":"2018-11-09T00:00:00.000Z","tags":["mineable"],"max_supply":21000000,"circulating_supply":18399477.0819233,"total_supply":18399477.0819233,"platform":null,"cmc_rank":6,"last_updated":"2020-05-07T22:29:13.000Z","quote":{"USD":{"price":210.045381795,"volume_24h":2364053730.1132,"percent_change_1h":0.184715,"percent_change_24h":0.471004,"percent_change_7d":0.376569,"market_cap":3864725188.5009317,"last_updated":"2020-05-07T22:29:13.000Z"}}},{"id":2,"name":"Litecoin","symbol":"LTC","slug":"litecoin","num_market_pairs":594,"date_added":"2013-04-28T00:00:00.000Z","tags":["mineable"],"max_supply":84000000,"circulating_supply":64679818.2265572,"total_supply":64679818.2265572,"platform":null,"cmc_rank":7,"last_updated":"2020-05-07T22:29:04.000Z","quote":{"USD":{"price":47.5118592612,"volume_24h":5125345009.06526,"percent_change_1h":0.088664,"percent_change_24h":1.09024,"percent_change_7d":0.729679,"market_cap":3073058420.6201844,"last_updated":"2020-05-07T22:29:04.000Z"}}},{"id":1839,"name":"Binance Coin","symbol":"BNB","slug":"binance-coin","num_market_pairs":309,"date_added":"2017-07-25T00:00:00.000Z","tags":[],"max_supply":187536713,"circulating_supply":155536713,"total_supply":187536713,"platform":null,"cmc_rank":8,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":17.1140649277,"volume_24h":420671647.902459,"percent_change_1h":0.30303,"percent_change_24h":0.906496,"percent_change_7d":-0.331755,"market_cap":2661865404.9230404,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":1765,"name":"EOS","symbol":"EOS","slug":"eos","num_market_pairs":369,"date_added":"2017-07-01T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":922432400.2454,"total_supply":1019132412.2453,"platform":null,"cmc_rank":9,"last_updated":"2020-05-07T22:29:06.000Z","quote":{"USD":{"price":2.76250457645,"volume_24h":4479875945.11953,"percent_change_1h":0.0911718,"percent_change_24h":-0.562061,"percent_change_7d":-3.14574,"market_cap":2548223727.1436753,"last_updated":"2020-05-07T22:29:06.000Z"}}},{"id":2011,"name":"Tezos","symbol":"XTZ","slug":"tezos","num_market_pairs":93,"date_added":"2017-10-06T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":709614503.656283,"total_supply":709614503.656283,"platform":null,"cmc_rank":10,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":2.78977181773,"volume_24h":221474054.711921,"percent_change_1h":0.118328,"percent_change_24h":1.86731,"percent_change_7d":0.115531,"market_cap":1979662543.7527604,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":512,"name":"Stellar","symbol":"XLM","slug":"stellar","num_market_pairs":298,"date_added":"2014-08-05T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":20232489123.5487,"total_supply":50001803861.8717,"platform":null,"cmc_rank":11,"last_updated":"2020-05-07T22:29:05.000Z","quote":{"USD":{"price":0.0728177647038,"volume_24h":677724312.270636,"percent_change_1h":0.0578241,"percent_change_24h":0.655066,"percent_change_7d":5.8243,"market_cap":1473284632.3707619,"last_updated":"2020-05-07T22:29:05.000Z"}}},{"id":2010,"name":"Cardano","symbol":"ADA","slug":"cardano","num_market_pairs":131,"date_added":"2017-10-01T00:00:00.000Z","tags":["mineable"],"max_supply":45000000000,"circulating_supply":25927070538,"total_supply":31112483745,"platform":null,"cmc_rank":12,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":0.0510633352483,"volume_24h":141986443.941559,"percent_change_1h":0.514759,"percent_change_24h":0.826166,"percent_change_7d":4.44178,"market_cap":1323922694.8882158,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":1975,"name":"Chainlink","symbol":"LINK","slug":"chainlink","num_market_pairs":176,"date_added":"2017-09-20T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":350000000,"total_supply":1000000000,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"0x514910771af9ca656af840dff83e8264ecf986ca"},"cmc_rank":13,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":3.77920034377,"volume_24h":418709132.774495,"percent_change_1h":0.533505,"percent_change_24h":2.4664,"percent_change_7d":0.784,"market_cap":1322720120.3195,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":328,"name":"Monero","symbol":"XMR","slug":"monero","num_market_pairs":141,"date_added":"2014-05-21T00:00:00.000Z","tags":["mineable"],"max_supply":null,"circulating_supply":17551706.2831166,"total_supply":17551706.2831166,"platform":null,"cmc_rank":14,"last_updated":"2020-05-07T22:29:03.000Z","quote":{"USD":{"price":63.6958970684,"volume_24h":138594438.229781,"percent_change_1h":-0.0144351,"percent_change_24h":6.70904,"percent_change_7d":1.20588,"market_cap":1117971676.7841847,"last_updated":"2020-05-07T22:29:03.000Z"}}},{"id":3635,"name":"Crypto.com Coin","symbol":"CRO","slug":"crypto-com-coin","num_market_pairs":45,"date_added":"2018-12-14T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":16603196347.032,"total_supply":100000000000,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"0xa0b73e1ff0b80914ab6fe0444e65848c4c34450b"},"cmc_rank":15,"last_updated":"2020-05-07T22:29:12.000Z","quote":{"USD":{"price":0.0661289969569,"volume_24h":13273646.2327268,"percent_change_1h":0.708087,"percent_change_24h":1.58403,"percent_change_7d":12.5725,"market_cap":1097952720.7076924,"last_updated":"2020-05-07T22:29:12.000Z"}}},{"id":1958,"name":"TRON","symbol":"TRX","slug":"tron","num_market_pairs":340,"date_added":"2017-09-13T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":66682072191.4,"total_supply":99281283754.3,"platform":null,"cmc_rank":16,"last_updated":"2020-05-07T22:29:07.000Z","quote":{"USD":{"price":0.0161019012186,"volume_24h":1688361576.97626,"percent_change_1h":0.326526,"percent_change_24h":0.328149,"percent_change_7d":4.27848,"market_cap":1073708139.477477,"last_updated":"2020-05-07T22:29:07.000Z"}}},{"id":3957,"name":"UNUS SED LEO","symbol":"LEO","slug":"unus-sed-leo","num_market_pairs":25,"date_added":"2019-05-21T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":999498892.9,"total_supply":999498892.9,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"0x2af5d2ad76741191d15dfe7bf6ac92d4bd912ca3"},"cmc_rank":17,"last_updated":"2020-05-07T22:29:13.000Z","quote":{"USD":{"price":1.0642871664,"volume_24h":27354757.8747149,"percent_change_1h":-0.607654,"percent_change_24h":-0.475802,"percent_change_7d":-0.413056,"market_cap":1063753844.544478,"last_updated":"2020-05-07T22:29:13.000Z"}}},{"id":2502,"name":"Huobi Token","symbol":"HT","slug":"huobi-token","num_market_pairs":137,"date_added":"2018-02-03T00:00:00.000Z","tags":[],"max_supply":null,"circulating_supply":222668092.971921,"total_supply":500000000,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"0x6f259637dcd74c767781e37bc6133cd6a68aa161"},"cmc_rank":18,"last_updated":"2020-05-07T22:29:10.000Z","quote":{"USD":{"price":4.21654427261,"volume_24h":168275952.111364,"percent_change_1h":0.086173,"percent_change_24h":-0.303191,"percent_change_7d":-0.958224,"market_cap":938889872.1137445,"last_updated":"2020-05-07T22:29:10.000Z"}}},{"id":1321,"name":"Ethereum Classic","symbol":"ETC","slug":"ethereum-classic","num_market_pairs":245,"date_added":"2016-07-24T00:00:00.000Z","tags":["mineable"],"max_supply":210700000,"circulating_supply":116313299,"total_supply":116313299,"platform":null,"cmc_rank":19,"last_updated":"2020-05-07T22:29:04.000Z","quote":{"USD":{"price":7.15171701655,"volume_24h":2363781454.77942,"percent_change_1h":0.629817,"percent_change_24h":-0.447855,"percent_change_7d":9.69334,"market_cap":831839799.7093681,"last_updated":"2020-05-07T22:29:04.000Z"}}},{"id":131,"name":"Dash","symbol":"DASH","slug":"dash","num_market_pairs":287,"date_added":"2014-02-14T00:00:00.000Z","tags":["mineable"],"max_supply":18900000,"circulating_supply":9481861.59422516,"total_supply":9481861.59422516,"platform":null,"cmc_rank":20,"last_updated":"2020-05-07T22:29:04.000Z","quote":{"USD":{"price":79.2516409705,"volume_24h":731145809.286151,"percent_change_1h":0.306309,"percent_change_24h":-0.531529,"percent_change_7d":-3.50444,"market_cap":751453090.7975053,"last_updated":"2020-05-07T22:29:04.000Z"}}}]}`
	wantedTickers  = tickers.Tickers([]tickers.Ticker{
		{
			Coin:     0,
			CoinName: "BTC",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     9862.53985763,
				Change24h: 5.47477,
				Currency:  "USD",
				Provider:  "coinmarketcap",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     60,
			CoinName: "ETH",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     213.544073721,
				Change24h: 2.48595,
				Currency:  "USD",
				Provider:  "coinmarketcap",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
		{
			Coin:     144,
			CoinName: "XRP",
			TokenId:  "",
			CoinType: tickers.Coin,
			Price: tickers.Price{
				Value:     0.219026524813,
				Change24h: 0.191592,
				Currency:  "USD",
				Provider:  "coinmarketcap",
			},
			LastUpdate: time.Now(),
			Error:      "",
		},
	})
)
