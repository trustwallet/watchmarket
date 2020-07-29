package tickerscontroller

import (
	"context"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

func (c Controller) getRateByPriority(currency string, ctx context.Context) (watchmarket.Rate, error) {
	rates, err := c.database.GetRates(currency, ctx)
	if err != nil {
		logger.Error(err, "getRateByPriority")
		return watchmarket.Rate{}, err
	}

	providers := c.ratesPriority

	var result models.Rate
ProvidersLoop:
	for _, p := range providers {
		for _, r := range rates {
			if p == r.Provider {
				result = r
				break ProvidersLoop
			}
		}
	}
	emptyRate := models.Rate{}
	if result == emptyRate || (watchmarket.IsFiatRate(result.Currency) && result.Provider != "fixer") {
		return watchmarket.Rate{}, errors.New(watchmarket.ErrNotFound)
	}

	return watchmarket.Rate{
		Currency:         result.Currency,
		PercentChange24h: result.PercentChange24h,
		Provider:         result.Provider,
		Rate:             result.Rate,
		Timestamp:        result.LastUpdated.Unix(),
	}, nil
}

func (c Controller) rateToDefaultCurrency(t watchmarket.Ticker, rate watchmarket.Rate, ctx context.Context) (watchmarket.Rate, bool) {
	if t.Price.Currency != watchmarket.DefaultCurrency {
		newRate, err := c.getRateByPriority(strings.ToUpper(t.Price.Currency), ctx)
		if err != nil {
			return watchmarket.Rate{}, false
		}
		rate.Rate /= newRate.Rate
		rate.PercentChange24h = newRate.PercentChange24h
	}
	return rate, true
}

func applyRateToTicker(t watchmarket.Ticker, rate watchmarket.Rate) watchmarket.Ticker {
	if t.Price.Currency == rate.Currency {
		return t
	}
	t.Price.Value *= 1 / rate.Rate
	t.Price.Currency = rate.Currency
	t.Volume *= 1 / rate.Rate
	t.MarketCap *= 1 / rate.Rate

	if rate.PercentChange24h != 0 {
		t.Price.Change24h -= rate.PercentChange24h // Look at it more detailed
	}
	return t
}
