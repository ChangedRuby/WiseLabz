# Design

## Context

See proposal.md for why. The current state that shapes the approach:

- Every third-party action except `actions/checkout@v6` is already SHA-pinned with a `# vX.Y.Z` comment. That is the format Dependabot recognizes and keeps in sync when it bumps a SHA.
- Composite actions live in `.github/actions/{bun-setup,go-cache-restore,go-cache-save}/action.yml`. They reference `actions/setup-go`, `actions/cache{,/restore,/save}` and `oven-sh/setup-bun`.
- `scripts/ci/changes.sh` classifies changed paths by an ordered rules table. Any path no rule matches selects every area. `scripts/ci/changes-fixtures.txt` currently uses `.github/dependabot.yml` as a fail-safe example of an unclassified path.
- `.github/workflows/ci.yml` selects every area (`bfwmp`). Any actions bump that touches it therefore runs the full suite.
- The lefthook `commit-msg` regex accepts `feat|fix|docs|chore|refactor|test|ci|perf|revert` with an optional scope. `release-please-config.json` hides `ci`, `chore` and `test`.
- The Go workspace is `go.work` at the repo root with `use ./backend`. The web lockfile is the text-format `web/bun.lock`.
- The only `pull_request_target` workflow is `cache-cleanup.yml`, and it never checks out code.
- Eleven pinned actions declare `runs.using: node20`: setup-go, upload/download-artifact, the four docker actions, codeql-action init/analyze, golangci-lint-action and release-please-action. The rest (cache, setup-bun, paths-filter, sbom-action) are already `node24`, and trivy-action is composite with node24 internals.
- `aquasecurity/trivy-action` is pinned to `a9c7b0f…`. That is the SHA of the annotated tag object for v0.36.0, not the commit (`ed142fd…`). Git resolves it, but it isn't a commit SHA.
- `release.yml` (tag push) and `release-please.yml` (push to main) don't run on pull requests, so a PR can't exercise their action upgrades. `release.yml` may run on a self-hosted runner (`vars.RELEASE_RUNNER`).

## Goals / Non-Goals

**Goals:**
- Pin every action reference immutably, and let Dependabot keep the pins current.
- Run every action on Node 24 so the Node 20 deprecation warnings disappear.
- Keep update PR volume low: at most five PRs a month (actions, plus minor/patch and majors for Go and web).
- Update PRs pass repository conventions (commit-msg, changelog visibility, labels) without any manual touch-up.

**Non-Goals:**
- Auto-merging Dependabot PRs. The `automerge` label exists, but we don't use it here. That is a separate decision.
- Changing any Node.js version used by the project itself. The web toolchain is Bun, no workflow runs `setup-node`, and Node 26 isn't an Actions runtime option. The deprecation is only about each action's `runs.using`.
- Dependabot for the Docker base images in `Dockerfile` (the `docker` ecosystem). Worth a follow-up, but out of the issue's scope.
- Any enforcement of SHA-pinning in CI (for example a lint that fails on `@v`), beyond the documented rule and the grep check.

## Decisions

### `directories` with a glob for github-actions, not `directory: "/"`
The issue assumed `directory: "/"` also covers `.github/actions/*`. It does not: for `github-actions`, `/` scans `.github/workflows/` and a root `action.yml`. Composite actions in subdirectories must be listed explicitly. We use `directories: ["/", "/.github/actions/*"]`, and the glob picks up future composite actions automatically.
- *Alternative:* list each composite action directory. Rejected, because a new composite action would silently go unwatched.
- The group is a plain `groups.actions.patterns: ["*"]`. When a config entry uses `directories`, Dependabot builds one grouped PR across all listed directories, so workflows and composite actions land together.

### Checkout at v7.0.1 rather than v6.1.0
v7 is the current major. Its only behavioral change blocks checking out fork PR code under `pull_request_target`/`workflow_run`. Our one `pull_request_target` workflow (`cache-cleanup.yml`) never checks out code, and no `workflow_run` workflow exists. Pinning v6 would just make Dependabot's first PR a major bump. We resolve the SHA from the tag at apply time (`gh api repos/actions/checkout/commits/v7.0.1`; at planning time it was `3d3c42e5aac5ba805825da76410c181273ba90b1`).

### Upgrade Node 20 actions manually now, to their latest release
Targets were resolved at planning time. Re-resolve them at apply time and re-check `runs.using`:

| Action | From | To | Commit | Breaking changes vs our usage |
|---|---|---|---|---|
| actions/setup-go | v5.4.0 | v7.0.0 | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` | v6: node24 only. v7: ESM/deps. None. |
| actions/upload-artifact | v4.6.2 | v7.0.1 | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` | node24 only. None. |
| actions/download-artifact | v4.3.0 | v8.0.1 | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` | v5 changes the path for *by-ID* single downloads; we download by `name`/`pattern`. v8 fails on a digest mismatch (safer). None. |
| github/codeql-action/{init,analyze} | v3.36.2 | v4.38.2 | `2892aa5e19bbd11bc0cff5427e3b750a04d9e3c2` | node24. Same inputs. |
| golangci/golangci-lint-action | v7.0.1 | v9.3.0 | `ba0d7d2ec06a0ea1cb5fa41b2e4a3ab91d21278a` | v8 needs golangci-lint ≥ v2.1 (we use v2.13.2) and uses absolute paths with `working-directory` (annotation paths only). v9: node24. |
| googleapis/release-please-action | v4.3.0 | v5.0.0 | `45996ed1f6d02564a971a2fa1b5860e934307cf7` | node24 only. |
| docker/setup-buildx-action | v3.10.0 | v4.4.1 | `f87e5991a6d7451dcb8d9637bfbc97413f497069` | Removes deprecated inputs/outputs; we pass none. |
| docker/build-push-action | v6.15.0 | v7.4.0 | `c3c9e263c25d99ce0380d002d59b67737d91b0dc` | Removes `DOCKER_BUILD_NO_SUMMARY`/`…EXPORT_RETENTION_DAYS` envs and the legacy export tool; we use neither. |
| docker/login-action | v3.3.0 | v4.6.0 | `dbcb813823bdd20940b903addbd779551569679f` | node24 only. |
| docker/metadata-action | v5.7.0 | v6.2.0 | `dc802804100637a589fabce1cb79ff13a1411302` | node24 only. |
| aquasecurity/trivy-action | v0.36.0 | v0.36.0 | `ed142fd0673e97e23eac54620cfb913e5ce36c25` | Re-pin only: tag object → commit. |

- *Alternative:* let Dependabot's first grouped PR do this. Rejected, because it would mix ten majors with routine bumps in one unreviewed PR. Doing it here gives each major a reviewed breaking-change check, and Dependabot then starts from a clean, current baseline.
- *Alternative:* the minimal node24 major for each action (for example setup-go v6). Rejected, because Dependabot would immediately propose the next major, so we go straight to latest.

### Grouping and commit conventions
| Ecosystem | Dir | Group | Commit | Label |
|---|---|---|---|---|
| github-actions | `/`, `/.github/actions/*` | `actions`: all updates | `ci: ...` (`prefix: ci`) | `area:platform` |
| gomod | `/backend` | `go-minor-patch`: `update-types: [minor, patch]` | `chore(deps): ...` (`prefix: chore`, `include: scope`) | `area:backend` |
| bun | `/web` | `web-minor-patch`: `update-types: [minor, patch]` | `chore(deps): ...` | `area:frontend` |

- Actions bumps get `ci` because they only affect CI, and the issue asked for it.
- Dependency bumps get `chore(deps)`, which matches `CONTRIBUTING.md` ("`chore/` … dependency bumps"). A bump that fixes a vulnerability can be relabelled `fix` by hand if it should show up in release notes.
- Majors stay ungrouped so each breaking upgrade can be reviewed, and reverted, on its own.
- Setting custom `labels` replaces Dependabot's default `dependencies`/ecosystem labels. Those don't exist in the repo, and Dependabot would create them. The area labels already exist and match the triage scheme.
- `open-pull-requests-limit: 5` per ecosystem caps the ungrouped major PRs.

### `.github/dependabot.yml` is ignored (`-`), not `w`
The `workflows` area runs actionlint, and actionlint does not validate the Dependabot config. GitHub validates it on push and reports errors in the Dependabot UI, so running any CI job for it buys nothing. The new rule goes with the other `.github/*` ignore rules, above the `.github/workflows/*` catch-alls. The fail-safe fixture that used this path moves to `.github/FUNDING.yml`, which is also unclassified.

### Reducing Dependabot CI load (added after the first run)
The first run opened 7 PRs: 1 Go group, 1 web group of 41, 4 separate web majors and 1 actions PR. `main` also required branches to be up to date, so merging N PRs one by one cost about N(N+1)/2 CI runs, because every merge put the rest behind. Changes:
- **Majors grouped per ecosystem** (`go-major`, `web-major`): at most 2 PRs per ecosystem per cycle. The trade-off is that one breaking major holds the group, so the fix is on the PR or an `ignore` entry. The minor/patch group is never blocked by a major.
- **Monthly schedule + `cooldown.default-days: 7`** for all three ecosystems. Security updates are unaffected, because they bypass both.
- **Merge queue instead of strict up-to-date checks.** `ci.yml` gains `merge_group: [checks_requested]`. `dorny/paths-filter` diffs `merge_group.base_sha..head_sha`, so classification works unchanged. Caches are saved only on push to main, so queue runs write none. The version-only release rule is `pull_request`-only, so a release PR runs the frontend jobs once in the queue (fail-safe). The ruleset switches after the workflow change is on `main`, because otherwise queue entries would wait for a `CI Status` that never comes.
- Rejected: `rebase-strategy: disabled`. With the merge queue, out-of-date branches no longer matter, and it would also stop conflict rebases.

## Risks / Trade-offs

- **Stale `go.work.sum`.** Dependabot updates `backend/go.mod`/`go.sum` but not the root `go.work.sum`, so a Go bump could leave it out of date → if a backend job then fails with a checksum error, a maintainer runs `go work sync` (or `go mod download` in workspace mode) and pushes the result onto the PR. We accept this and document it in `docs/TESTING.md`. If it becomes frequent, a follow-up can revisit whether `go.work.sum` needs to be committed.
- **Bun lockfile support.** Dependabot's `bun` ecosystem supports only the text `bun.lock` (which we use), not `bun.lockb` → no action needed. If its updates ever disagree with the local Bun version, the frontend CI `bun install --frozen-lockfile` step fails visibly.
- **Restricted token on Dependabot PRs.** Dependabot-triggered `pull_request` runs get a read-only `GITHUB_TOKEN` and only Dependabot secrets → `ci.yml` uses no repository secrets, but some jobs request `actions: write` (cache writes). Verify on the first Dependabot PR that cache-save steps degrade gracefully and don't fail the job. If they do fail, gate those steps on `github.actor != 'dependabot[bot]'`.
- **Full-suite cost of actions PRs.** Any bump touching `ci.yml` runs every area, about as expensive as a normal backend+frontend PR → acceptable at one PR a month, and it genuinely re-validates CI with the new action versions.
- **Release-only workflows aren't exercised by the PR.** The docker actions, download-artifact in `release.yml` and release-please-action run only on a tag push or a push to main → run actionlint on the PR, then watch the first `release-please` run after merge. For `release.yml`, either watch the next tagged release or run a prerelease tag (`vX.Y.Z-rc.1`, which the `latest` tag rule already excludes) to smoke-test it.
- **Self-hosted release runner.** Node 24 actions require Actions Runner v2.327.1 or later → if `RELEASE_RUNNER` is set, confirm the runner's version before the next release. The GitHub-hosted runners already qualify.
- **Major bumps of pinned actions.** Dependabot will also propose major action upgrades inside the grouped PR → reviewed like any PR, and a problematic one can be excluded with an `ignore` entry.

## Migration Plan

1. Merge the PR: pins, config, classifier rule, docs.
2. Dependabot runs on its first schedule, or trigger it from Insights → Dependency graph → Dependabot → "Check for updates".
3. Confirm that Dependabot recognizes all three ecosystems with no config errors, and that the first grouped PRs pass `CI Status`.

Rollback: delete `.github/dependabot.yml` to stop updates. The SHA pins are harmless to keep.
