# Contributing to WiseLabz

Thank you for contributing. This guide covers everything you need to set up the project,
follow our conventions, and get your changes merged.

---

## Prerequisites

- **Go 1.23+** — backend and CLI tooling
- **[Bun](https://bun.sh) 1.3+** — frontend (React + Vite)
- **Docker and Docker Compose** — running the full stack and testing connectors
- **[lefthook](https://github.com/evilmartians/lefthook)** — commit hooks (install with
  `go install github.com/evilmartians/lefthook@latest` or your package manager)

## Setting up the project

```bash
# Clone and enter the repo
git clone https://github.com/WiseLabz/WiseLabz.git
cd WiseLabz

# Install all dependencies and generate code
make setup

# Run the full stack in development mode
make dev
```

`make dev` starts the Go backend on port 8080 and the Vite dev server on port 5173 with
hot reload for both. Use `make help` to see all available targets.

## Branch naming

| Prefix      | Use for                                          |
|-------------|--------------------------------------------------|
| `feat/`     | New features and enhancements                    |
| `fix/`      | Bug fixes                                        |
| `docs/`     | Documentation-only changes                       |
| `chore/`    | CI, tooling, dependency bumps, housekeeping      |
| `refactor/` | Code changes that don't add features or fix bugs |

Examples: `feat/diff-engine-ai-update`, `fix/portainer-tls-verify`, `docs/connector-guide`.

## Commit messages

We use [Conventional Commits](https://www.conventionalcommits.org/). Every commit message
must follow this format:

```
<type>(<scope>): <short description>

<optional body>
```

### Types

- `feat` — a new feature
- `fix` — a bug fix
- `docs` — documentation
- `chore` — maintenance
- `refactor` — code restructuring without behavior change
- `test` — adding or updating tests
- `ci` — CI pipeline changes

### WiseLabz-specific examples

```
feat(connector): add Proxmox VM snapshot support
fix(diff): handle empty service snapshot on first sync
docs(config): document WISELABZ_ env var override order
chore(deps): bump golang-migrate to v4.18
refactor(api): extract service handler from router setup
test(connector): add integration tests for Portainer fetch
ci: add go vet to pre-commit hook
```

## Commit hooks

WiseLabz uses `lefthook` to run automated checks on every commit:

| Hook         | What runs                                |
|--------------|------------------------------------------|
| `pre-commit` | Lint (Go + frontend) and fast unit tests |
| `commit-msg` | Validates Conventional Commits format    |

Install the hooks after cloning:

```bash
lefthook install
```

If the pre-commit hook blocks you on a lint issue, run `make lint` to see the full
output and `make fmt` to auto-fix formatting. The commit-msg hook will reject
non-conforming messages — the error tells you exactly what's wrong.

## Releasing

Releases are automated with [release-please](https://github.com/googleapis/release-please):

1. Conventional Commits merged to `main` are picked up automatically.
2. release-please keeps an open "release PR" that bumps the version (root
   `CHANGELOG.md` and `web/package.json`) based on those commits.
3. Merging that PR pushes a `vX.Y.Z` tag.
4. The tag push triggers `release.yml`, which builds and publishes the
   multi-arch image to GHCR, runs a Trivy scan, generates release notes with
   git-cliff, and publishes the GitHub Release with an attached SBOM.

No manual tagging or changelog editing is needed — write a correct
Conventional Commit and the rest is automatic.

## Testing and coverage baseline

The frontend rewrite (#52) landed, so the deferral in the old policy (#12) no
longer applies — `web/` has an automated Vitest suite and a `test-frontend` CI
job, same as the backend's `test-backend` job.

Both suites run with coverage in CI and enforce a floor so coverage only goes
up over time:

- **Backend** (`go test -coverpkg=./... -coverprofile=coverage.out ./...`):
  floor is 60% total statement coverage (CI fails under it).
- **Frontend** (`bun run test:coverage`, v8 provider): floor is set in
  `web/vitest.config.ts` (`test.coverage.thresholds`) — currently 40%
  statements/lines, 25% functions/branches.

Raise these floors as coverage improves; never lower them. When adding a
feature with non-trivial branching (auth, destructive actions, WebSocket
handling), add regression tests for it rather than relying on the floor alone.

## Pull request process

1. **Branch from `main`** — create a branch using the naming convention above.
2. **Keep commits small and focused** — one concern per commit.
3. **Run hooks locally** — `lefthook` runs `pre-commit` before every commit. Ensure CI
   will pass by running `make test` before opening the PR.
4. **Open the PR** — use the pull request template. Link any related issues.
5. **One approval required** — a maintainer must approve before merge.
6. **CI must pass** — all checks (lint, test, build) must be green.
7. **Rebase, don't merge** — update your branch with `git rebase main`. Maintainers
   squash-and-rebase when merging. No merge commits land on `main`.

## Writing a connector

Connectors are how WiseLabz talks to services. If you want to add support for a new
platform, read the [Connector Guide](docs/connectors/CONNECTOR_GUIDE.md). It walks through
the `Connector` interface, the `ServiceSnapshot` data structure, registration, testing,
and mocking with a working example.

## Getting help

Open a [GitHub Discussion](https://github.com/WiseLabz/WiseLabz/discussions) or comment on
your open PR. We're a small team and we'll get back to you as soon as we can.
