# 0005 — Cross-replica WebSocket event relay for active/active

Status: accepted (implementation deferred until active/active is a concrete
deployment requirement)

## Context

[ADR 0004](0004-leader-election.md) supports PostgreSQL active/passive
operation: one elected leader runs background work, standbys return 503 from
`/readyz`, and the load balancer routes every client to the leader. The
WebSocket hub (`backend/internal/ws/ws.go`) is process-local, and so are the
tickets that authorize an upgrade. An event produced on one process reaches only
the clients connected to that process.

Events are produced in two places:

- HTTP handlers on whichever replica serves the request (connector sync
  triggers, lifecycle and config-push alerts, doc locks, AI suggestions,
  in-app notifications);
- leader-only background work (scheduled sync and quality checks, doc lock
  expiry, notification delivery retries).

Today's envelope is `{type, payload}` with no id or timestamp. The server keeps
no history and never replays. Global broadcasts go to every authenticated user,
connector events to that connector's readers, and
user-targeted events go to every connection that user has open.

Active/active means several replicas return 200 from `/readyz` behind a plain,
non-sticky load balancer. That breaks event delivery, and it also breaks other
process-local state that active/passive hides. This ADR records how events
would be relayed and what else must change first. It does not make
active/active a supported mode.

Target topology: N application replicas, one shared PostgreSQL primary, and no
new infrastructure (no Redis or NATS). SQLite remains single-instance.

## Decision

### Opt-in mode

Add `ha.mode: active_passive | active_active`, defaulting to
`active_passive`. `active_active` requires `db.driver: postgres` and
`ha.leader_election: true`. Config validation rejects it on SQLite, as it
already does for leader election. In `active_active` every replica reports
ready once started. Leader election still decides who runs background jobs, so
ADR 0004's leader-only guarantees are unchanged.

### Mechanism: PostgreSQL LISTEN/NOTIFY

Each replica listens on one channel (`wiselabz_ws`). The listener uses a
dedicated `pgx.Connect(dsn)` connection opened outside the `database/sql`
pool, so a long-lived LISTEN never takes one of the pool's slots. The listener
reconnects with capped exponential backoff.

### Event ownership: local first, then relay

The replica that produces an event delivers it to its own hub immediately and
then publishes it with `pg_notify`. Every relay message carries the producing
replica's `origin` ID, a random value generated at startup. A receiver drops
messages from its own origin, so no event is delivered locally twice and
nothing is re-relayed. Received events go straight to the local hub and are
never re-published, which prevents loops. Local delivery does not depend on
PostgreSQL being healthy.

### Relay payload

```json
{"v":1,"origin":"<replica>","id":"<uuid>","ts":"<RFC3339>","type":"alert.created",
 "audience":{"kind":"all"},"payload":{...}}
```

`audience` is `{"kind":"all"}` for `Broadcast`,
`{"kind":"user","userID":"..."}` for `BroadcastToUser`, and
`{"kind":"connector","connectorID":"..."}` for `BroadcastConnector`. A connector
audience carries only the connector ID, never the reader list: each replica
resolves the readers locally through `Hub.publish` when it delivers, so a grant
change is honored per replica and per event. `payload` and the envelope's
top-level `connectorId` are unchanged. PostgreSQL limits
NOTIFY payloads to 8000 bytes. A message larger than about 7.5 KB once encoded
is written to a short-lived `ws_relay` table, and the notification carries only
`{"v":1,"origin":..,"id":..,"ref":..}`. A periodic sweep deletes rows older
than a few minutes. A receiver that finds no matching row skips the event and
logs it.

### Duplicate suppression and ordering

`Envelope` gains `id` (UUID) and `ts`, as `docs/WS_CONTRACT.md` already
describes. Receivers keep a bounded set of recently seen IDs and drop repeats,
for example after a listener reconnect. Ordering is best-effort and holds only
per origin, because PostgreSQL delivers one session's notifications in commit
order. There is no global order across origins. Clients must treat events as
hints to refetch, not as an authoritative log.

### Reconnects and degraded relay

A relay failure never makes a replica unready. Local clients keep receiving
events produced locally. When the listener reconnects after a gap, the replica
sends a `system.resync` event to its local clients. The frontend invalidates
its volatile queries on `system.resync` and after every socket reconnect. Missed
events are recovered by refetching over REST, not by replay. A failed
`pg_notify` is logged and counted, and the producer does not retry.

### Authorization and secrets

The relay preserves each event's audience exactly: a user-targeted event is
delivered only to that user's connections on every replica, a connector event
only to that connector's readers, and a broadcast is never narrowed or widened. Relay payloads contain only what the event already
sends to browsers, never credentials, tokens or connector secrets.
PostgreSQL statement logging can record NOTIFY payloads, so this rule applies
to the relay too.

Per-connector filtering is implemented (#421): connector-scoped events are
delivered only to users holding a grant on the connector and to restricted API
keys covering it. The relay receiver delivers through the same `Hub.publish`
path, so it applies the same filter with that replica's database. It fails
closed: an event whose readers cannot be resolved is dropped.

## Prerequisites

Relaying events alone is **not** active/active support. Each item below must
be done before `ha.mode: active_active` ships.

| Process-local state | Location | Decision |
|---|---|---|
| WebSocket tickets | `backend/internal/ws/ws.go` (`IssueTicket`, `RedeemTicket`) | Move to a `ws_tickets` table that stores the ticket hash, the full identity (user, role, session, API key and its connector restriction) and a 30s expiry. Redeem once with `DELETE … RETURNING`, and sweep expired rows. A ticket issued on one replica must be redeemable on another. |
| User role/disabled cache (30s TTL) | `backend/internal/store/user.go` | Publish invalidations on the same NOTIFY channel so a disable or role change applies to every replica immediately, not up to 30s later. |
| HTTP-triggered sync de-duplication | `backend/internal/sync/engine.go` (`inFlight`) | Use the existing database sync claim lease for manual triggers as well, not just the in-process `sync.Map`. |
| Doc export, template version and backup/retention job serialization | `backend/internal/docexport/export.go`, `backend/internal/api/templates/handlers.go`, `backend/internal/api/system/handlers.go` | Replace the in-process mutexes with a database lease or PostgreSQL advisory lock. |
| Auth rate limiter | `backend/internal/api/middleware/ratelimit.go` | Accepted as per-replica. The effective limit is N times the configured limit. Documented in DEPLOYMENT.md. |
| Dashboard and attention TTL caches (5s), chat vector cache, OIDC provider cache | `backend/internal/ttlcache`, `backend/internal/chat/vectorcache.go`, `backend/internal/api/auth/handlers.go` | Accepted as per-replica. Data may be up to 5s stale; the other caches only affect performance. |
| Per-connector WebSocket event filtering | `backend/internal/ws/ws.go` | Done (#421). |

## Consequences

- **Active/passive (default, supported):** unchanged. One ready replica owns
  every WebSocket connection. Clients reconnect to the new leader after
  failover and must refetch then.
- **Active/active (future, opt-in):** a client connected to any replica
  receives events produced on any replica, provided the relay is healthy.
  Delivery is at most once per replica and best-effort. There is no ordering
  across origins, and a relay outage drops remote events until
  `system.resync` triggers a refetch. Background jobs still run only on the
  leader.
- One extra PostgreSQL connection per replica for LISTEN, and a small
  `ws_relay` table used only for oversized events.
- Implementation must prove, with two server processes against one
  PostgreSQL database, that a client on replica B receives an event produced
  on replica A, that no event is looped or delivered twice, that a
  user-targeted event reaches no other user, and that background jobs run
  only once.
