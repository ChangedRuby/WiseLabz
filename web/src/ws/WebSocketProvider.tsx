/**
 * Single WS connection for the app. In dev it transparently hits MockWebSocket
 * (installed by the mock layer), which replays the WS_CONTRACT timeline. Frames
 * update the live store (sync progress, status, alerts, activity feed) and
 * invalidate the relevant React Query caches so REST data stays fresh.
 */
import { useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { getGetDashboardOverviewQueryKey } from '../api/generated/dashboard/dashboard';
import { getGetChangesQueryKey } from '../api/generated/changes/changes';
import { getGetAlertsQueryKey } from '../api/generated/alerts/alerts';
import { getGetFindingsQueryKey } from '../api/generated/findings/findings';
import { getGetNotificationsQueryKey } from '../api/generated/notifications/notifications';
import { getGetDocsTreeQueryKey } from '../api/generated/docs/docs';
import { useLive } from '../store/live';
import { toast } from '../lib/toast';
import { navigateTo } from '../lib/navigation';
import i18n from '../i18n';
import { useAuth } from '../store/auth';
import { customInstance } from '../api/axios-instance';
import type { WsEvent } from '../types/ws';

const jump = (to: string) => ({
  label: i18n.t('notify.view'),
  onClick: () => navigateTo(to),
});

const wsUrl = (ticket: string) => {
  return `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api/ws?ticket=${encodeURIComponent(ticket)}`;
};

// The socket is authorized by a one-time, short-lived ticket minted with the
// normal access token, so the refresh cookie is never the WS credential.
const fetchTicket = () =>
  customInstance<{ ticket: string }>({ url: '/ws/ticket', method: 'POST' }).then((r) => r.ticket);

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const qc = useQueryClient();
  const live = useRef(useLive.getState()).current; // stable store handle
  const reconnects = useRef(0);
  const status = useAuth((s) => s.status);

  useEffect(() => {
    if (status !== 'authenticated') {
      useLive.getState().setWs('closed');
      return;
    }

    let socket: WebSocket | null = null;
    let retry: ReturnType<typeof setTimeout> | undefined;
    let closed = false;

    const connect = async () => {
      useLive.getState().setWs('connecting');
      let ticket: string;
      try {
        ticket = await fetchTicket();
      } catch {
        if (closed) return;
        const delay = Math.min(1000 * 2 ** reconnects.current++, 15000);
        retry = setTimeout(connect, delay);
        return;
      }
      if (closed) return;
      socket = new WebSocket(wsUrl(ticket));

      socket.onopen = () => {
        reconnects.current = 0;
        useLive.getState().setWs('open');
      };

      socket.onclose = () => {
        useLive.getState().setWs('closed');
        if (closed) return;
        const delay = Math.min(1000 * 2 ** reconnects.current++, 15000);
        retry = setTimeout(connect, delay);
      };

      socket.onmessage = (ev) => {
        let frame: WsEvent;
        try {
          frame = JSON.parse(ev.data as string) as WsEvent;
        } catch {
          return;
        }
        handle(frame, qc);
      };
    };

    void connect();
    return () => {
      closed = true;
      if (retry) clearTimeout(retry);
      socket?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status]);

  // touch `live` so the ref isn't flagged unused under noUnusedLocals
  void live;
  return <>{children}</>;
}

function handle(frame: WsEvent, qc: ReturnType<typeof useQueryClient>) {
  const s = useLive.getState();
  switch (frame.type) {
    case 'service.status': {
      const p = frame.payload;
      s.setStatus(p.serviceId, p.status);
      qc.invalidateQueries({ queryKey: getGetDashboardOverviewQueryKey() });
      break;
    }
    case 'sync.progress': {
      const p = frame.payload;
      s.upsertJob({
        jobId: p.jobId,
        serviceId: p.serviceId,
        phase: p.phase,
        percent: p.percent,
        message: p.message,
        startedAt: Date.now(),
      });
      break;
    }
    case 'sync.complete': {
      const p = frame.payload;
      s.clearJob(p.serviceId ?? 'global');
      s.pushActivity({
        id: frame.id ?? crypto.randomUUID(),
        kind: 'sync',
        label: i18n.t('notify.syncCompleteTitle'),
        detail: i18n.t('notify.syncCompleteDetail', {
          changes: p.changesDetected,
          alerts: p.alertsRaised,
        }),
        at: frame.ts,
        tone: p.alertsRaised > 0 ? 'warn' : 'ok',
      });
      qc.invalidateQueries({ queryKey: getGetDashboardOverviewQueryKey() });
      qc.invalidateQueries({ queryKey: getGetChangesQueryKey() });
      {
        const msg = i18n.t('notify.syncComplete', {
          changes: p.changesDetected,
          alerts: p.alertsRaised,
        });
        const opts = p.changesDetected > 0 ? { action: jump('/changes') } : undefined;
        if (p.alertsRaised > 0) toast.warning(msg, opts);
        else toast.success(msg, opts);
      }
      break;
    }
    case 'change.detected': {
      const p = frame.payload;
      s.pushActivity({
        id: p.changeId,
        kind: 'change',
        label: p.summary,
        detail: p.changeType,
        at: frame.ts,
        tone: p.severity === 'critical' ? 'err' : p.severity === 'warning' ? 'warn' : 'signal',
      });
      qc.invalidateQueries({ queryKey: getGetChangesQueryKey() });
      qc.invalidateQueries({ queryKey: getGetDashboardOverviewQueryKey() });
      {
        const opts = { action: jump(`/changes/${p.changeId}`) };
        if (p.severity === 'critical') toast.error(p.summary, opts);
        else if (p.severity === 'warning') toast.warning(p.summary, opts);
        else toast.info(p.summary, opts);
      }
      break;
    }
    case 'alert.created': {
      const p = frame.payload;
      s.bumpAlerts(1);
      s.pushActivity({
        id: p.alertId,
        kind: 'alert',
        label: p.title,
        detail: i18n.t('notify.newAlert'),
        at: frame.ts,
        tone: p.severity === 'critical' ? 'err' : 'warn',
      });
      qc.invalidateQueries({ queryKey: getGetAlertsQueryKey() });
      qc.invalidateQueries({ queryKey: getGetNotificationsQueryKey() });
      {
        const opts = { action: jump('/alerts') };
        if (p.severity === 'critical') toast.error(p.title, opts);
        else toast.warning(p.title, opts);
      }
      break;
    }
    case 'alert.resolved': {
      s.bumpAlerts(-1);
      qc.invalidateQueries({ queryKey: getGetAlertsQueryKey() });
      qc.invalidateQueries({ queryKey: getGetNotificationsQueryKey() });
      break;
    }
    case 'quality.finding.created':
    case 'quality.findings.changed': {
      qc.invalidateQueries({ queryKey: getGetFindingsQueryKey() });
      break;
    }
    case 'finding.created':
    case 'system.job_failed': {
      // Per-user notification dispatch (Dispatcher.NotifyFindingCreated /
      // NotifySystemEvent), not the broadcast-to-everyone
      // 'quality.finding.created' above — refresh the notification bell the
      // same way 'alert.resolved' does.
      qc.invalidateQueries({ queryKey: getGetNotificationsQueryKey() });
      break;
    }
    case 'doc.generated': {
      s.pushActivity({
        id: frame.id ?? crypto.randomUUID(),
        kind: 'doc',
        label: i18n.t('notify.documentationRegenerated'),
        detail: i18n.t('notify.documentationRegeneratedDetail', {
          trigger: frame.payload.trigger,
          version: frame.payload.newVersion,
        }),
        at: frame.ts,
        tone: 'signal',
      });
      qc.invalidateQueries({ queryKey: getGetDocsTreeQueryKey() });
      break;
    }
    case 'doc.lock.acquired': {
      const p = frame.payload;
      s.setDocLock(p.docId, { userId: p.userId, expiresAt: p.expiresAt });
      break;
    }
    case 'doc.lock.released': {
      const p = frame.payload;
      s.setDocLock(p.docId, undefined);
      break;
    }
    case 'doc.lock.expired': {
      const p = frame.payload;
      s.setDocLock(p.docId, undefined);
      break;
    }
    default:
      break;
  }
}
