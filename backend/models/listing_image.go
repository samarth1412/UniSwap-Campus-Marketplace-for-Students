package models

import "time"

type ListingImage struct {
	ID        int64     `json:"id"`
	ListingID int64     `json:"listing_id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}
