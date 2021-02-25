package tickerscontroller

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (c Controller) getRateByPriority(currency string) (watchmarket.Rate, error) {
	if c.configuration.RestAPI.UseMemoryCache {
		rawResult, err := c.cache.Get(currency)
		if err != nil {
			return watchmarket.Rate{}, err
		}
		var result watchmarket.Rate
		if err = json.Unmarshal(rawResult, &result); err != nil {
			return watchmarket.Rate{}, err
		}
		return result, nil
	}

	rates, err := c.database.GetRates(currency)
	if err != nil {
		log.Error(err, "getRateByPriority")
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

func (c Controller) rateToDefaultCurrency(t watchmarket.Ticker, rate watchmarket.Rate) (watchmarket.Rate, bool) {
	if t.Price.Currency != watchmarket.DefaultCurrency {
		newRate, err := c.getRateByPriority(strings.ToUpper(t.Price.Currency))
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
