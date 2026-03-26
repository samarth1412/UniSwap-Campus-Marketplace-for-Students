package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckGet(t *testing.T) {
	a := &app{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	a.healthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHealthCheckMethodNotAllowed(t *testing.T) {
	a := &app{}
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	a.healthCheck(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
