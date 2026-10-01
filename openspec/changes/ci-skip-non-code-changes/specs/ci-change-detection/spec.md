# Spec Delta

## Purpose

Decides which CI jobs a push or pull request must run from the files it changes. Changes with no CI value (documentation, agent tooling, generated knowledge graphs, release metadata) skip every build, lint and test job, while the required status check still reports a result.

## ADDED Requirements

### Requirement: Non-code changes skip all heavy CI jobs
CI SHALL classify every changed path, and SHALL skip every lint, test, build, vulnerability-scan and compose-smoke job when no changed path maps to a CI area. Paths with no CI value SHALL include at least:
- Markdown files (`*.md`) at any depth
- `docs/**` except `docs/openapi.yaml`
- `graphify-out/**` and `.graphifyignore`
- `openspec/**`, `.claude/**`, `.agents/**`, `.codex/**`, `.idea/**`, `prototypes/**`, `deploy/**`
- Repository metadata: `LICENSE`, `.github/CODEOWNERS`, `.github/ISSUE_TEMPLATE/**`, `.github/PULL_REQUEST_TEMPLATE.md`, `.gitignore`, `cliff.toml`, `skills-lock.json`
- Local developer tooling: `.air.toml`, `lefthook.yml`, `Makefile`
- Release metadata: `CHANGELOG.md`, `release-please-config.json`, `.release-please-manifest.json`

#### Scenario: Docs-only pull request
- **WHEN** a pull request changes only `README.md`, `docs/ARCHITECTURE.md` and `graphify-out/GRAPH_REPORT.md`
- **THEN** every backend, frontend, build and compose-smoke job is skipped
- **AND** the `CI Status` check completes successfully

#### Scenario: Markdown inside a code directory
- **WHEN** a pull request changes only `web/README.md`
- **THEN** no frontend job runs

#### Scenario: OpenAPI contract is not treated as docs
- **WHEN** a pull request changes only `docs/openapi.yaml`
- **THEN** the frontend jobs and compose-smoke run

### Requirement: Required status check always reports
The `CI Status` check SHALL be reported on every pull request to `main` and every push to `main`, whatever the changed paths, so a required-check ruleset never leaves a pull request waiting. It SHALL fail when any job that ran failed or was cancelled, and SHALL succeed when every job either passed or was skipped by classification.

#### Scenario: Skipped run still satisfies the ruleset
- **WHEN** a pull request changes only ignored paths
- **THEN** the CI workflow runs, and `CI Status` reports success without waiting on any skipped job

#### Scenario: Failure still blocks
- **WHEN** a pull request changes `backend/` and the backend lint job fails
- **THEN** `CI Status` reports failure

### Requirement: Code paths map to their CI areas
CI SHALL map changed paths to areas, and SHALL run each area's jobs when at least one changed path maps to that area, except where a requirement below narrows a specific job (Postgres tests, vulnerability scan, draft PRs):
- **backend**: `backend/**`, `go.work`, `go.work.sum`, `.golangci.yml`, `scripts/ci/**`, the Go cache composite actions
- **frontend**: `web/**`, `docs/openapi.yaml`, the Bun setup composite action
- **compose**: everything in backend and frontend, plus `Dockerfile`, `.dockerignore`, `docker-compose*.yml`, `.env.example`, `scripts/compose-smoke.*`
- **workflows**: `.github/workflows/**`, `.github/actions/**`. Selecting this area runs a workflow linter and nothing else.

A change to the CI workflow definition itself SHALL select every area. Changes to other workflow files, or to `.github/codeql/**`, SHALL select only the workflows area.

#### Scenario: Lint configuration change runs backend lint
- **WHEN** a pull request changes only `.golangci.yml`
- **THEN** the backend lint job runs

#### Scenario: Docker build-context change runs compose-smoke
- **WHEN** a pull request changes only `.dockerignore`
- **THEN** compose-smoke runs and the backend and frontend unit-test jobs do not

#### Scenario: Other workflow edit is linted only
- **WHEN** a pull request changes only `.github/workflows/release.yml`
- **THEN** only the workflow lint job runs, and the backend, frontend and compose areas are not selected

#### Scenario: Workflow edit runs everything
- **WHEN** a pull request changes only `.github/workflows/ci.yml`
- **THEN** the backend, frontend and compose areas are all selected

### Requirement: Unclassified paths fail safe
A changed path that matches neither an ignore rule nor an area rule SHALL select every area. Unrecognized files must never cause CI to be skipped.

#### Scenario: New top-level file
- **WHEN** a pull request adds `renovate.json`, which no rule covers
- **THEN** the backend, frontend and compose areas are all selected (not workflows)
- **AND** the classification output marks the path as unclassified

### Requirement: Release version bumps skip frontend CI
On pull requests, a change to `web/package.json` that differs from the pull request's base only in its top-level `version` field SHALL be treated as ignored. Any other difference in that file SHALL select the frontend and compose areas. On pushes to `main`, and whenever the base version of the file cannot be read, any `web/package.json` change SHALL select the frontend and compose areas.

#### Scenario: release-please pull request
- **WHEN** a pull request changes only `CHANGELOG.md`, `.release-please-manifest.json` and the `version` field of `web/package.json`
- **THEN** every heavy job is skipped and `CI Status` succeeds

#### Scenario: Release merged to main
- **WHEN** the release-please pull request is merged and the push to `main` includes the `web/package.json` version bump
- **THEN** the frontend jobs and compose-smoke run on `main`

#### Scenario: Dependency bump in package.json
- **WHEN** a pull request changes a dependency version in `web/package.json`
- **THEN** the frontend jobs and compose-smoke run

### Requirement: Classification is explained and self-tested
Every CI run SHALL publish, in the change-detection job summary, the selected areas and the rule that matched each changed path. The classification rules SHALL have a fixture table of paths and expected areas. CI SHALL check that table on every run and SHALL fail change detection when any fixture disagrees with the rules.

#### Scenario: Reading why a run was skipped
- **WHEN** a docs-only pull request's CI run finishes
- **THEN** the change-detection summary lists each changed path as ignored, with its matching rule

#### Scenario: Broken rule edit
- **WHEN** a rules edit makes `backend/main.go` stop mapping to the backend area
- **THEN** the fixture check fails and `CI Status` reports failure

### Requirement: CodeQL scans only on relevant changes
The CodeQL workflow SHALL run on pull requests and pushes to `main` only when Go sources, Go workspace files, its own workflow file or its CodeQL configuration (`.github/codeql/**`) change, and SHALL NOT run for Markdown-only changes under `backend/`. When a newer commit is pushed to the same pull request, the CodeQL run still in progress for the older commit SHALL be cancelled. The weekly scheduled scan SHALL continue unchanged.

#### Scenario: CodeQL config change
- **WHEN** a pull request changes only `.github/codeql/codeql-config.yml`
- **THEN** the CodeQL workflow runs

#### Scenario: Rapid pushes
- **WHEN** two commits are pushed to the same pull request while the first CodeQL run is still running
- **THEN** the first run is cancelled

### Requirement: Draft pull requests defer heavy jobs
While a pull request is a draft, CI SHALL run only change detection, the rules self-test and the workflow linter, and SHALL skip every other job. `CI Status` SHALL still report. When a draft is marked ready for review, CI SHALL run again on the same commit with the full job selection for its changed paths, and that run's `CI Status` SHALL be the one the ruleset evaluates. CodeQL SHALL likewise skip drafts and run when a draft is marked ready for review.

#### Scenario: Pushing to a draft
- **WHEN** a commit changing `backend/internal/api/router.go` is pushed to a draft pull request
- **THEN** no backend lint, test, build or compose-smoke job runs, and `CI Status` succeeds

#### Scenario: Marking ready for review
- **WHEN** that draft pull request is marked ready for review without a new commit
- **THEN** CI runs the backend, build and compose jobs for the same commit

#### Scenario: Broken workflow on a draft
- **WHEN** a draft pull request changes a workflow file so that the linter reports an error
- **THEN** `CI Status` reports failure

### Requirement: Postgres tests run only when their dependencies change
The Postgres-backed test job SHALL run when the backend area is selected and at least one changed path belongs to the Go dependency closure of the Postgres test packages and the migration command. A path belongs to that closure if it is a file in one of those packages' directories, a file those packages embed, or test data under those directories. The closure SHALL be computed from the code being tested, not from a maintained list. The job SHALL also run when Go module or workspace files, CI scripts, the Go cache actions or the CI workflow change, when a path is unclassified, or when the closure cannot be computed.

#### Scenario: Handler-only change
- **WHEN** a pull request changes only `backend/internal/api/handlers.go` (a package outside the closure)
- **THEN** backend lint, unit and race tests run, and the Postgres test job is skipped

#### Scenario: Store change
- **WHEN** a pull request changes `backend/internal/store/doc.go`
- **THEN** the Postgres test job runs

#### Scenario: Migration file change
- **WHEN** a pull request changes only a file under `backend/internal/store/migrations/postgres/`
- **THEN** the Postgres test job runs

#### Scenario: Transitive dependency change
- **WHEN** a pull request changes only `backend/internal/crypto/` (imported by the store)
- **THEN** the Postgres test job runs

### Requirement: Vulnerability scanning runs on dependency changes and nightly
On pull requests and pushes to `main`, the Go vulnerability scan SHALL run only when `backend/go.mod`, `backend/go.sum`, `go.work`, `go.work.sum` or the CI workflow change. The scan SHALL also run on a daily schedule against `main` and SHALL be manually dispatchable. A failing scheduled scan SHALL fail its workflow run, so it is visible in the Actions tab and to repository notifications.

#### Scenario: Code-only backend change
- **WHEN** a pull request changes only `backend/internal/api/router.go`
- **THEN** the vulnerability scan does not run in that pull request's CI

#### Scenario: Dependency bump
- **WHEN** a pull request changes `backend/go.sum`
- **THEN** the vulnerability scan runs

#### Scenario: New advisory against an unchanged dependency
- **WHEN** a vulnerability affecting a reachable function in an existing dependency is published
- **THEN** the next daily scheduled scan on `main` fails and reports it
