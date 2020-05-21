package models

import "time"

type Rate struct {
	BasicTimeModel
	Currency         string `gorm:"primary_key;"`
	PercentChange24h float64
	Provider         string `gorm:"primary_key;"`
	Rate             float64
	LastUpdated      time.Time
}
