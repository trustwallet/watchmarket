package controllers

import (
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type (
	TickersController interface {
		HandleTickersRequest(tr TickerRequest) (watchmarket.Tickers, error)
	}

	RatesController interface {
		HandleRatesRequest(r RateRequest) (RateResponse, error)
		GetFiatRates() (FiatRates, error)
	}

	ChartsController interface {
		HandleChartsRequest(cr ChartRequest) (watchmarket.Chart, error)
	}

	InfoController interface {
		HandleInfoRequest(dr DetailsRequest) (InfoResponse, error)
	}
)
