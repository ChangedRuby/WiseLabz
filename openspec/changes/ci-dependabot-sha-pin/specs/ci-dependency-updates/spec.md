# Spec Delta

## Purpose

Keeps the third-party code CI runs, and the project's Go and web dependencies, pinned to reviewed versions and updated automatically. Update PRs are grouped and labelled, use commit messages that pass the repository's checks, and show up in the changelog only where they should.

## ADDED Requirements

### Requirement: Every action reference is pinned to a commit SHA
Every `uses:` reference to an action outside this repository, in workflows (`.github/workflows/**`) and composite actions (`.github/actions/**`), SHALL name a full 40-character commit SHA followed by a `# vX.Y.Z` comment giving the release that SHA belongs to. References to local actions (`./...`) are exempt. No reference SHALL use a mutable tag or branch.

#### Scenario: No mutable tags remain
- **WHEN** `grep -rnE 'uses: [^.].*@v[0-9]' .github` runs on the repository
- **THEN** it prints nothing

#### Scenario: Checkout is pinned to v7
- **WHEN** any workflow checks out the repository
- **THEN** it references `actions/checkout@<40-char SHA> # v7.0.1`, or a later release that Dependabot has moved it to

#### Scenario: Pins name commits, not tag objects
- **WHEN** each pinned SHA is looked up as a commit (`gh api repos/<owner>/<repo>/commits/<sha>`)
- **THEN** every lookup succeeds, and none of the SHAs is an annotated-tag object

#### Scenario: Composite actions are included
- **WHEN** `.github/actions/bun-setup/action.yml` references `oven-sh/setup-bun`
- **THEN** that reference is SHA-pinned with a version comment

### Requirement: No action runs on a deprecated Node.js runtime
Every JavaScript action referenced by a workflow or composite action SHALL declare a Node.js runtime that GitHub Actions currently supports without deprecation (`node24` at the time of writing). Composite actions are checked through the actions they call. CI runs SHALL NOT emit the "Node.js 20 is deprecated … being forced to run on Node.js 24" annotation.

#### Scenario: Pinned actions declare node24
- **WHEN** the `runs.using` field of each pinned action's `action.yml`/`action.yaml` is read at its pinned SHA
- **THEN** every value is `node24`, `composite` or `docker`

#### Scenario: CI run is free of Node 20 warnings
- **WHEN** a pull request that selects every area runs CI
- **THEN** no job's annotations include "Node.js 20 is deprecated"

### Requirement: GitHub Actions are updated monthly in one grouped PR
The repository SHALL configure automated updates for the `github-actions` ecosystem. They cover both the workflows and every composite action under `.github/actions/*`, run monthly, and skip releases younger than 7 days. All pending action updates SHALL arrive in a single pull request. The update commits SHALL use the `ci` Conventional Commits type, so they pass the `commit-msg` hook and stay hidden in the release changelog. Update PRs SHALL be labelled `area:platform`.

#### Scenario: Several actions are outdated
- **WHEN** a monthly run finds new releases of `actions/cache` and `docker/build-push-action`
- **THEN** one pull request updates both, rewriting each SHA and its `# vX.Y.Z` comment

#### Scenario: A composite action's dependency is outdated
- **WHEN** a new `actions/setup-go` release exists and it is referenced only in `.github/actions/go-cache-restore/action.yml`
- **THEN** the grouped actions PR updates that file

#### Scenario: Commit message follows conventions
- **WHEN** an actions update PR is created
- **THEN** its commit subject matches `^ci(\(.+\))?: ` and the change does not appear in release notes

#### Scenario: Everything is current
- **WHEN** a monthly run finds no newer action releases
- **THEN** no pull request is opened

### Requirement: Go modules and web packages are updated monthly
The repository SHALL configure monthly automated updates for Go modules in `/backend` and Bun-managed packages in `/web`, skipping releases younger than 7 days. Security updates SHALL NOT wait for the schedule. For each ecosystem, minor and patch updates SHALL be grouped into one pull request, and major updates into a second one, so an ecosystem opens at most two update PRs per cycle. Commits SHALL use the `chore` type with a `deps` scope (`chore(deps): ...`). Go update PRs SHALL be labelled `area:backend` and web update PRs `area:frontend`.

#### Scenario: Several minor Go updates
- **WHEN** a monthly run finds patch releases of two Go modules
- **THEN** one pull request titled with `chore(deps)` updates `backend/go.mod` and `backend/go.sum` for both

#### Scenario: Several major web updates
- **WHEN** a monthly run finds new major versions of three web dependencies and minor updates of forty more
- **THEN** exactly two web PRs are opened: one with the three majors and one with the minor and patch updates

#### Scenario: Fresh release is deferred
- **WHEN** a dependency published a release three days before the run
- **THEN** that release is not proposed until it is at least 7 days old

#### Scenario: Update PR runs the matching CI area
- **WHEN** a web update PR changes only `web/package.json` and `web/bun.lock`
- **THEN** the frontend jobs and compose-smoke run, the backend jobs are skipped, and `CI Status` reports the result

### Requirement: The update configuration has no CI cost
A change that touches only the Dependabot configuration file SHALL be classified as having no CI value. It SHALL skip every lint, test, build and compose-smoke job, and `CI Status` SHALL still report success.

#### Scenario: Config-only pull request
- **WHEN** a pull request changes only `.github/dependabot.yml`
- **THEN** only change detection and `CI Status` run, and `CI Status` succeeds

#### Scenario: Unclassified paths still fail safe
- **WHEN** a pull request changes only a path no rule classifies (for example `.github/FUNDING.yml`)
- **THEN** every area is selected, as before

### Requirement: Merging several PRs does not re-run CI on every open PR
PRs to `main` SHALL merge through a merge queue rather than requiring each branch to be up to date with `main`. The CI workflow SHALL run on merge-queue entries and report `CI Status` for them, classifying the changes in the entry. Merging N queued PRs SHALL NOT require rebasing and re-running CI on the remaining open PRs.

#### Scenario: Three Dependabot PRs merged together
- **WHEN** three approved Dependabot PRs are added to the merge queue
- **THEN** each queue entry runs CI once, and the other open PRs are not rebased or re-run

#### Scenario: Queue entry is gated by CI
- **WHEN** a queue entry's `CI Status` fails
- **THEN** that PR is removed from the queue and not merged

### Requirement: Security updates arrive immediately and grouped
Dependabot security updates SHALL be enabled for the repository. They SHALL NOT wait for the monthly schedule or the cooldown, and SHALL be grouped into one pull request per ecosystem.

#### Scenario: Advisory against a Go module
- **WHEN** an advisory is published for a Go module in `backend/go.mod`
- **THEN** a grouped `chore(deps)` security PR opens without waiting for the monthly run

### Requirement: Container images are pinned and updated
Base images in `Dockerfile` and service images in `docker-compose.yml` SHALL be pinned as `tag@sha256:<multi-arch index digest>`. Dependabot SHALL propose image updates monthly, in one PR covering both files. It SHALL NOT propose a new Go minor or major for the `golang` image, or a new Postgres major.

#### Scenario: New digest for a pinned tag
- **WHEN** `oven/bun:1-alpine` is republished with a new digest
- **THEN** the monthly `images` PR updates the digest and keeps the tag

#### Scenario: New Postgres major
- **WHEN** `postgres:17-alpine` is released
- **THEN** no update is proposed for the `postgres:16-alpine` pin

### Requirement: Minor and patch Dependabot PRs merge without further action
A Dependabot PR whose largest change is a minor or patch update SHALL have auto-merge enabled when it is opened. It SHALL merge through the merge queue once it has the required approval and `CI Status` passes. Major updates SHALL NOT be auto-merged.

#### Scenario: Grouped minor/patch PR approved
- **WHEN** a code owner approves the `web-minor-patch` PR and its CI passes
- **THEN** the PR is queued and merged without anyone clicking merge

#### Scenario: Majors stay manual
- **WHEN** the `web-major` PR is opened
- **THEN** auto-merge is not enabled

### Requirement: release-please skips pushes that cannot change a release
The release workflow SHALL skip the release-please action when every pushed commit is of a type the changelog hides (`ci`, `chore`, `test`), has no breaking marker, contains no visible-type commit line, and is not the release PR's merge. When the commit list is missing or can't be read, it SHALL run. Release runs SHALL be serialized so that bursts of pushes collapse into one pending run, and a running release SHALL NOT be cancelled.

#### Scenario: Dependabot merge
- **WHEN** a push contains only `ci: bump the actions group ...`
- **THEN** the action is skipped and the job summary says so

#### Scenario: Release PR merge
- **WHEN** a push contains `chore(main): release 1.0.0`
- **THEN** release-please runs and publishes the release

#### Scenario: Burst of merges
- **WHEN** four PRs land on `main` while a release-please run is in progress
- **THEN** only one more run follows the current one
