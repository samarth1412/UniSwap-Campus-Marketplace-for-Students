package models

import "time"

type CreateListingRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type UpdateListingRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type Listing struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Category        string    `json:"category"`
	PrimaryImageURL string    `json:"primary_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
