// Package httputil provides HTTP response helpers used across the API layer.
package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/WiseLabz/wiselabz/internal/storeerr"
)

// ErrorResponse is the standard error envelope returned by all API endpoints.
// Matches the OpenAPI Error schema.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// JSON writes a JSON response with the given status code and body.
func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			slog.Error("failed to encode JSON response", "error", err)
		}
	}
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes a structured error response.
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// Errorf writes a 500 Internal Server Error with a generic message.
// Use for unexpected errors; the caller should log the actual error.
func Errorf(w http.ResponseWriter, err error) {
	slog.Error("internal server error", "error", err)
	Error(w, http.StatusInternalServerError, "internal_error", "An internal error occurred")
}

// HandleStoreError maps a store sentinel error to a structured HTTP response:
// ErrNotFound -> 404, ErrConflict/ErrVersionConflict -> 409, ErrUnauthorized
// -> 401, ErrForbidden -> 403. Anything else falls through to Errorf (500).
// Handlers that need a resource-specific message should keep their own
// errors.Is check; this is for the generic case.
func HandleStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storeerr.ErrNotFound):
		Error(w, http.StatusNotFound, "not_found", "Resource not found")
	case errors.Is(err, storeerr.ErrVersionConflict):
		Error(w, http.StatusConflict, "version_conflict", "Version conflict")
	case errors.Is(err, storeerr.ErrConflict):
		Error(w, http.StatusConflict, "conflict", "Resource already exists")
	case errors.Is(err, storeerr.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
	case errors.Is(err, storeerr.ErrForbidden):
		Error(w, http.StatusForbidden, "forbidden", "Forbidden")
	default:
		Errorf(w, err)
	}
}

// MaxJSONBodyBytes caps JSON request bodies read via DecodeJSON.
const MaxJSONBodyBytes = 4 << 20

// DecodeJSON decodes the request body into a T. On failure it writes the
// standard 400 "invalid_request" response and returns the zero value with
// ok=false; callers should return immediately when ok is false.
func DecodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var v T
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxJSONBodyBytes)).Decode(&v); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			Error(w, http.StatusRequestEntityTooLarge, "payload_too_large", "Request body too large")
			var zero T
			return zero, false
		}
		Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		var zero T
		return zero, false
	}
	return v, true
}

// PaginatedResponse wraps a paginated list response. Matches the AlertPage /
// ChangePage OpenAPI schemas: { items, total, page, pageSize }.
type PaginatedResponse struct {
	Items    any `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// WritePaginated writes a paginated response with pagination metadata.
func WritePaginated(w http.ResponseWriter, data any, page, pageSize, total int) {
	JSON(w, http.StatusOK, PaginatedResponse{
		Items:    data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

const (
	// DefaultPage is the fallback page number for paginated endpoints.
	DefaultPage = 1
	// DefaultPageSize is the default number of items per page.
	DefaultPageSize = 20
	// MaxPageSize is the maximum number of items per page.
	MaxPageSize = 100
	// MaxPage is the maximum page number, chosen so (MaxPage-1)*MaxPageSize
	// cannot overflow int and stays well within any reasonable dataset size.
	MaxPage = 1_000_000
	// MaxBulkIDs is the maximum number of IDs accepted by any bulk-action
	// endpoint (bulk-resolve, bulk-snooze, bulk-sync, bulk-reauth, bulk-restart).
	MaxBulkIDs = 500
)

// Paginate extracts pagination parameters from the request query string.
// Returns normalized page, pageSize, and offset values.
func Paginate(r *http.Request) (page, pageSize, offset int) {
	page = intQuery(r, "page", DefaultPage)
	pageSize = intQuery(r, "pageSize", DefaultPageSize)

	if page < 1 {
		page = 1
	}
	if page > MaxPage {
		page = MaxPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset = (page - 1) * pageSize
	return
}

func intQuery(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}
