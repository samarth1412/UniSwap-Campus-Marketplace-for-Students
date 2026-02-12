package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

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
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authService.Register(r.Context(), req)
	if err != nil {
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

	writeSuccess(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authService.Login(r.Context(), req)
	if err != nil {
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

	writeSuccess(w, http.StatusOK, result)
}
