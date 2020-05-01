package cmc

import (
	"fmt"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/watchmarket/services/tickers/models"
	"time"
)

type (
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
	CmcMapping  map[uint][]CoinMap

	CoinResult struct {
		Id       uint
		Coin     coin.Coin
		TokenId  string
		CoinType models.CoinType
	}
)

func (c *CmcSlice) coinToCmcMap() (m CoinMapping) {
	m = make(map[string]CoinMap)
	for _, cm := range *c {
		m[generateId(cm.Coin, cm.TokenId)] = cm
	}
	return
}

func (c *CmcSlice) cmcToCoinMap() (m CmcMapping) {
	m = make(map[uint][]CoinMap)
	for _, cm := range *c {
		_, ok := m[cm.Id]
		if !ok {
			m[cm.Id] = make([]CoinMap, 0)
		}
		m[cm.Id] = append(m[cm.Id], cm)
	}
	return
}

func (cm CmcMapping) GetCoins(coinId uint) ([]CoinResult, error) {
	cmcCoin, ok := cm[coinId]
	if !ok {
		return nil, errors.E("CmcMapping.getCoin: coinId notFound")
	}
	tokens := make([]CoinResult, 0)
	for _, cc := range cmcCoin {
		c, ok := coin.Coins[cc.Coin]
		if !ok {
			continue
		}
		tokens = append(tokens, CoinResult{Coin: c, Id: cc.Id, TokenId: cc.TokenId, CoinType: models.CoinType(cc.Type)})
	}
	return tokens, nil
}

func (cm CoinMapping) GetCoinByContract(coinId uint, contract string) (c CoinMap, err error) {
	c, ok := cm[generateId(coinId, contract)]
	if !ok {
		err = errors.E("No coin found", errors.Params{"coin": coinId, "token": contract})
	}

	return
}

func generateId(id uint, token string) string {
	return fmt.Sprintf("%d:%s", id, token)
}
