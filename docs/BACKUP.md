# Configuration & Documentation Backup (Export/Import)

This is a portable, application-level backup of WiseLabz configuration and
content — distinct from the full database backup described in
[DEPLOYMENT.md](DEPLOYMENT.md#backups) (which copies the entire SQLite/
Postgres database, secrets included, for disaster recovery of *this exact
instance*). This export is meant to move connectors, docs, and templates
between instances, or as a lightweight, secret-free config snapshot —
excludes secrets by default, and no separate DB tooling is needed.

## Endpoints

Both are operator-only (403 for `viewer`), not step-up-gated (no elevation
token required — this isn't a destructive action):

- `GET /api/system/backup/export` — returns the bundle as JSON, with
  `Content-Disposition: attachment; filename="wiselabz-backup-<timestamp>.json"`.
- `POST /api/system/backup/import` — accepts a bundle (as exported, or
  hand-edited) in the request body.

## What's included

- **Connectors** — all fields except secret configuration fields (see
  Redaction below).
- **Docs** — all documentation records, plus every historical version of
  each (`GET /api/docs/{id}/versions` equivalent).
- **Templates** — all template records, plus every section of each.
- **AI config summary** (`aiConfig`) — `enabled`, `provider`, `model`,
  `baseUrl`, `mode`. Informational only (see Exclusions).

## What's excluded, and why

- **Connector secret fields** — any config field a connector type's schema
  marks `password` (e.g. Proxmox `token_secret`, OPNsense `api_key` /
  `api_secret`) is stripped from `configData` before export. Non-secret
  fields (URLs, IDs, toggles) are kept.
- **AI provider API key** — `ai_config.api_key_encrypted` is never read by
  the export, let alone written to the bundle. `aiConfig` in the bundle is a
  read-only summary of the non-secret fields, exported for operator
  visibility; **it is never applied on import** — a working AI config in the
  target instance would otherwise be silently broken by importing a summary
  with no usable key.
- **Notification channel config** (`notification_config.config_json`) — not
  exported at all. It's an arbitrary JSON blob (SMTP credentials, webhook
  URLs/secrets) with no schema WiseLabz can use to redact just the secret
  parts. Reconfigure notification channels manually on the target instance.

## Bundle format

Top-level JSON fields: `version` (integer, currently `1`), `exportedAt`
(RFC3339 timestamp), `connectors`, `docs`, `docVersions`, `templates`,
`templateSections` (arrays mirroring the corresponding store records), and
optionally `aiConfig`. See `docs/openapi.yaml` (`BackupBundle` schema) for
the exact shape.

## Import behavior

- **Validates before writing anything.** `version` must match the format
  version this instance supports; every `docVersions[].docId` must reference
  a `docs[].id` present in the bundle; every `templateSections[].templateId`
  must reference a `templates[].id` present in the bundle; every
  `connectors[].category` must be one of `virtualization`, `containers_paas`,
  `networking`. The first validation failure aborts the whole import with a
  400 `invalid_backup` response — no partial writes.
- **Additive and idempotent.** Each record is created only if its ID doesn't
  already exist in the target instance; an existing ID is left untouched and
  counted as `skipped` rather than overwritten. Importing the same bundle
  twice is safe — the second run reports everything as skipped.
- **`aiConfig` is never imported** (see Exclusions above).
- The response reports per-entity `{ imported, skipped }` counts
  (`BackupImportResult` in `docs/openapi.yaml`).

## Manifest, checksum, and verification

Every bundle written by `ExportToFile` (the scheduled backup job, and
`POST /api/system/backup/run`) gets a `<bundle>.manifest.json` sidecar with
its sha256 checksum, per-entity row counts, and app/schema version.
`Import`ing from a file checks the checksum first; a scheduled job restores
the latest bundle into a scratch in-memory database nightly to confirm it
actually comes back intact; and `backup verify`/`backup restore` (a new
`cmd/backup` binary) do the same checks from the command line, plus an
actual restore. See [BACKUP_RECOVERY.md](BACKUP_RECOVERY.md) for the full
verify/restore workflow and, since secrets are redacted (see Exclusions
above), what has to be manually re-entered after a restore.
