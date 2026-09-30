# 0001 — Lab-mutating operation boundaries

Status: accepted (implemented; start/stop by ADR 0002, config push by ADR 0003)

## Context

`PRODUCT.md` locks v1's manager scope to trigger-sync, enable/disable a
connector, and add/remove connectors — explicitly excluding lab-mutating
operations (service start/stop/restart, config push), gated on "the
permission/confirmation model" (`PRODUCT.md:51-54,70-72`). `docs/MISSING.md`
carried the same deferral until it shipped (#236). This ADR is that model: it defines which
operation may go first, and how it reuses the authorization, step-up,
audit, and rollback machinery already in the codebase, so that when the
operation is actually implemented it has no open design questions left.

**This ADR ships no code.** It defines the model; implementation is a
separate, future change.

## Decision

### Eligible first operation: `service.restart` only

Restart, and nothing else, is in scope for the first lab-mutating
operation. Config-push and start/stop are explicitly deferred to a
follow-up ADR.

Reasoning:
- Restart is momentary and self-healing: the service either comes back or
  it doesn't, with no partial-state risk to reason about.
- Every connector category that exists today (`virtualization`,
  `containers_paas`, `networking`) already has an unambiguous "restart"
  concept operators already understand.
- Config-push is qualitatively riskier — it writes arbitrary content into
  a live system with no generic rollback (nothing today snapshots-and-diffs
  a target config before writing it). Start/stop as two separate verbs
  doubles the surface for what restart already covers in one action.
- Naming one narrow operation keeps this ADR's job cheap to validate:
  prove the model works, then extend it.

### Authorization

Reuse `auth.RequireRole("operator")` / `roleSatisfies`
(`backend/internal/auth/middleware.go:56-67`) as-is. No new role. Per the
existing role-model note (`docs/ARCHITECTURE.md`, "Role model amendment"),
a distinct admin role stays deferred to a multi-tenant future and this
ADR does not reopen that.

### Confirmation / step-up

Reuse `auth.RequireElevation(jwtSvc, action)`
(`backend/internal/auth/middleware.go:83-111`) with a new action string
`"connector.restart"`, wired exactly like the existing elevation-gated
routes (`connector.delete`, `template.delete`, `user.delete`,
`user.resetPassword` — `backend/internal/api/router.go:142,179,279,284`).
Reuse the existing 60s `elevationTTL` (`backend/internal/auth/jwt.go:52`)
and the existing `POST /api/auth/elevate` flow verbatim. No new elevation
mechanism.

### Audit

Reuse the existing audit writer and log described in `docs/AUDIT.md`. When
implemented, add one row to that document's action table:

| Action | Trigger | targetType / targetId |
|---|---|---|
| `connector.restart` | `POST /api/connectors/{id}/restart` | connector / id |

The existing "successes only, not attempts" scope cut (`docs/AUDIT.md`,
"What's not recorded") applies unchanged — a failed restart attempt is not
audited, only a completed one.

### Dry-run

No dry-run pattern exists anywhere in this codebase today — this is new
ground. Minimal contract for when restart is implemented:

- `POST /api/connectors/{id}/restart?dryRun=true` (or a dedicated
  `/restart/preview` path) returns what *would* happen — target service
  name, an estimated downtime window, and any other services with a
  snapshot dependency on this one — **without** calling any mutating
  connector method (none exists yet; `connector.Connector` today is
  strictly `Fetch`/`Validate`, per `backend/internal/connector/connector.go`).
- The "dependent services" list is where issue #31's `ServiceDependency`
  model (`Kind: host | network | storage | upstream_service` on
  `ServiceSnapshot`) pays off directly: a preview can report "N dependent
  services" once that data is populated.
- Dry-run is **mandatory** before the real action is exposed in the UI —
  mirroring how destructive deletes already require a confirm dialog.

### Rollback expectations

`service.restart` has no rollback in the traditional sense — it's an
action to re-attempt, not a state change to undo. If a restart fails, the
system leaves the service in whatever state the restart call reported,
surfaces it as a new `AlertRecord` (existing alerts model, no new
mechanism), and does **not** auto-retry.

Config-push's rollback story is unresolved and explicitly blocks that
operation's own future ADR — it is not answered here.

### Out of scope

- Config-push, service start/stop (deferred to a follow-up ADR).
- Any new role.
- Any new elevation mechanism beyond reusing `RequireElevation`.
- Any new audit mechanism beyond reusing the existing writer.
- Dry-run for any operation other than `service.restart`.

## Consequences

Once implemented, `service.restart` becomes the template for every
subsequent lab-mutating operation: same role gate, same elevation pattern
with a new action string, same audit writer with a new action name, and a
dry-run preview required before UI exposure. Config-push and start/stop
each need their own ADR because their rollback and blast-radius stories
differ from restart's and are not resolved by this decision.

### Addendum (#282): runbook steps are an additional entry point

Runbook steps (`runbook_steps`, `POST
/api/runbooks/{id}/steps/{stepId}/execute`) let a runbook point at a
connector's `restart`/`start`/`stop`, but they are not a new operation and
not a new guardrail model — they are a second entry point into the exact
same boundaries this ADR (and its start/stop follow-up) already defines:
the same `connector.<verb>` elevation action, the same operator-grant
check, the same dry-run-preview/elevation-gated-mutate split, the same
failure-alert-instead-of-rollback behavior, and the same audit writer
(`connector.<verb>`, now additionally carrying `runbookId`/`stepId` in
`detail`). Authoring a step is deliberately **not** gated by a connector
grant — only linking guidance to a connector, same as issue #42 established
for the pre-existing snapshot/doc pointers — so the operator-grant check
happens at execution time, not authoring time.
