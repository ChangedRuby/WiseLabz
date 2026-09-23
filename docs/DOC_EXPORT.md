# Scheduled Doc Export

Writes every generated doc to a local directory as Markdown on a schedule,
so documentation survives outside the tool — and doubles as a lightweight,
human-readable secondary backup alongside the JSON bundle described in
[BACKUP.md](BACKUP.md). Optionally, the directory is a clone of a Git
remote and every run commits and pushes the changes, giving the docs a real
Git history outside WiseLabz. Implemented in `backend/internal/docexport/`.

This is distinct from `GET /api/system/backup/export` (a single JSON bundle
meant for restoring into another WiseLabz instance): doc export writes one
`.md` file per doc, using exactly the content the doc engine already
generates and persists (`store.DocRecord.Content`) — no re-rendering.

## Configuration

Config-file (or env) only — there's no operator-facing API to change it at
runtime, matching the "quality"/"digest"/"backup-verify" scheduled jobs in
`cmd/server/main.go` rather than the schedule-via-API pattern used for the
primary backup job.

```yaml
doc_export:
  enabled: false            # opt-in; default false
  dir: ./data/docexport     # target directory (the persistent clone in Git mode)
  cron_expr: "0 2 * * *"    # daily at 2 AM by default
  git:                      # optional; Git mode is on when remote is set
    remote: https://github.com/acme/lab-docs.git
    branch: main            # default main
    path: docs              # repo subdirectory the docs go to; default docs
    author_name: WiseLabz   # default WiseLabz
    author_email: wiselabz@localhost  # default wiselabz@localhost
    token: ""               # HTTPS access token (prefer the env var)
    ssh_key_path: ""        # SSH private key file (ssh remotes only)
    ssh_known_hosts: ""     # known_hosts file (required for ssh remotes)
    insecure_skip_host_key: false     # skip SSH host key checks (not recommended)
```

| Key | Env var | Default | Notes |
|---|---|---|---|
| `doc_export.enabled` | `WISELABZ_DOC_EXPORT_ENABLED` | `false` | |
| `doc_export.dir` | `WISELABZ_DOC_EXPORT_DIR` | `./data/docexport` | In Git mode: the persistent clone |
| `doc_export.cron_expr` | `WISELABZ_DOC_EXPORT_CRON_EXPR` | `0 2 * * *` | 5- or 6-field cron |
| `doc_export.git.remote` | `WISELABZ_DOC_EXPORT_GIT_REMOTE` | empty (local mode) | `https://…`, `ssh://…` or scp-style `git@host:org/repo.git`. No credentials in the URL |
| `doc_export.git.branch` | `WISELABZ_DOC_EXPORT_GIT_BRANCH` | `main` | Created by the first push if missing |
| `doc_export.git.path` | `WISELABZ_DOC_EXPORT_GIT_PATH` | `docs` | Relative, inside the repo |
| `doc_export.git.author_name` | `WISELABZ_DOC_EXPORT_GIT_AUTHOR_NAME` | `WiseLabz` | Author and committer |
| `doc_export.git.author_email` | `WISELABZ_DOC_EXPORT_GIT_AUTHOR_EMAIL` | `wiselabz@localhost` | |
| `doc_export.git.token` | `WISELABZ_DOC_EXPORT_GIT_TOKEN` | empty | HTTPS only; secret, never logged |
| `doc_export.git.ssh_key_path` | `WISELABZ_DOC_EXPORT_GIT_SSH_KEY_PATH` | empty | Required for SSH remotes; unencrypted key |
| `doc_export.git.ssh_known_hosts` | `WISELABZ_DOC_EXPORT_GIT_SSH_KNOWN_HOSTS` | empty | Required for SSH remotes unless insecure |
| `doc_export.git.insecure_skip_host_key` | `WISELABZ_DOC_EXPORT_GIT_INSECURE_SKIP_HOST_KEY` | `false` | Logs a WARN at startup and on every run |

The server refuses to start (`server config validate` reports it) when the
Git settings are inconsistent: a scheme other than https/ssh, credentials
embedded in the remote URL, a token together with any `ssh_*` setting, an SSH
remote without `ssh_key_path`, or an SSH remote without `ssh_known_hosts`
while `insecure_skip_host_key` is off.

## Behavior

- Each run writes one file per doc, named `<slugified-title>-<first 8 chars
  of doc ID>.md`, to `dir` (local mode) or `dir/<git.path>` (Git mode), and
  mirrors that directory to the current set of docs.
- **Safe prune:** only files matching the generated name pattern
  (`<slug>-<8 hex>.md`, or `<uuid>.md` for the rare collision fallback) are
  ever deleted, so a `README.md` or any other hand-written file next to the
  export survives every run, in both modes.
- The directory is created if it doesn't exist.
- Failures are logged and skip that run; they don't crash the server (same
  convention as the other scheduled jobs).

## Git mode

The Git implementation is pure Go (`go-git`), because the runtime image is
`distroless/static` and has no `git` binary.

**Workdir.** `doc_export.dir` is a persistent clone. If it's empty or
missing, the first run initialises it with `origin` set to the remote. If
it's non-empty, it must already be a Git repository whose `origin` URL is
exactly `doc_export.git.remote`; anything else makes the run fail rather
than touch a directory the exporter doesn't own.

**Each run:**

1. Fetch `origin/<branch>` and hard-reset the clone to it (untracked files
   are cleaned). An empty remote, or a missing branch, is fine: the first
   push creates it.
2. Export the docs into `<dir>/<path>` (with the safe prune above).
3. Stage additions, edits and removals under `<path>`. If nothing changed,
   no commit is made.
4. Commit once with the configured author, message
   `docs: export N docs (+added ~modified -removed)` and the per-kind file
   lists in the body.
5. Push `<branch>`, **never forced**.

**Divergence.** Because every run starts from `origin/<branch>`, commits
other people push to the same branch (a README, CI config, manual edits
outside the generated files) are kept. If the remote moves between the fetch
and the push, the push is rejected as non-fast-forward: the run is logged as
failed and simply retried on the next schedule, which fetches the new head
and rebuilds the export on top of it. Local unpushed commits are never kept
between runs; the export is always regenerated from the store.

**Authentication.**

- HTTPS: set `token` (e.g. a GitHub fine-grained PAT with *Contents:
  read and write* on the one repo). It's sent as basic auth with the user
  `x-access-token`. Traffic goes through the shared hardened HTTP client
  (`internal/httpx`: TLS 1.2+, no redirects); the transfer is bounded by the
  job's context rather than a fixed client timeout. Without a token, the
  remote must accept anonymous pushes.
- SSH: set `ssh_key_path` to an unencrypted private key and
  `ssh_known_hosts` to a known_hosts file containing the server's host key
  (e.g. `ssh-keyscan github.com > known_hosts`, then verify the fingerprint).
  The user comes from the URL (`git` if absent). `insecure_skip_host_key:
  true` disables host key verification entirely; it logs a WARN at startup
  and on every run and should only be used on a trusted network.

## Failure notifications

The exporter emits the `system.job_failed` notification event through the
normal notification dispatcher:

- severity `warning` when the job goes from OK to failing (e.g. a revoked
  token, an unreachable remote, a rejected push);
- severity `info` when it recovers.

Only these transitions notify, so a job that keeps failing sends one
notification, not one per run. The state lives in memory: after a server
restart, the first failure notifies again. Route the event per channel under
**Settings → Notifications** (the `system job failed` row); it's also
included in digests.
