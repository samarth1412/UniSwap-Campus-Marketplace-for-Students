package services

import (
	"context"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type ReportService struct {
	reportRepo  repository.ReportRepository
	listingRepo repository.ListingRepository
}

func NewReportService(reportRepo repository.ReportRepository, listingRepo repository.ListingRepository) *ReportService {
	return &ReportService{
		reportRepo:  reportRepo,
		listingRepo: listingRepo,
	}
}

func (s *ReportService) Create(ctx context.Context, listingID, reporterUserID int64, req models.CreateReportRequest) (*models.Report, error) {
	if strings.TrimSpace(req.Reason) == "" {
		return nil, fmt.Errorf("%w: reason is required", ErrValidation)
	}

	// Confirm listing exists to return a clean 404 from handler logic.
	if _, err := s.listingRepo.GetByID(ctx, listingID); err != nil {
		return nil, err
	}

	report := &models.Report{
		ListingID:      listingID,
		ReporterUserID: reporterUserID,
		Reason:         strings.TrimSpace(req.Reason),
	}

	return s.reportRepo.Create(ctx, report)
}
