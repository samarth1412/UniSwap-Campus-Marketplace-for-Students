package handlers

import (
	"net/http"

	"uniswap-campus-marketplace/middleware"
)

// userIDFromContext reads the authenticated user ID injected by jwt middleware.
func userIDFromContext(r *http.Request) (int64, bool) {
	return middleware.UserIDFromContext(r.Context())
}
