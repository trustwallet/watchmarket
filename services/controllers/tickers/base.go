package tickerscontroller

import (
	"encoding/json"
	"errors"
	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
	"strings"
)

type Controller struct {
	cache           cache.Provider
	ratesPriority   []string
	tickersPriority []string
	configuration   config.Configuration
}

func NewController(
	database db.Instance,
	cache cache.Provider,
	ratesPriority, tickersPriority []string,
	configuration config.Configuration,
) Controller {
	return Controller{
		cache,
		ratesPriority,
		tickersPriority,
		configuration,
	}
}

func (c Controller) HandleTickersRequest(request controllers.TickerRequest) (watchmarket.Tickers, error) {
	rate, err := c.getCacheRate(strings.ToUpper(request.Currency))
	if err != nil {
		return watchmarket.Tickers{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getCacheTickers(request.Assets)
	if err != nil {
		return watchmarket.Tickers{}, errors.New(watchmarket.ErrInternal)
	}
	tickers = c.filterTickers(tickers, rate)
	return tickers, nil
}

func (c Controller) filterTickers(tickers watchmarket.Tickers, rate watchmarket.Rate) (result watchmarket.Tickers) {
	for _, ticker := range tickers {
		r, ok := c.rateToDefaultCurrency(ticker, rate)
		if !ok {
			continue
		}
		if !watchmarket.IsSuitableUpdateTime(ticker.LastUpdate, c.configuration.RestAPI.Tickers.RespectableUpdateTime) {
			continue
		}
		c.applyRateToTicker(&ticker, r)
		result = append(result, ticker)
	}
	return result
}

func (c Controller) rateToDefaultCurrency(t watchmarket.Ticker, rate watchmarket.Rate) (watchmarket.Rate, bool) {
	if t.Price.Currency != watchmarket.DefaultCurrency {
		newRate, err := c.getCacheRate(strings.ToUpper(t.Price.Currency))
		if err != nil {
			return watchmarket.Rate{}, false
		}
		rate.Rate /= newRate.Rate
		rate.PercentChange24h = newRate.PercentChange24h
	}
	return rate, true
}

func (c Controller) applyRateToTicker(ticker *watchmarket.Ticker, rate watchmarket.Rate) {
	if ticker.Price.Currency == rate.Currency {
		return
	}
	ticker.Price.Value *= 1 / rate.Rate
	ticker.Price.Currency = rate.Currency
	ticker.Volume *= 1 / rate.Rate
	ticker.MarketCap *= 1 / rate.Rate

	if rate.PercentChange24h != 0 {
		ticker.Price.Change24h -= rate.PercentChange24h // Look at it more detailed
	}
}

func (c Controller) getCacheTickers(assets []controllers.Asset) (watchmarket.Tickers, error) {
	if c.cache == nil {
		return watchmarket.Tickers{}, errors.New(watchmarket.ErrInternal)
	}
	var results watchmarket.Tickers
	for _, assetData := range assets {
		key := strings.ToLower(asset.BuildID(assetData.CoinId, assetData.TokenId))
		rawResult, err := c.cache.Get(key)
		if err != nil {
			continue
		}
		var result watchmarket.Ticker
		if err = json.Unmarshal(rawResult, &result); err != nil {
			continue
		}
		result.TokenId = assetData.TokenId
		results = append(results, result)
	}
	return results, nil
}

func (c Controller) getCacheRate(currency string) (result watchmarket.Rate, err error) {
	if c.cache == nil {
		return watchmarket.Rate{}, errors.New(watchmarket.ErrInternal)
	}
	rawResult, err := c.cache.Get(currency)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(rawResult, &result); err != nil {
		return result, err
	}
	return result, nil
}
