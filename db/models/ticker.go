package models

type Ticker struct {
	BasicTimeModel
	Coin      uint   `gorm:"primary_key; auto_increment:false" sql:"index"`
	CoinName  string `gorm:"primary_key;"sql:"index"`
	CoinType  string `gorm:"primary_key;"sql:"index"`
	TokenId   string `gorm:"primary_key;"sql:"index"`
	Change24h float64
	Currency  string `gorm:"primary_key;"sql:"index"`
	Provider  string `gorm:"primary_key;"sql:"index"`
	Value     float64
}
