package chart

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrices_IsEmpty(t *testing.T) {
	price := Price{
		Price: 0,
		Date:  0,
	}

	prices := make([]Price, 0)
	prices = append(prices, price)
	prices = append(prices, price)
	prices = append(prices, price)
	prices = append(prices, price)

	notEmptyChartData := Data{
		Prices: prices,
		Error:  "",
	}

	emptyChartData := Data{
		Prices: nil,
		Error:  "",
	}

	assert.True(t, !notEmptyChartData.IsEmpty())
	assert.True(t, emptyChartData.IsEmpty())
}
