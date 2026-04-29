package models

import "time"

const ContactRequestStatusPending = "pending"

type CreateContactRequestRequest struct {
	Message string `json:"message"`
}

type ContactRequest struct {
	ID        int64     `json:"id"`
	ListingID int64     `json:"listing_id"`
	BuyerID   int64     `json:"buyer_id"`
	SellerID  int64     `json:"seller_id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
