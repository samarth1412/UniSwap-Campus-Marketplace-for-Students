package models

import "time"

type CreateReportRequest struct {
	Reason string `json:"reason"`
}

type Report struct {
	ID             int64     `json:"id"`
	ListingID      int64     `json:"listing_id"`
	ReporterUserID int64     `json:"reporter_user_id"`
	Reason         string    `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}
