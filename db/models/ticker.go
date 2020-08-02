package models

import "time"

type TickerQuery struct {
	Coin    uint
	TokenId string
}

type Ticker struct {
	BasicTimeModel
	ID          string `sql:"index"`
	Coin        uint   `gorm:"primary_key" sql:"index"`
	CoinName    string `gorm:"primary_key"`
	CoinType    string `gorm:"primary_key"`
	TokenId     string `gorm:"primary_key" sql:"index"`
	Currency    string `gorm:"primary_key"`
	Provider    string `gorm:"primary_key"`
	Change24h   float64
	Value       float64
	Volume      float64
	MarketCap   float64
	ShowOption  ShowOption
	LastUpdated time.Time
}
