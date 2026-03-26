package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"uniswap-campus-marketplace/services"
)

type UserHandler struct {
	listingService *services.ListingService
}

func NewUserHandler(listingService *services.ListingService) *UserHandler {
	return &UserHandler{listingService: listingService}
}

func (h *UserHandler) UserRoutes(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUserListingsPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.getUserListings(w, r, userID)
}

func (h *UserHandler) getUserListings(w http.ResponseWriter, r *http.Request, userID int64) {
	listings, err := h.listingService.GetByUserID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to fetch user listings")
		}
		return
	}

	writeSuccess(w, http.StatusOK, listings)
}

func parseUserListingsPath(path string) (int64, bool) {
	trimmed := strings.TrimPrefix(path, "/api/users/")
	if trimmed == path || trimmed == "" {
		return 0, false
	}

	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) != 2 || parts[1] != "listings" || parts[0] == "" {
		return 0, false
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || userID <= 0 {
		return 0, false
	}

	return userID, true
}
