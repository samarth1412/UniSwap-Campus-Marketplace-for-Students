package services

import (
	"context"
	"errors"
	"testing"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type listingImageListingRepoStub struct {
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error)
	getByUser  func(ctx context.Context, userID int64) ([]models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *listingImageListingRepoStub) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	return r.createFn(ctx, listing)
}
func (r *listingImageListingRepoStub) GetAll(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error) {
	return r.getAllFn(ctx, query)
}
func (r *listingImageListingRepoStub) GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error) {
	return r.getByUser(ctx, userID)
}
func (r *listingImageListingRepoStub) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	return r.getByIDFn(ctx, id)
}
func (r *listingImageListingRepoStub) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	return r.updateByID(ctx, id, listing)
}
func (r *listingImageListingRepoStub) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return r.deleteByID(ctx, id)
}

type listingImageRepoStub struct {
	createManyFn func(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error)
}

func (r *listingImageRepoStub) CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
	return r.createManyFn(ctx, listingID, imageURLs)
}

func TestListingImageServiceValidation(t *testing.T) {
	svc := NewListingImageService(&listingImageListingRepoStub{}, &listingImageRepoStub{})
	_, err := svc.AddImages(context.Background(), 1, 1, nil)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestListingImageServiceForbidden(t *testing.T) {
	listingRepo := &listingImageListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 5}, nil
		},
	}
	svc := NewListingImageService(listingRepo, &listingImageRepoStub{})
	_, err := svc.AddImages(context.Background(), 2, 10, []string{"/uploads/x.jpg"})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestListingImageServiceListingNotFound(t *testing.T) {
	listingRepo := &listingImageListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return nil, repository.ErrListingNotFound
		},
	}
	svc := NewListingImageService(listingRepo, &listingImageRepoStub{})
	_, err := svc.AddImages(context.Background(), 1, 99, []string{"/uploads/x.jpg"})
	if !errors.Is(err, repository.ErrListingNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestListingImageServiceSuccess(t *testing.T) {
	listingRepo := &listingImageListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 3}, nil
		},
	}
	imageRepo := &listingImageRepoStub{
		createManyFn: func(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
			return []models.ListingImage{
				{ID: 1, ListingID: listingID, ImageURL: imageURLs[0]},
				{ID: 2, ListingID: listingID, ImageURL: imageURLs[1]},
			}, nil
		},
	}
	svc := NewListingImageService(listingRepo, imageRepo)
	out, err := svc.AddImages(context.Background(), 3, 11, []string{"/uploads/1.jpg", "/uploads/2.jpg"})
	if err != nil {
		t.Fatalf("add images: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 images, got %d", len(out))
	}
}
