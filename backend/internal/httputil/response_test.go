package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/storeerr"
)

// TestPaginateClampsOverflowingPage guards against issue #150: an
// unbounded ?page= value multiplied by pageSize overflowed int and
// produced a negative offset.
func TestPaginateClampsOverflowingPage(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?page=9223372036854775807&pageSize=100", nil)

	page, pageSize, offset := Paginate(r)

	if offset < 0 {
		t.Fatalf("Paginate() offset = %d, want non-negative", offset)
	}
	if page > MaxPage {
		t.Errorf("Paginate() page = %d, want <= MaxPage (%d)", page, MaxPage)
	}
	if pageSize != 100 {
		t.Errorf("Paginate() pageSize = %d, want 100", pageSize)
	}
}

func TestHandleStoreError(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"not found", storeerr.ErrNotFound, http.StatusNotFound, "not_found"},
		{"wrapped not found", fmt.Errorf("get: %w", storeerr.ErrNotFound), http.StatusNotFound, "not_found"},
		{"conflict", storeerr.ErrConflict, http.StatusConflict, "conflict"},
		{"version conflict", storeerr.ErrVersionConflict, http.StatusConflict, "version_conflict"},
		{"unauthorized", storeerr.ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{"forbidden", storeerr.ErrForbidden, http.StatusForbidden, "forbidden"},
		{"unknown", errors.New("boom"), http.StatusInternalServerError, "internal_error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			HandleStoreError(rec, tc.err)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			var body ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.Code != tc.wantCode {
				t.Errorf("code = %q, want %q", body.Code, tc.wantCode)
			}
		})
	}
}
