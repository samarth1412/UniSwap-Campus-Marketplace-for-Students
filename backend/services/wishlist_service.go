package services

import (
	"context"
	"fmt"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type WishlistService struct {
	wishlistRepo repository.WishlistRepository
	listingRepo  repository.ListingRepository
}

func NewWishlistService(wishlistRepo repository.WishlistRepository, listingRepo repository.ListingRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		listingRepo:  listingRepo,
	}
}

func (s *WishlistService) Create(ctx context.Context, userID int64, req models.CreateWishlistRequest) (*models.WishlistItem, error) {
	if req.ListingID <= 0 {
		return nil, fmt.Errorf("%w: listing_id is required", ErrValidation)
	}

	if _, err := s.listingRepo.GetByID(ctx, req.ListingID); err != nil {
		return nil, err
	}

	return s.wishlistRepo.Create(ctx, &models.WishlistItem{
		UserID:    userID,
		ListingID: req.ListingID,
	})
}

func (s *WishlistService) Delete(ctx context.Context, userID, wishlistID int64) error {
	if wishlistID <= 0 {
		return fmt.Errorf("%w: wishlist id is required", ErrValidation)
	}

	return s.wishlistRepo.Delete(ctx, wishlistID, userID)
}

func (s *WishlistService) ListByUser(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
	return s.wishlistRepo.ListByUser(ctx, userID)
}
