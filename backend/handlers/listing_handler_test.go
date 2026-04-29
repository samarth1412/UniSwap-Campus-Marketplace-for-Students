package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/services"
)

type listingRepoForHandlerTest struct {
	createFn  func(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	getAllFn  func(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error)
	getByUser func(ctx context.Context, userID int64) ([]models.Listing, error)
	getByIDFn func(ctx context.Context, id int64) (*models.Listing, error)
	updateFn  func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	deleteFn  func(ctx context.Context, id int64) (int64, error)
}

func (r *listingRepoForHandlerTest) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	if r.createFn == nil {
		return &models.Listing{}, nil
	}
	return r.createFn(ctx, listing)
}

func (r *listingRepoForHandlerTest) GetAll(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error) {
	if r.getAllFn == nil {
		return &models.PaginatedListings{
			Items:      []models.Listing{},
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      0,
			TotalPages: 1,
		}, nil
	}
	return r.getAllFn(ctx, query)
}

func (r *listingRepoForHandlerTest) GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error) {
	if r.getByUser == nil {
		return []models.Listing{}, nil
	}
	return r.getByUser(ctx, userID)
}

func (r *listingRepoForHandlerTest) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	if r.getByIDFn == nil {
		return &models.Listing{ID: id, UserID: 1}, nil
	}
	return r.getByIDFn(ctx, id)
}

func (r *listingRepoForHandlerTest) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	if r.updateFn == nil {
		return &models.Listing{ID: id, UserID: 1}, nil
	}
	return r.updateFn(ctx, id, listing)
}

func (r *listingRepoForHandlerTest) DeleteByID(ctx context.Context, id int64) (int64, error) {
	if r.deleteFn == nil {
		return 1, nil
	}
	return r.deleteFn(ctx, id)
}

type reportRepoForHandlerTest struct{}

func (r *reportRepoForHandlerTest) Create(ctx context.Context, report *models.Report) (*models.Report, error) {
	return &models.Report{ID: 1}, nil
}

type contactRequestRepoForHandlerTest struct {
	createFn func(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error)
}

func (r *contactRequestRepoForHandlerTest) Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
	if r.createFn == nil {
		return &models.ContactRequest{ID: 1}, nil
	}
	return r.createFn(ctx, request)
}

func (r *contactRequestRepoForHandlerTest) ListBySellerID(ctx context.Context, sellerID int64) ([]models.ContactRequest, error) {
	return []models.ContactRequest{}, nil
}

type parserForHandlerTest struct {
	userID int64
}

func (p *parserForHandlerTest) ParseToken(tokenString string) (int64, error) {
	return p.userID, nil
}

func newListingHandlerForTest(repo *listingRepoForHandlerTest) *ListingHandler {
	listingSvc := services.NewListingService(repo)
	reportSvc := services.NewReportService(&reportRepoForHandlerTest{}, repo)
	contactSvc := services.NewContactRequestService(&contactRequestRepoForHandlerTest{}, repo)
	return NewListingHandler(listingSvc, reportSvc, contactSvc)
}

func TestParseListingPath(t *testing.T) {
	id, isReport, isContact, ok := parseListingPath("/api/listings/10")
	if !ok || isReport || isContact || id != 10 {
		t.Fatalf("expected listing route parse success")
	}

	id, isReport, isContact, ok = parseListingPath("/api/listings/10/report")
	if !ok || !isReport || isContact || id != 10 {
		t.Fatalf("expected report route parse success")
	}

	id, isReport, isContact, ok = parseListingPath("/api/listings/10/contact")
	if !ok || isReport || !isContact || id != 10 {
		t.Fatalf("expected contact route parse success")
	}

	_, _, _, ok = parseListingPath("/api/listings/x")
	if ok {
		t.Fatalf("expected parse failure for invalid id")
	}
}

func TestListingByIDRoutesMethodNotAllowed(t *testing.T) {
	h := newListingHandlerForTest(&listingRepoForHandlerTest{})
	req := httptest.NewRequest(http.MethodPatch, "/api/listings/1", nil)
	rec := httptest.NewRecorder()

	h.ListingByIDRoutes(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestUpdateListingUnauthorized(t *testing.T) {
	h := newListingHandlerForTest(&listingRepoForHandlerTest{})
	req := httptest.NewRequest(http.MethodPut, "/api/listings/1", strings.NewReader(`{"title":"A","description":"B","price":1,"category":"Books"}`))
	rec := httptest.NewRecorder()

	h.ListingByIDRoutes(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUpdateListingForbidden(t *testing.T) {
	repo := &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 999}, nil
		},
		updateFn: func(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
			t.Fatalf("update should not be called for non-owner")
			return nil, nil
		},
	}
	h := newListingHandlerForTest(repo)

	req := httptest.NewRequest(http.MethodPut, "/api/listings/1", strings.NewReader(`{"title":"A","description":"B","price":1,"category":"Books"}`))
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 1})(http.HandlerFunc(h.ListingByIDRoutes)).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestDeleteListingSuccess(t *testing.T) {
	repo := &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 1}, nil
		},
		deleteFn: func(ctx context.Context, id int64) (int64, error) {
			return 1, nil
		},
	}
	h := newListingHandlerForTest(repo)
	req := httptest.NewRequest(http.MethodDelete, "/api/listings/1", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 1})(http.HandlerFunc(h.ListingByIDRoutes)).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestCreateContactRequestUnauthorized(t *testing.T) {
	h := newListingHandlerForTest(&listingRepoForHandlerTest{})
	req := httptest.NewRequest(http.MethodPost, "/api/listings/1/contact", strings.NewReader(`{"message":"hello"}`))
	rec := httptest.NewRecorder()

	h.ListingByIDRoutes(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestCreateContactRequestForbidden(t *testing.T) {
	repo := &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 1}, nil
		},
	}
	h := newListingHandlerForTest(repo)
	req := httptest.NewRequest(http.MethodPost, "/api/listings/1/contact", strings.NewReader(`{"message":"hello"}`))
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 1})(http.HandlerFunc(h.ListingByIDRoutes)).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestCreateContactRequestCreated(t *testing.T) {
	repo := &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id, UserID: 9}, nil
		},
	}
	h := newListingHandlerForTest(repo)
	req := httptest.NewRequest(http.MethodPost, "/api/listings/1/contact", strings.NewReader(`{"message":"hello"}`))
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 1})(http.HandlerFunc(h.ListingByIDRoutes)).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}
