package models

import (
	"github.com/jinzhu/gorm"
)

type (
	Condition uint

	Subscription struct {
		gorm.Model
		Coin  uint
		Token string `gorm:"type:varchar(128)"`
		Condition
		SubscriptionId uint `gorm:"unique_index"`
		Price          float64
	}
)

const (
	_    = iota
	Less = 0
	More
)
