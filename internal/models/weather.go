package models

import "time"

type Weather struct {
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}
