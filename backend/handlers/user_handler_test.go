package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/services"
)

func newUserHandlerForTest(repo *listingRepoForHandlerTest) *UserHandler {
	return NewUserHandler(services.NewListingService(repo))
}

func TestUserListingsUnauthorizedWithoutContextUser(t *testing.T) {
	h := newUserHandlerForTest(&listingRepoForHandlerTest{})
	req := httptest.NewRequest(http.MethodGet, "/api/users/1/listings", nil)
	rec := httptest.NewRecorder()

	h.UserListings(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestUserListingsForbiddenForDifferentUser(t *testing.T) {
	h := newUserHandlerForTest(&listingRepoForHandlerTest{})
	req := httptest.NewRequest(http.MethodGet, "/api/users/2/listings", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 1})(http.HandlerFunc(h.UserListings)).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestUserListingsOwnerSuccess(t *testing.T) {
	repo := &listingRepoForHandlerTest{
		getByUser: func(ctx context.Context, userID int64) ([]models.Listing, error) {
			return []models.Listing{{ID: 10, UserID: userID, Title: "Book"}}, nil
		},
	}
	h := newUserHandlerForTest(repo)
	req := httptest.NewRequest(http.MethodGet, "/api/users/3/listings", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForHandlerTest{userID: 3})(http.HandlerFunc(h.UserListings)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
