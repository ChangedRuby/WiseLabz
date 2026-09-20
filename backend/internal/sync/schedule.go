package sync

import (
	"time"
)

// retrySchedule is the backoff delay before each retry attempt (index 0 =
// delay before the 2nd consecutive attempt, etc), mirroring
// notifications.retrySchedule.
// ponytail: fixed schedule, not exponential-from-config; add jitter/config if a real deployment needs it.
var retrySchedule = []time.Duration{time.Minute, 5 * time.Minute, 15 * time.Minute, 30 * time.Minute, time.Hour}

// computeNextRun decides a connector's next scheduled run time and updated
// retry count after one sync attempt. Unlike notification delivery retries
// (which give up after a fixed number of attempts), a scheduled connector
// must keep trying forever: once retryCount exceeds the backoff schedule, it
// keeps retrying at the schedule's last (longest) step. The retry delay is
// also never allowed to exceed the connector's own scheduleSeconds cadence,
// so a short-cadence connector doesn't end up waiting an hour to recover.
//
// scheduleSeconds nil means manual-only: nextRunAt is always nil, but
// retryCount is still tracked for display purposes. retryCount is the
// connector's consecutive-failure count *before* this attempt.
func computeNextRun(scheduleSeconds *int, retryCount int, success bool, now time.Time) (nextRunAt *time.Time, newRetryCount int) {
	if success {
		newRetryCount = 0
	} else {
		newRetryCount = retryCount + 1
	}

	if scheduleSeconds == nil {
		return nil, newRetryCount
	}
	cadence := time.Duration(*scheduleSeconds) * time.Second

	delay := cadence
	if !success {
		idx := newRetryCount - 1
		if idx >= len(retrySchedule) {
			idx = len(retrySchedule) - 1
		}
		delay = retrySchedule[idx]
		if cadence < delay {
			delay = cadence
		}
	}

	next := now.Add(delay)
	return &next, newRetryCount
}
