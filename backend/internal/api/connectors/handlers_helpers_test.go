package connectors

import (
	"encoding/json"
	"testing"
)

func TestParseScheduleUpdates(t *testing.T) {
	raw := func(s string) json.RawMessage { return json.RawMessage(s) }

	t.Run("absent fields yield no updates", func(t *testing.T) {
		u, msg := parseScheduleUpdates(nil, nil, nil)
		if msg != "" || len(u) != 0 {
			t.Fatalf("got %v, %q", u, msg)
		}
	})
	t.Run("null clears", func(t *testing.T) {
		u, msg := parseScheduleUpdates(raw("null"), raw("null"), raw("null"))
		if msg != "" {
			t.Fatal(msg)
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
		if _, msg := parseScheduleUpdates(raw(`"x"`), nil, nil); msg != "Invalid scheduleSeconds" {
			t.Fatalf("msg = %q", msg)
		}
	})
	t.Run("invalid userExpiresAt", func(t *testing.T) {
		if _, msg := parseScheduleUpdates(nil, raw(`1`), nil); msg != "Invalid userExpiresAt" {
			t.Fatalf("msg = %q", msg)
		}
	})
	t.Run("invalid rotation", func(t *testing.T) {
		if _, msg := parseScheduleUpdates(nil, nil, raw(`"x"`)); msg != "Invalid rotationMaxAgeDays" {
			t.Fatalf("msg = %q", msg)
		}
	})
}
