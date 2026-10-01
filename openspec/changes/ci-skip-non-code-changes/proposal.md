# Proposal

## Why

PRs that change only docs, agent tooling, the knowledge graph or release metadata should cost close to nothing in CI. Today the picture is mixed:

- **Docs-only PRs already skip the heavy jobs.** `ci.yml` path-filters with `dorny/paths-filter`, so a PR like #428 (only `.gitignore`) runs just `Detect changes` + `CI Status`, about 11 seconds.
- **The filter is a hand-kept allowlist spread over three blocks and one extra workflow.** Nobody tests it. It has real gaps:
  - `.golangci.yml` is not in the `backend` filter, so a PR that changes lint rules never runs lint.
  - `.dockerignore` is not in `compose`.
  - `.github/codeql/**` is not in `codeql.yml`'s `paths`.
  - A `*.md` file under `web/` or `backend/` would start the whole suite.
- **The release-please PR (#248) re-runs frontend lint, tests, build and compose-smoke (~2.5 min) after every merge to main.** It only rewrites `CHANGELOG.md`, `.release-please-manifest.json` and the `version` field of `web/package.json`. That is the biggest recurring waste in the current setup.

We can't just skip the workflow for non-code PRs with workflow-level `paths-ignore`. The `main` ruleset makes `CI Status` a required check, and a workflow that never starts leaves that check pending forever, so the PR could not merge. The workflow still has to start, but for non-code changes it must do nothing beyond classifying the change and reporting success.

## What Changes

- Replace the three inline `dorny/paths-filter` lists with one tested classifier, `scripts/ci/changes.sh`, and a rules table. The classifier:
  - Sorts each changed path into **ignored** (no CI value), one or more **areas** (`backend`, `frontend`, `compose`), or **unknown**.
  - Treats **unknown paths as code for every area** (fail-safe: a new top-level file runs everything until someone classifies it).
  - Ignores these globally: `*.md` anywhere, `docs/**` except `docs/openapi.yaml`, `graphify-out/**`, `.graphifyignore`, `openspec/**`, `.claude/**`, `.agents/**`, `.codex/**`, `.idea/**`, `prototypes/**`, `deploy/**`, `LICENSE`, `.github/CODEOWNERS`, `.github/ISSUE_TEMPLATE/**`, `.github/PULL_REQUEST_TEMPLATE.md`, `cliff.toml`, `skills-lock.json`, `.gitignore`, `.air.toml`, `lefthook.yml`, `Makefile`, `release-please-config.json`, `.release-please-manifest.json`, `CHANGELOG.md`.
  - On pull requests, treats a `web/package.json` diff that only changes `"version"` as ignored. This makes the release-please PR skip every heavy job. The push to `main` after a release still runs the full frontend suite.
  - Has a `check` mode that runs a fixture table of path → expected areas, so a rule edit can't silently break gating.
- Fix the filter gaps: `.golangci.yml` → `backend`; `.dockerignore` → `compose`; `.github/codeql/**` and a `!**/*.md` exclusion → `codeql.yml` `paths`.
- Write the classification (areas plus each path's reason) to the job's step summary, so a skipped run explains itself.
- Add a `workflows` area and a cheap `actionlint` job, so edits to `.github/workflows/**` and `.github/actions/**` are validated. Workflow files other than `ci.yml` then no longer need the full suite.
- Smaller optimizations in scope:
  - Move the Bun setup and `node_modules` cache steps, which are repeated in `lint-frontend`, `test-frontend` and `build`, into one composite action, `.github/actions/bun-setup`.
  - Add a `concurrency` group with `cancel-in-progress` to `codeql.yml`. It has none, so every push to a PR queues another full CodeQL build.
- Skip heavy jobs on **draft PRs**. `ci.yml` and `codeql.yml` add `ready_for_review` to their `pull_request` types, so marking a PR ready runs the full suite on the same commit. Change detection, the fixture check and `actionlint` still run on drafts.
- Run **`test-backend-postgres` only when its Go dependency closure changes**. `internal/store/...` and `cmd/migrate` import 10 of the 28 `internal/` packages. A PR that touches only `internal/api`, `internal/chat` and similar no longer starts a Postgres service. `go list -deps -test` computes the closure at run time, so it never goes stale.
- Run **`govulncheck` on PRs only when Go module files change** (`go.mod`, `go.sum`, `go.work*`), plus a new **nightly scheduled run** on `main`. A new advisory against an unchanged dependency, or a code change that makes a known-vulnerable function reachable, is then found within 24 h instead of only when someone happens to open a backend PR.
- **"CodeQL – Code Quality" (GitHub's built-in Code Quality scan)**: check whether it supports path or trigger settings. If it does, configure it to skip docs-only PRs. If not, keep it enabled and document in `docs/TESTING.md` why it runs on every PR.
- Follow-up tracked separately: Dependabot for `github-actions` and SHA-pinning `actions/checkout@v6` → #429.

No breaking changes. `CI Status` stays the single required check, with the same name.

## Capabilities

### New Capabilities
- `ci-change-detection`: which changed paths start which CI jobs. Covers the ignore list, area mapping, the fail-safe for unknown paths, the release-please version-only rule, draft-PR behavior, Postgres and `govulncheck` gating, the nightly vulnerability scan, required-check behavior on skipped runs, and self-test of the rules.

### Modified Capabilities
<!-- none: openspec/specs/ is empty -->

## Impact

- `.github/workflows/ci.yml`: the `changes` job (classifier replaces the inline filters; `fetch-depth: 2`; Go set up only when backend changed, to compute the Postgres closure); new `postgres` and `gomod` outputs; draft gating on heavy jobs; the frontend jobs use the new composite action.
- `.github/workflows/codeql.yml`: `paths` fixes, `concurrency`, draft gating, `ready_for_review`.
- New `.github/workflows/govulncheck-nightly.yml` (schedule + `workflow_dispatch`).
- Repository settings: possibly the Code Quality configuration, depending on what the check finds.
- New: `scripts/ci/changes.sh`, `scripts/ci/changes-fixtures.txt`, `.github/actions/bun-setup/action.yml`, an `actionlint` job (pinned `rhysd/actionlint` release).
- `docs/TESTING.md` / `CONTRIBUTING.md`: a short section on how CI decides what to run and how to classify a new path.
- No application code, API, or runtime dependency changes. `dorny/paths-filter` stays, but only to list changed files (it already handles PR/push/force-push base resolution).
