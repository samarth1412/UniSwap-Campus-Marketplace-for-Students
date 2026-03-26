package services

import (
	"context"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type UploadService struct {
	listingRepo      repository.ListingRepository
	listingImageRepo repository.ListingImageRepository
}

func NewUploadService(listingRepo repository.ListingRepository, listingImageRepo repository.ListingImageRepository) *UploadService {
	return &UploadService{
		listingRepo:      listingRepo,
		listingImageRepo: listingImageRepo,
	}
}

func (s *UploadService) SaveListingImage(ctx context.Context, userID, listingID int64, imageURL string, isPrimary bool) (*models.ListingImage, error) {
	if listingID <= 0 {
		return nil, fmt.Errorf("%w: listing_id is required", ErrValidation)
	}
	if strings.TrimSpace(imageURL) == "" {
		return nil, fmt.Errorf("%w: image url is required", ErrValidation)
	}

	listing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if listing.UserID != userID {
		return nil, ErrForbidden
	}

	return s.listingImageRepo.Create(ctx, &models.ListingImage{
		ListingID: listingID,
		ImageURL:  strings.TrimSpace(imageURL),
		IsPrimary: isPrimary,
	})
}
