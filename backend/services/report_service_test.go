package services

import (
	"context"
	"errors"
	"testing"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type reportRepoStub struct {
	createFn func(ctx context.Context, report *models.Report) (*models.Report, error)
}

func (r *reportRepoStub) Create(ctx context.Context, report *models.Report) (*models.Report, error) {
	return r.createFn(ctx, report)
}

type reportListingRepoStub struct {
	getByIDFn  func(ctx context.Context, id int64) (*models.Listing, error)
	createFn   func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn   func(ctx context.Context, search string) ([]models.Listing, error)
	updateByID func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteByID func(ctx context.Context, id int64) (int64, error)
}

func (r *reportListingRepoStub) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	return r.createFn(ctx, listing)
}
func (r *reportListingRepoStub) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	return r.getAllFn(ctx, search)
}
func (r *reportListingRepoStub) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	return r.getByIDFn(ctx, id)
}
func (r *reportListingRepoStub) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	return r.updateByID(ctx, id, listing)
}
func (r *reportListingRepoStub) DeleteByID(ctx context.Context, id int64) (int64, error) {
	return r.deleteByID(ctx, id)
}

func TestReportServiceValidation(t *testing.T) {
	svc := NewReportService(&reportRepoStub{}, &reportListingRepoStub{})
	_, err := svc.Create(context.Background(), 1, 1, models.CreateReportRequest{Reason: ""})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestReportServiceListingNotFound(t *testing.T) {
	listingRepo := &reportListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return nil, repository.ErrListingNotFound
		},
	}
	svc := NewReportService(&reportRepoStub{}, listingRepo)
	_, err := svc.Create(context.Background(), 1, 1, models.CreateReportRequest{Reason: "spam"})
	if !errors.Is(err, repository.ErrListingNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestReportServiceSuccess(t *testing.T) {
	listingRepo := &reportListingRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 1}, nil
		},
	}
	reportRepo := &reportRepoStub{
		createFn: func(ctx context.Context, report *models.Report) (*models.Report, error) {
			return &models.Report{ID: 7, ListingID: report.ListingID, Reason: report.Reason}, nil
		},
	}
	svc := NewReportService(reportRepo, listingRepo)
	got, err := svc.Create(context.Background(), 5, 8, models.CreateReportRequest{Reason: "spam"})
	if err != nil {
		t.Fatalf("create report: %v", err)
	}
	if got.ID != 7 || got.ListingID != 5 {
		t.Fatalf("unexpected report returned: %+v", got)
	}
}
