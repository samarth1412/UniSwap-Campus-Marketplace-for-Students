package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"uniswap-campus-marketplace/models"
	"uniswap-campus-marketplace/repository"
	"uniswap-campus-marketplace/services"
)

type WishlistHandler struct {
	wishlistService *services.WishlistService
}

func NewWishlistHandler(wishlistService *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{wishlistService: wishlistService}
}

func (h *WishlistHandler) Wishlist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createWishlist(w, r)
	case http.MethodGet:
		h.getWishlist(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *WishlistHandler) WishlistByID(w http.ResponseWriter, r *http.Request) {
	wishlistID, ok := parseWishlistPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.deleteWishlist(w, r, wishlistID)
}

func (h *WishlistHandler) createWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.wishlistService.Create(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		case errors.Is(err, repository.ErrWishlistAlreadyExists):
			writeError(w, http.StatusConflict, "listing already exists in wishlist")
		default:
			writeError(w, http.StatusInternalServerError, "failed to add wishlist item")
		}
		return
	}

	writeSuccess(w, http.StatusCreated, item)
}

func (h *WishlistHandler) getWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.wishlistService.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch wishlist")
		return
	}

	writeSuccess(w, http.StatusOK, items)
}

func (h *WishlistHandler) deleteWishlist(w http.ResponseWriter, r *http.Request, wishlistID int64) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.wishlistService.Delete(r.Context(), userID, wishlistID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist item not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to delete wishlist item")
		}
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{
		"message": "wishlist item removed successfully",
	})
}

func parseWishlistPath(path string) (int64, bool) {
	trimmed := strings.TrimPrefix(path, "/api/wishlist/")
	if trimmed == path || trimmed == "" {
		return 0, false
	}

	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) != 1 || parts[0] == "" {
		return 0, false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
