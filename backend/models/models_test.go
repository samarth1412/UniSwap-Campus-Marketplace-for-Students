package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRegisterRequestJSONTags(t *testing.T) {
	var req RegisterRequest
	err := json.Unmarshal([]byte(`{"full_name":"Test User","email":"test@example.edu","password":"secret123","university":"UF"}`), &req)
	if err != nil {
		t.Fatalf("unmarshal register request: %v", err)
	}

	if req.FullName != "Test User" || req.Email != "test@example.edu" || req.University != "UF" {
		t.Fatalf("unexpected register request values: %+v", req)
	}
}

func TestListingAndReportJSONTags(t *testing.T) {
	listing := Listing{
		ID:          9,
		UserID:      4,
		Title:       "Book",
		Description: "Used book",
		Price:       12.5,
		Category:    "Books",
		CreatedAt:   time.Unix(100, 0),
	}
	report := Report{
		ID:             11,
		ListingID:      9,
		ReporterUserID: 2,
		Reason:         "Spam",
		CreatedAt:      time.Unix(200, 0),
	}

	listingJSON, err := json.Marshal(listing)
	if err != nil {
		t.Fatalf("marshal listing: %v", err)
	}
	reportJSON, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal report: %v", err)
	}

	var listingMap map[string]any
	if err := json.Unmarshal(listingJSON, &listingMap); err != nil {
		t.Fatalf("unmarshal listing json: %v", err)
	}
	if _, ok := listingMap["user_id"]; !ok {
		t.Fatalf("expected user_id json key in listing")
	}

	var reportMap map[string]any
	if err := json.Unmarshal(reportJSON, &reportMap); err != nil {
		t.Fatalf("unmarshal report json: %v", err)
	}
	if _, ok := reportMap["reporter_user_id"]; !ok {
		t.Fatalf("expected reporter_user_id json key in report")
	}
}

func TestUserJSONOmitsPasswordHash(t *testing.T) {
	user := User{
		ID:           1,
		FullName:     "User One",
		Email:        "user@example.edu",
		PasswordHash: "internal-only",
		University:   "UF",
		CreatedAt:    time.Unix(100, 0),
		UpdatedAt:    time.Unix(100, 0),
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatalf("unmarshal user json: %v", err)
	}
	if _, ok := body["password_hash"]; ok {
		t.Fatalf("password_hash must not be serialized")
	}
}

func TestListingImageJSONTags(t *testing.T) {
	image := ListingImage{
		ID:        2,
		ListingID: 7,
		ImageURL:  "/uploads/img.jpg",
		CreatedAt: time.Unix(300, 0),
	}

	b, err := json.Marshal(image)
	if err != nil {
		t.Fatalf("marshal listing image: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatalf("unmarshal listing image json: %v", err)
	}
	if _, ok := body["image_url"]; !ok {
		t.Fatalf("expected image_url json key")
	}
}

func TestContactRequestJSONTags(t *testing.T) {
	request := ContactRequest{
		ID:        3,
		ListingID: 7,
		BuyerID:   1,
		SellerID:  2,
		Message:   "Interested in this listing",
		Status:    ContactRequestStatusPending,
		CreatedAt: time.Unix(400, 0),
	}

	b, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal contact request: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(b, &body); err != nil {
		t.Fatalf("unmarshal contact request json: %v", err)
	}
	if _, ok := body["listing_id"]; !ok {
		t.Fatalf("expected listing_id json key")
	}
	if _, ok := body["buyer_id"]; !ok {
		t.Fatalf("expected buyer_id json key")
	}
	if _, ok := body["seller_id"]; !ok {
		t.Fatalf("expected seller_id json key")
	}
	if body["status"] != ContactRequestStatusPending {
		t.Fatalf("expected pending status, got %+v", body["status"])
	}
}
