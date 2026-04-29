package handlers

import (
	"errors"
	"net/http"

	"uniswap-campus-marketplace/services"
)

type ContactRequestHandler struct {
	contactRequestService *services.ContactRequestService
}

func NewContactRequestHandler(contactRequestService *services.ContactRequestService) *ContactRequestHandler {
	return &ContactRequestHandler{contactRequestService: contactRequestService}
}

func (h *ContactRequestHandler) Received(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sellerID, ok := userIDFromContext(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	requests, err := h.contactRequestService.ListReceivedBySellerID(r.Context(), sellerID)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch received contact requests")
		return
	}

	writeSuccess(w, http.StatusOK, requests)
}
