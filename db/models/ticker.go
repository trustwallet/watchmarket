package models

import "time"

type TickerQuery struct {
	Coin    uint
	TokenId string
}

type Ticker struct {
	BasicTimeModel
	ID                string `gorm:"index:,"`
	Coin              uint   `gorm:"primaryKey; autoIncrement:false; index:,"`
	CoinName          string `gorm:"primaryKey"`
	CoinType          string `gorm:"primaryKey"`
	TokenId           string `gorm:"primaryKey; index:,"`
	Currency          string `gorm:"primaryKey"`
	Provider          string `gorm:"primaryKey"`
	Change24h         float64
	Value             float64
	Volume            float64
	MarketCap         float64
	CirculatingSupply float64
	TotalSupply       float64
	ShowOption        ShowOption
	LastUpdated       time.Time
}
