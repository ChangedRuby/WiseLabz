# Diagnostics Bundle

A downloadable, secret-free snapshot of instance state for support requests —
distinct from the [backup export](BACKUP.md) (which is meant to move
config/content between instances). The diagnostics bundle is meant to be
attached to a support ticket or bug report: health, versions, sanitized
config, and recent failures, with nothing an operator would need to redact
by hand first.

## Endpoint

- `GET /api/system/diagnostics` — operator-only (403 for `viewer`), not
  step-up-gated. Returns the bundle as JSON, with
  `Content-Disposition: attachment; filename="wiselabz-diagnostics-<timestamp>.json"`.

## What's included

- **`health`** — the same DB-ping check as `GET /api/health` and `GET /readyz`.
- **`versions`** — Go/build version information collected for the operator-only bundle.
- **`sanitizedConfig`**:
  - `connectors` — name/type/category/enabled, plus `configData` with any
    `password`-typed field redacted (see Redaction below).
  - `syncSchedule` — the configured cron schedule.
  - `authProviders` — `local` (always true) and `oidc` (configured provider
    `id`/`displayName` only — no client ID or secret).
  - `aiConfig` — `enabled`, `provider`, `model`, `baseUrl`, `mode`. Same
    shape as the backup bundle's AI config summary.
- **`recentFailures`** — the most recent 50 (bounded, not configurable)
  failed sync runs (`connectorId`, `connectorName`, `startedAt`, `error`)
  and the most recent 50 failed notification deliveries (`channel`,
  `attempts`, `lastError`, timestamps — no channel *config*, just the
  delivery attempt record).
- **`generatedAt`** — RFC3339 timestamp.

## What's excluded, and why

- **Connector secret fields** — redacted from `configData` exactly like the
  backup export (`internal/backup.RedactConnectorConfig`): any field a
  connector type's schema marks `password` is stripped.
- **AI provider API key** — `ai_config.api_key_encrypted` is never read.
- **Notification channel config** — excluded entirely, not just redacted.
  It's an arbitrary JSON blob (SMTP credentials, webhook URLs/secrets) with
  no schema to redact just the secret parts, so the whole thing is left out.
  Only the delivery *attempt* records (channel name, status, error) are
  included, never the channel's own config.
- **OIDC client secrets** — only `id` and `displayName` are included per
  provider; `clientId`, `clientSecret`, and `issuerUrl` are omitted.
- **The app's own JWT signing secret** (`auth.secret`) — never touched.

## Bundle format

See `docs/openapi.yaml` (`DiagnosticsBundle` schema) for the exact shape.
