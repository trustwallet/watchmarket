package models

import "time"

type (
	IDModel struct {
		ID uint `gorm:"primary_key"`
	}

	CreatedAtModel struct {
		CreatedAt time.Time
	}

	BasicModel struct {
		IDModel
		CreatedAt time.Time
	}

	ExtendedModel struct {
		BasicModel
		UpdatedAt time.Time
	}

	BasicTimeModel struct {
		UpdatedAt time.Time
		CreatedAt time.Time
	}

	ShowOption int
)

const (
	Default    = 0
	AlwaysShow = 1
	NeverShow  = 2
)
