package tickerscontroller

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/trustwallet/golibs/asset"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	sortedTickersResponse struct {
		sync.Mutex
		tickers []models.Ticker
	}
)

func (c Controller) getTickersByPriority(tickerQueries []models.TickerQuery) (watchmarket.Tickers, error) {
	var results watchmarket.Tickers
	for _, tr := range tickerQueries {
		key := strings.ToLower(asset.BuildID(tr.Coin, tr.TokenId))
		rawResult, err := c.cache.Get(key)
		if err != nil {
			continue
		}
		var result watchmarket.Ticker
		if err = json.Unmarshal(rawResult, &result); err != nil {
			continue
		}
		results = append(results, result)
	}
	return results, nil
}
