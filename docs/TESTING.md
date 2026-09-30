# Backend test performance

This file records where backend Go test time went, what #401 and #405–#408 changed, what they
rejected and why, and the rules that keep the suite fast without weakening
determinism, race detection or coverage. Read it before adding a slow test,
a `time.Sleep`, a `t.Parallel()` or a new CI shard. [What CI runs](#what-ci-runs)
explains which jobs a change starts.

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
`workflow_dispatch` workflows must be available on the default branch before
they can be dispatched; then select the branch to measure when running them.

A change is adopted when it saves at least 5% of the affected CI job, or at
least 10s (median of 3 runs), and does not weaken determinism, race detection
or coverage (see the parity check below).

### Shuffled race stress gate

Before extending `t.Parallel()` to any new package, add its explicit package
path to the maintained allowlist in `.github/workflows/test-stress.yml` and
require a passing **Test stress** run on the rollout commit. The allowlist
currently contains `./internal/api` and `./internal/api/connectors`; it takes
no user-supplied package expressions. Keep fixtures isolated per test and
use local servers; the gate needs no public test services or shared fixtures.

After the workflow is available on the default branch, run it from Actions
→ Test stress → Run workflow, selecting the rollout branch, or use:

```bash
gh workflow run test-stress.yml --ref <rollout-branch>
```

Each package runs `go test -race -count=20 -shuffle=on -timeout=10m -v` in
its own job, with a 20-minute job timeout. This manual gate stays out of
routine PR CI. Failed tests fail the job; other package jobs still finish.
The `stress-*` artifacts are retained for 14 days, including on failure:
`metadata.txt` records the checked-out commit, package, flags, run URL and Go
version, and `tests.log` contains stdout/stderr and the `-test.shuffle` seed.
To reproduce, check out that commit and run from `backend/` with the recorded
package and Go version, replacing `-shuffle=on` with `-shuffle=<seed>`.

PR #404 already passed the temporary shuffled race stress runs for the two
current parallel packages (288s for `internal/api`, 155s for
`internal/api/connectors`). This retained gate supports future rollouts.

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
  `apitest`, `notifications` and `sync` harnesses, plus the #406 fixtures
  listed below. `internal/store` uses its own `_test.go` template to avoid
  an import cycle. Tests stay fully isolated:
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
  - The race suite is one shard covering `internal/api/...` and all eleven
    #405–#407 packages listed below.
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
- The `store` Postgres job runs with `-p 1` and one schema per test. It runs
  only when a change reaches the store's Go dependency closure (see
  [What CI runs](#what-ci-runs)).

**Time:**

- No fixed sleeps to "let something happen". Wait on the event: a channel,
  `r.Context().Done()`, or `synctest.Wait` inside a bubble.
- Code on a clock (cron, tickers, timeouts) runs under `testing/synctest`.
  Settle startup with `synctest.Wait()`, then use `synctest.Sleep` for bounded
  clock advances and assert completion immediately using channels or
  synchronized counts. Self-contained SQLite and filesystem work can run
  inside a bubble when all associated goroutines and resources finish there;
  real network I/O stays outside. Otherwise inject the clock (`now func()`,
  as in `internal/quality/checker.go`).
- Deadlines in polling loops are hang guards only, and generous: they must
  hold with `-race` on a slow runner. `waitForSyncRuns` uses 30s; it returns
  as soon as the condition is met.
- Test-only speedups go through `testing.Testing()` in one clearly named
  place, never an env var or config knob production could reach.

**Network:** tests never call public services. Point the code at an
`httptest` server (with `connector.AllowLoopbackForTest(t)` where the
guarded dialer applies), or at a closed local port such as
`http://127.0.0.1:1` when the test only needs the call to fail.

## What CI runs

`CI Status` is the only required check, and it reports on every PR and every
push to `main`. The workflow always starts: a `paths-ignore` trigger would
leave the required check pending forever on a docs-only PR. The `changes` job
decides what else runs, using `scripts/ci/changes.sh`, and every other job is
gated on its outputs.

**Areas.** Each changed path is matched against the ordered rule table at the
top of `scripts/ci/changes.sh`. The first match wins.

| Output | Selected by | Jobs |
|---|---|---|
| `backend` | `backend/**`, `go.work*`, `.golangci.yml`, `scripts/ci/**`, the Go cache actions | lint, staticcheck, unit, race, build |
| `frontend` | `web/**`, `docs/openapi.yaml`, the Bun setup action | lint, test, build |
| `compose` | backend or frontend, plus `Dockerfile`, `.dockerignore`, `docker-compose*.yml`, `.env.example`, `scripts/compose-smoke.*` | compose-smoke |
| `workflows` | `.github/workflows/**`, `.github/actions/**`, `.github/codeql/**` | `Lint (Workflows)` (actionlint + shellcheck) |
| `gomod` | `backend/go.mod`, `backend/go.sum`, `go.work*` | govulncheck |
| `postgres` | a backend change inside the Postgres tests' Go dependency closure (below) | Postgres shards |

Changing `.github/workflows/ci.yml` selects everything.

**Ignored.** These paths never start a job: `*.md` anywhere, `docs/**` (except
`openapi.yaml`), `graphify-out/**`, `openspec/**`, `.claude/**`, `.agents/**`,
`.codex/**`, `deploy/**`, `prototypes/**`, release metadata (`CHANGELOG.md`,
release-please config and manifest), local tooling (`Makefile`,
`lefthook.yml`, `.air.toml`) and repository metadata (`LICENSE`,
`CODEOWNERS`, issue and PR templates, `.github/dependabot.yml`). A docs-only
PR finishes in about 15 seconds.

**Fail-safe.** A path that matches no rule is `unclassified` and selects
backend, frontend, compose and gomod. Unknown files never cause a skip; they
cost a full run until someone classifies them.

**Release PRs.** On pull requests, a `web/package.json` change that only bumps
`version` is ignored, so the release-please PR skips the heavy jobs. The push
to `main` after a release still runs them.

**Postgres closure.** When backend changed, `changes` sets up Go and runs
`changes.sh go-closure`. It lists the dependencies of the Postgres test
packages (from `test-shards.json`) and `./cmd/migrate` with
`go list -deps -test`. It then maps each changed file to its nearest enclosing
Go package, so migrations, `testdata/` and deleted files count too. The shards
run when that package is in the closure. They also run for Go module files,
`scripts/ci/**`, the Go cache actions, `ci.yml`, unclassified paths, a path
with no enclosing package, or a `go list` failure.

**govulncheck.** PR CI runs it only when `gomod` is selected. The
`Govulncheck nightly` workflow scans `main` every day at 06:17 UTC (it can
also be run by hand). It catches new advisories against unchanged
dependencies, and code changes that made a known-vulnerable function
reachable. A red nightly run needs an owner: treat it like a failing `main`.

**Drafts.** On a draft PR only `changes`, the rule check and
`Lint (Workflows)` run. `CI Status` stays green, and the summary says heavy
jobs are deferred. Marking the PR ready for review re-runs CI on the same
commit with the full selection. CodeQL also skips drafts.

**CodeQL – Code Quality.** GitHub's Code Quality scan (runs named
`CodeQL - Code Quality`, event `dynamic`) is configured in repository
settings, not in a workflow file. It runs on every PR, docs-only ones
included, for about 2.5 minutes. It is not a required check. It can't be
path-filtered: its settings (`gh api repos/WiseLabz/WiseLabz/code-quality/setup`)
cover only on/off, languages, runner and AI findings.

**Reading a run.** The `Detect changes` job summary lists the selected areas.
It has a collapsible table of every changed path with the rule it matched,
followed by the Postgres decision and its reason.

**Adding or changing a rule.** Edit the table in `scripts/ci/changes.sh`
(order matters), then add a row to `scripts/ci/changes-fixtures.txt` for the
path. `changes` runs the fixtures on every run, and a mismatch fails CI.
Suppose a `.md` file under `backend/` or `web/` ever becomes embedded or
imported. The `*.md` rule would then hide changes to it, so add an exception
for it above that rule. To check locally:

```sh
scripts/ci/changes.sh check          # classification fixtures
scripts/ci/changes.sh check --go     # plus Postgres closure fixtures (needs Go)
git diff --name-only origin/main... | scripts/ci/changes.sh classify
git diff --name-only origin/main... | scripts/ci/changes.sh go-closure
```

### Dependency updates and action pinning

**Pinning.** Every `uses:` of an action outside this repository names a full
commit SHA followed by the release it belongs to, for example
`actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`. Pin the
commit, not an annotated tag object: resolve it with
`gh api repos/<owner>/<repo>/commits/<tag> --jq .sha`. The action must also
run on a supported Node runtime (`node24` today, never `node20`). Check with
`gh api "repos/<owner>/<repo>/contents/action.yml?ref=<sha>" --jq .content | base64 -d | grep using:`.
This must print nothing:

```sh
grep -rnE 'uses: [^.].*@v[0-9]' .github
```

**Dependabot.** `.github/dependabot.yml` opens update PRs every week:

| Ecosystem | Scope | PRs | Commit | Label |
|---|---|---|---|---|
| `github-actions` | workflows and `.github/actions/*` | one grouped PR | `ci: ...` | `area:platform` |
| `gomod` | `/backend` | minor+patch grouped, one PR per major | `chore(deps): ...` | `area:backend` |
| `bun` | `/web` | minor+patch grouped, one PR per major | `chore(deps): ...` | `area:frontend` |

Both commit types pass the `commit-msg` hook and stay out of the changelog.
Dependabot rewrites each SHA together with its `# vX.Y.Z` comment. An actions
PR that touches `ci.yml` runs the full suite, which re-validates CI with the
new versions; Go and web PRs run only their area.

Dependabot does not update the root `go.work.sum`. If a Go update PR fails
with a checksum error, run `go work sync` on its branch and push the result.

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

## Fixture reuse and lifecycle tests (#405–#407)

These changes are adopted **regardless of the 5% / 10s threshold**. This is
an explicitly chosen exception: schema reuse and deterministic lifecycle
checks make regular race coverage affordable. Correctness, isolation, race
detection and coverage remain acceptance gates. No parallel tests,
dependencies, public APIs or configuration knobs were added.

The fifteen ordinary fixture sites converted in this change are:

| Package | Fixture sites |
|---|---|
| `internal/store` | `newDocTestStore`, `newConcurrentQualityTestStore`, `newCascadeTestStore`, MFA cascade test, external `newBackupTestStore` |
| `internal/quality` | `newTestStore` |
| `internal/backup` | `newTestStore` |
| `internal/doc` | `newEngineTestStore` |
| `internal/docexport` | `newTestStore` (also used by Git export tests) |
| `internal/chat` | embedding failure test and retrieval/cache test |
| `internal/mcp` | `newTestHarness` |
| `internal/retention` | `newTestStore` |
| `internal/diagnostics` | `newTestStore` |
| `cmd/backup` | `newSeededStore`; export directories remain independently supplied |

The store package builds its template with `sync.Once`, closes the migrated
database before reading immutable bytes, and writes a separate `0600` file
under each test's temporary directory. `MigratedSQLiteForTest` exists only
in test binaries and lets external `store_test` fixtures share this cache.
Other packages use `storetest.MigratedSQLite`. Templates contain only the
migration result (including migration-defined defaults), with no `Store.Init`
or test seeding. Both helpers are checked for concurrent copy creation,
permissions and current migration status; seeding one copy must leave both
an existing copy and a later copy free of its users and docs.

Callers retain their own DSN options, pool sizes, foreign-key enforcement,
WAL and busy timeouts, per-test initialization, seeding and cleanup. Real
migrations remain in migration/upgrade/rollback tests, restore verification,
production restore paths, and isolated PostgreSQL schemas.

The lifecycle's unexported scheduler dependency now accepts `Start` and
`Stop`; production still supplies the real runner. The standby test advances
a scheduler double by two one-second intervals and checks zero starts and
callbacks, HTTP `/readyz` 503 and `WaitingForLeader` both before and after.
The acquired-leader positive control advances one second and executes its
callback with the live work context. Startup waits for a buffered
`http.Server.BaseContext` listener signal and a successful local HTTP
request. Real HTTP and SQLite stay outside fake-time bubbles.

The two real-runner shutdown tests remain. A gated `Stop` additionally
proves shutdown marks readiness false while the database is usable and the
work context is live, then checks cancellation, goroutine completion and
database closure after releasing the gate. Cleanup is registered before
assertions; bounded 30-second waits serve only as hang guards.

`race/all` now also runs `internal/store`, `quality`, `backup`, `doc`,
`docexport`, `chat`, `mcp`, `retention`, `diagnostics`, `cmd/backup` and
`cmd/server`. The single shard and its ten-minute CI timeout are unchanged.

### Local measurements

Fresh uncached runs on the same machine, before resource-limited verification:

| Package | Normal before / after (s) | Race before / after (s) |
|---|---|---|
| `internal/store` | 6.877 / 0.860 | 161.894 / 15.657 |
| `internal/quality` | 1.772 / 0.152 | 41.656 / 3.675 |
| `internal/backup` | 1.143 / 0.306 | 27.506 / 6.575 |
| `internal/doc` | 0.854 / 0.083 | 22.046 / 2.693 |
| `internal/docexport` | 1.303 / 0.952 | 16.831 / 3.593 |
| `internal/chat` | 0.099 / 0.053 | 3.265 / 2.116 |
| `internal/mcp` | 0.279 / 0.077 | 6.827 / 2.415 |
| `internal/retention` | 0.233 / 0.062 | 6.601 / 2.186 |
| `internal/diagnostics` | 0.190 / 0.057 | 5.512 / 2.185 |
| `cmd/backup` | 0.323 / 0.191 | 9.838 / 6.244 |
| `cmd/server` | 1.241 / 0.072 | 3.593 / 2.392 |

Wall clock for those eleven packages: **8s → 3s normal**, **165s → 18s race**.
Both profiles use `-count=1`; the candidate race run also uses `-shuffle=on`.
These are local observations, not CI medians or a claim about the unchanged
race job: previously those eleven packages were absent from that job.

### Validation and coverage

All eleven packages pass normally and with `-race -count=1 -shuffle=on`.
Lifecycle tests also pass `-race -count=20 -shuffle=on`, and the existing
real-cron scheduler suite passes with race detection. All three existing
PostgreSQL shards pass locally on PostgreSQL 17.11, with `migrate up` and
`migrate verify` reporting schema version 45, clean and current; CI continues
to use PostgreSQL 16.

Full backend coverage uses the unchanged `-coverpkg` strategy and helper
exclusions. The baseline and candidate both round to **80.4%** overall.
`coverage-parity.sh` keeps its 0.1pp package tolerance and 1pp `internal/ws`
tolerance unchanged. Its raw comparison exits nonzero for intentional gains:
`cmd/server` **20.8% → 22.5%** and `internal/api/system` **77.2% → 78.0%**,
from the acquired-leader and draining-readiness tests. No package loses
coverage beyond the existing tolerances in the confirmation run. The first comparison also observed one
uncovered `internal/sync` connector-list error statement (about 0.14pp);
the confirmation covered it without any source change. Raw profiles and both
comparison outputs are retained so this variation is visible.

### CI measurements

The existing Test profile workflow ran normal, race and cover modes three
times per configuration on `ubuntu-latest`, with `-count=1`, the same Go
version and each mode's usual restored cache kind. Each run links to raw
JSON/text artifacts (retained by the workflow for 14 days). Normal and cover
use all backend packages; race uses the shard scope at that commit.
The expanded-scope baseline changes only `test-shards.json`, separating
fixture/timing improvements from the extra race coverage.

| Configuration / source | Normal runs → median (s) | Race runs → median (s) | Cover runs → median (s) |
|---|---|---|---|
| [Unchanged baseline](https://github.com/WiseLabz/WiseLabz/actions/runs/36496216327) (`fb615e3`) | 143, 140, 139 → **140** | 50, 48, 49 → **49** | 61, 62, 61 → **61** |
| [Expanded-scope baseline](https://github.com/WiseLabz/WiseLabz/actions/runs/36496761938) (`b067e81`) | 139, 140, 146 → **140** | 424, 448, 439 → **439** | 53, 59, 67 → **59** |
| [Implementation](https://github.com/WiseLabz/WiseLabz/actions/runs/36497384852) (`2aeb6cc`) | 106, 107, 99 → **106** | 95, 106, 89 → **95** | 29, 40, 65 → **40** |

At the same expanded race scope, the median falls **439s → 95s**: **344s
(78%)** saved. Against the former smaller scope, the final race run costs
**46s more** (49s → 95s) while covering eleven additional packages. Normal
falls 140s → 106s and cover 61s → 40s. These wall clocks include compilation;
cache restoration, cold builds of newly included race binaries and runner
variance are distinct from package execution. The package medians below
show the test-time effect directly; the normal baseline has a cold profile
build cache while race/cover restore their regular CI cache kinds.

| Package | Normal baseline / implementation median (s) | Race expanded baseline / implementation median (s) |
|---|---|---|
| `internal/store` | 33.851 / 3.548 | 304.379 / 26.538 |
| `internal/quality` | 10.620 / 0.937 | 88.422 / 7.468 |
| `internal/backup` | 9.260 / 2.538 | 73.753 / 14.930 |
| `internal/doc` | 6.455 / 0.575 | 58.286 / 4.582 |
| `internal/docexport` | 5.961 / 1.240 | 43.799 / 6.626 |
| `internal/chat` | 0.736 / 0.361 | 7.272 / 3.676 |
| `internal/mcp` | 1.797 / 0.492 | 17.545 / 4.860 |
| `internal/retention` | 2.264 / 0.406 | 16.584 / 4.267 |
| `internal/diagnostics` | 1.469 / 0.428 | 13.940 / 3.880 |
| `cmd/backup` | 2.375 / 1.055 | 34.649 / 17.875 |
| `cmd/server` | 1.654 / 0.533 | 6.277 / 8.123 |

[Regular backend CI](https://github.com/WiseLabz/WiseLabz/actions/runs/36497434502)
passes on the implementation commit, including PostgreSQL 16, coverage,
the expanded race shard, lint/static/vulnerability checks, build and compose
smoke. Local checks were completed sequentially with bounded Go build
concurrency after a resource-heavy verification attempt; their logs and
baseline/candidate coverage profiles are retained with the local evidence.
On memory-limited machines, run checks sequentially with disk-backed
`GOTMPDIR` and `TMPDIR`, `GOFLAGS=-p=1`, `GOMAXPROCS=2` and
`GOMEMLIMIT=384MiB`. These limit verification processes; production and CI
configuration are unchanged.

## Deterministic scheduled exports (#408)

This change also uses the agreed **reliability exception to the 5% / 10s
adoption threshold**. Timing evidence is still required; saving less than
that threshold does not block deterministic tests or the added race coverage.

`TestScheduledExportRunsAndFires` creates its migrated SQLite fixture,
exporter, real cron runner, context and channels inside one `synctest.Test`
bubble, using the bubble's `t`. Its `AddJob` callback calls the real
`RunExportOnce` and sends its returned error through a buffered channel.
After startup settles, one fake second advances the actual cron trigger.
A nonblocking receive fails immediately if execution is missing. The runner
stops before the test checks exactly `runbook-0000000a.md` containing `steps`.
Cleanup cancels the context, calls blocking `Stop()` and settles the
cancellation watcher before database closure or temporary-directory removal,
including assertion failures. The migrated-template fixture is unchanged.

Scheduler tests use exact fake-time advances and settled counts/signals
instead of virtual polling or settling sleeps. Removal and cancellation
cross two subsequent cron ticks and prove no additional callbacks run.
A gated callback proves `Stop()` remains blocked until work is released;
cleanup releases the gate even on failure. Existing health and overlap tests
remain unchanged. `race/all` now includes `internal/scheduler`, with its
existing ten-minute CI timeout as the wall-clock hang guard. No production
APIs, clock seams, dependencies, configuration or parallel-test changes are
included.

### Measurements and validation

The planning baseline for the scheduled export test was **0.85, 1.01, 1.01s**
(median **1.01s**). Fresh local and three-run CI measurements follow. Local profiles use `-count=1`
with Go 1.27.1, warm build caches, disabled test-result caching, `GOFLAGS=-p=1`,
`GOMAXPROCS=2` and `GOMEMLIMIT=384MiB`. Each package process starts with a cold
in-memory migration template.

| Local measurement | Baseline runs → median (s) | Candidate runs → median (s) | Median change |
|---|---|---|---|
| Scheduled export, normal | 0.10, 0.70, 0.90 → **0.70** | <0.01, <0.01, <0.01 → **<0.01** | >0.69s saved (>98.6%) |
| Scheduled export, race | 0.83, 0.13, 0.43 → **0.43** | 0.03, 0.03, 0.03 → **0.03** | 0.40s saved (93.0%) |
| `internal/docexport`, normal | 0.304, 0.911, 1.110 → **0.911** | 0.208, 0.209, 0.209 → **0.209** | 0.702s saved (77.1%) |
| `internal/docexport`, race | 3.718, 3.017, 3.297 → **3.297** | 2.813, 2.797, 2.850 → **2.813** | 0.484s saved (14.7%) |
| `internal/scheduler`, normal | 0.004, 0.004, 0.004 → **0.004** | 0.004, 0.004, 0.004 → **0.004** | 0.000s (0%) |
| `internal/scheduler`, race | 1.014, 1.015, 1.016 → **1.015** | 1.014, 1.015, 1.014 → **1.014** | 0.001s saved (0.1%) |

The normal scheduled-test JSON rounds elapsed time to hundredths; zero is
reported as <0.01s here. Wall-clock alignment of the baseline cron tick
explains its spread. Race package times retain the standard one-second exit
sleep. Package savings do not establish a 5% whole-job saving.

Both affected packages pass normally, with `-race`, and with
`-race -count=20 -shuffle=on`.

Full backend tests and the expanded `race/all` shard pass with `-count=1`.
Full `-coverpkg` coverage stays **80.4% → 80.4%**;
`coverage-parity.sh` reports **zero packages outside the existing tolerances**.
Formatting, shard-definition validation, golangci-lint, `go vet ./...` and
staticcheck v0.8.1 pass. The export test passes alone with race detection,
and first in the entire shuffled package (`-race -count=1 -shuffle=10`),
proving cold template initialization works inside the bubble. Isolated
negative controls remove the export registration or fail while the callback
gate is held: each exits at the intended assertion and completes cleanup,
without a timeout or bubble deadlock. Raw local JSON, profiles and logs are
retained under `/tmp/issue408-evidence/` on the verification machine.

### Three-run CI comparison

The existing Test profile workflow ran three uncached test executions per
mode, with `-count=1`, Go 1.27.1 and `ubuntu-latest`. Baseline and candidate
coverage/race runs restore the **same** main build-cache snapshot
(`36507104776`, dependency hash `aa226528…`). Modified test binaries and the
new scheduler race binary still require compilation. Normal mode has no
profile cache in either configuration and builds cold. Workflow timings
below measure the test command including compilation; setup/cache download
time is outside that measurement. Raw JSON/text artifacts remain attached
to each linked run for the workflow's usual 14-day retention.

| Profile / source | Normal runs → median (s) | Race runs → median (s) | Coverage runs → median (s) |
|---|---|---|---|
| [Baseline](https://github.com/WiseLabz/WiseLabz/actions/runs/36507943499) (`91d073a`) | 101, 106, 78 → **101** | 95, 94, 154 → **95** | 38, 24, 30 → **30** |
| [Candidate](https://github.com/WiseLabz/WiseLabz/actions/runs/36508418726) (`329d208`) | 102, 105, 127 → **105** | 171, 83, 96 → **96** | 29, 37, 25 → **29** |

Coverage saves **1s (3.3%)**; race costs **1s more (1.1%)**, with
`internal/scheduler` newly included; normal costs **4s more (4.0%)**.
Neither affected CI job meets the **5% / 10s** savings threshold. Adoption
therefore uses the explicitly agreed reliability exception. These noisy
whole-job timings do not establish a whole-suite speedup. The race scope
changes intentionally; its candidate scheduler package median is **1.019s**
(including the unchanged one-second race exit sleep), versus no scheduler
coverage in the baseline shard.

| CI package/test median | Baseline (s) | Candidate (s) | Change |
|---|---|---|---|
| `docexport`, normal | 1.464 | 1.144 | 0.320s saved (21.9%) |
| `docexport`, coverage | 2.124 | 1.117 | 1.007s saved (47.4%) |
| `docexport`, race | 8.351 | 6.567 | 1.784s saved (21.4%) |
| Scheduled export, normal | 0.41 | 0.01 | 0.40s saved (97.6%) |
| Scheduled export, coverage | 0.92 | 0.01 | 0.91s saved (98.9%) |
| Scheduled export, race | 0.58 | 0.08 | 0.50s saved (86.2%) |
| `scheduler`, normal | 0.008 | 0.009 | 0.001s more (12.5%) |
| `scheduler`, coverage | 0.044 | 0.040 | 0.004s saved (9.1%) |

[Regular CI](https://github.com/WiseLabz/WiseLabz/actions/runs/36508410636)
passes on the implementation revision, including the expanded race shard,
coverage, PostgreSQL shards, static/lint/vulnerability checks, build and
compose smoke. The initial attempt failed in the unchanged
`TestConnectorsSyncAcceptsFieldsHint/fields` during `TempDir` removal
(`directory not empty`); the same-revision confirmation passed. A separate
100-run baseline check did not reproduce that cleanup failure. The initial
failure is retained in the run's first attempt and local evidence; this
patch does not modify the API fixture or sync cleanup.

## Follow-ups

- Further parallel test rollout still needs its own stress validation and
  adoption decision; this change adds no `t.Parallel()` calls.
