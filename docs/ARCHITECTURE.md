# WiseLabz — Architecture & Technical Decisions

This document records every significant technical decision made for the WiseLabz project.
It exists so contributors (and future maintainers) can understand **why** the project is structured
the way it is, not just **what** it uses.

---

## Project overview

WiseLabz is a self-hosted, open-source service that automatically generates and maintains
documentation for home labs. Users register their services (Proxmox, Docker/Portainer,
pfSense/OPNsense, self-hosted apps, etc.), connect their APIs, and WiseLabz fetches live data
to generate structured documentation from configurable templates. A diff engine detects changes
since the last sync and either triggers an AI-assisted section update or emits an alert. A web
dashboard displays live service data alongside the generated docs.

**Deployment target:** Docker Compose (self-hosted).  
**Repository:** `github.com/wiselabz/wiselabz`  
**Container registry:** `ghcr.io/wiselabz/wiselabz`

---

## Repository structure

### Monorepo with Go workspaces

All code lives in a single repository. Go workspaces (`go.work`) are configured from day one,
separating the core module from connector packages. This keeps the door open for extracting
connectors into importable standalone libraries later without breaking the internal API.

```
WiseLabz/
├── backend/
│   ├── cmd/server/          # entrypoint — main.go
│   └── internal/
│       ├── connector/       # Connector interface + implementations
│       │   ├── proxmox/
│       │   ├── docker/
│       │   ├── pfsense/
│       │   └── custom/      # example connector for contributors
│       ├── sync/            # scheduler, diff engine
│       ├── doc/             # template engine, renderer
│       ├── api/             # REST handlers + WebSocket
│       ├── auth/            # JWT + OIDC provider
│       └── store/           # repositories, migrations
├── web/
│   └── src/
│       ├── components/      # reusable UI
│       ├── pages/           # dashboard, settings, docs viewer
│       ├── hooks/           # custom React hooks (incl. useWebSocket)
│       ├── store/           # Zustand global state
│       └── api/             # generated REST client
├── deploy/
│   ├── docker-compose.yml
│   └── config.example.yaml
├── docs/                    # ADRs, connector guides, this file
├── .github/
│   └── workflows/           # CI/CD pipelines
├── Makefile
├── CONTRIBUTING.md
└── go.work
```

---

## Technology stack

### Backend

| Concern          | Decision                              | Rationale                                                                                                       |
|------------------|---------------------------------------|-----------------------------------------------------------------------------------------------------------------|
| Language         | Go                                    | performance, single binary, strong stdlib                                                                       |
| HTTP router      | `chi`                                 | idiomatic middleware composition, no magic, lightweight                                                         |
| WebSocket        | `gorilla/websocket`                   | de facto standard in Go, well-maintained                                                                        |
| Auth             | Local JWT + OIDC via `coreos/go-oidc` | single session logic covers both local and external providers                                                   |
| Database queries | `sqlc`                                | type-safe generated code, no ORM, fully auditable SQL                                                           |
| Migrations       | `golang-migrate`                      | simple, supports SQLite and PostgreSQL, low contributor friction                                                |
| Config           | `viper`                               | YAML file for traditional installs; all values overridable via WISELABZ_ env vars, friendly to PaaS deployments |
| Logging          | `slog` (stdlib)                       | no external dependency, JSON output in production                                                               |

### Frontend

| Concern          | Decision                                                  | Rationale                                                                                                                                                            |
|------------------|-----------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Framework        | React + Vite                                              | mature ecosystem, fast HMR in dev                                                                                                                                    |
| State management | Zustand                                                   | minimal boilerplate, works cleanly alongside React Query                                                                                                             |
| Server state     | React Query                                               | caching, background refetch, loading/error states for REST data                                                                                                      |
| Styling          | Tailwind CSS                                              | productive, dark mode built-in, low barrier for contributors                                                                                                         |
| API client       | generated via `orval` (OpenAPI → typed React Query hooks) | contract-first from `docs/openapi.yaml`; frontend/backend types stay in sync, and MSW mocks generated from the same spec let the frontend build ahead of the backend. Output (`web/src/api/generated`, `web/src/api/model`) is committed; CI regenerates and fails the build if it drifts from the spec, so run `bun run gen:api` after editing `docs/openapi.yaml` and commit the result |

### Infrastructure

| Concern            | Decision                                        | Rationale                                                      |
|--------------------|-------------------------------------------------|----------------------------------------------------------------|
| Deployment         | Docker Compose                                  | `docker compose up` is the entire install story                |
| Build              | Vite SPA embedded into Go binary via `go:embed` | single binary, single Docker service, no Nginx needed          |
| Container registry | GHCR (`ghcr.io`)                                | native GitHub integration, free for open-source                |
| CI                 | GitHub Actions                                  | lint + test + build on every PR; multi-arch Docker push on tag |

---

## Build pipeline

The production build uses a multi-stage Dockerfile:

1. **Stage 1 (`web`)** — Bun Alpine image runs `bun install` and `bun run build`, producing a
   static SPA in `/web/dist`.
2. **Stage 2 (`builder`)** — Go image copies the SPA output into `internal/web/dist/` (where
   `go:embed` expects it) and compiles a static binary with `CGO_ENABLED=0`.
3. **Stage 3 (`runtime`)** — Minimal scratch or distroless image receives only the compiled
   binary. No Node.js, no Go toolchain reaches the runtime image.

In development, Vite's dev server runs separately with a proxy forwarding `/api` and `/ws`
requests to the Go backend. The `go:embed` path only activates in the production build.

---

## Development workflow

### Branching

| Prefix      | Use                                   |
|-------------|---------------------------------------|
| `feat/`     | new features                          |
| `fix/`      | bug fixes                             |
| `docs/`     | documentation only                    |
| `chore/`    | tooling, dependencies, config         |
| `refactor/` | code changes with no behaviour change |

Pull requests are required for all work, including between the two core maintainers.
This creates a reviewable decision trail that external contributors can follow.

### Commit conventions

All commits follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
feat(connector): add Proxmox VM snapshot support
fix(auth): handle expired refresh token on WS reconnect
docs(connector): add guide for writing a custom connector
chore(ci): pin golangci-lint to v1.57
```

### Commit hooks (`lefthook`)

Two hooks are enforced locally via `lefthook`:

- **`pre-commit`** — runs `gofmt`, `golangci-lint`, and fast unit tests
- **`commit-msg`** — validates the Conventional Commits format

Both hooks also run in CI, so a passing local environment guarantees a passing pipeline.

### Code quality

- Go: `gofmt` + `golangci-lint` with an explicit `.golangci.yml` config checked into the repo
- TypeScript: ESLint + Prettier
- Both block merges if they fail — no exceptions

---

## Authentication design

WiseLabz supports two authentication modes, sharing the same session logic:

**Local mode** — username/password stored in the database (bcrypt). On login, issues a
short-lived JWT access token and a longer-lived refresh token. No external dependency required.

**OIDC mode** — the user sets `auth.oidc.issuer_url` in `config.yaml` pointing to their
provider (Authentik, Authelia, Keycloak, etc.). The login flow delegates to the provider;
WiseLabz receives and validates the ID token via `coreos/go-oidc`. The session mechanism
(JWT issuance, refresh) is identical to local mode.

Switching between modes requires only a config change and restart — no code path diverges.

### OIDC provider configuration (decided 2026-06-25: file-defined, app toggles only)

OIDC providers are defined **entirely in config file / env** — `issuer_url`,
`client_id`, and `client_secret` never transit the HTTP API, never live in the app
database, and are never editable through the UI. This is deliberate: an API that could
rewrite `issuer_url` would be an auth-bypass vector (an attacker repointing the issuer
to an IdP they control logs in as admin), and storing the secret in the DB only adds an
exfiltration surface without removing the env secret you need to encrypt it.

The in-app Settings → Auth page is therefore **read-only for providers plus a few safe
toggles**:
- View detected providers (issuer, clientId, `secretConfigured: yes/no`) — never the secret.
- Enable/disable a provider (`PUT /api/auth/providers/{id}/enabled`) — persisted as a DB flag; the definition stays in config.
- Toggle local login and set access/refresh token TTLs (`PUT /api/auth/config`).

This supersedes the frontend plan's original §8.6 ("add/edit OIDC provider via UI"),
which has been amended accordingly. Endpoint contract: `docs/openapi.yaml`
(`/auth/config`, `/auth/providers/{providerId}/enabled`).

### Permissions & step-up for mutating actions (decided 2026-06-27)

Two layers gate state-changing actions, on top of the session model above:

**Role-based access.** Two roles: `viewer` (read docs + live state) and `operator`
(everything viewer can do, plus mutating actions). Role is a property of the user
record, enforced server-side on every mutating endpoint — the UI hiding a control is
a convenience, never the security boundary.

**Per-action step-up.** Destructive actions require re-authentication (password, or
TOTP/2FA when configured) immediately before they execute, even for an operator. The
server issues a short-lived elevation token scoped to the single action.

Step-up is **scoped to genuinely destructive ops**, not every mutation:

| Action                              | Mutating                        | Step-up              |
|-------------------------------------|---------------------------------|----------------------|
| Trigger sync (global / per-service) | yes (WiseLabz-side, reversible) | no                   |
| Connector enable/disable            | yes (reversible toggle)         | no                   |
| Connector **add**                   | yes (creates state)             | no                   |
| Connector **remove**                | yes (cascades, destructive)     | **yes (default on)** |

Gating sync or a toggle behind a password is friction with no payoff and trains
users to disable step-up wholesale. Step-up therefore defaults **on for connector
removal** (the one destructive op in v1) and off elsewhere. A Settings toggle lets a
solo operator disable step-up entirely — honest default, escape hatch for the
single-admin homelab. As lab-mutating ops (lifecycle, config push) land post-v1,
each enters this table as `step-up: yes`.

### Destructive-action pattern: confirm + blast radius (decided 2026-06-27)

Destructive actions never fire on a bare click. The server computes and returns the
**blast radius** — the concrete dependents the action destroys (e.g. removing a
connector drops N tracked services and M doc sections) — and the UI states it
exactly before the user commits, then requires type-to-confirm (type the resource
name) for the highest-risk ops. This is the "machine-honest" principle applied to
mutation: show what changes and what depends on it, never imply the action is
smaller than it is.

---

## Manager actions (v1 scope)

WiseLabz is a homelab manager, not only a doc manager (see `PRODUCT.md`). The v1
manager layer is deliberately narrow and all of it is REST (per the API design
below — no mutations over WS):

- **Sync** — trigger a sync/refresh globally or per-service.
- **Connector enable/disable** — a persisted flag; the connector definition is
  untouched.
- **Connector add/remove (full CRUD via UI)** — see below.

**Out of scope for v1:** lab-mutating operations (service start/stop/restart,
config push back to the service). These are desired but deferred; each will require
its own row in the step-up table and a blast-radius definition before shipping.

### Connector management via UI (decided 2026-06-27)

Service connectors (Proxmox, Docker, pfSense, custom) are full CRUD through the UI —
distinct from OIDC providers, which remain config-file-only for the auth-bypass
reasons above. Adding a connector is **a dedicated flow, not a row action**: it
carries a per-connector-type config schema (each type exposes different fields),
credential entry, and a connection test (`Connector.Validate()`) before the
connector is persisted. Removal is the destructive op that drives the blast-radius
pattern above (it cascades to that connector's snapshots and generated doc
sections). Connectors also have an optional `owner` field used by documentation
quality checks. Credentials are write-only over the API — never returned in reads.

---

## Connector interface

Every service integration implements the same interface. The core knows nothing about
Proxmox or Docker specifically — only about `Connector`.

```go
type Connector interface {
    Name()     string
    Fetch(ctx context.Context) (ServiceSnapshot, error)
    Validate() error
}
```

Each connector lives in its own package under `internal/connector/`. Contributors adding a
new integration only need to implement this interface and register the connector — no changes
to the core sync or doc engine are required.

A `ServiceSnapshot` can also carry `Dependencies` — the hosts, networks, storage, and
upstream services a connector's service relies on. When present, generated documentation
renders a "Dependencies" section, so docs can describe a service's operational context, not
just its isolated state.

---

## API design

- **REST** — all data reads and state mutations. Schema exported as OpenAPI, used to generate the TypeScript REST client in `web/src/api/`.
The API itself is implemented entirely in Go.
- **WebSocket** — live updates only (service status changes, diff engine alerts, sync progress).
  No mutations over WS. A single `/ws` endpoint streams typed JSON events to the dashboard.

---

## Database

- **SQLite** — default for development and small single-user installs. Zero setup, embedded.
- **PostgreSQL** — recommended for multi-user or higher-reliability deployments. Configured via
  `config.yaml`.

`sqlc` generates type-safe Go code from plain `.sql` query files. Migrations are managed by `golang-migrate` and run automatically on startup.
Both SQLite and PostgreSQL dialects are supported.  
Database connection is configured via environment variables (recommended for PaaS and container platforms like Dokploy, Coolify, and Portainer)
or via `config.yaml` for traditional installs. Environment variables always take precedence over the config file. The `WISELABZ_` prefix is used
for all env overrides (e.g. `WISELABZ_DB_DSN`, `WISELABZ_DB_DRIVER`).  
`golang-migrate` and run automatically on startup. Both SQLite and PostgreSQL dialects are
supported.

Operators can roll back the most recently applied migration with the `migrate` CLI
(`go run ./cmd/migrate down`), which delegates to golang-migrate's `Steps(-1)` and supports
both SQLite and PostgreSQL.

---

## Template history and read-only preview (decided 2026-09-05)

Templates retain immutable save and restore snapshots, mirroring document history.
Restoring an older template creates a new current revision instead of rewriting history.

`POST /templates/{templateId}/preview` is read-only: it renders each matching connector
sequentially, reports whether its document would change, and optionally returns one detailed
render without creating or updating `docs` or `doc_versions`.

---

## Data retention (decided 2026-09-05)

A background job (`internal/retention`, registered with the shared cron
scheduler in `cmd/server/main.go` alongside the quality stale-sweep and sync
poll jobs — see `internal/scheduler`) bounds the growth of historical data:
`service_snapshots`, `doc_versions`, `alerts`, and `sync_runs`. Each category
has an independent, configurable retention
window (`retention.*_days` in `config.yaml` / `WISELABZ_RETENTION_*_DAYS` env
vars); a value of `0` disables cleanup for that category entirely.

The job always protects records an active workflow depends on, regardless of
age: a connector's latest snapshot (what `GetLatestSnapshot` serves as
"current data"), a doc's current revision (`docs.current_version`), and any
`pending` or `snoozed` (non-terminal) alert. Only `resolved`/`dismissed`
alerts are purged, and `sync_runs` has no such guard — a connector's
last-sync status lives on the `connectors` row itself, not derived from
`sync_runs` history.

The job runs on a configurable cron schedule (`retention.cron_expr`, default
`0 0 * * *` — daily at midnight) and is idempotent: each pass deletes only
rows newly past that category's cutoff, so running it repeatedly (or after
downtime) is safe and produces no duplicate effects.

---

## Documentation quality findings (decided 2026-09-05)

WiseLabz records actionable findings when documentation is stale, effectively
empty, backed by three consecutive failed syncs, or attached to a connector with
no owner. Only one open finding exists per connector and check type: repeated
detections atomically update its last-seen time and count, while fixing the
condition auto-resolves it. A later recurrence creates a new finding.

All five checks run after each non-skipped sync. A daily sweep checks staleness
for every connector as well, including connectors that have stopped syncing.
Findings appear on their own authenticated API and navigation page; operators can
resolve them, remediation links point to the affected document or connector, and
WebSocket events refresh the badge. Retention pruning remains a follow-up.

**Credential rotation (#239 PR1)** is the fifth check: a connector's
secret-typed config fields have a `secret_rotated_at` timestamp, bumped only
when their decrypted value actually changes (not on a rename or a resubmitted
identical secret). The check flags a connector whose secret is older than
`rotation.max_age_days` (`config.yaml` / `WISELABZ_ROTATION_MAX_AGE_DAYS`,
default 90), or past an optional per-connector `userExpiresAt` if that's
sooner — `warning` inside `rotation.warn_days` (default 14) of the due date,
`critical` past it. A connector can override the global max age
(`rotationMaxAgeDays`). Connectors whose implementation satisfies
`CredentialRefresher` (self-refreshing credentials, e.g. OAuth2) are always
skipped. Unlike the other four checks, this one also **notifies**: opening a
finding, or its severity escalating, dispatches through the same
`notifications.Dispatcher` alerts use (in-app, SMTP, webhook, Discord, Slack, ntfy, Telegram; see
[NOTIFICATIONS.md](NOTIFICATIONS.md)) via
`Dispatcher.NotifyFindingCreated` — deduplicated by `notified_severity` so a
repeat detection at the same severity doesn't re-notify, while a
resolve-then-reopen does. This notify hook is generic across every check
type, not just `credential_rotation`.

---

## AI module

The AI doc generation feature is **opt-in and provider-agnostic**. When enabled, the diff
engine passes changed service data and the current doc section to the configured LLM API.
The user can choose their provider (Anthropic, OpenAI, or a local model via Ollama) in
`config.yaml`. When disabled, the diff engine emits alerts instead of auto-updating sections.

---

## Frontend shell & theme (decided 2026-06)

**Shell — bottom dock, single variant.** The authenticated frame is `AppShell`
(`web/src/components/shell/AppShell.tsx`): a slim top bar (brand, ⌘K search,
global-sync, alerts, account), full-width routed content, and a floating bottom
nav dock (`ShellDock` + `Dock`). A sidebar variant and a live shell-switch toggle
were built for comparison and **removed** once the dock was chosen — there is no
runtime shell switching and no `shell` UI state.

**Theme engine — code default, user-overridable behind Settings.** Palettes and
fonts are generated from OKLCH builder options in `web/src/theme.ts`. The shipped
default is `ACTIVE = { palette: 'blueprint', font: 'geist' }` (a first-paint
fallback is inlined in `web/src/index.css`). Settings → Theme exposes presets,
font sets, and advanced OKLCH sliders for exploration, but the default is the
brand, not a neutral starting point. Motion is a user-controlled setting
(Full / Reduced / Off), seeded by `prefers-reduced-motion` but never silently
forced.

## Changes / diff contract (decided 2026-06)

The `Diff` schema in `docs/openapi.yaml` carries two shapes keyed by `format`:
`infra` uses structured key-path `hunks` (config/state drift); `doc` carries the
**full before/after revisions** (`baseText` / `headText`, with optional
`baseLabel` / `headLabel` / `language`). The frontend computes the line + word
diff (`web/src/lib/{linediff,worddiff,docdiffmodel}.ts`) and renders it
JetBrains-style — whole document with foldable context, intra-line word
highlighting, unified or side-by-side. Sending full revisions (rather than
pre-windowed hunks) keeps the contract and the diff engine simple; large-doc
perf is handled client-side by context-folding.

**Role model amendment.** The two roles `viewer` / `operator` (above) are now the
contract `Role` enum. For v1's single-admin homelab, the former admin-only
surfaces (user management, auth/system config) are gated on `operator`; a distinct
`admin` role can be reintroduced if multi-tenant separation is needed later.

**Runbook bindings.** A `runbooks` row attaches actionable guidance to a change
type or an alert severity (`target_type` / `target_value`), not to a specific
change or alert instance. Each runbook carries free-text `title`/`body` plus
optional `snapshot_id` and `doc_id` pointers, letting operators jump from "this
kind of change" or "this severity" to a known-good historical snapshot or a
relevant doc. Like the diff and retention models above, the binding is purely
a read-only reference: `snapshot_id`/`doc_id` are inert lookups with no
restore or apply path, and the CRUD surface (`GET`/`POST`/`PUT`/`DELETE
/api/runbooks`, viewer-read/operator-write) only ever manages the guidance
record itself, never the target it points at.

## ADR index

Architectural Decision Records are stored in `docs/adr/`. Each significant decision that
required evaluation of alternatives gets its own numbered ADR file (e.g. `0001-monorepo.md`).
This file records the _outcome_ of each decision; the ADRs record the _reasoning_.

- [`0001-lab-mutating-operation-boundaries.md`](adr/0001-lab-mutating-operation-boundaries.md) —
  permission, step-up, audit, dry-run, and rollback model for the first lab-mutating
  operation (`service.restart`), ahead of implementation.

---

_Last updated: 2026-06-27_
