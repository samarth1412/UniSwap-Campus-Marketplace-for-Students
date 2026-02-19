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

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.getListingByID(w, r, listingID)
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
	search := r.URL.Query().Get("search")
	listings, err := h.listingService.GetAll(r.Context(), search)
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

	// Allowed routes:
	// /api/listings/{id}
	// /api/listings/{id}/report
	if len(parts) == 1 {
		return id, false, true
	}
	if len(parts) == 2 && parts[1] == "report" {
		return id, true, true
	}

	return 0, false, false
}
