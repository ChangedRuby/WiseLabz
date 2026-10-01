# Design

## Context

See proposal.md, Why, for the motivation. The constraints that shape the approach:

- **The ruleset requires `CI Status` (integration 15368, strict).** The workflow must start on every PR and push to `main`, and `ci-status` must always report. Skipping at the trigger level (`paths-ignore`) is ruled out.
- **`ci.yml` already has a `changes` job feeding `if:` conditions on every heavy job.** Only *how* `changes` decides has to change. Its outputs (`backend`, `frontend`, `compose`, `race_shards`, `postgres_shards`) and the jobs' `if:` expressions stay as they are.
- **The repo already has a pattern for CI rules: `scripts/ci/test-shards.sh` + `test-shards.json`.** It is the single source of truth, has a `check` mode run in CI and is driven by a Bash script. The classifier follows the same pattern.
- **Verified inputs to the rules:**
  - `git ls-files` finds no `*.md` under `backend/` or `web/`.
  - The only Markdown-like file Go embeds is `backend/internal/report/templates/report.md.tmpl`, and `*.md` does not match it.
  - `web/` never reads `package.json`'s `version`.
  - Nothing in CI or compose reads `deploy/`. `config.go` searches `./deploy/config.yaml`, which is not tracked.

## Goals / Non-Goals

**Goals:**
- One place to edit when a new path appears, with a test that catches a bad edit.
- Unclassified paths fail safe (run everything).
- A skipped run explains itself in the job summary.

**Non-Goals:**
- Dependency-closure gating for anything other than the Postgres job. The race job's package list (`test-shards.json`) already spans most of `internal/`, so gating it would rarely skip anything.
- Changing the required-check setup or the ruleset.
- Dependabot and SHA-pinning `actions/checkout` (tracked in #429).

## Decisions

### D1. Keep `dorny/paths-filter` to list files; classify in `scripts/ci/changes.sh`
`dorny/paths-filter` runs with one catch-all filter (`all: '**'`, `list-files: json`) and emits the changed-file list. It already resolves the right base for `pull_request`, for `push` (via `event.before`) and for a new branch or force-push. Re-implementing that is error-prone. The list goes to `scripts/ci/changes.sh classify`. The script prints `key=value` lines for `$GITHUB_OUTPUT` and a Markdown table for `$GITHUB_STEP_SUMMARY`.

*Alternatives considered:*
- **Keep dorny's own filters with `predicate-quantifier: every` and negations (`!**/*.md`).** This works, but `every` needs each positive list folded into one brace glob. It still can't express "unknown → all" or the `package.json` version rule, and it can't be tested offline.
- **A pure `git diff` script with no dorny.** It needs `fetch-depth: 0` or manual base fetching for push events, and it duplicates dorny's edge cases.

### D2. Rules are an ordered Bash `case` table; first match wins
In Bash `case`, `*` matches `/`, so `backend/*` covers the whole subtree and `*.md` matches Markdown at any depth. No extglob or `globstar` is needed. Order:
1. **Explicit code exceptions:**
   - `.github/workflows/ci.yml` → all areas
   - `docs/openapi.yaml` → frontend, compose
   - `web/package.json` → the version-only check (D3)
2. **Ignore rules**, the spec's list. `*.md` sits here, ahead of the area rules, so `web/README.md` is ignored.
3. **Area rules:**
   - backend: `backend/*`, `go.work*`, `.golangci.yml`, `scripts/ci/*`, `.github/actions/go-cache-*`
   - frontend: `web/*`, `.github/actions/bun-setup/*`
   - compose-only: `Dockerfile`, `.dockerignore`, `docker-compose*.yml`, `.env.example`, `scripts/compose-smoke.*`
   - workflows: `.github/workflows/*`, `.github/actions/*`, `.github/codeql/*`
4. **Fallback:** `unclassified` → backend, frontend and compose (not workflows).

The rules also derive:
- `compose = backend || frontend || compose-only`, matching today's filter.
- `workflows = true` whenever any `.github/workflows/*` or `.github/actions/*` path changed. A path under `.github/actions/*` is both linted and routed to its area.

*Alternative:* a JSON rules file read with `jq`, like `test-shards.json`. It was rejected because glob-with-order semantics would have to be re-implemented in jq. A `case` table is shorter and reads top-down.

### D3. `web/package.json` version-only rule, pull requests only
On `pull_request`, `actions/checkout` checks out the merge commit, so `HEAD^1` is the base tip. The `changes` job checks out with `fetch-depth: 2` and compares `jq -S 'del(.version)'` of `HEAD^1:web/package.json` and `HEAD:web/package.json`. If they are equal, the file counts as ignored with reason `version-only`. If they differ, or `HEAD^1` is missing, it counts as frontend. On `push`, `HEAD^1` may not be the push's `before` (rebase merges push several commits), so the rule is off and the file is always frontend. The script takes the event name via an env var (`CI_EVENT`), and in `check` mode the version rule is stubbed through a fixture flag.

### D4. Fixture table and `check` mode
`scripts/ci/changes-fixtures.txt` holds one `path<TAB>expected-areas` per line, where expected-areas is a comma list or `-` for ignored. It includes a line for every spec scenario, plus guard rows:
- `backend/internal/report/templates/report.md.tmpl → backend,compose`
- `docs/openapi.yaml → frontend,compose`
- `renovate.json → backend,frontend,compose`

`changes.sh check` classifies each fixture path on its own and diffs the result. The `changes` job runs `check` before `classify` on every run; it is pure Bash and takes milliseconds. A mismatch fails `changes`, which `ci-status` treats as a failure.

### D5. `ci-status` keeps its current logic
It already fails on `failure`/`cancelled` and passes on `skipped`. The only addition is the new `actionlint` job in its `needs`. A failed `changes` job makes the downstream jobs `skipped`, but `changes` itself shows `failure` in `needs.*.result`, so the check still fails.

### D6. `actionlint` job
It runs when `needs.changes.outputs.workflows == 'true'` and downloads a pinned `rhysd/actionlint` release binary via its `download-actionlint.bash`, pinned to a commit SHA. That is under 15 seconds and avoids a Docker pull. Its `shellcheck` integration is on by default, since `ubuntu-latest` ships shellcheck, so `run:` blocks are linted too. Before landing the job, run it once locally and fix any existing findings in the same PR.

### D7. `.github/actions/bun-setup` composite action
The action wraps `oven-sh/setup-bun` (same pin, `bun-version-file: web/package.json`), the `node_modules` `actions/cache` (same key) and a conditional `bun install --frozen-lockfile`. It replaces the three copies in `lint-frontend`, `test-frontend` and `build`. The cache key is unchanged, so existing caches stay valid.

### D8. `codeql.yml`
- `paths` gains `.github/codeql/**` and a trailing `!backend/**/*.md` (GitHub evaluates `paths` in order, and a later negation excludes).
- New block: `concurrency: { group: codeql-${{ github.event.pull_request.number || github.run_id }}, cancel-in-progress: ${{ github.event_name == 'pull_request' }} }`. Main-branch and scheduled runs each get their own group (the run id), because runs sharing a group drop each other's queued runs even without `cancel-in-progress`. Every `main` commit keeps its SARIF upload.

### D9. Draft pull requests
- `ci.yml` and `codeql.yml` add `types: [opened, synchronize, reopened, ready_for_review]` to `pull_request`.
- The `changes` job exposes `draft=${{ github.event.pull_request.draft == true }}`.
- Every heavy job's `if:` gains `&& needs.changes.outputs.draft != 'true'`. `changes`, `actionlint` and `ci-status` are not gated.
- The CodeQL `analyze` job gets `if: github.event.pull_request.draft != true` (always true for push and schedule).

A `ready_for_review` event on an unchanged head produces a new run for the same SHA. The ruleset evaluates the latest `CI Status` for the head SHA, so the full run replaces the draft run's green result. That run may even fail, which is the intent.

*Alternative:* a workflow-level `if` that skips drafts entirely. It was rejected because `CI Status` would then not report, and the rules self-test and workflow lint would not run on drafts.

### D10. Postgres job gated on the Go dependency closure
When `backend=true` (and the PR is not a draft), the `changes` job runs `./.github/actions/go-cache-restore` with `kind: lint`, which is restore-only and already brings in module downloads. It then runs `scripts/ci/changes.sh go-closure`, which calls:

```
go list -deps -test -f '{{if not .Standard}}{{.Dir}}|{{join .EmbedFiles ","}}|{{join .TestEmbedFiles ","}}{{end}}' <postgres packages from test-shards.json> ./cmd/migrate
```

It maps each changed path under `backend/` to its **nearest enclosing Go package**: walk up the path's directories until one is a package from `go list ./...`. The path touches the closure when that package is one of the closure's packages. This one rule covers embedded files (`internal/store/migrations/{sqlite,postgres}/*.sql` belong to `internal/store`), `testdata/`, non-Go files and deleted files. A path with no enclosing package (for example, a removed package directory) counts as affected. (It replaces an earlier plan to match `EmbedFiles` lists. Those lists miss deleted files and embedded subdirectories.)

The following are always affected:
- `backend/go.mod` and `backend/go.sum`
- `go.work*`
- `scripts/ci/*`
- `.github/actions/go-cache-*`
- `.github/workflows/ci.yml`
- any unclassified path
- a non-zero exit from `go list`

The result goes out as the output `postgres=true|false`. The Postgres job's `if:` becomes `needs.changes.outputs.postgres == 'true'`.

Package roots come from `scripts/ci/test-shards.json` (`jq '[.postgres[].packages[]] | unique'`), so the set of packages that run against Postgres is still defined in one place.

Measured on `main`: the closure is `cmd/migrate` plus `internal/{auth,config,connector,crypto,httputil,httpx,logsafe,store,store/storetest,storeerr,ws}`, 11 of the 28 `internal/` packages. `go list` takes about 0.1 s with a warm module cache.

*Alternatives considered:*
- **A static path list in `changes.sh`.** It goes stale the first time `store` imports a new package, which is exactly the silent-skip failure the fail-safe exists to prevent.
- **An in-job gate step, with the job starting and then exiting early.** The Postgres service container would still start (about 15 s per shard × 3 shards) for nothing.

### D11. `govulncheck` gating and nightly workflow
`changes.sh` emits `gomod=true` when any of the following change:
- `backend/go.mod`
- `backend/go.sum`
- `go.work`
- `go.work.sum`
- `.github/workflows/ci.yml`
- an unclassified path

The `govulncheck` job in `ci.yml` switches its `if:` to `gomod`.

A new `.github/workflows/govulncheck-nightly.yml`:
- runs on `schedule: cron "17 6 * * *"` (off the hour, to avoid GitHub's top-of-hour load) and on `workflow_dispatch`;
- checks out `main` and uses `go-cache-restore` with `kind: vuln` (restore-only; `main` CI saves that cache whenever `gomod` changes);
- runs the same pinned `govulncheck@v1.8.0` command, with `permissions: contents: read` and `timeout-minutes: 10`.

GitHub marks a failed scheduled run red and notifies the workflow's last committer by default. That is the alerting mechanism, and no extra integration is added.

*Trade-off accepted:* a PR that makes an already-vulnerable dependency function reachable passes CI and is caught by the next nightly run, within 24 h, instead of at PR time.

### D12. "CodeQL – Code Quality" (GitHub's built-in Code Quality scan): check before acting
The scan appears in run lists with `event: dynamic` and is set up in repository settings, not in a workflow file. `gh api repos/{owner}/{repo}/code-scanning/default-setup` reports `not-configured`, so it is separate from CodeQL code scanning's default setup. It is not known whether GitHub exposes path or trigger options for it. The implementation starts with a short, time-boxed check:
- the repository's Settings → Code security → Code quality page;
- GitHub Docs for "code quality" configuration;
- the REST API (`gh api` on `code-quality` / `code-scanning` settings endpoints).

It then takes exactly one of two outcomes:
- **(a)** A path or trigger option exists: configure it to skip the same ignore set as `changes.sh` (at least `**/*.md`, `docs/**`, `graphify-out/**`, `openspec/**`). Record the setting and how to change it in docs/TESTING.md.
- **(b)** No option exists: leave it enabled. Document in docs/TESTING.md that it runs on every PR, takes about 2.5 min, is not required for merge, and can't be path-filtered.

Neither outcome changes the specs or the other tasks.

**Findings (2026-09-29): outcome (b).** Code Quality has no path or trigger options. Sources checked:
- `gh api repos/WiseLabz/WiseLabz/code-quality/setup` returns `state: configured`, `languages: [go, javascript-typescript]`, `schedule: weekly`, `runner_type: standard`, `ai_findings_option: on_push`.
- The REST reference ([docs.github.com/en/rest/code-quality/code-quality](https://docs.github.com/en/rest/code-quality/code-quality)) accepts only `state`, `languages`, `runner_type`, `runner_label` and `ai_findings_option` on `PATCH`.
- The enablement guide ([docs.github.com/en/code-security/code-quality/how-tos/enable-code-quality](https://docs.github.com/en/code-security/code-quality/how-tos/enable-code-quality)) lists languages, runner type and repository access, and no path or event filtering.
- The run's workflow is `dynamic/github-code-quality/codeql`. There is no file in the repository to edit.

It stays enabled and is documented in docs/TESTING.md. The only lever is turning it off entirely, which is out of scope.

## Risks / Trade-offs

- **[A future `go:embed` or Vite import of a `.md` file under `backend/` or `web/` would be ignored]** → A header comment in `changes.sh` states the assumption. The fallback does not help here, because `*.md` is an explicit rule. Mitigation: a fixture row for the known `.md.tmpl` template, plus a note in docs/TESTING.md. Anyone adding an embedded `.md` adds an exception above the ignore rules.
- **[Ignoring `Makefile`/`lefthook.yml` means CI never checks them]** → CI doesn't check them today either. Neither is invoked by any workflow.
- **[The `package.json` rule relies on the merge-commit parent]** → The fallback on missing `HEAD^1` is "frontend", so the failure mode is running extra jobs, never skipping.
- **[One more Bash script to maintain]** → It is under 150 lines, with self-tests, following the `test-shards.sh` precedent.
- **[The Postgres closure misses a runtime-only dependency, such as a file read via `os.ReadFile` outside `testdata/`]** → Postgres tests would be skipped for that change. The unit and race jobs still run on every backend change, and SQLite covers the same store code, so only Postgres-specific behavior is exposed. Mitigation: store tests use `embed` and `testdata/`, which the closure covers, and a fixture-style check in `changes.sh check` asserts that `internal/store/migrations/postgres/*.sql` and `internal/crypto/*.go` map to `postgres=true`.
- **[Setting up Go in `changes` slows every backend PR by about 10–20 s]** → Only when `backend=true`. Docs-only and frontend-only runs are unaffected. The time is paid back whenever the three Postgres shards and their service containers are skipped.
- **[A draft PR shows a green `CI Status` for untested code]** → Drafts can't be merged, and `ready_for_review` re-runs everything on the same SHA. The draft run's summary states "draft: heavy jobs deferred" so the green result isn't misread.
- **[Nightly `govulncheck` failures go unnoticed]** → A failed scheduled run notifies by email by default. docs/TESTING.md names the workflow as something to watch.
- **[A job-name or output-name change could silently disconnect an `if:`]** → Outputs keep their current names. The spec scenarios are checked once by hand on a draft PR (see tasks).

## Migration Plan

It ships as one PR. Because that PR edits `ci.yml`, its own run selects every area and exercises every job end to end. After merge, verify with throwaway draft PRs: docs-only, `.golangci.yml`-only, `release.yml`-only. Watch the next release-please PR update. To roll back, revert the PR. The ruleset and job names don't change, so a revert needs no settings changes.

## Open Questions

- The outcome of the D12 check (a or b). Both are planned, and neither changes the specs or the other tasks.
