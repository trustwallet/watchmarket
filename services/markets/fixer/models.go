package fixer

import (
	"time"
)

type Rate struct {
	Timestamp int64              `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
	UpdatedAt time.Time          `json:"updated_at"`
}
