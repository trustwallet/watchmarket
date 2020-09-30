package models

import "time"

type Rate struct {
	BasicTimeModel
	Currency         string `gorm:"primary_key" sql:"index"`
	Provider         string `gorm:"primary_key"`
	PercentChange24h float64
	Rate             float64
	ShowOption       ShowOption
	LastUpdated      time.Time
}
