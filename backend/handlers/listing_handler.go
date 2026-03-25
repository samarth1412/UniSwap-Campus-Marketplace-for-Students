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

type ListingHandler struct {
	listingService *services.ListingService
	reportService  *services.ReportService
}

func NewListingHandler(listingService *services.ListingService, reportService *services.ReportService) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
		reportService:  reportService,
	}
}

func (h *ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createListing(w, r)
	case http.MethodGet:
		h.getAllListings(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ListingHandler) ListingByIDRoutes(w http.ResponseWriter, r *http.Request) {
	listingID, isReportRoute, ok := parseListingPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	if isReportRoute {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.reportListing(w, r, listingID)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getListingByID(w, r, listingID)
	case http.MethodPut:
		h.updateListing(w, r, listingID)
	case http.MethodDelete:
		h.deleteListing(w, r, listingID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *ListingHandler) createListing(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.listingService.Create(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to create listing")
		}
		return
	}

	writeSuccess(w, http.StatusCreated, result)
}

func (h *ListingHandler) getAllListings(w http.ResponseWriter, r *http.Request) {
	page, err := parsePositiveIntQuery(r, "page")
	if err != nil {
		writeError(w, http.StatusBadRequest, "page must be a positive integer")
		return
	}

	limit, err := parsePositiveIntQuery(r, "limit")
	if err != nil {
		writeError(w, http.StatusBadRequest, "limit must be a positive integer")
		return
	}

	params := models.ListingListParams{
		Search: r.URL.Query().Get("search"),
		Page:   page,
		Limit:  limit,
	}

	listings, err := h.listingService.GetAll(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch listings")
		return
	}

	writeSuccess(w, http.StatusOK, listings)
}

func (h *ListingHandler) getListingByID(w http.ResponseWriter, r *http.Request, listingID int64) {
	listing, err := h.listingService.GetByID(r.Context(), listingID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to fetch listing")
		}
		return
	}

	writeSuccess(w, http.StatusOK, listing)
}

func (h *ListingHandler) updateListing(w http.ResponseWriter, r *http.Request, listingID int64) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.listingService.Update(r.Context(), userID, listingID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to update listing")
		}
		return
	}

	writeSuccess(w, http.StatusOK, updated)
}

func (h *ListingHandler) deleteListing(w http.ResponseWriter, r *http.Request, listingID int64) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.listingService.Delete(r.Context(), userID, listingID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, "forbidden")
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to delete listing")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ListingHandler) reportListing(w http.ResponseWriter, r *http.Request, listingID int64) {
	userID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.reportService.Create(r.Context(), listingID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrListingNotFound):
			writeError(w, http.StatusNotFound, "listing not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to report listing")
		}
		return
	}

	writeSuccess(w, http.StatusCreated, report)
}

func parseListingPath(path string) (int64, bool, bool) {
	trimmed := strings.TrimPrefix(path, "/api/listings/")
	if trimmed == path || trimmed == "" {
		return 0, false, false
	}

	parts := strings.Split(strings.Trim(trimmed, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, false, false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		return 0, false, false
	}

	if len(parts) == 1 {
		return id, false, true
	}
	if len(parts) == 2 && parts[1] == "report" {
		return id, true, true
	}

	return 0, false, false
}

func parsePositiveIntQuery(r *http.Request, key string) (int, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid positive integer")
	}

	return parsed, nil
}
