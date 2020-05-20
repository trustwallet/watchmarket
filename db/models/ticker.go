package models

type TickerQuery struct {
	Coin    uint
	TokenId string
}

type Ticker struct {
	BasicTimeModel
	Coin      string `gorm:"primary_key;"`
	CoinName  string `gorm:"primary_key;"`
	CoinType  string `gorm:"primary_key;"`
	TokenId   string `gorm:"primary_key;"`
	Change24h float64
	Currency  string `gorm:"primary_key;"`
	Provider  string `gorm:"primary_key;"`
	Value     float64
}
