package coinmarketcap

import (
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
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
		LastUpdated       time.Time `json:"last_updated"`
		CirculatingSupply float64   `json:"circulating_supply"`
		TotalSupply       float64   `json:"total_supply"`
		Platform          Platform  `json:"platform"`
		Quote             Quote     `json:"quote"`
	}

	Platform struct {
		Coin
		TokenAddress string `json:"token_address"`
	}

	Quote struct {
		USD USD `json:"USD"`
	}

	USD struct {
		Volume           float64 `json:"volume_24h"`
		MarketCap        float64 `json:"market_cap"`
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
		CoinType watchmarket.CoinType
	}

	Charts struct {
		Data ChartQuotes `json:"data"`
	}

	ChartQuotes map[string]ChartQuoteValues

	ChartQuoteValues map[string][]float64

	ChartInfo struct {
		Data map[int]ChartInfoData `json:"data"`
	}

	ChartInfoData struct {
		Rank              uint32                    `json:"cmc_rank"`
		CirculatingSupply float64                   `json:"circulating_supply"`
		TotalSupply       float64                   `json:"total_supply"`
		Slug              string                    `json:"slug"`
		Quotes            map[string]ChartInfoQuote `json:"quote"`
	}

	ChartInfoQuote struct {
		Price     float64 `json:"price"`
		Volume24  float64 `json:"volume_24h"`
		MarketCap float64 `json:"market_cap"`
	}

	CmcSlice []CoinMap

	CoinMapping map[string]CoinMap
)
