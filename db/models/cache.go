package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Cache struct {
	gorm.Model
	Key  string `gorm:"unique_index, type:varchar(64)"`
	Data []byte
}

type CachingInterval struct {
	gorm.Model
	Timestamp time.Time
	Duration  int64
	Key       string `gorm:"unique_index, type:varchar(64)"`
}
