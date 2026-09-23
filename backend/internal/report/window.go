package report

import "time"

// defaultWindow is the lookback used when a definition has never produced a
// scheduled report (plan Q12).
const defaultWindow = 7 * 24 * time.Hour

// maxWindow caps how far back a single report can reach (plan Q22), so a
// long-dead definition that's re-enabled doesn't try to summarize months of
// history in one report. When the raw window would exceed this, the start
// is clamped and the report is marked truncated.
const maxWindow = 31 * 24 * time.Hour

// ComputeWindow computes the [start, end) window a report should cover.
//
// end is always now. start is watermark when set, or now-7d on a
// definition's first run. The result is capped to 31 days: if end.Sub(start)
// would exceed that, start is pulled forward to end.Add(-maxWindow) and
// truncated is reported true so the caller can note it in the report header.
//
// watermark must be the period_end of the definition's most recent
// SCHEDULED report — never a manual run's. Manual runs (plan Q12) use this
// same watermark to compute their window but must not persist their own
// period_end as the new watermark afterward; that bookkeeping is the
// caller's responsibility (it lives in the store-backed Generator, not
// here), which is why this function takes no "manual" flag: the watermark
// value passed in is what decides everything.
func ComputeWindow(watermark *time.Time, now time.Time) (start, end time.Time, truncated bool) {
	end = now
	if watermark != nil {
		start = *watermark
	} else {
		start = now.Add(-defaultWindow)
	}

	if end.Sub(start) > maxWindow {
		start = end.Add(-maxWindow)
		truncated = true
	}

	return start, end, truncated
}
