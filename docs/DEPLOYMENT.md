# WiseLabz — Deployment Guide

Supported deployment modes: Docker, docker-compose, a bare Go binary with
systemd, and PostgreSQL active/passive replicas behind a readiness-aware load
balancer.

## Scaling & high availability

The default is one application instance. SQLite supports only this model.
With PostgreSQL, set `ha.leader_election: true` on every replica and direct
the load balancer to instances whose `/readyz` returns 200. Standbys keep
`/healthz` available but return 503 on `/readyz` until they acquire the
advisory lock. Only the leader runs background jobs. On lock loss it exits;
the process manager should restart it as a standby. Takeover starts at the
next `ha.lock_poll_interval` poll (5 seconds by default), plus any process
restart time.

| Component | Multiple replicas |
|---|---|
| Scheduler jobs (sync, digest, quality, reports, doc export, backups, retention) | Leader only |
| Notification delivery retrier | Leader only |
| Document lock sweep | Leader only |
| WebSocket hub and tickets | Process-local; route clients only to the ready leader (see [ADR 0005](adr/0005-cross-replica-websocket-relay.md)) |
| Rate limiter and TTL caches | Process-local; reset on failover |
| Backups on local disk | Leader only, but store backup files on shared or durable storage |
| Migrations | Safe to start concurrently; golang-migrate locks PostgreSQL migrations |
| Scheduled sync claims | Protected by a database lease even across replicas |

Tune connector fan-out with `sync.max_concurrency` (default 4),
`sync.due_batch_size` (default 50 per tick), and `sync.timeout` (default 5m
per connector). The claim lease lasts one minute beyond the timeout.

Active/passive is the only supported multi-replica mode. Routing every client
to the ready leader keeps WebSocket events, tickets, and per-process state on
one process. Do not mark standbys ready or load-balance across replicas.
Active/active is deferred:
[ADR 0005](adr/0005-cross-replica-websocket-relay.md) records the planned
PostgreSQL LISTEN/NOTIFY event relay and the process-local state (tickets, user
cache invalidation, in-process locks, rate limits) that must be addressed
first. A relay alone would not make active/active safe.

## PostgreSQL support

`db.driver: postgres` runs schema migrations (the SQL migration files in
`backend/internal/store/migrations/postgres/` are portable DDL) and then
runs the application's runtime queries normally. The `backend/internal/store/`
query layer is written using SQLite-style `?` placeholders; when
`db.driver` is `postgres`, `Store` wraps the connection
(`backend/internal/store/pgdb.go`) to rewrite `?` placeholders to
PostgreSQL's `$1, $2, ...` form transparently, so no query code needs to
differ between drivers.

Both the `docker-compose.yml` (Postgres) and `docker-compose.sqlite.yml`
reference stacks are fully functional end-to-end.

The Postgres reference stack is smoke-tested in CI with
`scripts/compose-smoke.sh`: it starts with an empty database, waits for the
application health check, authenticates the bootstrap operator, and creates
then lists a connector. Run that script locally when changing Compose or
container startup behavior; it uses an isolated environment, port, and
temporary Compose volumes that it removes on exit.

## WebSocket behind a reverse proxy

`/api/ws` shares the same port and origin as the rest of the HTTP API — no
separate service or port. Reverse proxy behavior varies:

- **Traefik / Caddy**: handle WebSocket upgrades automatically, no extra
  config needed.
- **nginx**: needs explicit `proxy_set_header Upgrade $http_upgrade;`,
  `proxy_set_header Connection "upgrade";`, `proxy_http_version 1.1;`, and a
  bumped `proxy_read_timeout` (nginx's default 60s will silently drop idle
  WebSocket connections).
- **Cloudflare (orange-cloud/proxied)**: has a ~100s idle timeout on
  WebSocket connections. Either send an application-level heartbeat/ping
  more frequently than that, or set the DNS record to "DNS only" (grey
  cloud) to bypass Cloudflare's proxy for this traffic.

## Backups

**SQLite**: either stop the container and copy the `.db` file directly, or
use `sqlite3 /data/wiselabz.db ".backup /data/backup.db"` while the app is
running (SQLite's online backup API, safe under concurrent access). Example
cron entry on the host (bare-metal/systemd mode):
```
0 3 * * * sqlite3 /opt/wiselabz/data/wiselabz.db ".backup /opt/wiselabz/backups/wiselabz-$(date +\%Y\%m\%d).db"
```

**PostgreSQL**: standard `pg_dump`, e.g. as a cron job against the
`postgres` compose service:
```
0 3 * * * docker compose exec -T postgres pg_dump -U wiselabz wiselabz | gzip > backups/wiselabz-$(date +\%Y\%m\%d).sql.gz
```

## systemd (bare binary)

A reference unit file is at `deploy/wiselabz.service`. Key points:
- `WorkingDirectory=` is pinned to avoid relative-path surprises (the
  default DSN is now an absolute path, `/data/wiselabz.db`, but pin this
  anyway for config file resolution).
- `EnvironmentFile=` should point at a `root:wiselabz`-owned, mode `0640`
  file holding `WISELABZ_AUTH_SECRET` / `WISELABZ_ENCRYPTION_KEY` /
  `WISELABZ_SERVER_ORIGIN` / `WISELABZ_ADMIN_PASSWORD` / etc.
- Default `StandardOutput=journal` is sufficient — no extra logging config
  needed; use `journalctl -u wiselabz`.
- `TimeoutStopSec=15s` gives the app's graceful shutdown
  (`Server.ShutdownTimeoutSeconds`, default 10s) room to finish before
  systemd sends `SIGKILL`.
