package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

var ErrForbidden = errors.New("forbidden")

type ListingService struct {
	listingRepo repository.ListingRepository
}

func NewListingService(listingRepo repository.ListingRepository) *ListingService {
	return &ListingService{listingRepo: listingRepo}
}

func (s *ListingService) Create(ctx context.Context, userID int64, req models.CreateListingRequest) (*models.Listing, error) {
	if err := validateListingFields(req.Title, req.Category, req.Price); err != nil {
		return nil, err
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

func (s *ListingService) Update(ctx context.Context, actorUserID, listingID int64, req models.UpdateListingRequest) (*models.Listing, error) {
	if err := validateListingFields(req.Title, req.Category, req.Price); err != nil {
		return nil, err
	}

	existing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != actorUserID {
		return nil, ErrForbidden
	}

	listing := &models.Listing{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		Category:    strings.TrimSpace(req.Category),
	}

	return s.listingRepo.UpdateByID(ctx, listingID, listing)
}

func (s *ListingService) Delete(ctx context.Context, actorUserID, listingID int64) error {
	existing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return err
	}
	if existing.UserID != actorUserID {
		return ErrForbidden
	}

	rowsAffected, err := s.listingRepo.DeleteByID(ctx, listingID)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrListingNotFound
	}

	return nil
}

func validateListingFields(title, category string, price float64) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	if strings.TrimSpace(category) == "" {
		return fmt.Errorf("%w: category is required", ErrValidation)
	}
	if price < 0 {
		return fmt.Errorf("%w: price must be greater than or equal to 0", ErrValidation)
	}

	return nil
}
