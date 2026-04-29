package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListingQueryFromRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/listings?search=desk&category=Furniture&min_price=10&max_price=50&page=2&limit=5", nil)
	query, err := listingQueryFromRequest(req)
	if err != nil {
		t.Fatalf("parse listing query: %v", err)
	}
	if query.Keyword != "desk" || query.Category != "Furniture" || query.Page != 2 || query.Limit != 5 {
		t.Fatalf("unexpected query: %+v", query)
	}
	if query.MinPrice == nil || *query.MinPrice != 10 || query.MaxPrice == nil || *query.MaxPrice != 50 {
		t.Fatalf("unexpected price filters: %+v", query)
	}
}

func TestListingQueryFromRequestValidation(t *testing.T) {
	cases := []string{
		"/api/listings?page=0",
		"/api/listings?limit=0",
		"/api/listings?min_price=-1",
		"/api/listings?max_price=-1",
		"/api/listings?min_price=50&max_price=10",
	}
	for _, path := range cases {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			if _, err := listingQueryFromRequest(req); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}
