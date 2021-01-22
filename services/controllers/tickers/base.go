package tickerscontroller

import (
	"errors"
	"strings"

	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
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

func (c Controller) HandleTickersRequestV2(tr controllers.TickerRequestV2) (controllers.TickerResponseV2, error) {
	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency))
	if err != nil {
		return controllers.TickerResponseV2{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueriesV2(tr.Ids))
	if err != nil {
		return controllers.TickerResponseV2{}, errors.New(watchmarket.ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate)

	return createResponseV2(tr, tickers), nil
}

func (c Controller) HandleTickersRequest(tr controllers.TickerRequest) (controllers.TickerResponse, error) {
	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency))
	if err != nil {
		return controllers.TickerResponse{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueries(tr.Assets))
	if err != nil {
		return controllers.TickerResponse{}, errors.New(watchmarket.ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate)

	return createResponse(tr, tickers), nil
}
