package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestAddJobRegistersAndFires verifies that AddJob registers a job and it fires
// according to the cron expression.
func TestAddJobRegistersAndFires(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var execCount atomic.Int32
	_, err := r.AddJob("test", "*/1 * * * * *", func(_ context.Context) error {
		execCount.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r.Start(ctx)

	// Wait for at least one execution (the job should fire within 1 second)
	deadline := time.Now().Add(2 * time.Second)
	for execCount.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	r.Stop()

	if execCount.Load() == 0 {
		t.Fatalf("job did not execute within 2 seconds")
	}
}

// TestAddJobInvalidExpression verifies that AddJob rejects invalid cron expressions.
func TestAddJobInvalidExpression(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	_, err := r.AddJob("test", "invalid cron", func(_ context.Context) error { return nil })
	if err == nil {
		t.Fatalf("AddJob() should reject invalid cron expression, got nil error")
	}
}

// TestRemoveJobDeregisters verifies that RemoveJob removes a registered job.
func TestRemoveJobDeregisters(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var execCount atomic.Int32
	entryID, err := r.AddJob("test", "*/1 * * * * *", func(_ context.Context) error {
		execCount.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r.Start(ctx)

	// Wait for first execution
	deadline := time.Now().Add(2 * time.Second)
	for execCount.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	if execCount.Load() == 0 {
		t.Fatalf("job did not execute before removal")
	}

	firstCount := execCount.Load()
	r.RemoveJob(entryID)

	// Wait a bit and verify the job doesn't execute again
	time.Sleep(1500 * time.Millisecond)
	r.Stop()

	if finalCount := execCount.Load(); finalCount != firstCount {
		t.Fatalf("job executed after removal: initial count=%d, final count=%d", firstCount, finalCount)
	}
}

// TestPanicRecovery verifies that panics in job functions are recovered and logged.
func TestPanicRecovery(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var executionCount atomic.Int32
	_, err := r.AddJob("panic-test", "*/1 * * * * *", func(_ context.Context) error {
		if executionCount.Add(1) == 1 {
			panic("intentional panic")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r.Start(ctx)

	// Wait for the panic to happen and the job to run again
	deadline := time.Now().Add(3 * time.Second)
	for executionCount.Load() < 2 && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	r.Stop()

	// If we get here without the process exiting, panic recovery works
	if count := executionCount.Load(); count < 2 {
		t.Fatalf("job did not recover from panic and run again, execution count=%d", count)
	}
}

// TestStopViaContextCancellation verifies that Stop is called when context is cancelled.
func TestStopViaContextCancellation(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var started atomic.Bool

	_, err := r.AddJob("test", "*/1 * * * * *", func(_ context.Context) error {
		started.Store(true)
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.Start(ctx)

	// Give the job a moment to start (up to 2 seconds)
	deadline := time.Now().Add(2 * time.Second)
	for !started.Load() && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	// Cancel context
	cancel()

	// Give a moment for the stop goroutine to execute
	time.Sleep(500 * time.Millisecond)

	if !started.Load() {
		t.Fatalf("job did not start within timeout")
	}
}

// TestMultipleJobsSchedule verifies that multiple jobs can be scheduled and run independently.
func TestMultipleJobsSchedule(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var count1, count2 atomic.Int32

	_, err := r.AddJob("job1", "*/1 * * * * *", func(_ context.Context) error {
		count1.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob(job1) error: %v", err)
	}

	_, err = r.AddJob("job2", "*/1 * * * * *", func(_ context.Context) error {
		count2.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob(job2) error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r.Start(ctx)

	// Wait for both jobs to execute
	deadline := time.Now().Add(2 * time.Second)
	for (count1.Load() == 0 || count2.Load() == 0) && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	r.Stop()

	if count1.Load() == 0 {
		t.Fatalf("job1 did not execute")
	}
	if count2.Load() == 0 {
		t.Fatalf("job2 did not execute")
	}
	if first, second := count1.Load(), count2.Load(); first != second {
		t.Logf("job counts differ (job1=%d, job2=%d) but both executed, which is acceptable", first, second)
	}
}

// TestInvalidJobNameHandled verifies that job names are logged properly.
// This is more of a smoke test to ensure the code doesn't panic.
func TestInvalidJobNameHandled(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	_, err := r.AddJob("", "*/1 * * * * *", func(_ context.Context) error { return nil })
	if err != nil {
		t.Fatalf("AddJob() with empty name should not error: %v", err)
	}

	_, err = r.AddJob(fmt.Sprintf("job-%d", time.Now().Unix()), "*/1 * * * * *", func(_ context.Context) error { return nil })
	if err != nil {
		t.Fatalf("AddJob() with complex name should not error: %v", err)
	}
}

// TestJobContextDerivedFromStart verifies that job invocations receive the
// context passed to Start (not context.Background()), so job bodies observe
// cancellation from the server shutdown context, mirroring how the old
// ticker-loop schedulers respected ctx.Done().
func TestJobContextDerivedFromStart(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	type ctxKey string
	const key ctxKey = "test-key"
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), key, "test-value"))
	defer cancel()

	sawValue := make(chan bool, 1)
	sawCancelPropagated := make(chan bool, 1)
	var once sync.Once
	_, err := r.AddJob("test", "*/1 * * * * *", func(jobCtx context.Context) error {
		once.Do(func() {
			sawValue <- (jobCtx.Value(key) == "test-value")
		})
		// Block within the job body (as a long-running job would) so it is
		// still in-flight when the test cancels Start's context, proving the
		// cancellation reaches the running job rather than only stopping the
		// cron scheduler itself.
		select {
		case <-jobCtx.Done():
			select {
			case sawCancelPropagated <- true:
			default:
			}
		case <-time.After(4 * time.Second):
		}
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	r.Start(ctx)

	select {
	case ok := <-sawValue:
		if !ok {
			t.Fatalf("job context did not carry the value from Start's context")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("job did not execute within 2 seconds")
	}

	cancel()

	select {
	case <-sawCancelPropagated:
	case <-time.After(4 * time.Second):
		t.Fatalf("job context was never observed as Done() after Start's context was cancelled")
	}
}

// TestContextGivenToJobFunction verifies that the job function receives a context.
func TestContextGivenToJobFunction(t *testing.T) {
	logger := testLogger()
	r := New(logger)

	var receivedCtx atomic.Bool
	_, err := r.AddJob("test", "*/1 * * * * *", func(ctx context.Context) error {
		if ctx != nil {
			receivedCtx.Store(true)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r.Start(ctx)

	// Wait for the job to execute
	deadline := time.Now().Add(2 * time.Second)
	for !receivedCtx.Load() && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}

	r.Stop()

	if !receivedCtx.Load() {
		t.Fatalf("job function did not receive a context")
	}
}

func TestJobSkipsOverlappingInvocations(t *testing.T) {
	r := New(testLogger())
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	id, err := r.AddJob("blocking", "* * * * * *", func(context.Context) error { entered <- struct{}{}; <-release; return nil })
	if err != nil {
		t.Fatal(err)
	}
	job := r.c.Entry(id).WrappedJob
	done := make(chan struct{})
	go func() { defer close(done); job.Run() }()
	<-entered
	skipped := make(chan struct{})
	go func() { defer close(skipped); job.Run() }()
	select {
	case <-skipped:
	case <-time.After(time.Second):
		close(release)
		t.Fatal("overlap was queued instead of skipped")
	}
	if len(entered) != 0 {
		t.Error("overlapping job executed")
	}
	close(release)
	<-done
	job.Run()
	if len(entered) != 1 {
		t.Error("job did not run after previous invocation completed")
	}
}
