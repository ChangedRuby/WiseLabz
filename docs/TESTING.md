# Backend test performance

This file records where backend Go test time went, what #401 changed, what it
rejected and why, and the rules that keep the suite fast without weakening
determinism, race detection or coverage. Read it before adding a slow test,
a `time.Sleep`, a `t.Parallel()` or a new CI shard.

## Measuring

```bash
# Ranked packages, tests and subtests; raw `go test -json` saved to out.json.
scripts/ci/test-shards.sh profile out.json ./...
scripts/ci/test-shards.sh profile out.json -race ./internal/api/...
PROFILE_TOP=50 scripts/ci/test-shards.sh profile out.json "-coverpkg=$(scripts/ci/coverpkg.sh)" ./...
```

The **Test profile** workflow (Actions → Test profile → Run workflow) runs the
profile in normal, race and cover mode three times each on `ubuntu-latest`,
uploads the JSON, and reports the median wall clock. CI is the source of
truth: local runs guide the work, CI medians decide what is adopted.
`workflow_dispatch` workflows only run from the default branch, so compare a
branch through its regular CI job times.

A change is adopted when it saves at least 5% of the affected CI job, or at
least 10s (median of 3 runs), and does not weaken determinism, race detection
or coverage (see the parity check below).

## Where the time went

The baseline, measured locally on 6 cores before #401:

| Run | Before | After |
|---|---|---|
| `go test ./...` | 109s | 14.5s |
| `internal/api` | 107.4s | 0.63s |
| `internal/api/connectors` with `-race` | 92s | 3.3s |
| `internal/notifications` with `-race` | 48s | 3.1s |
| `internal/sync` with `-race` | 27.6s | 3.2s |
| `internal/scheduler` | 8.3s | 0.003s |
| Coverage run (`-coverpkg`, unsharded) | 106s | 20s |

The race-suite cost broke down as follows:

- **Fixture setup, almost all of it.** Every test opened a new SQLite database
  and ran every migration. `modernc.org/sqlite` is pure Go, so `-race`
  instruments the whole engine: a migration run cost about 1s per test with
  `-race` and about 0.04s without. `Store.Init` also hashed the admin password
  with bcrypt at cost 12.
- **Serial execution, most of the rest.** No backend test called
  `t.Parallel()`.
- **Timers.** Real `time.Sleep` and cron ticks in the scheduler, retry backoff
  in the connector clients, fixed server-side sleeps in timeout tests.
- **Not significant:** async sync goroutines and the single-connection SQLite
  pool (`SetMaxOpenConns(1)`).
- **Fixed per-package cost.** With `-race`, every test binary sleeps 1s at
  exit (`GORACE` `atexit_sleep_ms`, default 1000), so a package never reports
  under ~1s. It is left alone: the sleep gives late goroutines time to report
  races.

## What changed

- **bcrypt.** `auth.bcryptCost` is `bcrypt.MinCost` inside test binaries only
  (`testing.Testing()`). Production stays at 12, guarded by
  `TestBcryptCostOnlyLoweredUnderTest`.
- **Migrated database template.** `storetest.MigratedSQLite(t)` runs the
  migrations once per test binary and gives each test its own copy of the
  file. It is used by the `internal/api`, `api/auth`, `api/notifications`,
  `apitest`, `notifications` and `sync` harnesses. Tests stay fully isolated:
  each gets its own file.
- **Parallel tests.** Every top-level test in `internal/api` and
  `internal/api/connectors` calls `t.Parallel()` (see the rules below).
- **Deterministic time.**
  - Scheduler tests run inside `testing/synctest` bubbles, so cron ticks use a
    fake clock.
  - `connector.NewHTTPClient` backs off 1ms instead of 250ms+500ms inside test
    binaries only. Requests are still retried; `httpx`'s own tests check the
    real delays.
  - Timeout tests block the handler on `r.Context().Done()` instead of
    sleeping.
  - The doc-lock renewal test moves the expiry in SQL instead of sleeping.
- **No public network.**
  - The ntfy and Telegram tests used to send real requests to ntfy.sh and
    api.telegram.org. They now use a local server and assert the exact
    request path.
  - The AI-config and OIDC tests that only need a failing call point at a
    closed local port.
- **CI layout.**
  - Coverage is one unsharded job.
  - The race suite is one shard that now also covers `internal/api/...`.
  - The Postgres shards are unchanged: they need a Postgres service, so they
    are measured in CI only.

## Coverage strategy

CI runs, from `backend/`:

```bash
go test -coverpkg="$(../scripts/ci/coverpkg.sh)" -coverprofile=coverage.out ./...
```

`coverpkg.sh` is `./...` minus the test-only helpers (`api/apitest`,
`store/storetest`, `connector/connectortest`), which have non-`_test` file
names but are not product code. `-coverpkg` is needed because most handler
packages are exercised through `internal/api`'s router tests.

To check a change to coverage collection, compare per-package results, not
just the total:

```bash
scripts/ci/coverage-parity.sh old.out new.out   # fails if any package moves > 0.1pp
VERBOSE=1 scripts/ci/coverage-parity.sh old.out new.out
```

`internal/ws` is allowed 1pp (`NOISY_PACKAGES`): its write-pump `select`
branches (`ws.go`, around lines 108-117) run only under some goroutine
interleavings, and it moves about ±0.6pp between identical runs.
`internal/chat` had a similar load-dependent branch. It is now covered by
`TestVectorCachePutExistingKeyUpdatesInPlace` instead of being allowed a
tolerance.

These strategies were benchmarked after the setup-cost fixes:

| Strategy | Wall clock | Result |
|---|---|---|
| A. 3 coverage shards + `covdata` merge (previous) | ~22s per shard + merge job | baseline, 80.4% |
| B. One unsharded run | ~22s | **adopted**: identical per package (ws within noise) |
| C. Package-local coverage (no `-coverpkg`) | ~22s | rejected: 80.5% → 72.1% total, 48 packages differ |
| D. `-coverpkg=./internal/...` | ~22s | rejected: drops `cmd/backup` and `cmd/server` |

Before the setup fixes, B took 106s and A's shards 44-54s each. That is why
the suite was sharded. B now matches A's slowest shard and saves a job and a
merge step.

## Rules for new tests

**Parallel tests** (`internal/api`, `internal/api/connectors`):

- Each test gets its own database, `ws.Hub` and scheduler from
  `newTestApp`/`apitest`. Never share one between tests.
- Process-wide state must be safe under parallel tests:
  - `connector.AllowLoopbackForTest` is a refcount: it allows loopback while
    any test holds it. So in these packages, never write a test that depends
    on loopback being *blocked*. Put it in a serial package (`connector`,
    `notifications`).
  - Register fake connector types under a per-test name
    (`"health_test_fake/"+t.Name()`) when their factory captures per-test
    state.
  - Set per-type behaviour (such as the degraded-latency threshold) on the
    fake type's `TypeSchema` instead of mutating a package variable.
- Goroutines a test starts must stop with the test: use `t.Context()`.

**Serial by choice:** every other package. They are fast now, so
`t.Parallel()` would not meet the adoption bar. Some also have real
blockers:

- `notifications` swaps package variables (`fakePublicAPI`) and asserts that
  loopback is blocked.
- The connector packages mix `AllowLoopbackForTest` with blocked-loopback
  tests.
- The `store` Postgres job runs with `-p 1` and one schema per test.

**Time:**

- No fixed sleeps to "let something happen". Wait on the event: a channel,
  `r.Context().Done()`, or `synctest.Wait` inside a bubble.
- Code on a clock (cron, tickers, timeouts) runs under `testing/synctest`
  when it does no real I/O. Otherwise inject the clock (`now func()`, as in
  `internal/quality/checker.go`).
- Deadlines in polling loops are hang guards only, and generous: they must
  hold with `-race` on a slow runner. `waitForSyncRuns` uses 30s; it returns
  as soon as the condition is met.
- Test-only speedups go through `testing.Testing()` in one clearly named
  place, never an env var or config knob production could reach.

**Network:** tests never call public services. Point the code at an
`httptest` server (with `connector.AllowLoopbackForTest(t)` where the
guarded dialer applies), or at a closed local port such as
`http://127.0.0.1:1` when the test only needs the call to fail.

## CI job times

The baseline is `main` at 358d4a9 (run 36355404045), with warm caches:

| Job group | Jobs | Longest job | Runner time |
|---|---|---|---|
| Coverage (3 shards + merge) | 4 | 156s + 25s merge | 396s |
| Race (core + 4 connector shards) | 5 | 156s | 537s |
| Postgres (3 shards) | 3 | 52s | 137s |
| Whole CI run | | | 227s wall |

Results after #401 (PR #404, median of 3 CI runs):

| Job group | Jobs | Longest job | Runner time |
|---|---|---|---|
| Coverage | 1 | 94s (91, 94, 125) | 94s |
| Race (one shard, now incl. `internal/api/...`) | 1 | 147s (142, 147, 150) | 147s |
| Postgres (unchanged) | 3 | 66s | 178s |
| Whole CI run | | | 165s wall (160, 165, 169) |

- **Coverage:** the critical path goes from 181s (slowest shard plus the
  merge job) to 94s.
- **Race:** runner time goes from 537s to 147s, although the suite now also
  covers `internal/api/...`. The PR runs had a cold race build cache (the
  single-shard `race` cache kind is saved on `main` only), so most of those
  147s went to compiling; the slowest package, `internal/api`, took 22s.
  Expect the job to get faster once `main` has saved the cache.
- **Postgres:** these shards did not change; their 137s → 178s is runner
  variance.

## Follow-ups

- **`internal/store`** still migrates per test (6.6s normal, about 160s with
  `-race`). It cannot use `storetest`, which imports `store`: that would be an
  import cycle. It needs an in-package template helper. It is not in the race
  suite today.
- **Other packages that still migrate per test:** `quality`, `backup`, `doc`,
  `docexport`, `chat`, `mcp`, `retention`, `diagnostics`, `cmd/backup`. Each
  takes 1-2s normally, so the gain is small until they join a race job.
- **Two remaining one-tick waits:**
  - `cmd/server` `TestStandbyIsUnreadyAndRunsNoScheduler` sleeps 1.1s to
    prove a standby never runs a job.
  - `docexport` `TestScheduledExportRunsAndFires` waits up to one cron tick.
  - Making either instant means running the server lifecycle or the
    exporter's store inside a synctest bubble. Both are under the adoption
    bar.
- **Stress validation:** `go test -race -count=20 -shuffle=on` on the
  parallelized packages should run in CI (not on small dev machines) before
  the parallel rollout is extended to more packages.
