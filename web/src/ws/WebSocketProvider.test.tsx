import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { getGetAlertsQueryKey } from '../api/generated/alerts/alerts';
import { getGetChangesQueryKey } from '../api/generated/changes/changes';
import { getGetConnectorsQueryKey } from '../api/generated/connectors/connectors';
import { getGetDashboardOverviewQueryKey } from '../api/generated/dashboard/dashboard';
import i18n from '../i18n';
import { useAuth } from '../store/auth';
import { useLive } from '../store/live';
import { WebSocketProvider } from './WebSocketProvider';

vi.mock('../api/axios-instance', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../api/axios-instance')>()),
  customInstance: vi.fn().mockResolvedValue({ ticket: 'test-ticket' }),
}));

class TestWebSocket {
  static last: TestWebSocket | undefined;
  static urls: string[] = [];
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;

  constructor(url: string) {
    TestWebSocket.urls.push(url);
    TestWebSocket.last = this;
    queueMicrotask(() => this.onopen?.(new Event('open')));
  }

  close() {
    this.onclose?.(new Event('close'));
  }
}

describe('WebSocketProvider', () => {
  const originalWebSocket = window.WebSocket;

  afterEach(() => {
    vi.useRealTimers();
    window.WebSocket = originalWebSocket;
    useAuth.setState({ status: 'unknown', user: null });
    useLive.setState({ activity: [], jobs: {}, pendingAlerts: 0, statusOverrides: {} });
    TestWebSocket.last = undefined;
    TestWebSocket.urls = [];
  });

  it('invalidates the dashboard and changes caches after a completed sync', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const queryClient = new QueryClient();
    const invalidateQueries = vi.spyOn(queryClient, 'invalidateQueries');
    const translate = vi.spyOn(i18n, 't');

    render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    await waitFor(() => expect(TestWebSocket.last).toBeDefined());
    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'sync.complete',
            id: 'sync-1',
            ts: '2026-09-06T12:00:00Z',
            payload: {
              serviceId: 'svc-1',
              jobId: 'job-1',
              changesDetected: 0,
              alertsRaised: 0,
              durationMs: 1,
            },
          }),
        })
      )
    );

    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: getGetDashboardOverviewQueryKey() });
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: getGetChangesQueryKey() });
    expect(useLive.getState().activity[0]).toMatchObject({
      label: 'Sync complete',
      detail: '0 change(s) · 0 alert(s)',
    });
    expect(translate).toHaveBeenCalledWith('notify.syncCompleteTitle');
    expect(translate).toHaveBeenCalledWith('notify.syncCompleteDetail', { changes: 0, alerts: 0 });
  });

  it('stores localized activity copy for alert and documentation notices', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const translate = vi.spyOn(i18n, 't');
    render(
      <QueryClientProvider client={new QueryClient()}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );
    await waitFor(() => expect(TestWebSocket.last).toBeDefined());

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'alert.created',
            ts: '2026-09-06T12:00:00Z',
            payload: {
              alertId: 'alert-1',
              serviceId: 'svc-1',
              severity: 'warning',
              title: 'Disk space low',
            },
          }),
        })
      )
    );
    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'doc.generated',
            ts: '2026-09-06T12:00:00Z',
            payload: { docId: 'doc-1', trigger: 'template', newVersion: 3 },
          }),
        })
      )
    );

    expect(useLive.getState().activity).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ label: 'Disk space low', detail: 'New alert' }),
        expect.objectContaining({
          label: 'Documentation regenerated',
          detail: 'template · v3',
        }),
      ])
    );
    expect(translate).toHaveBeenCalledWith('notify.newAlert');
    expect(translate).toHaveBeenCalledWith('notify.documentationRegenerated');
    expect(translate).toHaveBeenCalledWith('notify.documentationRegeneratedDetail', {
      trigger: 'template',
      version: 3,
    });
  });

  const renderProvider = async (queryClient = new QueryClient()) => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );
    await waitFor(() => expect(TestWebSocket.last).toBeDefined());
  };

  const send = (frame: Record<string, unknown>) =>
    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', { data: JSON.stringify(frame) })
      )
    );

  const alertFrame = (extra: Record<string, unknown>, alertId: string) => ({
    type: 'alert.created',
    payload: { alertId, serviceId: 'svc-1', severity: 'warning', title: 'Disk space low' },
    ...extra,
  });

  it('handles a frame with a repeated id only once', async () => {
    await renderProvider();
    const frame = alertFrame({ id: 'evt-1', ts: '2026-09-06T12:00:00Z' }, 'alert-1');

    send(frame);
    send(frame);

    expect(useLive.getState().pendingAlerts).toBe(1);
    expect(useLive.getState().activity).toHaveLength(1);
  });

  it('still handles frames without an id, each one', async () => {
    await renderProvider();

    send(alertFrame({ ts: '2026-09-06T12:00:00Z' }, 'alert-1'));
    send(alertFrame({ ts: '2026-09-06T12:00:00Z' }, 'alert-2'));

    expect(useLive.getState().pendingAlerts).toBe(2);
  });

  it('fills a missing ts with a valid ISO timestamp', async () => {
    await renderProvider();

    send(alertFrame({ id: 'evt-2' }, 'alert-1'));

    const at = useLive.getState().activity[0]?.at;
    expect(at).toBeDefined();
    expect(new Date(at as string).toISOString()).toBe(at);
  });

  it('never puts the access token in the WebSocket URL', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    render(
      <QueryClientProvider client={new QueryClient()}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );
    await waitFor(() => expect(TestWebSocket.urls).toHaveLength(1));
    expect(TestWebSocket.urls[0]).toMatch(/\/api\/ws\?ticket=test-ticket$/);
  });

  it('updates docLocks store on doc lock events', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    render(
      <QueryClientProvider client={new QueryClient()}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );
    await waitFor(() => expect(TestWebSocket.last).toBeDefined());

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'doc.lock.acquired',
            ts: '2026-09-06T12:00:00Z',
            payload: {
              docId: 'doc-1',
              userId: 'user-1',
              acquiredAt: '2026-09-06T12:00:00Z',
              expiresAt: '2026-09-06T12:05:00Z',
            },
          }),
        })
      )
    );

    expect(useLive.getState().docLocks['doc-1']).toEqual({
      userId: 'user-1',
      expiresAt: '2026-09-06T12:05:00Z',
    });

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'doc.lock.released',
            ts: '2026-09-06T12:01:00Z',
            payload: { docId: 'doc-1', userId: 'user-1' },
          }),
        })
      )
    );

    expect(useLive.getState().docLocks['doc-1']).toBeUndefined();

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'doc.lock.acquired',
            ts: '2026-09-06T12:02:00Z',
            payload: {
              docId: 'doc-2',
              userId: 'user-2',
              acquiredAt: '2026-09-06T12:02:00Z',
              expiresAt: '2026-09-06T12:07:00Z',
            },
          }),
        })
      )
    );

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({
            type: 'doc.lock.expired',
            ts: '2026-09-06T12:07:30Z',
            payload: { docId: 'doc-2', userId: 'user-2' },
          }),
        })
      )
    );

    expect(useLive.getState().docLocks['doc-2']).toBeUndefined();
  });

  it('gracefully drops malformed messages without crashing or mutating state', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const queryClient = new QueryClient();
    const invalidateQueries = vi.spyOn(queryClient, 'invalidateQueries');

    render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    await waitFor(() => expect(TestWebSocket.last).toBeDefined());

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', { data: 'not valid json{{{' })
      )
    );

    expect(useLive.getState().activity).toEqual([]);
    expect(invalidateQueries).not.toHaveBeenCalled();
  });

  it('reconnects with exponential backoff after socket closes', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });

    const { unmount } = render(
      <QueryClientProvider client={new QueryClient()}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    // Wait for initial connection with real timers
    vi.useRealTimers();
    await waitFor(() => expect(TestWebSocket.last).toBeDefined());
    vi.useFakeTimers({ shouldAdvanceTime: true });

    const initialUrlCount = TestWebSocket.urls.length;

    act(() => TestWebSocket.last?.onclose?.(new Event('close')));

    expect(useLive.getState().ws).toBe('closed');
    expect(TestWebSocket.urls.length).toBe(initialUrlCount);

    await vi.advanceTimersByTimeAsync(1000);

    expect(TestWebSocket.urls.length).toBe(initialUrlCount + 1);

    unmount();
  });

  const volatileKeys = () => [
    getGetAlertsQueryKey(),
    getGetChangesQueryKey(),
    getGetConnectorsQueryKey(),
    getGetDashboardOverviewQueryKey(),
  ];

  const invalidatedKeys = (spy: { mock: { calls: unknown[][] } }) =>
    spy.mock.calls.map((c) => (c[0] as { queryKey: unknown }).queryKey);

  it('does not resync on the first open', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const queryClient = new QueryClient();
    const invalidateQueries = vi.spyOn(queryClient, 'invalidateQueries');

    render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    await waitFor(() => expect(useLive.getState().ws).toBe('open'));
    expect(invalidateQueries).not.toHaveBeenCalled();
  });

  it('refetches volatile queries and clears sync jobs after a reconnect', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const queryClient = new QueryClient();
    const invalidateQueries = vi.spyOn(queryClient, 'invalidateQueries');

    const { unmount } = render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    await vi.advanceTimersByTimeAsync(0);
    expect(useLive.getState().ws).toBe('open');
    expect(invalidateQueries).not.toHaveBeenCalled();

    useLive.getState().upsertJob({
      jobId: 'job-1',
      serviceId: 'svc-1',
      phase: 'fetching',
      percent: 10,
      startedAt: 0,
    });

    act(() => TestWebSocket.last?.onclose?.(new Event('close')));
    await vi.advanceTimersByTimeAsync(1000);

    expect(useLive.getState().ws).toBe('open');
    expect(invalidatedKeys(invalidateQueries)).toEqual(volatileKeys());
    expect(useLive.getState().jobs).toEqual({});

    unmount();
  });

  it('refetches volatile queries and clears sync jobs on system.resync', async () => {
    window.WebSocket = TestWebSocket as unknown as typeof WebSocket;
    useAuth.setState({ status: 'authenticated' });
    const queryClient = new QueryClient();
    const invalidateQueries = vi.spyOn(queryClient, 'invalidateQueries');

    render(
      <QueryClientProvider client={queryClient}>
        <WebSocketProvider>
          <div />
        </WebSocketProvider>
      </QueryClientProvider>
    );

    await waitFor(() => expect(useLive.getState().ws).toBe('open'));
    useLive.getState().upsertJob({
      jobId: 'job-1',
      serviceId: null,
      phase: 'fetching',
      percent: 10,
      startedAt: 0,
    });

    act(() =>
      TestWebSocket.last?.onmessage?.(
        new MessageEvent('message', {
          data: JSON.stringify({ type: 'system.resync', ts: '2026-09-06T12:00:00Z', payload: {} }),
        })
      )
    );

    expect(invalidatedKeys(invalidateQueries)).toEqual(volatileKeys());
    expect(useLive.getState().jobs).toEqual({});
  });
});
