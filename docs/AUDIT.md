# Audit Trail

A record of who did what, and when, for the operator actions sensitive
enough to need attribution — connector changes and syncs, elevation,
auth settings changes, and document restores.

## Endpoint

`GET /api/system/audit` — instance-admin only (403 otherwise, `auth.RequireInstanceAdmin`; #240 replaced the flat operator/viewer role with per-connector grants), paginated
newest-first like other list endpoints (`page`, `pageSize`; see
`docs/openapi.yaml`'s `Page`/`PageSize` parameters). Optional filters:

- `action` — exact match, e.g. `connector.create`.
- `targetType` — exact match, e.g. `connector`, `doc`.

Response shape mirrors `ChangePage`/`AlertPage`: `{ items, total, page, pageSize }`.

### Keyset (cursor) pagination

Offset pagination re-counts and re-skips rows on every page, which gets
expensive as `audit_log` grows. Sending `cursor` opts into keyset pagination
instead: `?cursor=` for the first page, then the previous response's
`nextCursor` for each page after it. The server answers
`WHERE (created_at, id) < (cursor) ORDER BY created_at DESC, id DESC LIMIT
pageSize`, which reads straight out of an index (migration 000033) and cannot
skip or repeat a row when several records share a `created_at`.

The cursor is opaque — never build one client-side. `nextCursor` is absent on
the last page and on requests that did not send `cursor`, so offset clients see
the unchanged four-key envelope. While a cursor is in use `page` is ignored;
`total` still reports the full filtered count.

`GET /api/changes` and `GET /api/connectors/{id}/syncs` take the same `cursor`
parameter. Since the sync-run response body is a bare array, its cursor comes
back in the `X-Next-Cursor` response header rather than an envelope field.

## What's recorded

Each entry (`AuditRecord` in `docs/openapi.yaml`) has: `actorUserId`,
`actorRole`, `action`, `targetType`, `targetId`, `detail` (a small JSON
object, action-specific), and `createdAt`.

| Action | Trigger | targetType / targetId |
|---|---|---|
| `connector.create` | `POST /api/connectors` | connector / new ID |
| `connector.update` | `PUT /api/connectors/{id}` | connector / id |
| `connector.delete` | `DELETE /api/connectors/{id}` | connector / id |
| `connector.toggle_enabled` | `PUT /api/connectors/{id}/enabled` | connector / id |
| `connector.sync` | `POST /api/connectors/{id}/sync` | connector / id |
| `snapshot.diff.export` | `GET /api/connectors/{id}/snapshots/diff?from=…&to=…&format=json\|csv\|md\|html` | connector / id |
| `connector.sync_all` | `POST /api/sync` | connector / (none) |
| `connector.restart` | `POST /api/connectors/{id}/restart` (dryRun omitted/false), or `POST /api/runbooks/{id}/steps/{stepId}/execute` for a `restart` step | connector / id |
| `connector.start` | `POST /api/connectors/{id}/start` (dryRun omitted/false), or `POST /api/runbooks/{id}/steps/{stepId}/execute` for a `start` step | connector / id |
| `connector.stop` | `POST /api/connectors/{id}/stop` (dryRun omitted/false), or `POST /api/runbooks/{id}/steps/{stepId}/execute` for a `stop` step | connector / id |
| `connector.configPush` | `POST /api/connectors/{id}/config-push` (successful, verified push only) | connector / id |
| `connector.maintenanceWindow.open` | `POST /api/connectors/{id}/maintenance-window` | connector / id |
| `connector.maintenanceWindow.close` | `DELETE /api/connectors/{id}/maintenance-window` (only when a window was actually active) | connector / id |
| `connector.bulk_sync` | `POST /api/connectors/bulk-sync` | connector / id — one record per resolved item |
| `connector.bulk_reauth` | `POST /api/connectors/bulk-reauth` | connector / id — one record per resolved item |
| `connector.bulk_restart` | `POST /api/connectors/bulk-restart` | connector / id — one record per resolved item |
| `auth.elevate` | `POST /api/auth/elevate` | action / the elevated action name |
| `auth.elevation_requested` | Any step-up-gated endpoint receiving `X-Elevation-Token` | action / the required action name |
| `auth.elevation_denied` | Failed elevation-token validation on a step-up-gated endpoint | action / the required action name |
| `auth.config.update` | `PUT /api/auth/config` | auth_config / (none) |
| `auth.provider.enabled` | `PUT /api/auth/providers/{id}/enabled` | oidc_provider / provider id |
| `doc.restore` | `POST /api/docs/{id}/versions/{rev}/restore` | doc / doc id |
| `change.ack` | `POST /api/changes/{id}/ack` | change / id |
| `change.dismiss` | `POST /api/changes/{id}/dismiss` | change / id |
| `change.bulk_ack` | `POST /api/changes/bulk-resolve` (`status: acknowledged`) | change / id — one record per resolved item |
| `change.bulk_dismiss` | `POST /api/changes/bulk-resolve` (`status: dismissed`) | change / id — one record per resolved item |
| `alert.resolve` | `POST /api/alerts/{id}/resolve` | alert / id |
| `alert.dismiss` | `POST /api/alerts/{id}/dismiss` | alert / id |
| `alert.snooze` | `POST /api/alerts/{id}/snooze` | alert / id |
| `finding.resolve` | `POST /api/findings/{id}/resolve` | finding / id |
| `runbook.create` | `POST /api/runbooks` | runbook / new ID |
| `runbook.update` | `PUT /api/runbooks/{id}` | runbook / id |
| `runbook.delete` | `DELETE /api/runbooks/{id}` | runbook / id |
| `compliance_rule.create` | `POST /api/compliance/rules` | compliance_rule / new ID |
| `compliance_rule.update` | `PUT /api/compliance/rules/{id}` | compliance_rule / id |
| `compliance_rule.delete` | `DELETE /api/compliance/rules/{id}` | compliance_rule / id |
| `compliance_rule.enable` | `PUT /api/compliance/rules/{id}` (enabled true) | compliance_rule / id |
| `compliance_rule.disable` | `PUT /api/compliance/rules/{id}` (enabled false) | compliance_rule / id |

`detail` never carries secret values. `connector.update` and
`auth.config.update` record which *fields* changed (a name list), not
their values — connector config can hold credentials, and this keeps the
audit log safe to expose to any instance admin without redaction logic.
`connector.configPush` follows the same discipline: `detail` records the
pushed `fieldKey` (name) and `entityRef`, never the pushed `value`.

`connector.restart`/`connector.start`/`connector.stop` triggered by
`POST /api/runbooks/{id}/steps/{stepId}/execute` carry `runbookId` and
`stepId` in `detail` alongside `entityRef`, so an execution can be traced
back to the runbook step that triggered it — linking a step still grants no
mutation permission by itself (see `docs/adr/0001-lab-mutating-operation-boundaries.md`);
the caller still needs an operator grant on the step's connector and a
valid elevation token for `connector.<verb>`. `runbook.create`/
`runbook.update` record the runbook's `title` and its `steps` (each step's
`id`, `connectorId`, `verb`, `entityRef`); `runbook.update` additionally
records `changedFields`, the list of top-level keys present in the request
body. `runbook.delete` records the deleted runbook's `title`.

## What's not recorded

- **Failed attempts, except elevation denials.** An audit entry is written
  only after the action itself succeeds. A failed create/update/delete never
  reaches the audit log — this is a deliberate scope cut: this endpoint
  answers "what happened," not "what was attempted." Elevation-token denials
  are the exception because they are security-relevant attempts.
- **Reads.** Listing or viewing a resource is not an audited action. Downloading
  a snapshot diff with `format` is the exception; viewing the same diff as JSON
  without `format` is not audited.
- A write to the audit log failing is logged (`slog.Error`) but never
  fails the request — the audited action has already gone through by
  that point, so refusing the response would be misleading.

## Retention

Audit rows are pruned by the retention job (`DeleteOldAuditRecords` in
`backend/internal/store/retention.go`, called from
`backend/internal/retention/retention.go`) like other operational tables, using
the `retention.audit_days` setting (env `WISELABZ_RETENTION_AUDIT_DAYS`,
default `180`). Rows with `created_at` older than that many days are deleted;
there is no "keep latest" guard. Setting `retention.audit_days` to `0` disables
audit pruning entirely, so rows are kept indefinitely. The elevation
request/denial rows are ordinary `audit_log` rows and follow the same policy.

## Filtering and export

`GET /api/system/audit` accepts these filters (all optional, combinable):

- `action` — exact match.
- `targetType` — exact match.
- `createdAfter` — inclusive lower bound on `createdAt` (RFC 3339 date-time).
- `createdBefore` — inclusive upper bound on `createdAt`.

`GET /api/system/audit/export?format=json|csv` (instance-admin only) returns
every record matching the same four filters, newest first. It is **not
paginated** and has no row limit, so narrow it with filters on large logs.
`format` defaults to `json`; any other value is a 400 (`invalid_format`).
Both formats are served as an attachment named
`wiselabz-audit-<UTC timestamp>.json|csv`. JSON is an array of `AuditRecord`
objects. CSV has a header row with the columns `id`, `actorUserId`,
`actorRole`, `action`, `targetType`, `targetId`, `detail` (the raw JSON string),
`createdAt`.
