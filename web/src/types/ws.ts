/**
 * WebSocket contract types — the typed mirror of docs/WS_CONTRACT.md.
 *
 * App code (a future `useWebSocket` + `WebSocketProvider`) and the mock emitter
 * both import from here, so the contract stays in one place. Enums are reused from
 * the orval-generated models where they already exist (Severity, ServiceStatus).
 */
import type { ServiceStatus, Severity } from '../api/model';

/** `domain.action`, see WS_CONTRACT.md §naming. */
export type WsEventType =
  | 'service.status'
  | 'sync.progress'
  | 'sync.complete'
  | 'change.detected'
  | 'alert.created'
  | 'alert.resolved'
  | 'quality.finding.created'
  | 'quality.findings.changed'
  | 'finding.created'
  | 'system.job_failed'
  | 'doc.generated'
  | 'doc.ai_suggestion'
  | 'doc.lock.acquired'
  | 'doc.lock.released'
  | 'doc.lock.expired'
  | 'system.health'
  | 'system.notice';

/** Uniform envelope wrapping every frame. */
export interface WsEnvelope<T = unknown> {
  type: WsEventType;
  ts: string; // ISO-8601
  payload: T;
  id?: string;
}

// ─── payloads (1:1 with WS_CONTRACT.md) ─────────────────────────────────────

export interface ServiceStatusPayload {
  serviceId: string;
  status: ServiceStatus;
  message?: string;
  lastChecked: string;
}

export type SyncPhase = 'queued' | 'fetching' | 'diffing' | 'generating' | 'done' | 'error';

export interface SyncProgressPayload {
  serviceId: string | null;
  jobId: string;
  phase: SyncPhase;
  percent: number;
  message?: string;
}

export interface SyncCompletePayload {
  serviceId: string | null;
  jobId: string;
  changesDetected: number;
  alertsRaised: number;
  durationMs: number;
}

export interface ChangeDetectedPayload {
  changeId: string;
  serviceId: string;
  changeType: string;
  severity: Severity;
  summary: string;
  willTriggerAi: boolean;
}

export interface AlertCreatedPayload {
  alertId: string;
  changeId?: string;
  serviceId: string;
  severity: Severity;
  title: string;
}

export interface AlertResolvedPayload {
  alertId: string;
  resolvedBy: string;
  resolution: 'resolved' | 'dismissed' | 'snoozed';
}

/** Per-user finding notification dispatch (Dispatcher.NotifyFindingCreated),
 * not the broadcast-to-everyone quality.finding.created above. alertId is
 * always "" — findings have no deep-link target yet. */
export interface FindingNotificationPayload {
  alertId: string;
  title: string;
  message: string;
}

/** Per-user system event dispatch (Dispatcher.NotifySystemEvent): a scheduled
 * job started failing, or recovered. alertId is always "". */
export type SystemJobFailedPayload = FindingNotificationPayload;

export interface QualityFindingCreatedPayload {
  findingId: string;
  connectorId: string;
  checkType: string;
  severity: Severity;
}

export interface QualityFindingsChangedPayload {
  connectorId: string;
}

export interface DocGeneratedPayload {
  docId: string;
  serviceId?: string;
  trigger: 'ai' | 'template' | 'manual';
  newVersion: number;
}

export interface DocAiSuggestionPayload {
  docId: string;
  requestId: string;
  status: 'streaming' | 'complete' | 'error';
  contentDelta?: string;
  fullContent?: string;
  /** Present on 'complete': which provider answered. */
  provider?: string;
  /** Present on 'complete': true if the primary provider failed over to a fallback. */
  fallbackUsed?: boolean;
  error?: string;
}

export interface DocLockAcquiredPayload {
  docId: string;
  userId: string;
  acquiredAt: string;
  expiresAt: string;
}

export interface DocLockReleasedPayload {
  docId: string;
  userId: string;
}

export interface DocLockExpiredPayload {
  docId: string;
  userId: string;
}

export interface SystemHealthPayload {
  status: 'ok' | 'degraded' | 'down';
  components: { name: string; status: 'ok' | 'degraded' | 'down'; detail?: string }[];
}

export interface SystemNoticePayload {
  level: 'info' | 'warning' | 'critical';
  message: string;
  action?: 'reauth' | 'reload' | 'none';
}

/** Maps each event type to its payload — lets consumers switch exhaustively. */
export interface WsEventMap {
  'service.status': ServiceStatusPayload;
  'sync.progress': SyncProgressPayload;
  'sync.complete': SyncCompletePayload;
  'change.detected': ChangeDetectedPayload;
  'alert.created': AlertCreatedPayload;
  'alert.resolved': AlertResolvedPayload;
  'quality.finding.created': QualityFindingCreatedPayload;
  'quality.findings.changed': QualityFindingsChangedPayload;
  'finding.created': FindingNotificationPayload;
  'system.job_failed': SystemJobFailedPayload;
  'doc.generated': DocGeneratedPayload;
  'doc.ai_suggestion': DocAiSuggestionPayload;
  'doc.lock.acquired': DocLockAcquiredPayload;
  'doc.lock.released': DocLockReleasedPayload;
  'doc.lock.expired': DocLockExpiredPayload;
  'system.health': SystemHealthPayload;
  'system.notice': SystemNoticePayload;
}

/** Discriminated union of all possible frames. */
export type WsEvent = {
  [K in WsEventType]: WsEnvelope<WsEventMap[K]> & { type: K };
}[WsEventType];
