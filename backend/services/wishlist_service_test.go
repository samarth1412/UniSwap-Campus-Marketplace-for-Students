package services

import (
	"context"
	"errors"
	"testing"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type wishlistRepoStub struct {
	createFn     func(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error)
	deleteFn     func(ctx context.Context, wishlistID, userID int64) error
	listByUserFn func(ctx context.Context, userID int64) ([]models.WishlistListing, error)
}

func (r *wishlistRepoStub) Create(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	return r.createFn(ctx, item)
}

func (r *wishlistRepoStub) Delete(ctx context.Context, wishlistID, userID int64) error {
	return r.deleteFn(ctx, wishlistID, userID)
}

func (r *wishlistRepoStub) ListByUser(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
	return r.listByUserFn(ctx, userID)
}

func TestWishlistServiceCreateValidation(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{}, &listingRepoStub{})
	_, err := svc.Create(context.Background(), 1, models.CreateWishlistRequest{})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestWishlistServiceCreateListingNotFound(t *testing.T) {
	listingRepo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return nil, repository.ErrListingNotFound
		},
	}
	svc := NewWishlistService(&wishlistRepoStub{}, listingRepo)
	_, err := svc.Create(context.Background(), 1, models.CreateWishlistRequest{ListingID: 10})
	if !errors.Is(err, repository.ErrListingNotFound) {
		t.Fatalf("expected listing not found, got %v", err)
	}
}

func TestWishlistServiceCreateSuccess(t *testing.T) {
	listingRepo := &listingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id}, nil
		},
	}
	wishlistRepo := &wishlistRepoStub{
		createFn: func(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
			if item.UserID != 7 || item.ListingID != 10 {
				t.Fatalf("unexpected wishlist item: %+v", item)
			}
			return &models.WishlistItem{ID: 3, UserID: item.UserID, ListingID: item.ListingID}, nil
		},
	}
	svc := NewWishlistService(wishlistRepo, listingRepo)
	item, err := svc.Create(context.Background(), 7, models.CreateWishlistRequest{ListingID: 10})
	if err != nil {
		t.Fatalf("create wishlist: %v", err)
	}
	if item.ID != 3 {
		t.Fatalf("expected item id 3, got %d", item.ID)
	}
}

func TestWishlistServiceDeleteValidation(t *testing.T) {
	svc := NewWishlistService(&wishlistRepoStub{}, &listingRepoStub{})
	err := svc.Delete(context.Background(), 1, 0)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestWishlistServiceDeleteSuccess(t *testing.T) {
	wishlistRepo := &wishlistRepoStub{
		deleteFn: func(ctx context.Context, wishlistID, userID int64) error {
			if wishlistID != 4 || userID != 9 {
				t.Fatalf("unexpected delete args wishlistID=%d userID=%d", wishlistID, userID)
			}
			return nil
		},
	}
	svc := NewWishlistService(wishlistRepo, &listingRepoStub{})
	if err := svc.Delete(context.Background(), 9, 4); err != nil {
		t.Fatalf("delete wishlist: %v", err)
	}
}

func TestWishlistServiceListByUser(t *testing.T) {
	wishlistRepo := &wishlistRepoStub{
		listByUserFn: func(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
			return []models.WishlistListing{{WishlistID: 2, Listing: models.Listing{ID: 8, UserID: userID}}}, nil
		},
	}
	svc := NewWishlistService(wishlistRepo, &listingRepoStub{})
	items, err := svc.ListByUser(context.Background(), 5)
	if err != nil {
		t.Fatalf("list wishlist: %v", err)
	}
	if len(items) != 1 || items[0].Listing.UserID != 5 {
		t.Fatalf("unexpected wishlist items: %+v", items)
	}
}
