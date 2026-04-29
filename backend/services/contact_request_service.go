package services

import (
	"context"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type ContactRequestService struct {
	contactRequestRepo repository.ContactRequestRepository
	listingRepo        repository.ListingRepository
}

func NewContactRequestService(contactRequestRepo repository.ContactRequestRepository, listingRepo repository.ListingRepository) *ContactRequestService {
	return &ContactRequestService{
		contactRequestRepo: contactRequestRepo,
		listingRepo:        listingRepo,
	}
}

func (s *ContactRequestService) Create(ctx context.Context, listingID, buyerID int64, req models.CreateContactRequestRequest) (*models.ContactRequest, error) {
	message := strings.TrimSpace(req.Message)
	if message == "" {
		return nil, fmt.Errorf("%w: message is required", ErrValidation)
	}

	listing, err := s.listingRepo.GetByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.UserID == buyerID {
		return nil, ErrForbidden
	}

	contactRequest := &models.ContactRequest{
		ListingID: listingID,
		BuyerID:   buyerID,
		SellerID:  listing.UserID,
		Message:   message,
		Status:    models.ContactRequestStatusPending,
	}

	return s.contactRequestRepo.Create(ctx, contactRequest)
}
