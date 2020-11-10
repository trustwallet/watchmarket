package models

import "time"

type Interval string

type Alert struct {
	UpdatedAt  time.Time
	AssetID    string   `gorm:"primaryKey; index:,"`
	Interval   Interval `gorm:"primaryKey; index:,"`
	Price      float64
	Difference float64
}

const (
	Hour Interval = "1h"
	Day  Interval = "1d"
	Week Interval = "1w"
)
