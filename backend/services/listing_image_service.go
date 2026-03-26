package services

import (
	"context"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type ListingImageService struct {
	listingRepo      repository.ListingRepository
	listingImageRepo repository.ListingImageRepository
}

func NewListingImageService(listingRepo repository.ListingRepository, listingImageRepo repository.ListingImageRepository) *ListingImageService {
	return &ListingImageService{
		listingRepo:      listingRepo,
		listingImageRepo: listingImageRepo,
	}
}

func (s *ListingImageService) AddImages(ctx context.Context, actorUserID, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
	if len(imageURLs) == 0 {
		return nil, fmt.Errorf("%w: at least one image is required", ErrValidation)
	}

	normalized := make([]string, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		trimmed := strings.TrimSpace(imageURL)
		if trimmed == "" {
			return nil, fmt.Errorf("%w: image url is required", ErrValidation)
		}
		normalized = append(normalized, trimmed)
	}

	existing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != actorUserID {
		return nil, ErrForbidden
	}

	return s.listingImageRepo.CreateMany(ctx, listingID, normalized)
}
