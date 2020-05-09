package fixer

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/rates"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api", "key", "USD")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
	assert.Equal(t, "key", client.key)
	assert.Equal(t, "USD", client.currency)
}

func TestInitProvider(t *testing.T) {
	provider := InitProvider("demo.api", "key", "USD")
	assert.NotNil(t, provider)
	assert.Equal(t, "demo.api", provider.client.BaseUrl)
	assert.Equal(t, "key", provider.client.key)
	assert.Equal(t, "fixer", provider.ID)
	assert.Equal(t, "USD", provider.currency)
}

func TestProvider_GetData(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()
	provider := InitProvider(server.URL, "", "USD")
	data, err := provider.GetData()
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Currency < data[j].Currency
	})
	assert.Nil(t, err)
	rawData, err := json.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, wantedRates, string(rawData))
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, mockedResponse); err != nil {
			panic(err)
		}
	})

	return r
}

func Test_normalizeRates(t *testing.T) {
	provider := "binancedex"
	tests := []struct {
		name      string
		latest    Rate
		wantRates rates.Rates
	}{
		{
			"test normalize fixer rate 1",
			Rate{
				Timestamp: 123,
				Rates:     map[string]float64{"USD": 22.111, "BRL": 33.2, "BTC": 44.99},
				UpdatedAt: time.Now(),
			},
			rates.Rates{
				rates.Rate{Currency: "USD", Rate: 22.111, Timestamp: 123, Provider: provider},
				rates.Rate{Currency: "BRL", Rate: 33.2, Timestamp: 123, Provider: provider},
				rates.Rate{Currency: "BTC", Rate: 44.99, Timestamp: 123, Provider: provider},
			},
		},
		{
			"test normalize fixer rate 2",
			Rate{
				Timestamp: 333,
				Rates:     map[string]float64{"LSK": 123.321, "IFC": 34.973, "DUO": 998.3},
				UpdatedAt: time.Now(),
			},
			rates.Rates{
				rates.Rate{Currency: "IFC", Rate: 34.973, Timestamp: 333, Provider: provider},
				rates.Rate{Currency: "LSK", Rate: 123.321, Timestamp: 333, Provider: provider},
				rates.Rate{Currency: "DUO", Rate: 998.3, Timestamp: 333, Provider: provider},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRates := normalizeRates(tt.latest, provider)
			sort.SliceStable(gotRates, func(i, j int) bool {
				return gotRates[i].Rate < gotRates[j].Rate
			})
			if !assert.ObjectsAreEqualValues(gotRates, tt.wantRates) {
				t.Errorf("normalizeRates() = %v, want %v", gotRates, tt.wantRates)
			}
		})
	}
}

var (
	mockedResponse = `{ "success": true, "timestamp": 1589030105, "base": "USD", "date": "2020-05-09", "rates": { "AED": 3.67315, "AFN": 76.203991, "ALL": 114.903989, "AMD": 484.110403, "ANG": 1.795199, "AOA": 551.905041, "ARS": 66.472813, "AUD": 1.530402, "AWG": 1.8, "AZN": 1.70397, "BAM": 1.803839, "BBD": 2.019276, "BDT": 84.99451, "BGN": 1.80549, "BHD": 0.378178, "BIF": 1906, "BMD": 1, "BND": 1.413261, "BOB": 6.895677, "BRL": 5.732404, "BSD": 1.000046, "BTC": 0.000103, "BTN": 75.58479, "BWP": 12.144339, "BYN": 2.435844, "BYR": 19600, "BZD": 2.015864, "CAD": 1.40175, "CDF": 1810.000362, "CHF": 0.970982, "CLF": 0.029928, "CLP": 825.803912, "CNY": 7.074204, "COP": 3896.5, "CRC": 568.93705, "CUC": 1, "CUP": 26.5, "CVE": 102.00036, "CZK": 25.132604, "DJF": 177.720394, "DKK": 6.880904, "DOP": 55.040393, "DZD": 128.33601, "EGP": 15.563556, "ERN": 15.000358, "ETB": 33.703876, "EUR": 0.911536, "FJD": 2.25304, "FKP": 0.805997, "GBP": 0.806062, "GEL": 3.220391, "GGP": 0.805997, "GHS": 5.74504, "GIP": 0.805997, "GMD": 51.18039, "GNF": 9450.000355, "GTQ": 7.705603, "GYD": 209.34286, "HKD": 7.76695, "HNL": 25.030389, "HRK": 6.938593, "HTG": 107.33289, "HUF": 323.140388, "IDR": 14773.435, "ILS": 3.506704, "IMP": 0.805997, "INR": 75.50404, "IQD": 1190, "IRR": 42105.000352, "ISK": 146.240386, "JEP": 0.805997, "JMD": 142.83448, "JOD": 0.709504, "JPY": 106.67604, "KES": 106.050385, "KGS": 78.903801, "KHR": 4110.000351, "KMF": 453.750384, "KPW": 900.000306, "KRW": 1219.803792, "KWD": 0.30935, "KYD": 0.833556, "KZT": 421.99093, "LAK": 9005.000349, "LBP": 1512.763039, "LKR": 186.51808, "LRD": 198.503775, "LSL": 18.390382, "LTL": 2.95274, "LVL": 0.60489, "LYD": 1.420381, "MAD": 9.825039, "MDL": 17.831157, "MGA": 3800.000347, "MKD": 56.826776, "MMK": 1395.135304, "MNT": 2794.123058, "MOP": 7.984837, "MRO": 357.00003, "MUR": 39.710379, "MVR": 15.503741, "MWK": 737.503739, "MXN": 23.672604, "MYR": 4.334039, "MZN": 68.080377, "NAD": 18.530377, "NGN": 390.000344, "NIO": 34.403725, "NOK": 10.217039, "NPR": 120.93559, "NZD": 1.629196, "OMR": 0.383447, "PAB": 1.000138, "PEN": 3.399039, "PGK": 3.430375, "PHP": 50.495039, "PKR": 159.650375, "PLN": 4.20475, "PYG": 6531.670104, "QAR": 3.641038, "RON": 4.453404, "RSD": 108.455038, "RUB": 73.40369, "RWF": 937.5, "SAR": 3.756262, "SBD": 8.267992, "SCR": 17.168052, "SDG": 55.325038, "SEK": 9.771904, "SGD": 1.412704, "SHP": 0.805997, "SLL": 9860.000339, "SOS": 583.000338, "SRD": 7.458038, "STD": 22051.386135, "SVC": 8.751671, "SYP": 514.451644, "SZL": 18.52037, "THB": 32.02037, "TJS": 10.265663, "TMT": 3.51, "TND": 2.912504, "TOP": 2.31435, "TRY": 7.089104, "TTD": 6.757574, "TWD": 29.857038, "TZS": 2314.203635, "UAH": 26.838273, "UGX": 3800.322804, "USD": 1, "UYU": 43.139569, "UZS": 10109.000335, "VEF": 9.987504, "VND": 23400.5, "VUV": 119.848296, "WST": 2.78472, "XAF": 604.99802, "XAG": 0.06474, "XAU": 0.000587, "XCD": 2.70255, "XDR": 0.734807, "XOF": 605.000332, "XPF": 110.4036, "YER": 250.375037, "ZAR": 18.350904, "ZMK": 9001.203593, "ZMW": 18.576435, "ZWL": 322.000001 } }`
	wantedRates    = `[{"currency":"AED","percent_change_24h":"0","provider":"fixer","rate":3.67315,"timestamp":1589030105},{"currency":"AFN","percent_change_24h":"0","provider":"fixer","rate":76.203991,"timestamp":1589030105},{"currency":"ALL","percent_change_24h":"0","provider":"fixer","rate":114.903989,"timestamp":1589030105},{"currency":"AMD","percent_change_24h":"0","provider":"fixer","rate":484.110403,"timestamp":1589030105},{"currency":"ANG","percent_change_24h":"0","provider":"fixer","rate":1.795199,"timestamp":1589030105},{"currency":"AOA","percent_change_24h":"0","provider":"fixer","rate":551.905041,"timestamp":1589030105},{"currency":"ARS","percent_change_24h":"0","provider":"fixer","rate":66.472813,"timestamp":1589030105},{"currency":"AUD","percent_change_24h":"0","provider":"fixer","rate":1.530402,"timestamp":1589030105},{"currency":"AWG","percent_change_24h":"0","provider":"fixer","rate":1.8,"timestamp":1589030105},{"currency":"AZN","percent_change_24h":"0","provider":"fixer","rate":1.70397,"timestamp":1589030105},{"currency":"BAM","percent_change_24h":"0","provider":"fixer","rate":1.803839,"timestamp":1589030105},{"currency":"BBD","percent_change_24h":"0","provider":"fixer","rate":2.019276,"timestamp":1589030105},{"currency":"BDT","percent_change_24h":"0","provider":"fixer","rate":84.99451,"timestamp":1589030105},{"currency":"BGN","percent_change_24h":"0","provider":"fixer","rate":1.80549,"timestamp":1589030105},{"currency":"BHD","percent_change_24h":"0","provider":"fixer","rate":0.378178,"timestamp":1589030105},{"currency":"BIF","percent_change_24h":"0","provider":"fixer","rate":1906,"timestamp":1589030105},{"currency":"BMD","percent_change_24h":"0","provider":"fixer","rate":1,"timestamp":1589030105},{"currency":"BND","percent_change_24h":"0","provider":"fixer","rate":1.413261,"timestamp":1589030105},{"currency":"BOB","percent_change_24h":"0","provider":"fixer","rate":6.895677,"timestamp":1589030105},{"currency":"BRL","percent_change_24h":"0","provider":"fixer","rate":5.732404,"timestamp":1589030105},{"currency":"BSD","percent_change_24h":"0","provider":"fixer","rate":1.000046,"timestamp":1589030105},{"currency":"BTC","percent_change_24h":"0","provider":"fixer","rate":0.000103,"timestamp":1589030105},{"currency":"BTN","percent_change_24h":"0","provider":"fixer","rate":75.58479,"timestamp":1589030105},{"currency":"BWP","percent_change_24h":"0","provider":"fixer","rate":12.144339,"timestamp":1589030105},{"currency":"BYN","percent_change_24h":"0","provider":"fixer","rate":2.435844,"timestamp":1589030105},{"currency":"BYR","percent_change_24h":"0","provider":"fixer","rate":19600,"timestamp":1589030105},{"currency":"BZD","percent_change_24h":"0","provider":"fixer","rate":2.015864,"timestamp":1589030105},{"currency":"CAD","percent_change_24h":"0","provider":"fixer","rate":1.40175,"timestamp":1589030105},{"currency":"CDF","percent_change_24h":"0","provider":"fixer","rate":1810.000362,"timestamp":1589030105},{"currency":"CHF","percent_change_24h":"0","provider":"fixer","rate":0.970982,"timestamp":1589030105},{"currency":"CLF","percent_change_24h":"0","provider":"fixer","rate":0.029928,"timestamp":1589030105},{"currency":"CLP","percent_change_24h":"0","provider":"fixer","rate":825.803912,"timestamp":1589030105},{"currency":"CNY","percent_change_24h":"0","provider":"fixer","rate":7.074204,"timestamp":1589030105},{"currency":"COP","percent_change_24h":"0","provider":"fixer","rate":3896.5,"timestamp":1589030105},{"currency":"CRC","percent_change_24h":"0","provider":"fixer","rate":568.93705,"timestamp":1589030105},{"currency":"CUC","percent_change_24h":"0","provider":"fixer","rate":1,"timestamp":1589030105},{"currency":"CUP","percent_change_24h":"0","provider":"fixer","rate":26.5,"timestamp":1589030105},{"currency":"CVE","percent_change_24h":"0","provider":"fixer","rate":102.00036,"timestamp":1589030105},{"currency":"CZK","percent_change_24h":"0","provider":"fixer","rate":25.132604,"timestamp":1589030105},{"currency":"DJF","percent_change_24h":"0","provider":"fixer","rate":177.720394,"timestamp":1589030105},{"currency":"DKK","percent_change_24h":"0","provider":"fixer","rate":6.880904,"timestamp":1589030105},{"currency":"DOP","percent_change_24h":"0","provider":"fixer","rate":55.040393,"timestamp":1589030105},{"currency":"DZD","percent_change_24h":"0","provider":"fixer","rate":128.33601,"timestamp":1589030105},{"currency":"EGP","percent_change_24h":"0","provider":"fixer","rate":15.563556,"timestamp":1589030105},{"currency":"ERN","percent_change_24h":"0","provider":"fixer","rate":15.000358,"timestamp":1589030105},{"currency":"ETB","percent_change_24h":"0","provider":"fixer","rate":33.703876,"timestamp":1589030105},{"currency":"EUR","percent_change_24h":"0","provider":"fixer","rate":0.911536,"timestamp":1589030105},{"currency":"FJD","percent_change_24h":"0","provider":"fixer","rate":2.25304,"timestamp":1589030105},{"currency":"FKP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"GBP","percent_change_24h":"0","provider":"fixer","rate":0.806062,"timestamp":1589030105},{"currency":"GEL","percent_change_24h":"0","provider":"fixer","rate":3.220391,"timestamp":1589030105},{"currency":"GGP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"GHS","percent_change_24h":"0","provider":"fixer","rate":5.74504,"timestamp":1589030105},{"currency":"GIP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"GMD","percent_change_24h":"0","provider":"fixer","rate":51.18039,"timestamp":1589030105},{"currency":"GNF","percent_change_24h":"0","provider":"fixer","rate":9450.000355,"timestamp":1589030105},{"currency":"GTQ","percent_change_24h":"0","provider":"fixer","rate":7.705603,"timestamp":1589030105},{"currency":"GYD","percent_change_24h":"0","provider":"fixer","rate":209.34286,"timestamp":1589030105},{"currency":"HKD","percent_change_24h":"0","provider":"fixer","rate":7.76695,"timestamp":1589030105},{"currency":"HNL","percent_change_24h":"0","provider":"fixer","rate":25.030389,"timestamp":1589030105},{"currency":"HRK","percent_change_24h":"0","provider":"fixer","rate":6.938593,"timestamp":1589030105},{"currency":"HTG","percent_change_24h":"0","provider":"fixer","rate":107.33289,"timestamp":1589030105},{"currency":"HUF","percent_change_24h":"0","provider":"fixer","rate":323.140388,"timestamp":1589030105},{"currency":"IDR","percent_change_24h":"0","provider":"fixer","rate":14773.435,"timestamp":1589030105},{"currency":"ILS","percent_change_24h":"0","provider":"fixer","rate":3.506704,"timestamp":1589030105},{"currency":"IMP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"INR","percent_change_24h":"0","provider":"fixer","rate":75.50404,"timestamp":1589030105},{"currency":"IQD","percent_change_24h":"0","provider":"fixer","rate":1190,"timestamp":1589030105},{"currency":"IRR","percent_change_24h":"0","provider":"fixer","rate":42105.000352,"timestamp":1589030105},{"currency":"ISK","percent_change_24h":"0","provider":"fixer","rate":146.240386,"timestamp":1589030105},{"currency":"JEP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"JMD","percent_change_24h":"0","provider":"fixer","rate":142.83448,"timestamp":1589030105},{"currency":"JOD","percent_change_24h":"0","provider":"fixer","rate":0.709504,"timestamp":1589030105},{"currency":"JPY","percent_change_24h":"0","provider":"fixer","rate":106.67604,"timestamp":1589030105},{"currency":"KES","percent_change_24h":"0","provider":"fixer","rate":106.050385,"timestamp":1589030105},{"currency":"KGS","percent_change_24h":"0","provider":"fixer","rate":78.903801,"timestamp":1589030105},{"currency":"KHR","percent_change_24h":"0","provider":"fixer","rate":4110.000351,"timestamp":1589030105},{"currency":"KMF","percent_change_24h":"0","provider":"fixer","rate":453.750384,"timestamp":1589030105},{"currency":"KPW","percent_change_24h":"0","provider":"fixer","rate":900.000306,"timestamp":1589030105},{"currency":"KRW","percent_change_24h":"0","provider":"fixer","rate":1219.803792,"timestamp":1589030105},{"currency":"KWD","percent_change_24h":"0","provider":"fixer","rate":0.30935,"timestamp":1589030105},{"currency":"KYD","percent_change_24h":"0","provider":"fixer","rate":0.833556,"timestamp":1589030105},{"currency":"KZT","percent_change_24h":"0","provider":"fixer","rate":421.99093,"timestamp":1589030105},{"currency":"LAK","percent_change_24h":"0","provider":"fixer","rate":9005.000349,"timestamp":1589030105},{"currency":"LBP","percent_change_24h":"0","provider":"fixer","rate":1512.763039,"timestamp":1589030105},{"currency":"LKR","percent_change_24h":"0","provider":"fixer","rate":186.51808,"timestamp":1589030105},{"currency":"LRD","percent_change_24h":"0","provider":"fixer","rate":198.503775,"timestamp":1589030105},{"currency":"LSL","percent_change_24h":"0","provider":"fixer","rate":18.390382,"timestamp":1589030105},{"currency":"LTL","percent_change_24h":"0","provider":"fixer","rate":2.95274,"timestamp":1589030105},{"currency":"LVL","percent_change_24h":"0","provider":"fixer","rate":0.60489,"timestamp":1589030105},{"currency":"LYD","percent_change_24h":"0","provider":"fixer","rate":1.420381,"timestamp":1589030105},{"currency":"MAD","percent_change_24h":"0","provider":"fixer","rate":9.825039,"timestamp":1589030105},{"currency":"MDL","percent_change_24h":"0","provider":"fixer","rate":17.831157,"timestamp":1589030105},{"currency":"MGA","percent_change_24h":"0","provider":"fixer","rate":3800.000347,"timestamp":1589030105},{"currency":"MKD","percent_change_24h":"0","provider":"fixer","rate":56.826776,"timestamp":1589030105},{"currency":"MMK","percent_change_24h":"0","provider":"fixer","rate":1395.135304,"timestamp":1589030105},{"currency":"MNT","percent_change_24h":"0","provider":"fixer","rate":2794.123058,"timestamp":1589030105},{"currency":"MOP","percent_change_24h":"0","provider":"fixer","rate":7.984837,"timestamp":1589030105},{"currency":"MRO","percent_change_24h":"0","provider":"fixer","rate":357.00003,"timestamp":1589030105},{"currency":"MUR","percent_change_24h":"0","provider":"fixer","rate":39.710379,"timestamp":1589030105},{"currency":"MVR","percent_change_24h":"0","provider":"fixer","rate":15.503741,"timestamp":1589030105},{"currency":"MWK","percent_change_24h":"0","provider":"fixer","rate":737.503739,"timestamp":1589030105},{"currency":"MXN","percent_change_24h":"0","provider":"fixer","rate":23.672604,"timestamp":1589030105},{"currency":"MYR","percent_change_24h":"0","provider":"fixer","rate":4.334039,"timestamp":1589030105},{"currency":"MZN","percent_change_24h":"0","provider":"fixer","rate":68.080377,"timestamp":1589030105},{"currency":"NAD","percent_change_24h":"0","provider":"fixer","rate":18.530377,"timestamp":1589030105},{"currency":"NGN","percent_change_24h":"0","provider":"fixer","rate":390.000344,"timestamp":1589030105},{"currency":"NIO","percent_change_24h":"0","provider":"fixer","rate":34.403725,"timestamp":1589030105},{"currency":"NOK","percent_change_24h":"0","provider":"fixer","rate":10.217039,"timestamp":1589030105},{"currency":"NPR","percent_change_24h":"0","provider":"fixer","rate":120.93559,"timestamp":1589030105},{"currency":"NZD","percent_change_24h":"0","provider":"fixer","rate":1.629196,"timestamp":1589030105},{"currency":"OMR","percent_change_24h":"0","provider":"fixer","rate":0.383447,"timestamp":1589030105},{"currency":"PAB","percent_change_24h":"0","provider":"fixer","rate":1.000138,"timestamp":1589030105},{"currency":"PEN","percent_change_24h":"0","provider":"fixer","rate":3.399039,"timestamp":1589030105},{"currency":"PGK","percent_change_24h":"0","provider":"fixer","rate":3.430375,"timestamp":1589030105},{"currency":"PHP","percent_change_24h":"0","provider":"fixer","rate":50.495039,"timestamp":1589030105},{"currency":"PKR","percent_change_24h":"0","provider":"fixer","rate":159.650375,"timestamp":1589030105},{"currency":"PLN","percent_change_24h":"0","provider":"fixer","rate":4.20475,"timestamp":1589030105},{"currency":"PYG","percent_change_24h":"0","provider":"fixer","rate":6531.670104,"timestamp":1589030105},{"currency":"QAR","percent_change_24h":"0","provider":"fixer","rate":3.641038,"timestamp":1589030105},{"currency":"RON","percent_change_24h":"0","provider":"fixer","rate":4.453404,"timestamp":1589030105},{"currency":"RSD","percent_change_24h":"0","provider":"fixer","rate":108.455038,"timestamp":1589030105},{"currency":"RUB","percent_change_24h":"0","provider":"fixer","rate":73.40369,"timestamp":1589030105},{"currency":"RWF","percent_change_24h":"0","provider":"fixer","rate":937.5,"timestamp":1589030105},{"currency":"SAR","percent_change_24h":"0","provider":"fixer","rate":3.756262,"timestamp":1589030105},{"currency":"SBD","percent_change_24h":"0","provider":"fixer","rate":8.267992,"timestamp":1589030105},{"currency":"SCR","percent_change_24h":"0","provider":"fixer","rate":17.168052,"timestamp":1589030105},{"currency":"SDG","percent_change_24h":"0","provider":"fixer","rate":55.325038,"timestamp":1589030105},{"currency":"SEK","percent_change_24h":"0","provider":"fixer","rate":9.771904,"timestamp":1589030105},{"currency":"SGD","percent_change_24h":"0","provider":"fixer","rate":1.412704,"timestamp":1589030105},{"currency":"SHP","percent_change_24h":"0","provider":"fixer","rate":0.805997,"timestamp":1589030105},{"currency":"SLL","percent_change_24h":"0","provider":"fixer","rate":9860.000339,"timestamp":1589030105},{"currency":"SOS","percent_change_24h":"0","provider":"fixer","rate":583.000338,"timestamp":1589030105},{"currency":"SRD","percent_change_24h":"0","provider":"fixer","rate":7.458038,"timestamp":1589030105},{"currency":"STD","percent_change_24h":"0","provider":"fixer","rate":22051.386135,"timestamp":1589030105},{"currency":"SVC","percent_change_24h":"0","provider":"fixer","rate":8.751671,"timestamp":1589030105},{"currency":"SYP","percent_change_24h":"0","provider":"fixer","rate":514.451644,"timestamp":1589030105},{"currency":"SZL","percent_change_24h":"0","provider":"fixer","rate":18.52037,"timestamp":1589030105},{"currency":"THB","percent_change_24h":"0","provider":"fixer","rate":32.02037,"timestamp":1589030105},{"currency":"TJS","percent_change_24h":"0","provider":"fixer","rate":10.265663,"timestamp":1589030105},{"currency":"TMT","percent_change_24h":"0","provider":"fixer","rate":3.51,"timestamp":1589030105},{"currency":"TND","percent_change_24h":"0","provider":"fixer","rate":2.912504,"timestamp":1589030105},{"currency":"TOP","percent_change_24h":"0","provider":"fixer","rate":2.31435,"timestamp":1589030105},{"currency":"TRY","percent_change_24h":"0","provider":"fixer","rate":7.089104,"timestamp":1589030105},{"currency":"TTD","percent_change_24h":"0","provider":"fixer","rate":6.757574,"timestamp":1589030105},{"currency":"TWD","percent_change_24h":"0","provider":"fixer","rate":29.857038,"timestamp":1589030105},{"currency":"TZS","percent_change_24h":"0","provider":"fixer","rate":2314.203635,"timestamp":1589030105},{"currency":"UAH","percent_change_24h":"0","provider":"fixer","rate":26.838273,"timestamp":1589030105},{"currency":"UGX","percent_change_24h":"0","provider":"fixer","rate":3800.322804,"timestamp":1589030105},{"currency":"USD","percent_change_24h":"0","provider":"fixer","rate":1,"timestamp":1589030105},{"currency":"UYU","percent_change_24h":"0","provider":"fixer","rate":43.139569,"timestamp":1589030105},{"currency":"UZS","percent_change_24h":"0","provider":"fixer","rate":10109.000335,"timestamp":1589030105},{"currency":"VEF","percent_change_24h":"0","provider":"fixer","rate":9.987504,"timestamp":1589030105},{"currency":"VND","percent_change_24h":"0","provider":"fixer","rate":23400.5,"timestamp":1589030105},{"currency":"VUV","percent_change_24h":"0","provider":"fixer","rate":119.848296,"timestamp":1589030105},{"currency":"WST","percent_change_24h":"0","provider":"fixer","rate":2.78472,"timestamp":1589030105},{"currency":"XAF","percent_change_24h":"0","provider":"fixer","rate":604.99802,"timestamp":1589030105},{"currency":"XAG","percent_change_24h":"0","provider":"fixer","rate":0.06474,"timestamp":1589030105},{"currency":"XAU","percent_change_24h":"0","provider":"fixer","rate":0.000587,"timestamp":1589030105},{"currency":"XCD","percent_change_24h":"0","provider":"fixer","rate":2.70255,"timestamp":1589030105},{"currency":"XDR","percent_change_24h":"0","provider":"fixer","rate":0.734807,"timestamp":1589030105},{"currency":"XOF","percent_change_24h":"0","provider":"fixer","rate":605.000332,"timestamp":1589030105},{"currency":"XPF","percent_change_24h":"0","provider":"fixer","rate":110.4036,"timestamp":1589030105},{"currency":"YER","percent_change_24h":"0","provider":"fixer","rate":250.375037,"timestamp":1589030105},{"currency":"ZAR","percent_change_24h":"0","provider":"fixer","rate":18.350904,"timestamp":1589030105},{"currency":"ZMK","percent_change_24h":"0","provider":"fixer","rate":9001.203593,"timestamp":1589030105},{"currency":"ZMW","percent_change_24h":"0","provider":"fixer","rate":18.576435,"timestamp":1589030105},{"currency":"ZWL","percent_change_24h":"0","provider":"fixer","rate":322.000001,"timestamp":1589030105}]`
)
