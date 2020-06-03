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
	if result == emptyRate || (isFiatRate(result.Currency) && result.Provider != "fixer") {
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

func isFiatRate(currency string) bool {
	switch currency {
	case "AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN", "BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BYR", "BZD", "CAD", "CDF", "CHF", "CLF", "CLP", "CNY", "COP", "CRC", "CUC", "CUP", "CVE", "CZK", "DJF", "DKK", "DOP", "DZD", "EGP", "ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL", "GGP", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD", "HKD", "HNL", "HRK", "HTG", "HUF", "IDR", "ILS", "IMP", "INR", "IQD", "IRR", "ISK", "JEP", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LVL", "LYD", "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRO", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN", "NAD", "NGN", "NIO", "NOK", "NPR", "NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG", "QAR", "RON", "RSD", "RUB", "RWF", "SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLL", "SOS", "SRD", "STD", "SVC", "SYP", "SZL", "THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD", "UYU", "UZS", "VEF", "VND", "VUV", "WST", "XAF", "XAG", "XAU", "XCD", "XDR", "XOF", "XPF", "YER", "ZAR", "ZMK", "ZMW", "ZWL":
		return true
	default:
	}
	return false
}
