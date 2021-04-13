package tickerscontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
)

func TestController_getTickersByPriority(t *testing.T) {
	ticker60ACMC := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinMarketCap,
		Value:     100,
	}

	ticker60ACG := models.Ticker{
		Coin:      60,
		CoinName:  "ETH",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinGecko,
		Value:     100,
	}

	ticker714ACG := models.Ticker{
		Coin:      714,
		CoinName:  "BNB",
		TokenId:   "a",
		Change24h: 10,
		Currency:  "USD",
		Provider:  watchmarket.CoinGecko,
		Value:     100,
	}

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG, ticker714ACG}
	c := setupController(t, db, false)
	assert.NotNil(t, c)

	tickers, err := c.getTickersByPriority([]controllers.Asset{{CoinId: 60, TokenId: "A"}, {CoinId: 714, TokenId: "A"}})
	assert.Nil(t, err)
	assert.NotNil(t, tickers)
	assert.Equal(t, 2, len(tickers))

	wantedTicker1 := watchmarket.Ticker{
		Coin:     60,
		CoinName: "ETH",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  watchmarket.CoinMarketCap,
			Value:     100,
		},
		TokenId: "A",
	}
	wantedTicker2 := watchmarket.Ticker{
		Coin:     714,
		CoinName: "BNB",
		CoinType: "",
		Price: watchmarket.Price{
			Change24h: 10,
			Currency:  "USD",
			Provider:  watchmarket.CoinGecko,
			Value:     100,
		},
		TokenId: "a",
	}
	var counter int
	for _, t := range tickers {
		if t == wantedTicker1 || t == wantedTicker2 {
			counter++
		}
	}
	assert.Equal(t, 1, counter)
	db2 := getDbMock()
	db2.WantedTickers = []models.Ticker{ticker60ACMC, ticker60ACG}
	c2 := setupController(t, db2, false)
	tickers2, err := c2.getTickersByPriority([]controllers.Asset{{CoinId: 60, TokenId: "A"}})
	assert.Nil(t, err)
	assert.NotNil(t, tickers2)
	assert.Equal(t, 1, len(tickers2))
	assert.Equal(t, wantedTicker1, tickers2[0])
}
