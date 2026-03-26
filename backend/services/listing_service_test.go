package services

import (
	"context"
	"errors"
	"testing"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type listingRepoStub struct {
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, search string) ([]models.Listing, error)
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *listingRepoStub) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	return r.createFn(ctx, listing)
}

func (r *listingRepoStub) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	return r.getAllFn(ctx, search)
}

func (r *listingRepoStub) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	return r.getByIDFn(ctx, id)
}

func (r *listingRepoStub) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	return r.updateByID(ctx, id, listing)
}

func (r *listingRepoStub) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return r.deleteByID(ctx, id)
}

func TestListingServiceCreateValidation(t *testing.T) {
	svc := NewListingService(&listingRepoStub{})
	_, err := svc.Create(context.Background(), 1, models.CreateListingRequest{
		Title:    "",
		Category: "Books",
		Price:    10,
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestListingServiceUpdateForbidden(t *testing.T) {
	repo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 2}, nil
		},
		updateByID: func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
			t.Fatalf("update should not be called")
			return nil, nil
		},
	}
	svc := NewListingService(repo)
	_, err := svc.Update(context.Background(), 1, 10, models.UpdateListingRequest{
		Title:    "x",
		Category: "Books",
		Price:    1,
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestListingServiceUpdateSuccess(t *testing.T) {
	repo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 7}, nil
		},
		updateByID: func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 7, Title: listing.Title}, nil
		},
	}
	svc := NewListingService(repo)
	out, err := svc.Update(context.Background(), 7, 15, models.UpdateListingRequest{
		Title:    "Updated",
		Category: "Books",
		Price:    2,
	})
	if err != nil {
		t.Fatalf("update error: %v", err)
	}
	if out.Title != "Updated" {
		t.Fatalf("expected updated title")
	}
}

func TestListingServiceDeleteNotFound(t *testing.T) {
	repo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 1}, nil
		},
		deleteByID: func(ctx context.Context, id int64) (int64, error) {
			return 0, nil
		},
	}
	svc := NewListingService(repo)
	err := svc.Delete(context.Background(), 1, 100)
	if !errors.Is(err, repository.ErrListingNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestListingServiceDeleteForbidden(t *testing.T) {
	repo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 9}, nil
		},
	}
	svc := NewListingService(repo)
	err := svc.Delete(context.Background(), 1, 9)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
