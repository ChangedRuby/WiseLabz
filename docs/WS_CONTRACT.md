# WiseLabz WebSocket Contract (`/ws`)

The single source of truth for the live-update channel. The Go backend implements
this; the frontend consumes it (and mocks it via a local emitter before the backend
exists). Mirrors `docs/FRONTEND_PLAN.md` §4 — keep them in sync.

**Scope:** live updates only. **No mutations over WS** (per `ARCHITECTURE.md`). Every
state change still goes through REST (`docs/openapi.yaml`); WS only pushes notifications.

---

## Transport

- Endpoint: `GET /ws` (same origin; Vite proxies in dev).
- Protocol: WebSocket, text frames, one JSON object per frame.
- Auth: the connection is opened by an authenticated session. The access token is
  passed at connect time (query param `?access_token=` or `Sec-WebSocket-Protocol`
  — backend's choice; frontend `useWebSocket` adapts). The server closes with code
  `4401` if the token is missing/expired, prompting the client to refresh then
  reconnect.
- Direction: server → client only for the events below. The client sends nothing
  except WS-level pong frames.

## Envelope

Every frame uses this envelope:

```ts
interface WsEnvelope<T = unknown> {
  type: string;   // "domain.action", see naming below
  ts: string;     // ISO-8601 server timestamp
  payload: T;     // event-specific, typed below
  id?: string;    // optional unique event id, used for client-side dedupe
}
```

## Naming convention

`domain.action`, lowercase, dot-separated. Domains: `service`, `sync`, `change`,
`alert`, `quality`, `doc`, `system`. Actions are past-tense/state nouns (`status`, `progress`,
`complete`, `detected`, `created`, `resolved`, `generated`, `ai_suggestion`,
`health`, `notice`).

## Client dispatch model

`useWebSocket` parses the envelope, updates `wsStore` liveness on every frame, and
hands the event to `WebSocketProvider`, which owns the `type → handler` map. Handlers
do exactly one of: write a Zustand store, push a toast, `setQueryData`, or
`invalidateQueries`. Store/query names below reference `FRONTEND_PLAN.md` §3.

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

---

## Reconnect behavior

`useWebSocket` auto-reconnects with exponential backoff and updates `wsStore.status`
(`connecting` → `open` → `closed` → `reconnecting`). On a `4401` close it triggers a
silent token refresh (see `openapi.yaml` `/auth/refresh`) before reconnecting.

On every successful **reconnect**, `WebSocketProvider` invalidates the volatile
queries to recover anything missed while disconnected:
`['alerts']`, `['changes']`, `['connectors']`, `['dashboard','overview']`.

WS failure is never fatal: it degrades to `wsStore.status='closed'`, shows a
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
