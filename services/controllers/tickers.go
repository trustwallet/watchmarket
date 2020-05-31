package controllers

import (
	"context"
	"errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
	"sync"
	"time"
)

func (c Controller) HandleTickersRequest(tr TickerRequest, ctx context.Context) (TickerResponse, error) {
	if tr.Assets == nil {
		return TickerResponse{}, errors.New(ErrBadRequest)
	}

	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency), ctx)
	if err != nil {
		return TickerResponse{}, errors.New(ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueries(tr.Assets), ctx)
	if err != nil {
		return TickerResponse{}, errors.New(ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate, ctx)

	return createResponse(tr, tickers), nil
}

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
		return watchmarket.Rate{}, errors.New(ErrNotFound)
	}

	return watchmarket.Rate{
		Currency:         result.Currency,
		PercentChange24h: result.PercentChange24h,
		Provider:         result.Provider,
		Rate:             result.Rate,
		Timestamp:        result.LastUpdated.Unix(),
	}, nil
}

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery, ctx context.Context) (watchmarket.Tickers, error) {
	dbTickers, err := c.database.GetTickersByQueries(tickerQueries, ctx)
	if err != nil {
		logger.Error(err, "getTickersByPriority")
		return nil, err
	}
	providers := c.tickersPriority

	res := new(sortedTickersResponse)
	wg := new(sync.WaitGroup)
	for _, q := range tickerQueries {
		wg.Add(1)
		go findBestProviderForQuery(q.Coin, q.TokenId, dbTickers, providers, wg, res, c.configuration)
	}

	wg.Wait()

	sortedTickers := res.tickers

	result := make(watchmarket.Tickers, len(sortedTickers))

	for i, sr := range sortedTickers {
		result[i] = watchmarket.Ticker{
			Coin:       sr.Coin,
			CoinName:   sr.CoinName,
			CoinType:   watchmarket.CoinType(sr.CoinType),
			LastUpdate: sr.LastUpdated,
			Price: watchmarket.Price{
				Change24h: sr.Change24h,
				Currency:  sr.Currency,
				Provider:  sr.Provider,
				Value:     sr.Value,
			},
			TokenId: sr.TokenId,
		}
	}

	return result, nil
}

func (c Controller) normalizeTickers(tickers watchmarket.Tickers, rate watchmarket.Rate, ctx context.Context) watchmarket.Tickers {
	result := make(watchmarket.Tickers, 0, len(tickers))
	for _, t := range tickers {
		r, ok := c.convertRateToDefaultCurrency(t, rate, ctx)
		if !ok {
			continue
		}
		if !c.isSuitableUpdateTime(t) {
			continue
		}
		result = append(result, applyRateToTicker(t, r))
	}
	return result
}

func (c Controller) convertRateToDefaultCurrency(t watchmarket.Ticker, rate watchmarket.Rate, ctx context.Context) (watchmarket.Rate, bool) {
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

func (c Controller) isSuitableUpdateTime(ticker watchmarket.Ticker) bool {
	if ticker.Price.Provider == "coinmarketcap" {
		return true
	}
	now := time.Now().Unix()
	last := ticker.LastUpdate.Unix()
	if now < last {
		return true
	}
	diff := now - last
	if diff < 0 {
		return true
	}
	respectableTime := watchmarket.DurationToUnix(c.configuration.RestAPI.Tickers.RespectableUpdateTime)
	return uint(diff) <= respectableTime
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

func createResponse(tr TickerRequest, tickers watchmarket.Tickers) TickerResponse {
	mergedTickers := make(watchmarket.Tickers, 0, len(tickers))
	for _, t := range tickers {
		newTicker, ok := foundTickerInAssets(tr.Assets, t)
		if !ok {
			continue
		}
		mergedTickers = append(mergedTickers, newTicker)
	}

	return TickerResponse{tr.Currency, mergedTickers}
}

func foundTickerInAssets(assets []Coin, t watchmarket.Ticker) (watchmarket.Ticker, bool) {
	for _, c := range assets {
		if c.Coin == t.Coin && strings.ToLower(c.TokenId) == t.TokenId {
			t.TokenId = c.TokenId
			return t, true
		}
	}
	return watchmarket.Ticker{}, false
}

func findBestProviderForQuery(coin uint, token string, sliceToFind []models.Ticker, providers []string, wg *sync.WaitGroup, res *sortedTickersResponse, configuration config.Configuration) {
	defer wg.Done()

	for _, p := range providers {
		for _, t := range sliceToFind {
			baseCheck := coin == t.Coin && strings.ToLower(token) == t.TokenId

			if baseCheck && t.ShowOption == models.AlwaysShow {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				return
			}
			if baseCheck && p == t.Provider && t.ShowOption != models.NeverShow &&
				isRespectableMarketCap(t.Provider, t.MarketCap, configuration) && isRespectableVolume(t.Provider, t.Volume, configuration) {
				res.Lock()
				res.tickers = append(res.tickers, t)
				res.Unlock()
				return
			}
		}
	}
}

func isRespectableMarketCap(provider string, marketCap float64, configuration config.Configuration) bool {
	if provider != "coingecko" {
		return true
	}
	return marketCap >= configuration.RestAPI.Tickers.RespsectableMarketCap
}

func isRespectableVolume(provider string, volume float64, configuration config.Configuration) bool {
	if provider != "coingecko" {
		return true
	}
	return volume >= configuration.RestAPI.Tickers.RespsectableVolume
}

func makeTickerQueries(coins []Coin) []models.TickerQuery {
	tickerQueries := make([]models.TickerQuery, 0, len(coins))
	for _, c := range coins {
		tickerQueries = append(tickerQueries, models.TickerQuery{
			Coin:    c.Coin,
			TokenId: strings.ToLower(c.TokenId),
		})
	}
	return tickerQueries
}

func isFiatRate(currency string) bool {
	switch currency {
	case "AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN", "BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL", "BSD", "BTN", "BWP", "BYN", "BYR", "BZD", "CAD", "CDF", "CHF", "CLF", "CLP", "CNY", "COP", "CRC", "CUC", "CUP", "CVE", "CZK", "DJF", "DKK", "DOP", "DZD", "EGP", "ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL", "GGP", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD", "HKD", "HNL", "HRK", "HTG", "HUF", "IDR", "ILS", "IMP", "INR", "IQD", "IRR", "ISK", "JEP", "JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LVL", "LYD", "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRO", "MUR", "MVR", "MWK", "MXN", "MYR", "MZN", "NAD", "NGN", "NIO", "NOK", "NPR", "NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG", "QAR", "RON", "RSD", "RUB", "RWF", "SAR", "SBD", "SCR", "SDG", "SEK", "SGD", "SHP", "SLL", "SOS", "SRD", "STD", "SVC", "SYP", "SZL", "THB", "TJS", "TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD", "UYU", "UZS", "VEF", "VND", "VUV", "WST", "XAF", "XAG", "XAU", "XCD", "XDR", "XOF", "XPF", "YER", "ZAR", "ZMK", "ZMW", "ZWL":
		return true
	default:
	}
	return false
}
