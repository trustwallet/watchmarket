package postgres

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/models"
	"testing"
)

func Test_normalizeRates(t *testing.T) {
	rates := []models.Rate{
		{
			Rate:             10,
			Currency:         "10",
			Provider:         "10",
			PercentChange24h: 10,
		},
		{
			Rate:             20,
			Currency:         "20",
			Provider:         "20",
			PercentChange24h: 20,
		},
		{
			Rate:             10,
			Currency:         "10",
			Provider:         "10",
			PercentChange24h: 10,
		},
		{
			Rate:             20,
			Currency:         "20",
			Provider:         "20",
			PercentChange24h: 20,
		},
	}
	result := normalizeRates(rates)
	assert.Len(t, result, 2)
	assert.NotEqual(t, result[0], result[1])

	rates = []models.Rate{
		{
			Rate:             10,
			Currency:         "10",
			Provider:         "10",
			PercentChange24h: 10,
		},
		{
			Rate:             100,
			Currency:         "10",
			Provider:         "10",
			PercentChange24h: 100,
		},
	}

	result = normalizeRates(rates)
	assert.Len(t, result, 1)
	assert.Equal(t, float64(10), result[0].Rate)
	assert.Equal(t, float64(10), result[0].PercentChange24h)
}
