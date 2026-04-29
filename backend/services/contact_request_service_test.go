package services

import (
	"context"
	"errors"
	"testing"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type contactRequestRepoStub struct {
	createFn func(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error)
	listFn   func(ctx context.Context, sellerID int64) ([]models.ContactRequest, error)
}

func (r *contactRequestRepoStub) Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
	return r.createFn(ctx, request)
}

func (r *contactRequestRepoStub) ListBySellerID(ctx context.Context, sellerID int64) ([]models.ContactRequest, error) {
	if r.listFn == nil {
		return []models.ContactRequest{}, nil
	}
	return r.listFn(ctx, sellerID)
}

type contactListingRepoStub struct {
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, search string) ([]models.Listing, error)
	getByUser  func(ctx context.Context, userID int64) ([]models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *contactListingRepoStub) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	return r.createFn(ctx, listing)
}
func (r *contactListingRepoStub) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	return r.getAllFn(ctx, search)
}
func (r *contactListingRepoStub) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	return r.getByIDFn(ctx, id)
}
func (r *contactListingRepoStub) GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error) {
	return r.getByUser(ctx, userID)
}
func (r *contactListingRepoStub) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	return r.updateByID(ctx, id, listing)
}
func (r *contactListingRepoStub) DeleteByID(ctx context.Context, id int64) (int64, error) {
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
