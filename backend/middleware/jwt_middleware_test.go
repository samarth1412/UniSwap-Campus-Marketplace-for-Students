package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type parserStub struct {
	userID int64
	err    error
	token  string
}

func (p *parserStub) ParseToken(tokenString string) (int64, error) {
	p.token = tokenString
	return p.userID, p.err
}

func TestAuthMissingHeader(t *testing.T) {
	p := &parserStub{}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	Auth(p)(next).ServeHTTP(rec, req)

	if called {
		t.Fatalf("next handler should not be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthInvalidHeaderFormat(t *testing.T) {
	p := &parserStub{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "invalidtoken")
	rec := httptest.NewRecorder()
	Auth(p)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthInvalidToken(t *testing.T) {
	p := &parserStub{err: errors.New("bad token")}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer tkn")
	rec := httptest.NewRecorder()
	Auth(p)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthSuccessSetsContext(t *testing.T) {
	p := &parserStub{userID: 101}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		if !ok || userID != 101 {
			t.Fatalf("expected user id to be set in context")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer abc")
	rec := httptest.NewRecorder()
	Auth(p)(next).ServeHTTP(rec, req)

	if p.token != "abc" {
		t.Fatalf("expected parser to receive token")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestAuthErrorBodyShape(t *testing.T) {
	p := &parserStub{}
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	Auth(p)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rec, req)

	var payload map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if payload["success"] != false {
		t.Fatalf("expected success false")
	}
}
