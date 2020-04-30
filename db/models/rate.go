package models

import (
	"github.com/jinzhu/gorm"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
)

type Rate struct {
	gorm.Model
	watchmarket.Rate
	Provider string
}
