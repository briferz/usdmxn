package model

import (
	"time"
)

type Provider struct {
	ProviderName string    `json:"provider_name"`
	Value        float32   `json:"value"`
	LastUpdated  time.Time `json:"last_updated"`
	Error        string    `json:"error,omitempty"`
}
