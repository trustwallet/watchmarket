package controllers

type (
	ChartRequest struct {
		coinQuery, token, currency, timeStartRaw, maxItems string
	}
	ChartsNormalizedRequest struct {
		coin            uint
		token, currency string
		timeStart       int64
		maxItems        int
	}
)
