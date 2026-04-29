package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"uniswap-campus-marketplace/models"
)

var (
	_ UserRepository           = (*PostgresUserRepository)(nil)
	_ ListingRepository        = (*PostgresListingRepository)(nil)
	_ ReportRepository         = (*PostgresReportRepository)(nil)
	_ ContactRequestRepository = (*PostgresContactRequestRepository)(nil)
	_ ListingImageRepository   = (*PostgresListingImageRepository)(nil)
)

func newClosedPostgresDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("postgres", "host=127.0.0.1 port=1 user=x password=x dbname=x sslmode=disable")
	if err != nil {
		t.Fatalf("open postgres db handle: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}
	return db
}

func TestConstructors(t *testing.T) {
	db := newClosedPostgresDB(t)

	userRepo := NewPostgresUserRepository(db)
	if userRepo == nil || userRepo.db != db {
		t.Fatalf("unexpected user repo constructor output")
	}

	listingRepo := NewPostgresListingRepository(db)
	if listingRepo == nil || listingRepo.db != db {
		t.Fatalf("unexpected listing repo constructor output")
	}

	reportRepo := NewPostgresReportRepository(db)
	if reportRepo == nil || reportRepo.db != db {
		t.Fatalf("unexpected report repo constructor output")
	}

	contactRequestRepo := NewPostgresContactRequestRepository(db)
	if contactRequestRepo == nil || contactRequestRepo.db != db {
		t.Fatalf("unexpected contact request repo constructor output")
	}

	imageRepo := NewPostgresListingImageRepository(db)
	if imageRepo == nil || imageRepo.db != db {
		t.Fatalf("unexpected listing image repo constructor output")
	}
}

func TestUserRepositoryMethodsReturnErrorWhenDBClosed(t *testing.T) {
	repo := NewPostgresUserRepository(newClosedPostgresDB(t))

	_, err := repo.Create(context.Background(), &models.User{FullName: "A", Email: "a@b.com", PasswordHash: "x"})
	if err == nil || !strings.Contains(err.Error(), "create user") {
		t.Fatalf("expected create user wrapped error, got %v", err)
	}

	_, err = repo.GetByEmail(context.Background(), "a@b.com")
	if err == nil || !strings.Contains(err.Error(), "get user by email") {
		t.Fatalf("expected get user by email wrapped error, got %v", err)
	}

	_, err = repo.GetByID(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "get user by id") {
		t.Fatalf("expected get user by id wrapped error, got %v", err)
	}
}

func TestListingRepositoryMethodsReturnErrorWhenDBClosed(t *testing.T) {
	repo := NewPostgresListingRepository(newClosedPostgresDB(t))

	_, err := repo.Create(context.Background(), &models.Listing{
		UserID:      1,
		Title:       "Book",
		Description: "Used",
		Category:    "Books",
		Price:       10,
	})
	if err == nil || !strings.Contains(err.Error(), "create listing") {
		t.Fatalf("expected create listing wrapped error, got %v", err)
	}

	_, err = repo.GetAll(context.Background(), models.ListingQuery{Page: 1, Limit: 10})
	if err == nil || !strings.Contains(err.Error(), "get listings") {
		t.Fatalf("expected get listings wrapped error, got %v", err)
	}

	_, err = repo.GetByID(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "get listing by id") {
		t.Fatalf("expected get listing by id wrapped error, got %v", err)
	}

	_, err = repo.UpdateByID(context.Background(), 1, &models.Listing{
		Title:       "New",
		Description: "Desc",
		Category:    "Books",
		Price:       12,
	})
	if err == nil || !strings.Contains(err.Error(), "update listing by id") {
		t.Fatalf("expected update listing by id wrapped error, got %v", err)
	}

	_, err = repo.DeleteByID(context.Background(), 1)
	if err == nil || !strings.Contains(err.Error(), "delete listing by id") {
		t.Fatalf("expected delete listing by id wrapped error, got %v", err)
	}
}

func TestReportRepositoryCreateReturnsErrorWhenDBClosed(t *testing.T) {
	repo := NewPostgresReportRepository(newClosedPostgresDB(t))

	_, err := repo.Create(context.Background(), &models.Report{
		ListingID:      1,
		ReporterUserID: 2,
		Reason:         "Spam",
	})
	if err == nil || !strings.Contains(err.Error(), "create report") {
		t.Fatalf("expected create report wrapped error, got %v", err)
	}
}

func TestContactRequestRepositoryMethodsReturnErrorWhenDBClosed(t *testing.T) {
	repo := NewPostgresContactRequestRepository(newClosedPostgresDB(t))

	_, err := repo.Create(context.Background(), &models.ContactRequest{
		ListingID: 1,
		BuyerID:   2,
		SellerID:  3,
		Message:   "Interested",
		Status:    models.ContactRequestStatusPending,
	})
	if err == nil || !strings.Contains(err.Error(), "create contact request") {
		t.Fatalf("expected create contact request wrapped error, got %v", err)
	}

	_, err = repo.ListReceivedBySellerID(context.Background(), 3)
	if err == nil || !strings.Contains(err.Error(), "list received contact requests by seller id") {
		t.Fatalf("expected list contact requests wrapped error, got %v", err)
	}
}

func TestListingImageRepositoryCreateManyReturnsErrorWhenDBClosed(t *testing.T) {
	repo := NewPostgresListingImageRepository(newClosedPostgresDB(t))

	_, err := repo.CreateMany(context.Background(), 1, []string{"/uploads/a.jpg"})
	if err == nil || !strings.Contains(err.Error(), "begin listing images transaction") {
		t.Fatalf("expected begin transaction wrapped error, got %v", err)
	}
}
