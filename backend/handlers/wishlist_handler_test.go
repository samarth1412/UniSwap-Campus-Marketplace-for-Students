package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"
)

type wishlistRepoForHandlerTest struct {
	createFn     func(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error)
	deleteFn     func(ctx context.Context, wishlistID, userID int64) error
	listByUserFn func(ctx context.Context, userID int64) ([]models.WishlistListing, error)
}

func (r *wishlistRepoForHandlerTest) Create(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	if r.createFn == nil {
		return &models.WishlistItem{ID: 1, UserID: item.UserID, ListingID: item.ListingID}, nil
	}
	return r.createFn(ctx, item)
}

func (r *wishlistRepoForHandlerTest) Delete(ctx context.Context, wishlistID, userID int64) error {
	if r.deleteFn == nil {
		return nil
	}
	return r.deleteFn(ctx, wishlistID, userID)
}

func (r *wishlistRepoForHandlerTest) ListByUser(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
	if r.listByUserFn == nil {
		return []models.WishlistListing{}, nil
	}
	return r.listByUserFn(ctx, userID)
}

func newWishlistHandlerForTest(wishlistRepo *wishlistRepoForHandlerTest, listingRepo *listingRepoForHandlerTest) *WishlistHandler {
	if wishlistRepo == nil {
		wishlistRepo = &wishlistRepoForHandlerTest{}
	}
	if listingRepo == nil {
		listingRepo = &listingRepoForHandlerTest{}
	}
	return NewWishlistHandler(services.NewWishlistService(wishlistRepo, listingRepo))
}

func TestParseWishlistPath(t *testing.T) {
	id, ok := parseWishlistPath("/api/wishlist/12")
	if !ok || id != 12 {
		t.Fatalf("expected wishlist path parse success")
	}
	if _, ok := parseWishlistPath("/api/wishlist/not-a-number"); ok {
		t.Fatalf("expected invalid wishlist path to fail")
	}
}

func TestWishlistMethodNotAllowed(t *testing.T) {
	h := newWishlistHandlerForTest(nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/wishlist", nil)
	rec := httptest.NewRecorder()

	h.Wishlist(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestWishlistCreateUnauthorized(t *testing.T) {
	h := newWishlistHandlerForTest(nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/wishlist", strings.NewReader(`{"listing_id":1}`))
	rec := httptest.NewRecorder()

	h.Wishlist(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestWishlistCreateSuccess(t *testing.T) {
	h := newWishlistHandlerForTest(nil, &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id}, nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/wishlist", strings.NewReader(`{"listing_id":5}`))
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 2})(http.HandlerFunc(h.Wishlist)).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestWishlistCreateConflict(t *testing.T) {
	h := newWishlistHandlerForTest(&wishlistRepoForHandlerTest{
		createFn: func(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
			return nil, repository.ErrWishlistAlreadyExists
		},
	}, &listingRepoForHandlerTest{
		getByIDFn: func(ctx context.Context, id int64) (*models.Listing, error) {
			return &models.Listing{ID: id}, nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/wishlist", strings.NewReader(`{"listing_id":5}`))
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 2})(http.HandlerFunc(h.Wishlist)).ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}

func TestWishlistGetSuccess(t *testing.T) {
	h := newWishlistHandlerForTest(&wishlistRepoForHandlerTest{
		listByUserFn: func(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
			return []models.WishlistListing{{WishlistID: 1, Listing: models.Listing{ID: 4, UserID: userID}}}, nil
		},
	}, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/wishlist", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 6})(http.HandlerFunc(h.Wishlist)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestWishlistDeleteNotFound(t *testing.T) {
	h := newWishlistHandlerForTest(&wishlistRepoForHandlerTest{
		deleteFn: func(ctx context.Context, wishlistID, userID int64) error {
			return repository.ErrWishlistNotFound
		},
	}, nil)
	req := httptest.NewRequest(http.MethodDelete, "/api/wishlist/9", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 6})(http.HandlerFunc(h.WishlistByID)).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
