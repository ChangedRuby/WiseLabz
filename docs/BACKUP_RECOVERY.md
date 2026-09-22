# Backup Recovery: What Comes Back, and What Doesn't

This covers the application-level backup (`internal/backup`, documented in
full in [BACKUP.md](BACKUP.md)): how to confirm a bundle actually restores,
and — since [BACKUP.md's Exclusions](BACKUP.md#whats-excluded-and-why)
redact every secret — what an operator has to manually re-enter after a
restore. It does not cover full database backups
([DEPLOYMENT.md#backups](DEPLOYMENT.md#backups)), which include secrets and
restore a byte-identical instance.

Encryption at rest for backup files and shipping backups to a remote target
(S3/SFTP/WebDAV) are explicitly out of scope here — see the issue this doc
was written for. Backups are local files, protected only by filesystem
permissions (0600/0700; see BACKUP.md).

## 1. Every export gets a manifest and a checksum

`backup export` (and the scheduled backup job) writes two files:

- `wiselabz-backup-<timestamp>.json` — the bundle itself.
- `wiselabz-backup-<timestamp>.manifest.json` — a sidecar recording the
  bundle's sha256 checksum, per-entity row counts, the app version
  (`go install`-style module version, or `dev`), and the schema (migration)
  version at export time.

`backup.Import`/`ImportFromFile` and `backup verify`/`backup restore` (below)
all check the bundle's checksum against the manifest before trusting its
contents. A bundle without a manifest (hand-edited, or exported before this
feature existed) is still accepted — the checksum check is simply skipped —
matching Import's existing tolerance for externally-authored bundles.

## 2. Verifying a backup actually restores

Row counts and a checksum in a manifest only prove the file wasn't corrupted
on disk — they don't prove the bundle *imports* cleanly. Two things do that:

**The scheduled verify job** (`backup.RunVerifyOnce`, registered in
`cmd/server/main.go` as the `backup-verify` cron job, daily by default —
see `backup.DefaultVerifyCronExpr`) restores the most recent backup bundle
into a throwaway in-memory SQLite database — never the real one — and
compares the imported row counts against the manifest. Results are appended
to `{backupDir}/verifications.jsonl` (newest-first via
`backup.ListVerifications`); a failure is logged at `error` level the same
way other scheduled jobs (retention, sync) surface failures, so it shows up
in whatever log aggregation already watches those.

**On demand**, run the same check from the command line:

```sh
# verify the most recent bundle in a directory
./backup verify -dir /opt/wiselabz/data/backups

# verify a specific file
./backup verify -file /opt/wiselabz/data/backups/wiselabz-backup-20260101-030000.json
```

Exit code 0 means the bundle's checksum matches its manifest and every
entity's row count matches after a full restore into a scratch database.
Non-zero means don't trust that backup — check the next-most-recent one.

## 3. Restoring for real

```sh
./backup restore -file /opt/wiselabz/data/backups/wiselabz-backup-20260101-030000.json
```

This runs the same checksum + row-count verification as `backup verify`
first and aborts if it fails, then imports into the database from the
running config (`WISELABZ_DB_DRIVER`/`WISELABZ_DB_DSN`, or `config.yaml`) —
the same additive, idempotent `Import` the API's
`POST /api/system/backup/import` uses: records whose ID already exists are
left untouched, never overwritten. Pass `-yes` to skip the confirmation
prompt (e.g. in a scripted recovery runbook).

## 4. What a restore does *not* bring back

Because export redacts every secret (see BACKUP.md's Exclusions), a restore
onto a fresh instance leaves these non-functional until reconfigured:

- **Connector credentials** — each connector row comes back with its secret
  fields (Proxmox `token_secret`, OPNsense `api_key`/`api_secret`, etc.)
  empty. The connector will show as offline/degraded on its next sync
  attempt. Re-enter credentials for each connector via
  `PUT /api/connectors/{id}` (or the UI's connector edit form) — the
  non-secret fields (URL, name, category) are already restored, so this is
  filling in one or two fields per connector, not reconfiguring from
  scratch.
- **AI provider API key** — `aiConfig` in the bundle is informational only
  and is never applied on import (a restored instance would otherwise show
  AI as "enabled" with no usable key). Re-enable and re-enter the key via
  `PUT /api/settings/ai` if AI features were in use.
- **Notification channels** — not exported at all (their config is an
  unstructured blob that can't be safely redacted). Reconfigure each
  channel (email/webhook/etc.) manually via `PUT /api/settings/notifications`.

Docs, doc versions, and templates/template sections carry no secrets and are
restored as-is.

## Recovery runbook (suggested order)

1. Provision the new instance; point it at a fresh (or restored) database.
2. `./backup verify -file <bundle>` — confirm it's good before touching
   anything.
3. `./backup restore -file <bundle>` (or `POST /api/system/backup/import`
   from the UI, if the API is already reachable).
4. Re-enter connector secrets, the AI key, and notification channel config
   (Section 4).
5. Trigger a manual sync per connector (or wait for the next scheduled one)
   to confirm credentials work.
