package controllers

import (
	"context"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	TickersController interface {
		HandleTickersRequest(tr TickerRequest, ctx context.Context) (TickerResponse, error)
		HandleTickersRequestV2(tr TickerRequestV2, ctx context.Context) (TickerResponseV2, error)
	}

	ChartsController interface {
		HandleChartsRequest(cr ChartRequest, ctx context.Context) (watchmarket.Chart, error)
	}

	InfoController interface {
		HandleInfoRequest(dr DetailsRequest, ctx context.Context) (watchmarket.CoinDetails, error)
	}
)
