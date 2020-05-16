package controllers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController_HandleTickersRequest(t *testing.T) {
	c := setupController(t)
	assert.NotNil(t, c)

	tickers, err := c.HandleTickersRequest(TickerRequest{
		Currency: "USD",
		Assets:   []Coin{{Coin: 60, TokenId: "A"}, {Coin: 714, TokenId: "A"}},
	})
	assert.Nil(t, err)
	assert.NotNil(t, tickers)
}
