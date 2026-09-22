# Scheduled Doc Export

Writes every generated doc to a local directory as Markdown on a schedule,
so documentation survives outside the tool — and doubles as a lightweight,
human-readable secondary backup alongside the JSON bundle described in
[BACKUP.md](BACKUP.md). Implemented in `backend/internal/docexport/`.

This is distinct from `GET /api/system/backup/export` (a single JSON bundle
meant for restoring into another WiseLabz instance): doc export writes one
`.md` file per doc, using exactly the content the doc engine already
generates and persists (`store.DocRecord.Content`) — no re-rendering.

## Configuration

Config-file only for this first cut — there's no operator-facing API to
change it at runtime, matching the "quality"/"digest"/"backup-verify"
scheduled jobs in `cmd/server/main.go` rather than the schedule-via-API
pattern used for the primary backup job.

```yaml
doc_export:
  enabled: false            # opt-in; default false
  dir: ./data/docexport     # target directory
  cron_expr: "0 2 * * *"    # daily at 2 AM by default
```

Or via environment variables: `WISELABZ_DOC_EXPORT_ENABLED`,
`WISELABZ_DOC_EXPORT_DIR`, `WISELABZ_DOC_EXPORT_CRON_EXPR`.

## Behavior

- Each run writes one file per doc to `dir`, named `<slugified-title>-<first
  8 chars of doc ID>.md`, and mirrors `dir` to the current set of docs: any
  `.md` file left over from a doc that was since renamed or deleted is
  removed on the next run.
- `dir` is created if it doesn't exist.
- Files outside the `.md` set the exporter manages (e.g. a README an
  operator drops into the directory) are left untouched.
- Failures are logged and skip that run; they don't crash the server (same
  convention as the other scheduled jobs).

## Scope and follow-up

This first cut only writes to a local directory. Exporting straight to a
Git remote (clone/pull, commit, push) was considered but deferred: it would
require adding a Git library (e.g. `go-git`) that isn't currently a
dependency, which is out of scope for this change's effort. Point `dir` at
a directory that's itself a Git working copy and drive the commit/push step
externally (a cron job, CI, `git pull` on a remote) as a workaround until
native Git remote support lands — tracked as a follow-up issue.
