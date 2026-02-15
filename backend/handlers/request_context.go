package handlers

import (
	"net/http"
	"strconv"
)

// userIDFromContext expects auth middleware to inject "user_id" into context.
func userIDFromContext(r *http.Request) (int64, bool) {
	value := r.Context().Value("user_id")
	if value == nil {
		return 0, false
	}

	switch v := value.(type) {
	case int64:
		return v, v > 0
	case int:
		id := int64(v)
		return id, id > 0
	case float64:
		id := int64(v)
		return id, id > 0
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			return 0, false
		}
		return id, true
	default:
		return 0, false
	}
}
