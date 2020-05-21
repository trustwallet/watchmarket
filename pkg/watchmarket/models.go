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
	}

	Rates []Rate

	CoinType string

	Response struct {
		Currency string  `json:"currency"`
		Docs     Tickers `json:"docs"`
	}

	Ticker struct {
		Coin       uint      `json:"coin"`
		CoinName   string    `json:"coin_name"`
		CoinType   CoinType  `json:"type"`
		Error      string    `json:"error,omitempty"`
		LastUpdate time.Time `json:"last_update"`
		Price      Price     `json:"price"`
		TokenId    string    `json:"token_id"`
		Volume     float64   `json:"_"`
		MarketCap  float64   `json:"_"`
	}

	Price struct {
		Change24h float64 `json:"change_24h"`
		Currency  string  `json:"currency"`
		Provider  string  `json:"provider"`
		Value     float64 `json:"value"`
	}

	Tickers []Ticker

	Chart struct {
		Prices []ChartPrice `json:"prices,omitempty"`
		Error  string       `json:"error,omitempty"`
	}

	ChartPrice struct {
		Price float64 `json:"price"`
		Date  int64   `json:"date"`
	}

	CoinDetails struct {
		Vol24             float64 `json:"volume_24"`
		MarketCap         float64 `json:"market_cap"`
		CirculatingSupply float64 `json:"circulating_supply"`
		TotalSupply       float64 `json:"total_supply"`
		Info              Info    `json:"assets,omitempty"`
	}

	Info struct {
		Name             string       `json:"name,omitempty"`
		Website          string       `json:"website,omitempty"`
		SourceCode       string       `json:"source_code,omitempty"`
		WhitePaper       string       `json:"white_paper,omitempty"`
		Description      string       `json:"description,omitempty"`
		ShortDescription string       `json:"short_description,omitempty"`
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
	DefaultPrecision              = 10
	Coin                 CoinType = "coin"
	Token                CoinType = "token"
	DefaultCurrency               = "USD"
	DefaultMaxChartItems          = 64
)

func (d Chart) IsEmpty() bool {
	return len(d.Prices) == 0
}

func (i CoinDetails) IsEmpty() bool {
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
