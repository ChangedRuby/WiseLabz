# WiseLabz WebSocket Contract (`/ws`)

The single source of truth for the live-update channel. The Go backend implements
this; the frontend consumes it (and mocks it via a local emitter before the backend
exists). Keep it in sync with `docs/openapi.yaml`.

**Scope:** live updates only. **No mutations over WS** (per `ARCHITECTURE.md`). Every
state change still goes through REST (`docs/openapi.yaml`); WS only pushes notifications.

---

## Transport

- Endpoint: `GET /api/ws` (same origin; Vite proxies in dev).
- Protocol: WebSocket, text frames, one JSON object per frame.
- Auth: a one-time ticket, not the access token. The client first calls
  `POST /api/ws/ticket` with its normal authenticated session and receives
  `{ "ticket": "..." }`, then connects to `/api/ws?ticket=<ticket>`. A ticket is
  valid for 30 seconds, works once, and is held in the memory of the process
  that issued it (so it must be redeemed on the same replica; see
  [ADR 0005](adr/0005-cross-replica-websocket-relay.md)). A missing, expired, or
  reused ticket fails the upgrade with HTTP 401. Restricted API keys are refused
  a ticket because the stream is not filtered per connector.
- Open connections are re-validated on every ping (about every 25 seconds). If
  the user is disabled, their role changed, or the bound refresh session ended,
  the server closes with code `1008` (policy violation).
- Direction: server → client only for the events below. The client sends nothing
  except WS-level pong frames.

## Envelope

Every frame uses this envelope:

```ts
interface WsEnvelope<T = unknown> {
  id: string;     // UUID v4, unique per emitted event
  ts: string;     // server emit time, UTC, millisecond RFC 3339 (2026-09-06T12:00:00.000Z)
  type: string;   // "domain.action", see naming below
  payload: T;     // event-specific, typed below
}
```

`id` and `ts` are always present, including on the `system.health` heartbeat. The
`id` is unique per emitted event, and clients use it to drop duplicate deliveries
(`WebSocketProvider` keeps a bounded set of recently seen ids). For compatibility
with older servers the client keeps a local fallback: a frame missing `id` or `ts`
gets a locally generated value and is never treated as a duplicate.

## Naming convention

`domain.action`, lowercase, dot-separated. Domains: `service`, `sync`, `change`,
`alert`, `quality`, `doc`, `system`. Actions are past-tense/state nouns (`status`, `progress`,
`complete`, `detected`, `created`, `resolved`, `generated`, `ai_suggestion`,
`health`, `notice`, `resync`).

## Client dispatch model

`WebSocketProvider` (`web/src/ws/WebSocketProvider.tsx`) owns the socket, tracks
its status in the `useLive` store, parses each frame, and dispatches it through its
`type → handler` map. Handlers
do exactly one of: write a Zustand store, push a toast, `setQueryData`, or
`invalidateQueries`. Query keys below are the frontend's React Query keys.

---

## Events

### 1. `service.status`
A connector's reachability/health changed.

```ts
interface ServiceStatusPayload {
  serviceId: string;
  status: 'online' | 'degraded' | 'offline' | 'unknown';
  message?: string;
  lastChecked: string; // ISO
}
```
- **Consumers:** `ServiceStatusWidget`, `ServicesListPage`, `ServiceDetailPage`.
- **Reaction:** `setQueryData(['connectors'])` + `['connector', id]` to patch status;
  `StatusPill` re-renders. Toast **only** on transition into `offline`/`degraded`.

### 2. `sync.progress`
Progress of an in-flight sync job.

```ts
interface SyncProgressPayload {
  serviceId: string | null; // null = global sync
  jobId: string;
  phase: 'queued' | 'fetching' | 'diffing' | 'generating' | 'done' | 'error';
  percent: number;          // 0–100
  message?: string;
}
```
- **Consumers:** `SyncActivityWidget`, `SyncStatusBanner`, `FirstSyncProgress`, `ServiceDetailPage`.
- **Reaction:** update `syncStore[jobId]`. On `phase: 'error'` → inline error + toast.
  Throttle/coalesce rapid frames client-side.

### 3. `sync.complete`
Terminal sync result (separate from `sync.progress` so consumers get a clean summary).

```ts
interface SyncCompletePayload {
  serviceId: string | null;
  jobId: string;
  changesDetected: number;
  alertsRaised: number;
  durationMs: number;
}
```
- **Consumers:** `WebSocketProvider` (orchestration), `SyncActivityWidget`.
- **Reaction:** clear `syncStore[jobId]`; invalidate `['connector', id, 'data']`,
  `['dashboard','overview']`, `['changes']`; toast `"Sync complete — N changes"`.

### 4. `change.detected`
Diff engine found an infrastructure change.

```ts
interface ChangeDetectedPayload {
  changeId: string;
  serviceId: string;
  changeType: string;       // e.g. 'vm.created', 'firewall.rule.modified'
  severity: 'info' | 'warning' | 'critical';
  summary: string;
  willTriggerAi: boolean;   // routed to AI update vs alert
}
```
- **Consumers:** `RecentChangesWidget`, `ChangesFeedPage`, `ServiceDetailPage` (history).
- **Reaction:** invalidate `['changes']` + `['dashboard','overview']`; toast for
  `warning`/`critical`.

### 5. `alert.created`
Diff engine raised an alert needing action.

```ts
interface AlertCreatedPayload {
  alertId: string;
  changeId?: string;
  serviceId: string;
  severity: 'info' | 'warning' | 'critical';
  title: string;
}
```
- **Consumers:** `AlertBell`, sidebar `AlertBadge`, `alertsStore`, `AlertCenterPage`,
  `AlertSummaryWidget`.
- **Reaction:** `alertsStore.increment` + prepend; invalidate `['alerts']`,
  `['dashboard','overview']`; toast by severity.

### 6. `alert.resolved`
An alert was resolved — possibly by **another** user/session (multi-user sync).

```ts
interface AlertResolvedPayload {
  alertId: string;
  resolvedBy: string;       // user id
  resolution: 'resolved' | 'dismissed' | 'snoozed';
}
```
- **Consumers:** `alertsStore`, `AlertCenterPage`, `AlertDetailPage`.
- **Reaction:** `alertsStore.decrement`/remove; invalidate `['alerts']`,
  `['alert', id]`, `['dashboard','overview']`.

### 7. `quality.finding.created` and `quality.findings.changed`

Quality checks emit `quality.finding.created` once when a new open finding is
inserted. Every completed connector check (including the stale sweep) also emits
`quality.findings.changed` so clients refresh after repeated detections or
automatic resolutions.

```ts
interface QualityFindingCreatedPayload {
  findingId: string;
  connectorId: string;
  checkType: 'stale' | 'empty' | 'failing' | 'ownership_incomplete' | 'credential_rotation';
  severity: 'info' | 'warning' | 'critical';
}

interface QualityFindingsChangedPayload {
  connectorId: string;
}
```
- **Consumers:** `FindingsPage` and the Findings navigation badge.
- **Reaction:** invalidate the findings query.

### 7a. `finding.created`

Per-user notification dispatch (`Dispatcher.NotifyFindingCreated`), sent when a
finding is newly opened or its severity escalates — not the broadcast-to-everyone
`quality.finding.created` above, and not deduplicated the same way (re-detection
at the same severity never sends this). Generic across every check type,
including `credential_rotation`.

```ts
interface FindingNotificationPayload {
  alertId: string;  // always ""; client navigates to /findings via eventType instead
  title: string;
  message: string;
}
```
- **Consumers:** `NotificationCenter` (topbar bell).
- **Reaction:** invalidate `['notifications']`. When clicked, `NotificationCenter` navigates to `/findings` (event type `finding.*` is used as navigation target since `alertId` is empty).

### 7b. `digest.summary`

Per-user digest notification emitted by the hourly digest sweep (`Dispatcher.RunDigestSweep`)
when there are new notifications accumulated since the user's last digest. Only sent during
the user's configured local-time window (fixed at 8am in their `digestTimezone`) and when
eligible per their `digestCadence` (off/daily/weekly). Never rate-limited or deduplicated
within a window; one digest fires per eligible sweep per user (hourly cron, at most once
per cadence period).

```ts
interface DigestSummaryPayload {
  alertId: string;  // always ""; no deep-link target, same as finding.created
  title: string;    // e.g. "Digest: 5 new notification(s)"
  message: string;  // plain-text summary: count + up to 10 notification titles + "…and N more" if truncated
}
```
- **Consumers:** `NotificationCenter` (topbar bell).
- **Reaction:** invalidate `['notifications']`. When clicked, `NotificationCenter` has no `alertId`
  and eventType doesn't start with `finding.`, so it falls into the "no navigation" case (see §7a Reaction) —
  clicking just marks it read, no route change. Digest summaries are informational; users follow up by
  navigating to the Alerts or Findings pages themselves if needed.

### 7c. `system.job_failed`

Per-user system notification (`Dispatcher.NotifySystemEvent`), sent when a
scheduled background job (today: the doc export, see `docs/DOC_EXPORT.md`)
starts failing (severity `warning`) and once more when it recovers (severity
`info`). Only state transitions notify; repeated failures stay silent. Routable
per channel like `alert.created`, and included in digests.

```ts
interface SystemJobFailedPayload {
  alertId: string;  // always ""; no deep-link target
  title: string;    // e.g. "Doc export failing" / "Doc export recovered"
  message: string;  // the error, or a recovery note
}
```
- **Consumers:** `NotificationCenter` (topbar bell).
- **Reaction:** invalidate `['notifications']`. Clicking just marks it read (no navigation), same as §7b.

### 8. `doc.generated`
A doc node was (re)generated by the engine/AI/template.

```ts
interface DocGeneratedPayload {
  docId: string;
  serviceId?: string;
  trigger: 'ai' | 'template' | 'manual';
  newVersion: number;
}
```
- **Consumers:** `DocsPage`, `DocContent`, `DocHistoryPage`, `DocsHealthWidget`.
- **Reaction:** invalidate `['doc', docId]`, `['doc', docId, 'versions']`,
  `['docs','tree']`. If the user is editing that doc → non-destructive
  "newer version available" banner (last-write-wins, plan §8.3).

### 9. `doc.ai_suggestion`
Streamed AI suggestion for the open editor. **Only emitted when the AI module is
enabled.** Correlates to the `requestId` returned by `POST /docs/{docId}/ai-suggest`
or `POST /changes/{id}/ai-update`.

```ts
interface DocAiSuggestionPayload {
  docId: string;
  requestId: string;
  status: 'streaming' | 'complete' | 'error';
  contentDelta?: string;    // incremental text while streaming
  fullContent?: string;     // present on 'complete'
  provider?: string;        // present on 'complete': which provider answered
  fallbackUsed?: boolean;   // present on 'complete': true if the primary provider failed over
  error?: string;           // present on 'error'
}
```
- **Consumers:** `AiAssistPanel` (in `DocEditorPage`) only.
- **Reaction:** append `contentDelta` to the local AI staging buffer; on `complete`
  enable accept/insert; on `error` show inline error. Ignored if no editor open for
  that `docId`.

### 10. `doc.lock.acquired`
An advisory edit lock was acquired on a doc by a user (presence signalling).

```ts
interface DocLockAcquiredPayload {
  docId: string;
  userId: string;
  acquiredAt: string;  // ISO
  expiresAt: string;   // ISO
}
```
- **Consumers:** `DocEditorPage`, presence banner.
- **Reaction:** update `useLive.docLocks[docId]` with holder info; show/update presence
  banner (e.g. "X is editing this doc"). Lock expires after 5 minutes of inactivity.

### 11. `doc.lock.released`
An advisory edit lock was explicitly released by its holder.

```ts
interface DocLockReleasedPayload {
  docId: string;
  userId: string;
}
```
- **Consumers:** `DocEditorPage`, presence banner.
- **Reaction:** clear `useLive.docLocks[docId]`; hide presence banner.

### 12. `doc.lock.expired`
An advisory edit lock expired (no renewal/heartbeat for 5 minutes).

```ts
interface DocLockExpiredPayload {
  docId: string;
  userId: string;
}
```
- **Consumers:** `DocEditorPage`, presence banner.
- **Reaction:** clear `useLive.docLocks[docId]`; hide presence banner (lazy cleanup).

### 13. `system.health`
Backend/integration health + diff-engine heartbeat.

```ts
interface SystemHealthPayload {
  status: 'ok' | 'degraded' | 'down';
  components: { name: string; status: 'ok' | 'degraded' | 'down'; detail?: string }[];
}
```
- **Consumers:** `SystemSettings` `HealthPanel`, `WsStatusDot` (indirect).
- **Reaction:** `setQueryData(['health'])`; `degraded`/`down` → admin toast + dot color.

### 14. `system.notice`
Server-pushed broadcast (maintenance, forced logout, config change affecting sessions).

```ts
interface SystemNoticePayload {
  level: 'info' | 'warning' | 'critical';
  message: string;
  action?: 'reauth' | 'reload' | 'none';
}
```
- **Consumers:** global `Toaster`, `AppRoot`.
- **Reaction:** toast; `action: 'reauth'` → clear `authStore` + redirect `/login`;
  `'reload'` → prompt refresh. Covers an admin toggling auth/AI config mid-session.

### 15. `system.resync`
Broadcast by a replica after its cross-replica relay listener reconnects following a
gap (ADR 0005), meaning frames may have been missed. Not emitted yet; the relay adds it.

```ts
type SystemResyncPayload = {};
```
- **Consumers:** `WebSocketProvider`.
- **Reaction:** invalidate `['/alerts']`, `['/changes']`, `['/connectors']`,
  `['/dashboard/overview']` and clear the sync-job store (`useLive.jobs`).

---

## Reconnect behavior

`WebSocketProvider` auto-reconnects: on any close it fetches a new ticket and
reconnects with exponential backoff (`min(1000·2^n, 15000)` ms, reset after a
successful open). It tracks `connecting` → `open` → `closed` in the `useLive`
store. Close codes are not distinguished.

The server never replays frames, so recovery is by REST refetch. On every reconnect
(not the first connect after mount or login) the client invalidates the volatile
queries `['/alerts']`, `['/changes']`, `['/connectors']` and `['/dashboard/overview']`
(prefix match, so filtered variants are covered) and clears the sync-job store; a
running sync re-adds itself on its next `sync.progress`. The same happens on a
`system.resync` event (ADR 0005).

WS failure is never fatal: it degrades to a `closed` socket status in `useLive`, shows a
reconnecting indicator, and the UI falls back to React Query's normal
refetch-on-focus/interval for freshness.

---

## Mock emitter (frontend-first)

Until the backend serves `/ws`, a local mock emitter replays these events so every
consumer can be verified. The mock must:
- emit the exact envelope + payload shapes above,
- support a scripted timeline (e.g. trigger a connector sync → emit `sync.progress`
  ramp → `sync.complete` → `change.detected` → `alert.created`),
- be toggled by the same `USE_MOCKS` flag that gates MSW.

Lives at `web/src/mocks/ws/` (next scaffold step, alongside MSW).
