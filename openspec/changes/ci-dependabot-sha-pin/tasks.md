# Tasks

## 1. Pin actions/checkout

- [x] 1.1 Resolve the v7.0.1 commit with `gh api repos/actions/checkout/commits/v7.0.1 --jq .sha` (expected `3d3c42e5aac5ba805825da76410c181273ba90b1`), and confirm that `gh api repos/actions/checkout/git/ref/tags/v7.0.1` points at the same commit (dereferencing an annotated tag if needed)
- [x] 1.2 Replace every `uses: actions/checkout@v6` with `uses: actions/checkout@<sha> # v7.0.1` in `ci.yml`, `codeql.yml`, `release.yml`, `test-profile.yml`, `test-stress.yml` and `govulncheck-nightly.yml`, keeping each `with:` block unchanged. Verify that `grep -rn 'actions/checkout@' .github | grep -vc '<sha> # v7.0.1'` prints `0` and that 17 references remain
- [x] 1.3 Verify that `grep -rnE 'uses: [^.].*@v[0-9]' .github` prints nothing, and that actionlint passes locally (`actionlint -color`, or the download script pinned in `ci.yml`)

## 2. Upgrade Node 20 actions to Node 24

- [x] 2.1 For each action in the design's upgrade table, re-resolve the target: take the latest release tag (for codeql-action, the latest `v4.x.y` tag, not a `codeql-bundle-*` release), get its commit with `gh api repos/<owner>/<repo>/commits/<tag> --jq .sha`, and confirm that `runs.using` in its `action.yml` at that SHA is `node24`. Update the table in design.md if any target has moved
- [x] 2.2 Bump `actions/setup-go` in `.github/actions/go-cache-restore/action.yml` and `.github/workflows/codeql.yml`, and `github/codeql-action/{init,analyze}` in `codeql.yml`, to the new `<sha> # vX.Y.Z` pins. Verify that the CodeQL workflow and the Go jobs pass on the PR
- [x] 2.3 Bump `actions/upload-artifact` (in `ci.yml` ×2, `test-profile.yml` and `test-stress.yml`) and `actions/download-artifact` (in `test-profile.yml` and `release.yml`). Verify that the coverage artifacts still upload in the PR's CI run, and that `download-artifact` is only called with `name:`/`pattern:` (`grep -n -A3 download-artifact .github/workflows/*.yml`)
- [x] 2.4 Bump `golangci/golangci-lint-action` in `ci.yml`, keeping `version: v2.13.2`. Verify that `Lint (Backend)` passes on the PR
- [x] 2.5 Bump `docker/setup-buildx-action` and `docker/build-push-action` (in `ci.yml` and `release.yml`), `docker/login-action` and `docker/metadata-action` (in `release.yml`), and `googleapis/release-please-action` (in `release-please.yml`). Verify that `compose-smoke` passes on the PR (buildx + build-push with the gha cache)
- [x] 2.6 Re-pin `aquasecurity/trivy-action` (both uses in `release.yml`) from `a9c7b0f…` to commit `ed142fd0673e97e23eac54620cfb913e5ce36c25 # v0.36.0`. Verify that every pinned SHA in `.github` resolves as a commit: `grep -rhoE '[a-z0-9_.-]+/[a-z0-9_.-]+(/[a-z0-9_./-]+)?@[0-9a-f]{40}' .github | sort -u | while IFS=@ read -r r s; do gh api "repos/$(cut -d/ -f1-2 <<<"$r")/commits/$s" --silent || echo "BAD $r@$s"; done` prints nothing
- [x] 2.7 Verify that no pinned action is left on Node 20: for each unique pin, read `runs.using` from `action.yml`/`action.yaml` at its SHA and confirm that it is `node24`, `composite` or `docker`
- [x] 2.8 If the `RELEASE_RUNNER` or `CODEQL_RUNNER` repository variable is set (`gh variable list`; `CODEQL_RUNNER` is used for scheduled CodeQL scans), confirm that the self-hosted runner is Actions Runner v2.327.1 or later before the next release

## 3. Dependabot configuration

- [x] 3.1 Create `.github/dependabot.yml` (`version: 2`) with the three `updates` entries from the design's table:
  - `github-actions`: `directories: ["/", "/.github/actions/*"]`, a weekly schedule, one group `actions` with `patterns: ["*"]`, `commit-message.prefix: ci`, `labels: [area:platform]`, `open-pull-requests-limit: 5`
  - `gomod`: `directory: /backend`, weekly, group `go-minor-patch` with `update-types: [minor, patch]`, `commit-message: {prefix: chore, include: scope}`, `labels: [area:backend]`, `open-pull-requests-limit: 5`
  - `bun`: `directory: /web`, the same settings as gomod but with group `web-minor-patch` and `labels: [area:frontend]`

  Verify that the file validates against the published Dependabot schema: `curl -sSfL https://json.schemastore.org/dependabot-2.0.json -o /tmp/dependabot.schema.json && uvx check-jsonschema --schemafile /tmp/dependabot.schema.json .github/dependabot.yml`.
- [x] 3.2 Verify that `commit-message` produces a subject the lefthook regex accepts: `echo 'ci: bump the actions group with 2 updates' | grep -qE '^(feat|fix|docs|chore|refactor|test|ci|perf|revert)(\(.+\))?: .{1,100}'`, and the same for `chore(deps): bump the go-minor-patch group ...`

## 4. Change classification

- [x] 4.1 In `scripts/ci/changes.sh`, add the rule `.github/dependabot.yml<TAB>-` next to the other `.github/*` ignore rules (after `.github/ISSUE_TEMPLATE/*`). Verify that `echo .github/dependabot.yml | scripts/ci/changes.sh classify` prints `false` for every area
- [x] 4.2 In `scripts/ci/changes-fixtures.txt`:
  - Move `.github/dependabot.yml` from the "Unclassified: fail safe" block to the ignored block, with no areas, in the same format as `.github/CODEOWNERS`.
  - Add `.github/FUNDING.yml	backend,frontend,compose,gomod` as the new unclassified example.

  Verify that `scripts/ci/changes.sh check` passes, and that `scripts/ci/changes.sh check --go` passes as well if Go is available.

## 5. Documentation

- [x] 5.1 Add a "Dependency updates and action pinning" subsection to `docs/TESTING.md` under "What CI runs". Cover:
  - the rule that every non-local `uses:` must be a full SHA plus a `# vX.Y.Z` comment, and the grep that checks it
  - the three Dependabot ecosystems, their grouping, commit types and labels
  - that actions PRs touching `ci.yml` run the full suite
  - the `go work sync` fix-up for a stale `go.work.sum`
  - the Node 24 runtime rule: a new action must not declare `node20`

  Also add `.github/dependabot.yml` to the ignored-paths list there, if that section lists paths. Verify that every command in the new text runs as written.

## 6. Reduce Dependabot CI load

- [x] 6.1 Switch all three ecosystems in `.github/dependabot.yml` to `interval: monthly` with `cooldown.default-days: 7`, and add `go-major`/`web-major` groups (`update-types: [major]`). Verify that the file passes the dependabot-2.0 schema
- [x] 6.2 Add `merge_group: types: [checks_requested]` to `ci.yml`. Verify that actionlint passes and `scripts/ci/changes.sh check` still passes
- [x] 6.3 Document the monthly/cooldown/majors behavior and the merge queue in `docs/TESTING.md` ("What CI runs") and CONTRIBUTING.md step 7. Verify that the tables match `dependabot.yml`
- [ ] 6.4 Open and merge the PR. Verify that `CI Status` passes on the PR
- [ ] 6.5 After it merges, update the `main` ruleset: add a `merge_queue` rule (squash, `CI Status` required) and set `strict_required_status_checks_policy: false`. Verify by queueing a PR: its `merge_group` CI run reports `CI Status` and the PR merges
- [ ] 6.6 Close the superseded individual major PRs (#443–#446) once Dependabot's next run opens the grouped `web-major` PR, or run "Check for updates" to trigger it

## 7. Release job and Dependabot automation

- [x] 7.1 Gate `release-please.yml` on release-relevant commits, and add a `release-please` concurrency group without cancel-in-progress. Verify that the gate script passes the 12 payload cases (hidden-only, Dependabot squash, feat, docs, release merge, `!`, BREAKING footer, visible squash-body line, mixed, empty, non-conventional, corrupt JSON) under `bash -e -o pipefail`, and that actionlint passes
- [x] 7.2 Enable Dependabot security updates (`gh api -X PUT repos/{owner}/{repo}/automated-security-fixes`) and add an `applies-to: security-updates` group per ecosystem. Verify `{"enabled":true}` and that the schema check passes
- [x] 7.3 Add the `docker` and `docker-compose` ecosystems in one `images` multi-ecosystem group (monthly, 7-day cooldown), ignoring golang minor/major and postgres major. Verify that the schema check passes
- [x] 7.4 Pin the Dockerfile (`oven/bun`, `golang`, `distroless/static-debian12`) and compose (`postgres`) images as `tag@sha256:<index digest>`, and confirm that each index covers linux/amd64 and linux/arm64. Verify that compose-smoke passes on the PR
- [x] 7.5 Add `dependabot-auto-merge.yml` (fetch-metadata v3.1.0, `gh pr merge --auto --squash` for minor/patch only). Verify that actionlint passes
- [x] 7.6 Document image pinning, security updates, the images group, auto-merge and the release-please gate in `docs/TESTING.md`
- [ ] 7.7 After merge, confirm on the next Dependabot merge to `main` that release-please logs `run=false`, and on the next Dependabot minor/patch PR that auto-merge is enabled

## 8. Integration checks

- [x] 8.1 Open the PR with a `ci:` title (for example `ci: add Dependabot and SHA-pin actions/checkout (#429)`). Verify that `CI Status` is green (the full suite runs because `ci.yml` changed), and that no job's annotations contain "Node.js 20 is deprecated" (`gh run view <id> --log | grep -c 'Node.js 20 is deprecated'` prints `0`)
- [ ] 8.2 After merging, open Insights → Dependency graph → Dependabot. Confirm that the three ecosystems are listed with no config errors, and trigger "Check for updates" on each
- [ ] 8.3 On the first Dependabot actions PR, confirm:
  - one grouped PR with a `ci:` subject and the `area:platform` label
  - SHAs and `# vX.Y.Z` comments updated together, including under `.github/actions/*`
  - `CI Status` passes under the Dependabot token, with cache-save steps not failing the job

  If a cache-save step does fail, gate it on `github.actor != 'dependabot[bot]'` in a follow-up commit.
- [ ] 8.4 After merge, confirm that the next `release-please` run succeeds on v5. On the next release tag, or a `vX.Y.Z-rc.1` prerelease, confirm that `release.yml` builds, scans, pushes and attaches the SBOM without Node 20 warnings
- [ ] 8.5 Check the acceptance criteria in #429 and close the issue, noting the v7 pin and the added gomod/bun ecosystems
