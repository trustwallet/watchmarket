package watchmarket

import (
	"math"
	"time"
)

type (
	Rate struct {
		Currency         string  `json:"currency"`
		PercentChange24h float64 `json:"percent_change_24h,omitempty"`
		Provider         string  `json:"provider,omitempty"`
		Rate             float64 `json:"rate"`
		Timestamp        int64   `json:"timestamp"`
		ShowOption       int     `json:"-"`
	}

	Rates []Rate

	CoinType string

	Ticker struct {
		Coin              uint      `json:"coin"`
		CoinName          string    `json:"coin_name"`
		TokenId           string    `json:"token_id,omitempty"`
		CoinType          CoinType  `json:"type,omitempty"`
		Price             Price     `json:"price,omitempty"`
		LastUpdate        time.Time `json:"last_update"`
		Error             string    `json:"error,omitempty"`
		Volume            float64   `json:"volume"`
		MarketCap         float64   `json:"market_cap"`
		CirculatingSupply float64   `json:"-"`
		TotalSupply       float64   `json:"-"`
		ShowOption        int       `json:"-"`
	}

	Price struct {
		Change24h float64 `json:"change_24h"`
		Currency  string  `json:"currency"`
		Provider  string  `json:"provider,omitempty"`
		Value     float64 `json:"value"`
	}

	Tickers []Ticker

	Chart struct {
		Provider string       `json:"provider,omitempty"`
		Prices   []ChartPrice `json:"prices,omitempty"`
		Error    string       `json:"error,omitempty"`
	}

	ChartPrice struct {
		Price float64 `json:"price"`
		Date  int64   `json:"date"`
	}

	CoinDetails struct {
		Provider    string `json:"provider,omitempty"`
		ProviderURL string `json:"provider_url,omitempty"`
		Info        *Info  `json:"info,omitempty"`
	}

	Info struct {
		Name             string       `json:"name,omitempty"`
		Website          string       `json:"website,omitempty"`
		SourceCode       string       `json:"source_code,omitempty"`
		WhitePaper       string       `json:"white_paper,omitempty"`
		Description      string       `json:"description,omitempty"`
		ShortDescription string       `json:"short_description,omitempty"`
		Research         string       `json:"research,omitempty"`
		Explorer         string       `json:"explorer,omitempty"`
		Socials          []SocialLink `json:"socials,omitempty"`
	}

	SocialLink struct {
		Name   string `json:"name"`
		Url    string `json:"url"`
		Handle string `json:"handle"`
	}
)

const (
	UnknownCoinID                 = 111111
	DefaultPrecision              = 10
	Coin                 CoinType = "coin"
	Token                CoinType = "token"
	DefaultCurrency               = "USD"
	DefaultMaxChartItems          = 64

	ErrNotFound   = "not found"
	ErrBadRequest = "bad request"
	ErrInternal   = "internal"

	CoinMarketCap = "coinmarketcap"
	CoinGecko     = "coingecko"
	Fixer         = "fixer"
)

func (d Chart) IsEmpty() bool {
	return len(d.Prices) == 0
}

func (i CoinDetails) IsEmpty() bool {
	if i.Info == nil {
		return true
	}
	return i.Info.Name == ""
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func TruncateWithPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UnixToDuration(unixTime uint) time.Duration {
	return time.Duration(unixTime * 1000000000)
}

func DurationToUnix(duration time.Duration) uint {
	return uint(duration.Seconds())
}

func IsRespectableValue(value, respValue float64) bool {
	return value >= respValue
}

func IsSuitableUpdateTime(LastUpdate time.Time, maxDuration time.Duration) bool {
	now := time.Now().Unix()
	last := LastUpdate.Unix()
	if now < last {
		return true
	}
	diff := now - last
	if diff < 0 {
		return true
	}
	respectableTime := DurationToUnix(maxDuration)
	return uint(diff) <= respectableTime
}

func IsFiatRate(currency string) bool {
	switch currency {
	case "AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN", "BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BYR", "BZD", "CAD", "CDF", "CHF", "CLF", "CLP", "CNY", "COP", "CRC", "CUC", "CUP", "CVE", "CZK", "DJF", "DKK", "DOP", "DZD", "EGP", "ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL", "GGP", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD", "HKD", "HNL", "HRK", "HTG", "HUF", "IDR", "ILS", "IMP", "INR", "IQD", "IRR", "ISK", "JEP", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LVL", "LYD", "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRO", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN", "NAD", "NGN", "NIO", "NOK", "NPR", "NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG", "QAR", "RON", "RSD", "RUB", "RWF", "SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLL", "SOS", "SRD", "STD", "SVC", "SYP", "SZL", "THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD", "UYU", "UZS", "VEF", "VND", "VUV", "WST", "XAF", "XAG", "XAU", "XCD", "XDR", "XOF", "XPF", "YER", "ZAR", "ZMK", "ZMW", "ZWL":
		return true
	default:
		return false
	}
}
