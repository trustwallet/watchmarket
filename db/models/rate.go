package models

type Rate struct {
	BasicTimeModel
	Currency         string `gorm:"primary_key;"sql:"index"`
	PercentChange24h float64
	Provider         string `gorm:"primary_key;"sql:"index"`
	Rate             float64
	Timestamp        int64
}
