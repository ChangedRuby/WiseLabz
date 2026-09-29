package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestAddJobRegistersAndFires verifies that AddJob registers a job and it fires
// according to the cron expression.
func TestAddJobRegistersAndFires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		r.Start(ctx)
		synctest.Wait()

		synctest.Sleep(time.Second)

		r.Stop()

		if count := execCount.Load(); count != 1 {
			t.Fatalf("execution count = %d, want 1 after the first cron tick", count)
		}
	})
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
	synctest.Test(t, func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		r.Start(ctx)
		synctest.Wait()

		synctest.Sleep(time.Second)

		if count := execCount.Load(); count != 1 {
			t.Fatalf("execution count before removal = %d, want 1", count)
		}

		firstCount := execCount.Load()
		r.RemoveJob(entryID)

		// Cross two more cron ticks to prove removal prevents execution.
		synctest.Wait()
		synctest.Sleep(2 * time.Second)
		r.Stop()

		if finalCount := execCount.Load(); finalCount != firstCount {
			t.Fatalf("job executed after removal: initial count=%d, final count=%d", firstCount, finalCount)
		}
	})
}

// TestPanicRecovery verifies that panics in job functions are recovered and logged.
func TestPanicRecovery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		r.Start(ctx)
		synctest.Wait()

		synctest.Sleep(2 * time.Second)

		r.Stop()

		// If we get here without the process exiting, panic recovery works
		if count := executionCount.Load(); count != 2 {
			t.Fatalf("job did not recover from panic and run again, execution count=%d", count)
		}
	})
}

// TestStopViaContextCancellation verifies that Stop is called when context is cancelled.
func TestStopViaContextCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		logger := testLogger()
		r := New(logger)

		var count atomic.Int32
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})
		_, err := r.AddJob("test", "*/1 * * * * *", func(_ context.Context) error {
			count.Add(1)
			return nil
		})
		if err != nil {
			t.Fatalf("AddJob() error: %v", err)
		}
		r.Start(ctx)
		synctest.Wait()
		synctest.Sleep(time.Second)
		if got := count.Load(); got != 1 {
			t.Fatalf("execution count = %d, want 1 before cancellation", got)
		}

		cancel()
		synctest.Wait()
		synctest.Sleep(2 * time.Second)
		if got := count.Load(); got != 1 {
			t.Fatalf("job executed after cancellation: count = %d, want 1", got)
		}
	})
}

// TestMultipleJobsSchedule verifies that multiple jobs can be scheduled and run independently.
func TestMultipleJobsSchedule(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		r.Start(ctx)
		synctest.Wait()

		synctest.Sleep(time.Second)

		r.Stop()

		if first, second := count1.Load(), count2.Load(); first != 1 || second != 1 {
			t.Fatalf("job execution counts = (%d, %d), want (1, 1)", first, second)
		}
	})
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
	synctest.Test(t, func(t *testing.T) {
		logger := testLogger()
		r := New(logger)

		type ctxKey string
		const key ctxKey = "test-key"
		ctx, cancel := context.WithCancel(context.WithValue(context.Background(), key, "test-value"))
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		sawValue := make(chan bool, 1)
		sawCancelPropagated := make(chan struct{}, 1)
		var count atomic.Int32
		_, err := r.AddJob("test", "*/1 * * * * *", func(jobCtx context.Context) error {
			count.Add(1)
			sawValue <- (jobCtx.Value(key) == "test-value")
			// Remain in flight until Start's context is cancelled.
			<-jobCtx.Done()
			sawCancelPropagated <- struct{}{}
			return nil
		})
		if err != nil {
			t.Fatalf("AddJob() error: %v", err)
		}
		r.Start(ctx)
		synctest.Wait()
		synctest.Sleep(time.Second)
		select {
		case ok := <-sawValue:
			if !ok {
				t.Fatal("job context did not carry the value from Start's context")
			}
		default:
			t.Fatal("job did not execute on the first cron tick")
		}

		cancel()
		synctest.Wait()
		select {
		case <-sawCancelPropagated:
		default:
			t.Fatal("job did not observe Start's context cancellation")
		}
		synctest.Sleep(2 * time.Second)
		if got := count.Load(); got != 1 {
			t.Fatalf("job executed after cancellation: count = %d, want 1", got)
		}
	})
}

// TestContextGivenToJobFunction verifies that the job function receives a context.
func TestContextGivenToJobFunction(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
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

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(func() {
			cancel()
			r.Stop()
			synctest.Wait()
		})

		r.Start(ctx)
		synctest.Wait()

		synctest.Sleep(time.Second)

		r.Stop()

		if !receivedCtx.Load() {
			t.Fatalf("job function did not receive a context")
		}
	})
}

// TestStopWaitsForInFlightJob proves Stop joins work before returning.
func TestStopWaitsForInFlightJob(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := New(testLogger())
		ctx, cancel := context.WithCancel(context.Background())
		entered := make(chan struct{}, 1)
		release := make(chan struct{})
		finished := make(chan struct{})
		stopped := make(chan struct{})
		var releaseOnce sync.Once
		// Release first, including on assertion failure, so cleanup can join.
		t.Cleanup(func() {
			releaseOnce.Do(func() { close(release) })
			cancel()
			r.Stop()
			synctest.Wait()
		})
		_, err := r.AddJob("blocking", "*/1 * * * * *", func(context.Context) error {
			entered <- struct{}{}
			<-release
			close(finished)
			return nil
		})
		if err != nil {
			t.Fatalf("AddJob() error: %v", err)
		}
		r.Start(ctx)
		synctest.Wait()
		synctest.Sleep(time.Second)
		select {
		case <-entered:
		default:
			t.Fatal("job did not enter on the first cron tick")
		}

		go func() {
			r.Stop()
			close(stopped)
		}()
		synctest.Wait()
		select {
		case <-stopped:
			t.Fatal("Stop returned while the job was still in flight")
		default:
		}
		releaseOnce.Do(func() { close(release) })
		synctest.Wait()
		select {
		case <-stopped:
		default:
			t.Fatal("Stop did not return after the job was released")
		}
		select {
		case <-finished:
		default:
			t.Fatal("Stop returned before the job finished")
		}
	})
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
