package services

import (
	"context"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type ListingService struct {
	listingRepo repository.ListingRepository
}

func NewListingService(listingRepo repository.ListingRepository) *ListingService {
	return &ListingService{listingRepo: listingRepo}
}

func (s *ListingService) Create(ctx context.Context, userID int64, req models.CreateListingRequest) (*models.Listing, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if strings.TrimSpace(req.Category) == "" {
		return nil, fmt.Errorf("%w: category is required", ErrValidation)
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("%w: price must be greater than or equal to 0", ErrValidation)
	}

	listing := &models.Listing{
		UserID:      userID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Category:    strings.TrimSpace(req.Category),
	}

	return s.listingRepo.Create(ctx, listing)
}

func (s *ListingService) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	return s.listingRepo.GetAll(ctx, search)
}

func (s *ListingService) GetByID(ctx context.Context, listingID int64) (*models.Listing, error) {
	return s.listingRepo.GetByID(ctx, listingID)
}
