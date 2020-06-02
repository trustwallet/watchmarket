package tickerscontroller

import (
	"context"
	"errors"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"github.com/trustwallet/watchmarket/services/controllers"
	"strings"
)

type Controller struct {
	database        db.Instance
	ratesPriority   []string
	tickersPriority []string
	configuration   config.Configuration
}

func NewController(
	database db.Instance,
	ratesPriority, tickersPriority []string,
	configuration config.Configuration,
) Controller {
	return Controller{
		database,
		ratesPriority,
		tickersPriority,
		configuration,
	}
}

func (c Controller) HandleTickersRequestV2(tr controllers.TickerRequestV2, ctx context.Context) (controllers.TickerResponseV2, error) {
	if tr.Ids == nil {
		return controllers.TickerResponseV2{}, errors.New(watchmarket.ErrBadRequest)
	}

	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency), ctx)
	if err != nil {
		return controllers.TickerResponseV2{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueriesV2(tr.Ids), ctx)
	if err != nil {
		return controllers.TickerResponseV2{}, errors.New(watchmarket.ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate, ctx)

	return createResponseV2(tr, tickers), nil
}

func (c Controller) HandleTickersRequest(tr controllers.TickerRequest, ctx context.Context) (controllers.TickerResponse, error) {
	if tr.Assets == nil {
		return controllers.TickerResponse{}, errors.New(watchmarket.ErrBadRequest)
	}

	rate, err := c.getRateByPriority(strings.ToUpper(tr.Currency), ctx)
	if err != nil {
		return controllers.TickerResponse{}, errors.New(watchmarket.ErrNotFound)
	}

	tickers, err := c.getTickersByPriority(makeTickerQueries(tr.Assets), ctx)
	if err != nil {
		return controllers.TickerResponse{}, errors.New(watchmarket.ErrInternal)
	}

	tickers = c.normalizeTickers(tickers, rate, ctx)

	return createResponse(tr, tickers), nil
}
