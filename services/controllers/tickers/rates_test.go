package tickerscontroller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func TestController_getRateByPriority(t *testing.T) {
	now := time.Now()
	rate := models.Rate{
		Currency:         "USD",
		PercentChange24h: 1,
		Provider:         watchmarket.CoinMarketCap,
		Rate:             1,
		LastUpdated:      now,
	}
	rate2 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 2,
		Provider:         watchmarket.CoinGecko,
		Rate:             2,
		LastUpdated:      now,
	}
	rate3 := models.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         watchmarket.Fixer,
		Rate:             6,
		LastUpdated:      now,
	}

	db := getDbMock()

	db.WantedTickersError = nil
	db.WantedRatesError = nil
	db.WantedRates = []models.Rate{rate, rate2, rate3}

	c := setupController(t, db, false)
	assert.NotNil(t, c)

	r, err := c.getRateByPriority("USD")
	assert.Nil(t, err)

	assert.Equal(t, watchmarket.Rate{
		Currency:         "USD",
		PercentChange24h: 4,
		Provider:         watchmarket.Fixer,
		Rate:             6,
		Timestamp:        now.Unix(),
	}, r)
}
