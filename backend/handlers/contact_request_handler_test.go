package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"uniswap-campus-marketplace/apiresponse"
	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/services"
)

type contactRequestRepoForReceivedHandlerTest struct {
	listReceivedBySeller func(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error)
}

func (r *contactRequestRepoForReceivedHandlerTest) Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
	return &models.ContactRequest{ID: 1}, nil
}

func (r *contactRequestRepoForReceivedHandlerTest) ListReceivedBySellerID(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error) {
	if r.listReceivedBySeller == nil {
		return []models.ReceivedContactRequest{}, nil
	}
	return r.listReceivedBySeller(ctx, sellerID)
}

func newContactRequestHandlerForTest(repo *contactRequestRepoForReceivedHandlerTest) *ContactRequestHandler {
	return NewContactRequestHandler(services.NewContactRequestService(repo, &listingRepoForHandlerTest{}))
}

func TestReceivedContactRequestsMethodNotAllowed(t *testing.T) {
	h := newContactRequestHandlerForTest(&contactRequestRepoForReceivedHandlerTest{})
	req := httptest.NewRequest(http.MethodPost, "/api/contact-requests/received", nil)
	rec := httptest.NewRecorder()

	h.Received(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestReceivedContactRequestsUnauthorized(t *testing.T) {
	h := newContactRequestHandlerForTest(&contactRequestRepoForReceivedHandlerTest{})
	req := httptest.NewRequest(http.MethodGet, "/api/contact-requests/received", nil)
	rec := httptest.NewRecorder()

	h.Received(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestReceivedContactRequestsSuccess(t *testing.T) {
	now := time.Unix(2000, 0).UTC()
	repo := &contactRequestRepoForReceivedHandlerTest{
		listReceivedBySeller: func(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error) {
			if sellerID != 12 {
				t.Fatalf("expected seller id 12, got %d", sellerID)
			}
			return []models.ReceivedContactRequest{{
				ID:           10,
				ListingID:    4,
				ListingTitle: "Bike",
				BuyerID:      7,
				BuyerName:    "Buyer One",
				BuyerEmail:   "buyer1@example.edu",
				Message:      "Can I pick this up today?",
				Status:       models.ContactRequestStatusPending,
				CreatedAt:    now,
			}}, nil
		},
	}
	h := newContactRequestHandlerForTest(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/contact-requests/received", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 12})(http.HandlerFunc(h.Received)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body apiresponse.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !body.Success {
		t.Fatalf("expected success response")
	}

	dataBytes, err := json.Marshal(body.Data)
	if err != nil {
		t.Fatalf("marshal body data: %v", err)
	}

	var requests []models.ReceivedContactRequest
	if err := json.Unmarshal(dataBytes, &requests); err != nil {
		t.Fatalf("unmarshal requests: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}
	if requests[0].ListingTitle != "Bike" || requests[0].BuyerEmail != "buyer1@example.edu" {
		t.Fatalf("unexpected request payload: %+v", requests[0])
	}
}
