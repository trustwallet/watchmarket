package rates

import "math/big"

type Rate struct {
	Currency         string    `json:"currency"`
	Rate             float64   `json:"rate"`
	Timestamp        int64     `json:"timestamp"`
	PercentChange24h big.Float `json:"percent_change_24h,omitempty"`
	Provider         string    `json:"provider,omitempty"`
}

type Rates []Rate
