package models

import "time"

type CreateWishlistRequest struct {
	ListingID int64 `json:"listing_id"`
}

type WishlistItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ListingID int64     `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
}

type WishlistListing struct {
	WishlistID int64     `json:"wishlist_id"`
	SavedAt    time.Time `json:"saved_at"`
	Listing    Listing   `json:"listing"`
}
