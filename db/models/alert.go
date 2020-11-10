package models

import "time"

type Alert struct {
	UpdatedAt  time.Time
	AssetID    string
	Type       string
	Difference float64
	Processed  bool
	Display    bool
}
