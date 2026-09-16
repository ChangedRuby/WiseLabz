# 0003 — Config-push lab-mutating operation

Status: accepted

## Context

ADR 0001 covered restart; ADR 0002 covered start/stop. Both are runtime
control with no persistent state change and a "re-attempt, not undo"
rollback story. Config-push is qualitatively different: it writes
persistent configuration into a live system. ADR 0001 deferred it
specifically because "it writes arbitrary content into a live system with
no generic rollback." This ADR is that rollback design.

## Decision

### Field-level partial update via a per-connector whitelist

Config-push does not accept an arbitrary config document. Each connector
that supports it declares a curated `[]connector.ConfigField` — a
deliberately narrower surface than what `Fetch`/`Validate` accept — and a
push targets exactly one `(entityRef, fieldKey, value)` triple per call.
This is the direct answer to ADR 0001's own risk framing: the risk is
arbitrary writes, so the fix is not exposing arbitrary write surface.

Representative starting fields (illustrative, not exhaustive — more can be
added to a connector's `WritableFields()` later without touching the
handler):

- **proxmox**: VM/container `memory` (MB), `cores`.
- **docker**: container `restartPolicy`. (Not image tag: Docker's Engine
  API has no live image swap for a running container — that needs a full
  stop/remove/recreate, which is a materially different and riskier
  operation than every other field-level push here. `restartPolicy` is
  what `/containers/{id}/update` can actually patch in place. Image-tag
  push is left for a future extension once a recreate-with-rollback shape
  is designed on its own.)
- **opnsense / pfsense**: a firewall rule's `enabled` flag, keyed by rule
  ID/UUID as `entityRef`.
- **pihole / dnsresolver**: a DNS record's `ip`, keyed by hostname as
  `entityRef`.

`custom` does not implement `ConfigPusher` — it has no fixed config schema
to whitelist against.

### Snapshot-before-write + verify-diff-after

The handler calls `Fetch` before the push (`pre`) and after (`post`), and
runs `sync.Compare(pre, post)` (`backend/internal/sync/diff.go:36`). A
non-empty diff is "the write landed as *some* change" — the field-scoped
verification ADR 0001 asked for. A stricter per-field-key value check
would need a fixed field→section mapping maintained per connector on top
of the whitelist; the diff-based check is the smaller, sufficient version
for this connector set (see `configPushLanded` in
`backend/internal/api/connectors/handlers.go`).

### Auto-revert-then-alert on mismatch

An empty diff (mismatch) does **not** stop at a manual-revert-only flow:

1. The handler calls `ConfigPush` again with the field's pre-push value —
   the revert. That value comes from the request's `previousValue`, which
   the frontend already holds (it read the field's current value from
   `GET /{id}/data` to populate the edit form before the user changed it).
   There is no generic way to recover "the old value of an arbitrary
   `fieldKey`" from a rendered `ServiceSnapshot` (markdown sections, not
   structured key/value pairs), so the caller supplies it rather than the
   handler inferring it.
2. A `critical` `AlertRecord` is raised describing the mismatch and that
   an auto-revert was attempted.
3. **If the revert call itself errors**, that is the one hard stop in this
   whole feature set: the alert is escalated (still `critical`, body notes
   the revert failed and manual intervention is required), and no further
   automation is attempted — no retry loop, no second revert attempt.

The endpoint returns `409 config_push_mismatch` either way; the alert (and
its severity/description) is what distinguishes "auto-revert succeeded"
from "manual intervention needed."

### Authorization / confirmation / audit

Same pattern as restart/start/stop: `auth.RequireRole("operator")`, new
elevation action `"connector.configPush"`, existing audit writer.

| Action | Trigger | targetType / targetId |
|---|---|---|
| `connector.configPush` | `POST /api/connectors/{id}/config-push` (successful, verified push only) | connector / id |

Pushed field *names* are always recorded in `detail`; values are not,
following the same discipline `connector.update` already uses for
`config_data` fields — this endpoint writes into the same kind of
potentially-sensitive connector config surface, so the same redaction
default applies without new logic.

### Rollback expectations

Auto-revert on mismatch, described above, *is* the rollback story — there
is no separate manual-revert UI, because the automatic path already
attempts it every time a mismatch occurs. A revert-failure alert is the
signal a human needs to step in.

### Out of scope

- A generic full-document config replace (explicitly rejected — this is
  what ADR 0001 flagged as too risky).
- Multi-field batch pushes in one call — `ConfigPush`/the revert path both
  handle exactly one field per call, by design (see the `ponytail:`
  comment at the revert call site).
- Any new role or elevation mechanism beyond reusing `RequireElevation`.

## Consequences

This closes the gap ADR 0001 left open: config-push now has a resolved
rollback story (auto-revert-then-alert, with a hard escalation stop when
revert itself fails), gated by the same authorization/elevation/audit
machinery every other lab-mutating operation uses. The whitelist mechanism
(`ConfigField`/`WritableFields()`) is the extension point for adding more
pushable fields per connector later, without touching the handler, the
elevation wiring, or the revert logic.
