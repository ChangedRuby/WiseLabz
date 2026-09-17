package notifications

import (
	"testing"
	"time"
)

func TestDigestDue(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 8, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		cadence   string
		lastSent  string
		localHour int
		now       time.Time
		want      bool
	}{
		// cadence == "off" always returns false
		{
			name:      "off cadence",
			cadence:   "off",
			lastSent:  "",
			localHour: 8,
			now:       baseTime,
			want:      false,
		},
		// localHour != 8 returns false
		{
			name:      "wrong hour daily",
			cadence:   "daily",
			lastSent:  "",
			localHour: 7,
			now:       baseTime,
			want:      false,
		},
		{
			name:      "wrong hour weekly",
			cadence:   "weekly",
			lastSent:  "",
			localHour: 9,
			now:       baseTime,
			want:      false,
		},
		// daily cadence at hour 8 always returns true
		{
			name:      "daily at hour 8 no last sent",
			cadence:   "daily",
			lastSent:  "",
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		{
			name:      "daily at hour 8 with last sent",
			cadence:   "daily",
			lastSent:  baseTime.Add(-24 * time.Hour).Format(time.RFC3339),
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		// weekly with no lastSent returns true
		{
			name:      "weekly at hour 8 no last sent",
			cadence:   "weekly",
			lastSent:  "",
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		// weekly with lastSent 3 days ago returns false
		{
			name:      "weekly 3 days since last",
			cadence:   "weekly",
			lastSent:  baseTime.Add(-3 * 24 * time.Hour).Format(time.RFC3339),
			localHour: 8,
			now:       baseTime,
			want:      false,
		},
		// weekly with lastSent 8 days ago returns true
		{
			name:      "weekly 8 days since last",
			cadence:   "weekly",
			lastSent:  baseTime.Add(-8 * 24 * time.Hour).Format(time.RFC3339),
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		// weekly with lastSent exactly 7 days ago returns true
		{
			name:      "weekly exactly 7 days",
			cadence:   "weekly",
			lastSent:  baseTime.Add(-7 * 24 * time.Hour).Format(time.RFC3339),
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		// weekly with malformed lastSent returns true (error case)
		{
			name:      "weekly malformed last sent",
			cadence:   "weekly",
			lastSent:  "not-a-date",
			localHour: 8,
			now:       baseTime,
			want:      true,
		},
		// invalid cadence returns false
		{
			name:      "invalid cadence",
			cadence:   "never",
			lastSent:  "",
			localHour: 8,
			now:       baseTime,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := digestDue(tt.cadence, tt.lastSent, tt.localHour, tt.now)
			if got != tt.want {
				t.Errorf("digestDue(%q, %q, %d, %v) = %v, want %v",
					tt.cadence, tt.lastSent, tt.localHour, tt.now, got, tt.want)
			}
		})
	}
}
