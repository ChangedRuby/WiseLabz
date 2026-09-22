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
| `connector.sync_all` | `POST /api/sync` | connector / (none) |
| `connector.restart` | `POST /api/connectors/{id}/restart` (dryRun omitted/false) | connector / id |
| `connector.start` | `POST /api/connectors/{id}/start` (dryRun omitted/false) | connector / id |
| `connector.stop` | `POST /api/connectors/{id}/stop` (dryRun omitted/false) | connector / id |
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

## What's not recorded

- **Failed attempts, except elevation denials.** An audit entry is written
  only after the action itself succeeds. A failed create/update/delete never
  reaches the audit log — this is a deliberate scope cut: this endpoint
  answers "what happened," not "what was attempted." Elevation-token denials
  are the exception because they are security-relevant attempts.
- **Reads.** Listing or viewing a resource is not an audited action.
- A write to the audit log failing is logged (`slog.Error`) but never
  fails the request — the audited action has already gone through by
  that point, so refusing the response would be misleading.

## Retention

Audit rows are not touched by the retention policies in
`backend/internal/store/retention.go` — they're kept indefinitely, since
their purpose is historical accountability rather than operational data. The
elevation request/denial rows inherit the same exemption through `audit_log`.
