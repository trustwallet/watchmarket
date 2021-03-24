package worker

import (
	"sync"
	"time"

	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	tickers struct {
		tickers watchmarket.Tickers
		sync.Mutex
	}

	rates struct {
		rates watchmarket.Rates
		sync.Mutex
	}
)

func (t *tickers) Add(tickers watchmarket.Tickers) {
	t.Lock()
	t.tickers = append(t.tickers, tickers...)
	t.Unlock()
}

func (r *rates) Add(rates watchmarket.Rates) {
	r.Lock()
	r.rates = append(r.rates, rates...)
	r.Unlock()
}

func isHigherPriority(priorities []string, current, new string) bool {
	for _, p := range priorities {
		if p == current {
			return false
		} else if p == new {
			return true
		}
	}
	return false
}

func toTickersModel(tickers watchmarket.Tickers) []models.Ticker {
	result := make([]models.Ticker, 0, len(tickers))
	for _, t := range tickers {
		result = append(result, models.Ticker{
			ID:                asset.BuildID(t.Coin, t.TokenId),
			Coin:              t.Coin,
			CoinName:          t.CoinName,
			CoinType:          string(t.CoinType),
			TokenId:           t.TokenId,
			Change24h:         t.Price.Change24h,
			Currency:          t.Price.Currency,
			Provider:          t.Price.Provider,
			Value:             t.Price.Value,
			Volume:            t.Volume,
			MarketCap:         t.MarketCap,
			CirculatingSupply: t.CirculatingSupply,
			TotalSupply:       t.TotalSupply,
			LastUpdated:       t.LastUpdate,
		})
	}

	return result
}

func toRatesModel(rates watchmarket.Rates) []models.Rate {
	result := make([]models.Rate, 0)
	for _, r := range rates {
		result = append(result, models.Rate{
			Currency:         r.Currency,
			PercentChange24h: r.PercentChange24h,
			Provider:         r.Provider,
			Rate:             r.Rate,
			LastUpdated:      time.Unix(r.Timestamp, 0),
		})
	}
	return result
}

func fromModelToTicker(m models.Ticker) watchmarket.Ticker {
	return watchmarket.Ticker{
		Coin:       m.Coin,
		CoinName:   m.CoinName,
		CoinType:   watchmarket.CoinType(m.CoinType),
		LastUpdate: m.LastUpdated,
		Price: watchmarket.Price{
			Change24h: m.Change24h,
			Currency:  m.Currency,
			Provider:  m.Provider,
			Value:     m.Value,
		},
		TokenId:    m.TokenId,
		ShowOption: int(m.ShowOption),
	}
}

func fromModelToRate(m models.Rate) watchmarket.Rate {
	return watchmarket.Rate{
		Currency:         m.Currency,
		PercentChange24h: m.PercentChange24h,
		Provider:         m.Provider,
		Rate:             m.Rate,
		Timestamp:        m.LastUpdated.Unix(),
		ShowOption:       int(m.ShowOption),
	}
}
