# 0002 — Start/stop lab-mutating operations

Status: accepted

## Context

ADR 0001 shipped `service.restart` as the first lab-mutating operation and
explicitly deferred start/stop and config-push to follow-up ADRs. This ADR
is start/stop's turn (config-push is ADR 0003). Restart already proved the
model — role gate, elevation step-up, audit writer, mandatory dry-run
preview — works end to end; start/stop reuse it verbatim rather than
inventing anything new.

## Decision

### Eligible ops: `service.start`, `service.stop`

Same connector categories PR1 scoped restart across: proxmox, docker,
opnsense, pfsense, pihole. `dnsresolver` has no restart either — it
monitors a DNS Resolver service it doesn't own the lifecycle of (the
pfSense/OPNsense connectors own that) — so it gets neither start/stop nor
restart, mirrored from PR1's scope decision, not re-litigated here.

### Authorization

Reuse `auth.RequireRole("operator")` as-is. No new role.

### Confirmation / step-up

Reuse `auth.RequireElevation(jwtSvc, action)` with two new action strings,
`"connector.start"` and `"connector.stop"`, wired the same way as
`"connector.restart"`. Same 60s `elevationTTL`, same `POST
/api/auth/elevate` flow. No new elevation mechanism.

### Audit

Reuse the existing audit writer. New rows:

| Action | Trigger | targetType / targetId |
|---|---|---|
| `connector.start` | `POST /api/connectors/{id}/start` (dryRun absent/false) | connector / id |
| `connector.stop` | `POST /api/connectors/{id}/stop` (dryRun absent/false) | connector / id |

Successes only, same as restart.

### Dry-run

Same preview contract as restart: `POST /api/connectors/{id}/start(or
/stop)?dryRun=true` returns `targetService`, `estimatedDowntimeSeconds`,
and `dependentServices` from the latest stored snapshot, without touching
the connector. Mandatory before the real action is exposed in the UI.

`estimatedDowntimeSeconds` is less meaningful for `stop`: a restart or
start has a roughly bounded window, but a stop's downtime is indefinite
until an explicit start. Rather than invent a new field for that, the stop
preview reports `estimatedDowntimeSeconds: 0`.

### Rollback expectations

Same "re-attempt, not undo" story as restart. `start` is idempotent-safe —
calling it against an already-running service is a no-op success at the
connector/vendor-API level, not an error the caller needs to special-case.
`stop` isn't destructive to configuration — nothing is deleted or
rewritten, only a running process is halted — so there is nothing to roll
back beyond starting it again. A failed start/stop surfaces as a new
`AlertRecord` (critical), not an automatic retry.

### Out of scope

- Config-push (ADR 0003).
- Any new role.
- Any new elevation mechanism beyond reusing `RequireElevation`.
- Any new audit mechanism beyond reusing the existing writer.

## Consequences

Start/stop add no new machinery — same role, same elevation pattern with
two new action strings, same audit writer, same dry-run contract minus one
field's meaning for stop. The one genuine decision left for
implementation was how to represent "start"/"stop" for connectors whose
vendor API has no literal start/stop verb (Pi-hole's FTL service exposes
only `restartdns`): those map to the closest equivalent the API actually
exposes (Pi-hole's DNS blocking toggle) rather than being left
unimplemented, documented at the call site.
