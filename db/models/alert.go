package models

import "time"

type Interval string

type Alert struct {
	UpdatedAt time.Time
	AssetID   string
	Interval
	Price      float64
	Difference float64
}

const (
	Hour Interval = "1h"
	Day  Interval = "1d"
	Week Interval = "1w"
)
