package connectors

import (
	"encoding/json"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

func TestParseScheduleUpdates(t *testing.T) {
	raw := func(s string) json.RawMessage { return json.RawMessage(s) }

	// onlyField asserts errs is a single FieldError on the named field and
	// returns its message, so the invalid-input cases stay one-liners.
	onlyField := func(t *testing.T, errs []httputil.FieldError, field string) string {
		t.Helper()
		if len(errs) != 1 || errs[0].Field != field {
			t.Fatalf("errs = %+v, want one error on %q", errs, field)
		}
		return errs[0].Msg
	}

	t.Run("absent fields yield no updates", func(t *testing.T) {
		u, errs := parseScheduleUpdates(nil, nil, nil)
		if len(errs) != 0 || len(u) != 0 {
			t.Fatalf("got %v, %+v", u, errs)
		}
	})
	t.Run("null clears", func(t *testing.T) {
		u, errs := parseScheduleUpdates(raw("null"), raw("null"), raw("null"))
		if len(errs) != 0 {
			t.Fatalf("errs = %+v", errs)
		}
		if v, ok := u["user_expires_at"]; !ok || v != nil {
			t.Errorf("user_expires_at = %v, %v", v, ok)
		}
		if _, ok := u["schedule_seconds"]; !ok {
			t.Error("schedule_seconds missing")
		}
		if _, ok := u["rotation_max_age_days"]; !ok {
			t.Error("rotation_max_age_days missing")
		}
	})
	t.Run("invalid schedule", func(t *testing.T) {
		_, errs := parseScheduleUpdates(raw(`"x"`), nil, nil)
		if msg := onlyField(t, errs, "scheduleSeconds"); msg == "" {
			t.Error("empty msg")
		}
	})
	t.Run("invalid userExpiresAt", func(t *testing.T) {
		_, errs := parseScheduleUpdates(nil, raw(`1`), nil)
		if msg := onlyField(t, errs, "userExpiresAt"); msg == "" {
			t.Error("empty msg")
		}
	})
	t.Run("invalid rotation", func(t *testing.T) {
		_, errs := parseScheduleUpdates(nil, nil, raw(`"x"`))
		if msg := onlyField(t, errs, "rotationMaxAgeDays"); msg == "" {
			t.Error("empty msg")
		}
	})
	t.Run("rotation values validated per field", func(t *testing.T) {
		_, errs := parseScheduleUpdates(nil, raw(`"not-a-time"`), raw(`0`))
		if len(errs) != 2 {
			t.Fatalf("errs = %+v, want one per field", errs)
		}
	})
}
