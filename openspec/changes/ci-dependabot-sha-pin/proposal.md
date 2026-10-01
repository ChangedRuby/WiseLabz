# Proposal

## Why

Dependencies only move when someone remembers to bump them. There is no `.github/dependabot.yml`, so SHA-pinned actions (`oven-sh/setup-bun`, `actions/cache`, `actions/setup-go`, `docker/build-push-action`, `golangci/golangci-lint-action`, `github/codeql-action`, ...), Go modules and web packages drift until a security advisory or a breaking runner change forces a rushed bump. And `actions/checkout@v6` is the only third-party action still on a mutable tag: 17 references in 6 workflow files (`ci.yml` ×12, `codeql.yml`, `release.yml`, `test-profile.yml`, `test-stress.yml`, `govulncheck-nightly.yml`). If that tag is retargeted, every job runs code nobody reviewed.

Separately, GitHub has deprecated Node.js 20 for action runtimes, and it now force-runs Node 20 actions on Node 24 with a warning on every job (for example "actions/setup-go@0aaccfd… targets Node.js 20 but is being forced to run on Node.js 24"). Eleven of our pinned actions still declare `runs.using: node20`. When the runner drops Node 20 entirely, the forced run is no longer guaranteed to work.

Tracked as #429, split out of `ci-skip-non-code-changes`. The Node 24 action upgrade is added to this change because it touches the same pins and would otherwise arrive as an unreviewed batch of majors in Dependabot's first PR.

## What Changes

- Add `.github/dependabot.yml` with three ecosystems, all on a weekly schedule:
  - **github-actions**: covers workflows and the composite actions under `.github/actions/*`. All bumps come in one grouped PR, and commits use the `ci` Conventional Commits type.
  - **gomod** (`/backend`) and **bun** (`/web`): minor and patch bumps are grouped into one PR per ecosystem, and each major bump gets its own PR. Commits use `chore(deps)`.
  - Every Dependabot PR is labelled with its area (`area:platform`, `area:backend` or `area:frontend`).
- Pin every `actions/checkout` reference to the commit SHA of **v7.0.1** with a `# v7.0.1` comment. This goes straight to the current major instead of pinning v6 and taking the major bump in Dependabot's first PR. After this, every non-local action in `.github/` is SHA-pinned, and Dependabot keeps each SHA and its comment in sync.
- Upgrade every action that declares a Node 20 runtime to the first release line on Node 24, and pin each at its latest release SHA:
  - `actions/setup-go` v5 → v7
  - `actions/upload-artifact` v4 → v7
  - `actions/download-artifact` v4 → v8
  - `github/codeql-action` v3 → v4
  - `golangci/golangci-lint-action` v7 → v9
  - `googleapis/release-please-action` v4 → v5
  - `docker/setup-buildx-action` v3 → v4, `docker/build-push-action` v6 → v7, `docker/login-action` v3 → v4, `docker/metadata-action` v5 → v6

  None of the breaking changes in these majors affect the inputs we use. (The project's own code needs no Node change. The web app uses Bun, and no workflow installs Node.)
- Re-pin `aquasecurity/trivy-action` v0.36.0 from its annotated-tag object SHA (`a9c7b0f…`) to the commit that tag points to (`ed142fd…`). The version doesn't change, but the pin becomes a real commit SHA that Dependabot can track.
- Classify `.github/dependabot.yml` as having no CI value in `scripts/ci/changes.sh`, with a fixture for it. It is currently unclassified, so editing the config would run the whole suite. The fixture that uses it as the unclassified fail-safe example switches to another path.
- Document the pinning rule and the Dependabot setup in `docs/TESTING.md` ("What CI runs").

No breaking changes. `CI Status` stays the single required check.

## Capabilities

### New Capabilities
- `ci-dependency-updates`: CI dependencies stay pinned and current. Covers SHA-pinning of every action reference, a supported (non-deprecated) Node runtime for every action, automated weekly update PRs per ecosystem with grouping and commit conventions, and how CI treats the update configuration and update PRs.

### Modified Capabilities
<!-- none: openspec/specs/ is empty; ci-change-detection is still an unarchived change -->

## Impact

- New: `.github/dependabot.yml`.
- `.github/workflows/{ci,codeql,release,test-profile,test-stress,govulncheck-nightly}.yml`: checkout pin (v6 → v7.0.1).
- `.github/workflows/{ci,codeql,release,release-please,test-profile,test-stress}.yml` and `.github/actions/go-cache-restore/action.yml`: Node 24 action upgrades and the trivy re-pin.
- The self-hosted release runner, if `RELEASE_RUNNER` is set, must be Actions Runner v2.327.1 or later.
- `scripts/ci/changes.sh` (rules table) and `scripts/ci/changes-fixtures.txt`.
- `docs/TESTING.md`.
- Ongoing effect: up to three weekly Dependabot PRs, plus a separate PR for each major Go or web bump. An actions PR that touches `ci.yml` runs the full suite, because `ci.yml` selects every area.
- No application code, API or runtime behavior changes.
