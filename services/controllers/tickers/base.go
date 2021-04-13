package tickerscontroller

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/cache"
	"github.com/trustwallet/watchmarket/services/controllers"
)

type Controller struct {
	database        db.Instance
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
		database,
		cache,
		ratesPriority,
		tickersPriority,
		configuration,
	}
}

func (c Controller) HandleTickersRequest(request controllers.TickerRequest) (watchmarket.Tickers, error) {
	rate, err := c.getRateByPriority(strings.ToUpper(request.Currency))
	if err != nil {
		return watchmarket.Tickers{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(request.Assets)
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

func (c Controller) getTickersByPriority(assets []controllers.Asset) (watchmarket.Tickers, error) {
	if result, err := c.getCachedTickers(assets); err == nil {
		return result, nil
	}
	tickers, err := c.database.GetTickers(assets)
	if err != nil {
		return nil, err
	}
	result := make(watchmarket.Tickers, 0)
	for _, assetData := range assets {
		ticker := findBestTicker(assetData, tickers, c.tickersPriority, c.configuration)
		if ticker == nil {
			continue
		}
		result = append(result, watchmarket.Ticker{
			Coin:       ticker.Coin,
			CoinName:   ticker.CoinName,
			CoinType:   watchmarket.CoinType(ticker.CoinType),
			LastUpdate: ticker.LastUpdated,
			Price: watchmarket.Price{
				Change24h: ticker.Change24h,
				Currency:  ticker.Currency,
				Provider:  ticker.Provider,
				Value:     ticker.Value,
			},
			TokenId: assetData.TokenId,
		})
	}

	return result, nil
}

func findBestTicker(assetData controllers.Asset, tickers []models.Ticker, providers []string, configuration config.Configuration) *models.Ticker {
	for _, p := range providers {
		for _, ticker := range tickers {
			baseCheck := assetData.CoinId == ticker.Coin && strings.ToLower(assetData.TokenId) == ticker.TokenId

			if baseCheck && ticker.ShowOption == models.AlwaysShow {
				return &ticker
			}

			if baseCheck && p == ticker.Provider && ticker.ShowOption != models.NeverShow &&
				(watchmarket.IsRespectableValue(ticker.MarketCap, configuration.RestAPI.Tickers.RespsectableMarketCap) || ticker.Provider != "coingecko") &&
				(watchmarket.IsRespectableValue(ticker.Volume, configuration.RestAPI.Tickers.RespsectableVolume) || ticker.Provider != "coingecko") {
				return &ticker
			}
		}
	}
	return nil
}

func (c Controller) getRateByPriority(currency string) (result watchmarket.Rate, err error) {
	if result, err := c.getCachedRate(currency); err == nil {
		return result, nil
	}

	rates, err := c.database.GetRates(currency)
	if err != nil {
		return result, err
	}
	isFiat := !watchmarket.IsFiatRate(currency)

	for _, p := range c.ratesPriority {
		if isFiat && p != watchmarket.Fixer {
			continue
		}
		for _, r := range rates {
			if p == r.Provider {
				return watchmarket.Rate{
					Currency:         r.Currency,
					PercentChange24h: r.PercentChange24h,
					Provider:         r.Provider,
					Rate:             r.Rate,
					Timestamp:        r.LastUpdated.Unix(),
				}, nil
			}
		}
	}
	return result, errors.New(watchmarket.ErrNotFound)
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

func (c Controller) getCachedTickers(assets []controllers.Asset) (watchmarket.Tickers, error) {
	if c.cache == nil {
		return watchmarket.Tickers{}, errors.New("cache isn't available")
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

// TODO: Remove duplicates or make method
func (c Controller) getCachedRate(currency string) (result watchmarket.Rate, err error) {
	if c.cache == nil {
		return watchmarket.Rate{}, errors.New("cache isn't available")
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
