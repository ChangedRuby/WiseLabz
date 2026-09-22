package httputil

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// CursorParam is the query parameter that opts a list endpoint into keyset
// (cursor) pagination. Its value is opaque to clients: it is only ever echoed
// back from a previous response's nextCursor.
const CursorParam = "cursor"

// cursorSep separates the sort key from the row id inside a cursor. It cannot
// appear in either half (both are an RFC3339 timestamp and a UUID).
const cursorSep = "\x00"

// EncodeCursor packs a (sortKey, id) keyset position into the opaque string
// handed to clients as nextCursor.
func EncodeCursor(sortKey, id string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(sortKey + cursorSep + id))
}

// DecodeCursor unpacks a cursor produced by EncodeCursor. ok is false for
// anything that is not a cursor this server issued, so callers can answer 400
// rather than silently returning a wrong page.
func DecodeCursor(s string) (sortKey, id string, ok bool) {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", "", false
	}
	sortKey, id, found := strings.Cut(string(raw), cursorSep)
	if !found || sortKey == "" || id == "" {
		return "", "", false
	}
	return sortKey, id, true
}

// Cursor reads the keyset position from the request.
//
// Keyset mode is opt-in and driven by the *presence* of ?cursor=, not its
// value: `?cursor=` (empty) asks for the first keyset page, and a non-empty
// value asks for the page after that position. Requests with no cursor
// parameter at all keep the historical offset behaviour untouched.
//
// On a malformed cursor it writes the standard 400 and returns ok=false;
// callers should return immediately.
func Cursor(w http.ResponseWriter, r *http.Request) (keyset bool, sortKey, id string, ok bool) {
	if !r.URL.Query().Has(CursorParam) {
		return false, "", "", true
	}
	raw := r.URL.Query().Get(CursorParam)
	if raw == "" {
		return true, "", "", true
	}
	sortKey, id, valid := DecodeCursor(raw)
	if !valid {
		Error(w, http.StatusBadRequest, "invalid_request", "Invalid cursor")
		return false, "", "", false
	}
	return true, sortKey, id, true
}

// NextCursor returns the cursor a client passes back to fetch the page after
// items, or "" when this page is the last one. key extracts the sort-column
// value and id of one item. A short page means there is nothing after it, so
// no cursor is issued and clients stop.
func NextCursor[T any](items []T, pageSize int, key func(T) (sortKey, id string)) string {
	if len(items) == 0 || len(items) < pageSize {
		return ""
	}
	return EncodeCursor(key(items[len(items)-1]))
}

// NextCursorHeader carries the next keyset cursor for list endpoints whose
// published response body is a bare array, leaving no envelope to put
// nextCursor in.
const NextCursorHeader = "X-Next-Cursor"
