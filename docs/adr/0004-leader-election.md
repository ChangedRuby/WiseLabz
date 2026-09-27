# 0004 — PostgreSQL leader election for background workers

Status: accepted

## Context

Scheduler jobs, notification retries, and document lock sweeping run in each
server process. Running two processes without coordination duplicates work and
notifications. The WebSocket hub, tickets, rate limiter, and TTL caches are
also process-local.

## Decision

PostgreSQL deployments may opt into active/passive operation with
`ha.leader_election`. Each server starts HTTP and its WebSocket hub, then polls
for one fixed PostgreSQL session advisory lock. A pinned database connection
holds the lock. Only its owner starts the scheduler, delivery retrier, and
document lock sweep. A standby returns 503 from `/readyz` while `/healthz`
remains live, so a readiness-aware load balancer routes users to the leader.
SQLite remains single-instance.

The leader pings its pinned connection every `ha.lock_poll_interval`. If the
connection fails, the lock may have been released. The process immediately
reports unready and exits non-zero after ordered shutdown; its supervisor
restarts it to campaign again. This prevents background work from continuing
under uncertain ownership. Normal shutdown closes the pinned connection
before the database pool.

## Consequences

Failover latency is about one poll interval plus any process restart and
readiness time. WebSocket connections reconnect to the newly ready leader.
Backups written to local disk need durable/shared storage for continuity.
LISTEN/NOTIFY fan-out for active/active WebSocket operation is deferred.
