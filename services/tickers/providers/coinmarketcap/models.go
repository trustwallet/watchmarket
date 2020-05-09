package coinmarketcap

import (
	"github.com/trustwallet/blockatlas/coin"
	ticker "github.com/trustwallet/watchmarket/services/tickers"
	"time"
)

type (
	CoinPrices struct {
		Status `json:"status"`
		Data   []Data `json:"data"`
	}

	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
	}

	Coin struct {
		Id     uint   `json:"id"`
		Symbol string `json:"symbol"`
	}

	Data struct {
		Coin
		LastUpdated time.Time `json:"last_updated"`
		Platform    Platform  `json:"platform"`
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

	CoinResult struct {
		Id       uint
		Coin     coin.Coin
		TokenId  string
		CoinType ticker.CoinType
	}
)
