package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"
)

type authHandlerUserRepoStub struct {
	createFn     func(ctx context.Context, user *models.User) (*models.User, error)
	getByEmailFn func(ctx context.Context, email string) (*models.User, error)
	getByIDFn    func(ctx context.Context, id int64) (*models.User, error)
}

func (r *authHandlerUserRepoStub) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if r.createFn == nil {
		return &models.User{ID: 1, Email: user.Email}, nil
	}
	return r.createFn(ctx, user)
}

func (r *authHandlerUserRepoStub) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if r.getByEmailFn == nil {
		return nil, repository.ErrUserNotFound
	}
	return r.getByEmailFn(ctx, email)
}

func (r *authHandlerUserRepoStub) GetByID(ctx context.Context, id int64) (*models.User, error) {
	if r.getByIDFn == nil {
		return nil, repository.ErrUserNotFound
	}
	return r.getByIDFn(ctx, id)
}

type authHandlerParserStub struct {
	userID int64
}

func (p *authHandlerParserStub) ParseToken(tokenString string) (int64, error) {
	return p.userID, nil
}

func newAuthHandlerForTest(repo *authHandlerUserRepoStub) *AuthHandler {
	return NewAuthHandler(services.NewAuthService(repo, "test-secret"))
}

func decodeHandlerResponseBody(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("unmarshal response body: %v", err)
	}
	return out
}

func TestRegisterMethodNotAllowed(t *testing.T) {
	h := newAuthHandlerForTest(&authHandlerUserRepoStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/register", nil)
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestRegisterSuccess(t *testing.T) {
	h := newAuthHandlerForTest(&authHandlerUserRepoStub{
		createFn: func(ctx context.Context, user *models.User) (*models.User, error) {
			return &models.User{
				ID:           55,
				FullName:     user.FullName,
				Email:        user.Email,
				PasswordHash: user.PasswordHash,
				University:   user.University,
			}, nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"full_name":"Test User","email":"test@example.edu","password":"secret123","university":"UF"}`))
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	body := decodeHandlerResponseBody(t, rec.Body.Bytes())
	if body["success"] != true {
		t.Fatalf("expected success=true, got %v", body["success"])
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	h := newAuthHandlerForTest(&authHandlerUserRepoStub{
		getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return nil, repository.ErrUserNotFound
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"missing@example.edu","password":"secret123"}`))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMeUnauthorizedWithoutContextUser(t *testing.T) {
	h := newAuthHandlerForTest(&authHandlerUserRepoStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestMeSuccess(t *testing.T) {
	h := newAuthHandlerForTest(&authHandlerUserRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.User, error) {
			return &models.User{ID: id, FullName: "User One", Email: "user1@example.edu"}, nil
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	middleware.Auth(&authHandlerParserStub{userID: 10})(http.HandlerFunc(h.Me)).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := decodeHandlerResponseBody(t, rec.Body.Bytes())
	if body["success"] != true {
		t.Fatalf("expected success=true")
	}
}
