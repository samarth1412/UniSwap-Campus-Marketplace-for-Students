package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"uniswap-campus-marketplace/middleware"
)

type parserForContextTest struct{}

func (p *parserForContextTest) ParseToken(tokenString string) (int64, error) {
	return 55, nil
}

func TestUserIDFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&parserForContextTest{})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r)
		if !ok || userID != 55 {
			t.Fatalf("expected user id 55 from context")
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
