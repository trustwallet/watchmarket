package models

import "time"

type Rate struct {
	BasicTimeModel
	Currency         string `gorm:"primaryKey; index:,"`
	Provider         string `gorm:"primaryKey"`
	PercentChange24h float64
	Rate             float64
	ShowOption       ShowOption
	LastUpdated      time.Time
}
