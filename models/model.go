package models

import (
	"time"
)

type RequestInfo struct {
	OriginalURL string `json:"url"`
	Expiry      uint64 `json:"expiry,omitempty"` // In hour
}

type URL struct {
	OriginalURL    string    `bson:"original_url"`
	ShortURL       string    `bson:"short_url"`
	ExpirationDate time.Time `bson:"expiration_date"`
}

type ResponseInfo struct {
	OriginalURL    string    `json:"original_url"`
	ShortURL       string    `json:"short_url"`
	ExpirationDate time.Time `json:"expiration_date"`
}
