package models

import (
	"github.com/jinzhu/gorm"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Ticker struct {
	gorm.Model
	watchmarket.Ticker
	Provider string `gorm:"type:varchar(64)"`
}
