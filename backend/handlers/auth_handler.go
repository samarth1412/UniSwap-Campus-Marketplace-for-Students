package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"uniswap-campus-marketplace/middleware"
	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("auth_handler.register: request method=%s path=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("auth_handler.register: invalid method method=%s", r.Method)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("auth_handler.register: decode failed err=%v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	_, err := h.authService.Register(r.Context(), req)
	if err != nil {
		log.Printf("auth_handler.register: service failed email=%s err=%v", req.Email, err)
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrEmailAlreadyExists):
			writeError(w, http.StatusConflict, "email already exists")
		default:
			writeError(w, http.StatusInternalServerError, "failed to register user")
		}
		return
	}

	log.Printf("auth_handler.register: success email=%s", req.Email)
	writeSuccess(w, http.StatusCreated, map[string]string{
		"message": "registration successful",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("auth_handler.login: request method=%s path=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("auth_handler.login: invalid method method=%s", r.Method)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("auth_handler.login: decode failed err=%v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authService.Login(r.Context(), req)
	if err != nil {
		log.Printf("auth_handler.login: service failed email=%s err=%v", req.Email, err)
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
		default:
			writeError(w, http.StatusInternalServerError, "failed to login")
		}
		return
	}

	log.Printf("auth_handler.login: success email=%s", req.Email)
	writeSuccess(w, http.StatusOK, map[string]string{
		"token": result.Token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	log.Printf("auth_handler.me: request method=%s path=%s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		log.Printf("auth_handler.me: invalid method method=%s", r.Method)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID <= 0 {
		log.Printf("auth_handler.me: missing or invalid user id in context")
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("auth_handler.me: service failed user_id=%d err=%v", userID, err)
		if errors.Is(err, repository.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	log.Printf("auth_handler.me: success user_id=%d", userID)
	writeSuccess(w, http.StatusOK, user)
}
