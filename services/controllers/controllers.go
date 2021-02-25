package controllers

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	TickersController interface {
		HandleTickersRequest(tr TickerRequest) (TickerResponse, error)
		HandleTickersRequestV2(tr TickerRequestV2) (TickerResponseV2, error)
	}

	RatesController interface {
		HandleRatesRequest(r RateRequest) (RateResponse, error)
	}

	ChartsController interface {
		HandleChartsRequest(cr ChartRequest) (watchmarket.Chart, error)
	}

	InfoController interface {
		HandleInfoRequest(dr DetailsRequest) (InfoResponse, error)
	}
)
