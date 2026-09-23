package report

import (
	"testing"
	"time"
)

func TestComputeWindow_FirstRun(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)

	start, end, truncated := ComputeWindow(nil, now)

	if !end.Equal(now) {
		t.Errorf("end = %v, want %v", end, now)
	}
	wantStart := now.Add(-7 * 24 * time.Hour)
	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v (7d lookback)", start, wantStart)
	}
	if truncated {
		t.Errorf("truncated = true, want false for a 7d first-run window")
	}
}

func TestComputeWindow_Watermark(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	watermark := now.Add(-3 * 24 * time.Hour)

	start, end, truncated := ComputeWindow(&watermark, now)

	if !start.Equal(watermark) {
		t.Errorf("start = %v, want watermark %v", start, watermark)
	}
	if !end.Equal(now) {
		t.Errorf("end = %v, want %v", end, now)
	}
	if truncated {
		t.Errorf("truncated = true, want false for a 3d window")
	}
}

func TestComputeWindow_CappedAt31Days(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	// A definition that's been disabled for months: the naive window would
	// be ~90 days, which must be clamped to 31.
	watermark := now.Add(-90 * 24 * time.Hour)

	start, end, truncated := ComputeWindow(&watermark, now)

	if !truncated {
		t.Fatalf("truncated = false, want true when window exceeds 31 days")
	}
	wantStart := now.Add(-31 * 24 * time.Hour)
	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v (31d cap)", start, wantStart)
	}
	if !end.Equal(now) {
		t.Errorf("end = %v, want %v", end, now)
	}
}

func TestComputeWindow_ExactlyAtCap(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	watermark := now.Add(-31 * 24 * time.Hour)

	start, _, truncated := ComputeWindow(&watermark, now)

	if truncated {
		t.Errorf("truncated = true, want false when window is exactly 31 days")
	}
	if !start.Equal(watermark) {
		t.Errorf("start = %v, want watermark %v (no clamping needed)", start, watermark)
	}
}

func TestComputeWindow_ManualRunUsesLastScheduledWatermarkUnchanged(t *testing.T) {
	// Manual runs don't get special handling in ComputeWindow itself — the
	// "manual doesn't advance the watermark" rule (plan Q12) is enforced by
	// the caller never persisting a manual run's period_end as the new
	// watermark, not by this function. Calling it twice with the same
	// watermark (as a manual run followed by the next scheduled tick would)
	// must yield the same start each time.
	now1 := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	now2 := now1.Add(2 * time.Hour) // manual run triggered a bit later
	watermark := now1.Add(-2 * 24 * time.Hour)

	start1, _, _ := ComputeWindow(&watermark, now1)
	start2, _, _ := ComputeWindow(&watermark, now2)

	if !start1.Equal(start2) {
		t.Errorf("start1 = %v, start2 = %v; want equal since watermark is unchanged between calls", start1, start2)
	}
}
