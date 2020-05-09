package coinmarketcap

import (
	"time"
)

type (
	Charts struct {
		Data ChartQuotes `json:"data"`
	}

	ChartQuotes map[string]ChartQuoteValues

	ChartQuoteValues map[string][]float64

	ChartInfo struct {
		Data ChartInfoData `json:"data"`
	}

	ChartInfoData struct {
		Rank              uint32                    `json:"rank"`
		CirculatingSupply float64                   `json:"circulating_supply"`
		TotalSupply       float64                   `json:"total_supply"`
		Quotes            map[string]ChartInfoQuote `json:"quotes"`
	}

	ChartInfoQuote struct {
		Price     float64 `json:"price"`
		Volume24  float64 `json:"volume_24h"`
		MarketCap float64 `json:"market_cap"`
	}

	CoinPrices struct {
		Status struct {
			Timestamp    time.Time   `json:"timestamp"`
			ErrorCode    int         `json:"error_code"`
			ErrorMessage interface{} `json:"error_message"`
		} `json:"status"`
		Data []Data `json:"data"`
	}

	Coin struct {
		Id     uint   `json:"id"`
		Symbol string `json:"symbol"`
	}

	Data struct {
		Coin
		LastUpdated time.Time `json:"last_updated"`
		Platform    *Platform `json:"platform"`
		Quote       Quote     `json:"quote"`
	}

	Platform struct {
		Coin
		TokenAddress string `json:"token_address"`
	}

	Quote struct {
		USD USD `json:"USD"`
	}

	USD struct {
		Price            float64 `json:"price"`
		PercentChange24h float64 `json:"percent_change_24h"`
	}

	CoinMap struct {
		Coin    uint   `json:"coin"`
		Id      uint   `json:"id"`
		Type    string `json:"type"`
		TokenId string `json:"token_id"`
	}

	CmcSlice    []CoinMap
	CoinMapping map[string]CoinMap
)
