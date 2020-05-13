package models

import "time"

type IDModel struct {
	ID uint `gorm:"primary_key"`
}

type CreatedAtModel struct {
	CreatedAt time.Time
}

type BasicModel struct {
	IDModel
	CreatedAt time.Time
}

type ExtendedModel struct {
	BasicModel
	UpdatedAt time.Time
}

type BasicTimeModel struct {
	UpdatedAt time.Time
	CreatedAt time.Time
}
