package fixer

import "time"

type responseModel struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      time.Time          `json:"date"`
	Rates     map[string]float32 `json:"rates"`
}
