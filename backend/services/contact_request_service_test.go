package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type contactRequestRepoStub struct {
	createFn             func(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error)
	listReceivedBySeller func(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error)
}

func (r *contactRequestRepoStub) Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
	if r.createFn == nil {
		return &models.ContactRequest{}, nil
	}
	return r.createFn(ctx, request)
}

func (r *contactRequestRepoStub) ListReceivedBySellerID(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error) {
	if r.listReceivedBySeller == nil {
		return []models.ReceivedContactRequest{}, nil
	}
	return r.listReceivedBySeller(ctx, sellerID)
}

type contactListingRepoStub struct {
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error)
	getByUser  func(ctx context.Context, userID int64) ([]models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *contactListingRepoStub) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	if r.createFn == nil {
		return &models.Listing{}, nil
	}
	return r.createFn(ctx, listing)
}

func (r *contactListingRepoStub) GetAll(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error) {
	if r.getAllFn == nil {
		return &models.PaginatedListings{Items: []models.Listing{}, Page: query.Page, Limit: query.Limit, TotalPages: 1}, nil
	}
	return r.getAllFn(ctx, query)
}

func (r *contactListingRepoStub) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	if r.getByIDFn == nil {
		return &models.Listing{ID: id}, nil
	}
	return r.getByIDFn(ctx, id)
}

func (r *contactListingRepoStub) GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error) {
	if r.getByUser == nil {
		return []models.Listing{}, nil
	}
	return r.getByUser(ctx, userID)
}

func (r *contactListingRepoStub) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	if r.updateByID == nil {
		return &models.Listing{ID: id}, nil
	}
	return r.updateByID(ctx, id, listing)
}

func (r *contactListingRepoStub) DeleteByID(ctx context.Context, id int64) (int64, error) {
	if r.deleteByID == nil {
		return 1, nil
	}
	return r.deleteByID(ctx, id)
}

func TestContactRequestServiceValidation(t *testing.T) {
	svc := NewContactRequestService(&contactRequestRepoStub{}, &contactListingRepoStub{})
	_, err := svc.Create(context.Background(), 1, 2, models.CreateContactRequestRequest{Message: "   "})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestContactRequestServiceListingNotFound(t *testing.T) {
	listingRepo := &contactListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return nil, repository.ErrListingNotFound
		},
	}
	svc := NewContactRequestService(&contactRequestRepoStub{}, listingRepo)
	_, err := svc.Create(context.Background(), 10, 2, models.CreateContactRequestRequest{Message: "hi"})
	if !errors.Is(err, repository.ErrListingNotFound) {
		t.Fatalf("expected listing not found, got %v", err)
	}
}

func TestContactRequestServiceForbiddenForOwnListing(t *testing.T) {
	listingRepo := &contactListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 7}, nil
		},
	}
	repo := &contactRequestRepoStub{
		createFn: func(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
			t.Fatalf("create should not be called for own listing contact")
			return nil, nil
		},
	}

	svc := NewContactRequestService(repo, listingRepo)
	_, err := svc.Create(context.Background(), 10, 7, models.CreateContactRequestRequest{Message: "hi"})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestContactRequestServiceSuccess(t *testing.T) {
	listingRepo := &contactListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 9}, nil
		},
	}
	repo := &contactRequestRepoStub{
		createFn: func(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
			return &models.ContactRequest{
				ID:        22,
				ListingID: request.ListingID,
				BuyerID:   request.BuyerID,
				SellerID:  request.SellerID,
				Message:   request.Message,
				Status:    request.Status,
			}, nil
		},
	}

	svc := NewContactRequestService(repo, listingRepo)
	result, err := svc.Create(context.Background(), 10, 2, models.CreateContactRequestRequest{Message: "  Interested in pickup today?  "})
	if err != nil {
		t.Fatalf("create contact request: %v", err)
	}
	if result.ID != 22 || result.SellerID != 9 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Message != "Interested in pickup today?" {
		t.Fatalf("expected trimmed message, got %q", result.Message)
	}
	if result.Status != models.ContactRequestStatusPending {
		t.Fatalf("expected pending status, got %q", result.Status)
	}
}

func TestContactRequestServiceListReceivedValidation(t *testing.T) {
	svc := NewContactRequestService(&contactRequestRepoStub{}, &contactListingRepoStub{})
	_, err := svc.ListReceivedBySellerID(context.Background(), 0)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestContactRequestServiceListReceivedBySellerID(t *testing.T) {
	now := time.Unix(1000, 0)
	repo := &contactRequestRepoStub{
		listReceivedBySeller: func(ctx context.Context, sellerID int64) ([]models.ReceivedContactRequest, error) {
			if sellerID != 9 {
				t.Fatalf("expected seller id 9, got %d", sellerID)
			}
			return []models.ReceivedContactRequest{{
				ID:           4,
				ListingID:    11,
				ListingTitle: "Desk Lamp",
				BuyerID:      2,
				BuyerName:    "Buyer One",
				BuyerEmail:   "buyer@example.edu",
				Message:      "Still available?",
				Status:       models.ContactRequestStatusPending,
				CreatedAt:    now,
			}}, nil
		},
	}

	svc := NewContactRequestService(repo, &contactListingRepoStub{})
	requests, err := svc.ListReceivedBySellerID(context.Background(), 9)
	if err != nil {
		t.Fatalf("list received contact requests: %v", err)
	}
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}
	if requests[0].ListingTitle != "Desk Lamp" || requests[0].BuyerEmail != "buyer@example.edu" {
		t.Fatalf("unexpected request: %+v", requests[0])
	}
}
