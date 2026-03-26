package services

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

type authUserRepoStub struct {
	createFn     func(ctx context.Context, user *models.User) (*models.User, error)
	getByEmailFn func(ctx context.Context, email string) (*models.User, error)
	getByIDFn    func(ctx context.Context, id int64) (*models.User, error)
}

func (r *authUserRepoStub) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return r.createFn(ctx, user)
}

func (r *authUserRepoStub) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getByEmailFn(ctx, email)
}

func (r *authUserRepoStub) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return r.getByIDFn(ctx, id)
}

func TestAuthServiceRegisterValidation(t *testing.T) {
	svc := NewAuthService(&authUserRepoStub{}, "secret")
	_, err := svc.Register(context.Background(), models.RegisterRequest{
		FullName: "",
		Email:    "",
		Password: "",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestAuthServiceRegisterSuccess(t *testing.T) {
	repo := &authUserRepoStub{
		createFn: func(ctx context.Context, user *models.User) (*models.User, error) {
			if user.Email != "user@example.com" {
				t.Fatalf("expected normalized email, got %s", user.Email)
			}
			return &models.User{ID: 5, Email: user.Email}, nil
		},
	}
	svc := NewAuthService(repo, "secret")

	resp, err := svc.Register(context.Background(), models.RegisterRequest{
		FullName: "User",
		Email:    " User@Example.com ",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("expected token")
	}
}

func TestAuthServiceLoginInvalidCredentials(t *testing.T) {
	repo := &authUserRepoStub{
		getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return nil, repository.ErrUserNotFound
		},
	}
	svc := NewAuthService(repo, "secret")
	_, err := svc.Login(context.Background(), models.LoginRequest{
		Email:    "missing@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &authUserRepoStub{
		getByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				ID:           10,
				Email:        "user@example.com",
				PasswordHash: string(hash),
			}, nil
		},
	}
	svc := NewAuthService(repo, "secret")

	resp, err := svc.Login(context.Background(), models.LoginRequest{
		Email:    "user@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("login error: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("expected token")
	}
}

func TestAuthServiceParseToken(t *testing.T) {
	svc := NewAuthService(&authUserRepoStub{}, "secret")
	token, err := svc.generateToken(&models.User{ID: 77, Email: "x@y.com"})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	userID, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if userID != 77 {
		t.Fatalf("expected user id 77, got %d", userID)
	}
}

func TestAuthServiceGetUserByID(t *testing.T) {
	repo := &authUserRepoStub{
		getByIDFn: func(ctx context.Context, id int64) (*models.User, error) {
			return &models.User{ID: id, Email: "u@u.com"}, nil
		},
	}
	svc := NewAuthService(repo, "secret")

	user, err := svc.GetUserByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}
	if user.ID != 3 {
		t.Fatalf("expected id 3, got %d", user.ID)
	}
}
