package connectors

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestWriteConfigRejection covers the connector-config branch of the details
// envelope directly, rather than through a route: no registered connector
// schema declares a Pattern/MinLength/MaxLength/Options rule today, so
// connector.ValidateConfig cannot be made to fail over the wire.
func TestWriteConfigRejection(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantField string
		wantMsg   string
	}{
		{
			name:      "top-level url field",
			err:       &connector.ConfigValidationError{Field: "url", Message: "does not match the required format"},
			wantField: "url",
			wantMsg:   "does not match the required format",
		},
		{
			name:      "verify_tls folds back to its JSON spelling",
			err:       &connector.ConfigValidationError{Field: "verify_tls", Message: "must be one of [true false]"},
			wantField: "verifyTls",
			wantMsg:   "must be one of [true false]",
		},
		{
			name:      "schema key is nested under config",
			err:       &connector.ConfigValidationError{Field: "token_id", Message: "must be at least 8 characters"},
			wantField: "config.token_id",
			wantMsg:   "must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeConfigRejection(rec, tt.err)

			if rec.Code != 400 {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			var got struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Details []struct {
					Field string `json:"field"`
					Msg   string `json:"msg"`
				} `json:"details"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode %q: %v", rec.Body.String(), err)
			}
			if got.Code != "invalid_request" || got.Message != tt.err.Error() {
				t.Errorf("code/message = %q/%q, want invalid_request/%q", got.Code, got.Message, tt.err.Error())
			}
			if len(got.Details) != 1 {
				t.Fatalf("details = %+v, want one entry", got.Details)
			}
			if got.Details[0].Field != tt.wantField || got.Details[0].Msg != tt.wantMsg {
				t.Errorf("details[0] = %+v, want {%s %s}", got.Details[0], tt.wantField, tt.wantMsg)
			}
		})
	}

	// A malformed schema pattern is a server-side defect, not a field the
	// caller can correct, so it must not grow a details array.
	t.Run("unattributed error keeps the plain envelope", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeConfigRejection(rec, errors.New(`field "url": invalid schema pattern: bad`))

		var got map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode %q: %v", rec.Body.String(), err)
		}
		if _, ok := got["details"]; ok {
			t.Errorf("details present on an unattributed error: %v", got)
		}
	})
}
