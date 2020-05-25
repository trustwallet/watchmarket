package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"testing"
)

func TestController_HandleDetailsRequest(t *testing.T) {
	cm := getChartsMock()
	wantedD := watchmarket.CoinDetails{
		Provider:          "coinmarketcap",
		Vol24:             1,
		MarketCap:         2,
		CirculatingSupply: 3,
		TotalSupply:       4,
		Info: watchmarket.Info{
			Name:             "2",
			Website:          "2",
			SourceCode:       "2",
			WhitePaper:       "2",
			Description:      "2",
			ShortDescription: "2",
			Explorer:         "2",
			Socials:          nil,
		},
	}
	cm.wantedDetails = wantedD
	c := setupController(t, getDbMock(), getCacheMock(), cm)
	assert.NotNil(t, c)
	details, err := c.HandleDetailsRequest(DetailsRequest{
		CoinQuery: "0",
		Token:     "2",
		Currency:  "3",
	}, context.Background())
	assert.Nil(t, err)
	assert.Equal(t, wantedD, details)
}
