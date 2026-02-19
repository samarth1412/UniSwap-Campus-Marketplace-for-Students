package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrValidation = errors.New("validation failed")

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	log.Printf("auth_service.register: validating request email=%s", req.Email)
	if strings.TrimSpace(req.FullName) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		log.Printf("auth_service.register: validation failed missing required fields")
		return nil, fmt.Errorf("%w: full_name, email, and password are required", ErrValidation)
	}
	if len(req.Password) < 6 {
		log.Printf("auth_service.register: validation failed short password email=%s", req.Email)
		return nil, fmt.Errorf("%w: password must be at least 6 characters", ErrValidation)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("auth_service.register: password hashing failed email=%s err=%v", req.Email, err)
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		FullName:     strings.TrimSpace(req.FullName),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hashedPassword),
		University:   strings.TrimSpace(req.University),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("auth_service.register: user creation failed email=%s err=%v", user.Email, err)
		return nil, err
	}

	token, err := s.generateToken(createdUser)
	if err != nil {
		log.Printf("auth_service.register: token generation failed user_id=%d err=%v", createdUser.ID, err)
		return nil, err
	}

	log.Printf("auth_service.register: success user_id=%d email=%s", createdUser.ID, createdUser.Email)
	return &models.AuthResponse{
		Token: token,
		User:  *createdUser,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	log.Printf("auth_service.login: validating request email=%s", req.Email)
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		log.Printf("auth_service.login: validation failed missing email or password")
		return nil, fmt.Errorf("%w: email and password are required", ErrValidation)
	}

	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		log.Printf("auth_service.login: user lookup failed email=%s err=%v", req.Email, err)
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("auth_service.login: password mismatch user_id=%d email=%s", user.ID, user.Email)
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		log.Printf("auth_service.login: token generation failed user_id=%d err=%v", user.ID, err)
		return nil, err
	}

	log.Printf("auth_service.login: success user_id=%d email=%s", user.ID, user.Email)
	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	log.Printf("auth_service.generate_token: creating token user_id=%d", user.ID)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Printf("auth_service.generate_token: signing failed user_id=%d err=%v", user.ID, err)
		return "", fmt.Errorf("sign token: %w", err)
	}

	log.Printf("auth_service.generate_token: success user_id=%d", user.ID)
	return signedToken, nil
}

func (s *AuthService) ParseToken(tokenString string) (int64, error) {
	log.Printf("auth_service.parse_token: parsing token")
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			log.Printf("auth_service.parse_token: unexpected signing method method=%v", token.Method)
			return nil, ErrInvalidCredentials
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		log.Printf("auth_service.parse_token: invalid token err=%v", err)
		return 0, ErrInvalidCredentials
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("auth_service.parse_token: invalid claims type")
		return 0, ErrInvalidCredentials
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		log.Printf("auth_service.parse_token: user_id claim missing")
		return 0, ErrInvalidCredentials
	}

	userIDFloat, ok := userIDValue.(float64)
	if !ok || userIDFloat <= 0 || userIDFloat > math.MaxInt64 {
		log.Printf("auth_service.parse_token: invalid user_id claim value=%v", userIDValue)
		return 0, ErrInvalidCredentials
	}

	log.Printf("auth_service.parse_token: success user_id=%d", int64(userIDFloat))
	return int64(userIDFloat), nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	log.Printf("auth_service.get_user_by_id: fetching user user_id=%d", id)
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("auth_service.get_user_by_id: lookup failed user_id=%d err=%v", id, err)
		return nil, err
	}
	log.Printf("auth_service.get_user_by_id: success user_id=%d", id)
	return user, nil
}
