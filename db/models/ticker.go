package models

import "time"

type TickerQuery struct {
	Coin    uint
	TokenId string
}

type Ticker struct {
	BasicTimeModel
	Coin        uint   `gorm:"primary_key;" sql:"index"`
	CoinName    string `gorm:"primary_key;" sql:"index"`
	CoinType    string `gorm:"primary_key;"`
	TokenId     string `gorm:"primary_key;" sql:"index"`
	Currency    string `gorm:"primary_key;" sql:"index"`
	Provider    string `gorm:"primary_key;" sql:"index"`
	Change24h   float64
	Value       float64
	Volume      float64
	MarketCap   float64
	ShowOption  ShowOption
	LastUpdated time.Time
}
