package tickerscontroller

import (
	"encoding/json"
	"strings"

	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

func (c Controller) getRateByPriority(currency string) (watchmarket.Rate, error) {
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
