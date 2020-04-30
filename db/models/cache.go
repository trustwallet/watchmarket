package models

import "github.com/jinzhu/gorm"

type Cache struct {
	gorm.Model
	Key  string `gorm:"unique_index"`
	Data []byte
}

type CachingInterval struct {
	gorm.Model
	Timestamp int64
	Duration  int64
	Key       string `gorm:"unique_index"`
}
