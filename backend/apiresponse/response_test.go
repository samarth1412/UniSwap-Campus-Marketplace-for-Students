package apiresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteSuccess(rec, http.StatusCreated, map[string]string{"ok": "yes"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var payload APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if !payload.Success {
		t.Fatalf("expected success=true")
	}
}

func TestWriteErrorHelpers(t *testing.T) {
	cases := []struct {
		name   string
		status int
		write  func(http.ResponseWriter)
	}{
		{
			name:   "validation",
			status: http.StatusBadRequest,
			write: func(w http.ResponseWriter) {
				WriteValidationError(w, "bad input")
			},
		},
		{
			name:   "unauthorized",
			status: http.StatusUnauthorized,
			write: func(w http.ResponseWriter) {
				WriteUnauthorized(w, "unauthorized")
			},
		},
		{
			name:   "forbidden",
			status: http.StatusForbidden,
			write: func(w http.ResponseWriter) {
				WriteForbidden(w, "forbidden")
			},
		},
		{
			name:   "not_found",
			status: http.StatusNotFound,
			write: func(w http.ResponseWriter) {
				WriteNotFound(w, "not found")
			},
		},
		{
			name:   "internal",
			status: http.StatusInternalServerError,
			write: func(w http.ResponseWriter) {
				WriteInternalError(w, "internal")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.write(rec)
			if rec.Code != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, rec.Code)
			}

			var payload APIResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("unmarshal response: %v", err)
			}
			if payload.Success {
				t.Fatalf("expected success=false")
			}
			if payload.Error == "" {
				t.Fatalf("expected non-empty error")
			}
		})
	}
}
