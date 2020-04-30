package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrices_IsEmpty(t *testing.T) {
	price := ChartPrice{
		Price: 0,
		Date:  0,
	}

	prices := make([]ChartPrice, 0)
	prices = append(prices, price)
	prices = append(prices, price)
	prices = append(prices, price)
	prices = append(prices, price)

	notEmptyChartData := ChartData{
		Prices: prices,
		Error:  "",
	}

	emptyChartData := ChartData{
		Prices: nil,
		Error:  "",
	}

	assert.True(t, !notEmptyChartData.IsEmpty())
	assert.True(t, emptyChartData.IsEmpty())
}
