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

func (h *UserHandler) UserListings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := parseUserListingsPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	actorUserID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if actorUserID != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	listings, err := h.listingService.GetByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch user listings")
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

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}
