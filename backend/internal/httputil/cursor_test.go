package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCursorRoundTrip(t *testing.T) {
	const (
		sortKey = "2026-09-21T12:00:00Z"
		id      = "0f8fad5b-d9cb-469f-a165-70867728950e"
	)

	encoded := EncodeCursor(sortKey, id)
	if encoded == sortKey+"|"+id {
		t.Error("cursor should be opaque, not the raw key pair")
	}

	gotSort, gotID, ok := DecodeCursor(encoded)
	if !ok {
		t.Fatalf("DecodeCursor(%q) not ok", encoded)
	}
	if gotSort != sortKey || gotID != id {
		t.Errorf("DecodeCursor() = %q, %q; want %q, %q", gotSort, gotID, sortKey, id)
	}
}

func TestDecodeCursorRejectsGarbage(t *testing.T) {
	for _, in := range []string{
		"not base64 !!",
		EncodeCursor("", "id"),
		EncodeCursor("sort", ""),
		"aGVsbG8", // valid base64, no separator
	} {
		if _, _, ok := DecodeCursor(in); ok {
			t.Errorf("DecodeCursor(%q) = ok, want rejected", in)
		}
	}
}

func TestCursorRequestModes(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantKeyset bool
		wantOK     bool
		wantStatus int
	}{
		{name: "absent means offset", query: "", wantKeyset: false, wantOK: true},
		{name: "empty means first keyset page", query: "?cursor=", wantKeyset: true, wantOK: true},
		{
			name:       "valid cursor",
			query:      "?cursor=" + EncodeCursor("2026-09-21T12:00:00Z", "id-1"),
			wantKeyset: true, wantOK: true,
		},
		{name: "garbage is a 400", query: "?cursor=not-a-real-cursor", wantOK: false, wantStatus: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/list"+tc.query, nil)

			keyset, _, _, ok := Cursor(rec, r)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if keyset != tc.wantKeyset {
				t.Errorf("keyset = %v, want %v", keyset, tc.wantKeyset)
			}
			if !tc.wantOK && rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestNextCursorStopsOnShortPage(t *testing.T) {
	key := func(s string) (string, string) { return "t-" + s, s }

	if got := NextCursor([]string{"a", "b"}, 3, key); got != "" {
		t.Errorf("short page cursor = %q, want empty", got)
	}
	if got := NextCursor([]string{}, 3, key); got != "" {
		t.Errorf("empty page cursor = %q, want empty", got)
	}

	got := NextCursor([]string{"a", "b", "c"}, 3, key)
	sort, id, ok := DecodeCursor(got)
	if !ok || sort != "t-c" || id != "c" {
		t.Errorf("full page cursor decoded to %q, %q (ok=%v); want the last item", sort, id, ok)
	}
}

// TestWritePaginatedOmitsNextCursor pins the additive part of the contract:
// offset clients must keep seeing exactly the four keys they always have.
func TestWritePaginatedOmitsNextCursor(t *testing.T) {
	rec := httptest.NewRecorder()
	WritePaginated(rec, []string{"a"}, 2, 10, 42)

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 4 {
		t.Errorf("envelope has %d keys (%v), want exactly items/total/page/pageSize", len(got), got)
	}
	if _, present := got["nextCursor"]; present {
		t.Error("nextCursor must be omitted when empty")
	}

	rec = httptest.NewRecorder()
	WritePaginatedCursor(rec, []string{"a"}, 1, 10, 42, "abc")
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode cursor body: %v", err)
	}
	if got["nextCursor"] != "abc" {
		t.Errorf("nextCursor = %v, want abc", got["nextCursor"])
	}
}
