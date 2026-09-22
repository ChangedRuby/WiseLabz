import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { ServiceDetailPage } from './ServiceDetailPage';

const { restart, start, stop, health, configPush, configFields } = vi.hoisted(() => ({
  restart: vi.fn(),
  start: vi.fn(),
  stop: vi.fn(),
  health: vi.fn(),
  configPush: vi.fn(),
  configFields: vi.fn().mockReturnValue({ data: [] }),
}));

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectorsConnectorId: () => ({
    data: {
      id: 'svc-pve1',
      name: 'pve1',
      type: 'proxmox',
      url: 'https://pve1.example',
      enabled: true,
      status: 'online',
      scheduleSeconds: null,
      nextRunAt: '',
      lastSyncAt: '',
    },
    isLoading: false,
    isError: false,
    refetch: vi.fn(),
  }),
  useGetConnectorsConnectorIdData: () => ({
    data: { sections: [], fetchedAt: new Date().toISOString() },
  }),
  useGetConnectorsConnectorIdSyncs: () => ({ data: [] }),
  useGetConnectorsConnectorIdConfigFields: configFields,
  useGetConnectorsSchema: () => ({ data: [] }),
  postConnectorsConnectorIdRestart: restart,
  postConnectorsConnectorIdStart: start,
  postConnectorsConnectorIdStop: stop,
  postConnectorsConnectorIdConfigPush: configPush,
  postConnectorsConnectorIdHealth: health,
  putConnectorsConnectorIdEnabled: vi.fn(),
  putConnectorsConnectorId: vi.fn(),
  getGetConnectorsQueryKey: () => [],
}));

vi.mock('../../api/generated/changes/changes', () => ({
  useGetChanges: () => ({ data: { items: [] } }),
}));
vi.mock('../../api/generated/docs/docs', () => ({
  useGetDocsServiceConnectorId: () => ({ isError: true }),
  getGetDocsServiceConnectorIdQueryKey: () => [],
  postDocsGenerate: vi.fn(),
}));
vi.mock('../../api/generated/templates/templates', () => ({
  useGetTemplates: () => ({ data: [] }),
}));
vi.mock('../../hooks/useRole', () => ({ useConnectorRole: () => 'operator' }));
vi.mock('../../store/live', () => ({
  useLive: (selector: (state: object) => unknown) =>
    selector({ statusOverrides: {}, activity: [] }),
}));
vi.mock('../../lib/runSync', () => ({ runSync: vi.fn() }));
const { toastWarning, toastSuccess } = vi.hoisted(() => ({
  toastWarning: vi.fn(),
  toastSuccess: vi.fn(),
}));
vi.mock('../../lib/toast', () => ({
  toast: { warning: toastWarning, success: toastSuccess, error: vi.fn() },
}));
vi.mock('../../components/manager/ConfirmDestructive', () => ({ ConfirmDestructive: () => null }));
vi.mock('../../components/manager/ElevationConfirm', () => ({
  ElevationConfirm: ({
    open,
    title,
    onConfirm,
  }: {
    open: boolean;
    title: string;
    onConfirm: (token: string | null) => void;
  }) =>
    open ? (
      <div role="dialog" aria-label={title}>
        <button onClick={() => onConfirm(null)}>confirm-elevation</button>
      </div>
    ) : null,
}));
vi.mock('../../components/ui/Dialog', () => ({
  Dialog: ({
    open,
    title,
    children,
  }: {
    open: boolean;
    title: string;
    children: React.ReactNode;
  }) =>
    open ? (
      <div role="dialog" aria-label={title}>
        {children}
      </div>
    ) : null,
}));

function renderPage() {
  render(
    <QueryClientProvider
      client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}
    >
      <MemoryRouter initialEntries={['/services/svc-pve1']}>
        <Routes>
          <Route path="/services/:id" element={<ServiceDetailPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('ServiceDetailPage restart preview', () => {
  it('shows the read-only restart impact returned by the connector API', async () => {
    restart.mockResolvedValue({
      targetService: 'pve1',
      estimatedDowntimeSeconds: 30,
      dependentServices: [{ kind: 'vm', name: 'home-assistant' }],
    });

    renderPage();
    fireEvent.click(await screen.findByRole('button', { name: 'Restart' }));

    await waitFor(() =>
      expect(restart).toHaveBeenCalledWith('svc-pve1', undefined, { dryRun: true })
    );
    expect(await screen.findByRole('dialog', { name: 'Restart impact' })).toHaveTextContent(
      'Review the impact below, then confirm to restart.'
    );
    expect(screen.getByText('30 seconds')).toBeInTheDocument();
    expect(screen.getByText('home-assistant')).toBeInTheDocument();
  });
});

describe('ServiceDetailPage health check', () => {
  it('runs a health check and shows the result, disabling the button while pending', async () => {
    let resolveHealth: (value: { status: string; message: string; latencyMs: number }) => void =
      () => {};
    health.mockReturnValue(
      new Promise((resolve) => {
        resolveHealth = resolve;
      })
    );

    renderPage();
    const button = await screen.findByRole('button', { name: /Health check/ });
    fireEvent.click(button);

    await waitFor(() => expect(health).toHaveBeenCalledWith('svc-pve1'));
    expect(await screen.findByRole('button', { name: 'Checking…' })).toBeDisabled();

    resolveHealth({ status: 'degraded', message: 'Slow response', latencyMs: 820 });

    expect(await screen.findByRole('status')).toHaveTextContent('Slow response');
    expect(screen.getByRole('status')).toHaveTextContent('took 820ms');
  });

  it('shows an error message when the health check fails', async () => {
    health.mockRejectedValue(new Error('boom'));

    renderPage();
    fireEvent.click(await screen.findByRole('button', { name: /Health check/ }));

    expect(await screen.findByRole('status')).toHaveTextContent(
      "Couldn't run the health check. Try again."
    );
  });
});

describe('ServiceDetailPage start/stop preview', () => {
  it('shows the read-only start impact and estimatedDowntimeSeconds', async () => {
    start.mockResolvedValue({
      targetService: 'pve1',
      estimatedDowntimeSeconds: 15,
      dependentServices: [],
    });

    renderPage();
    fireEvent.click(await screen.findByRole('button', { name: 'Start' }));

    await waitFor(() => expect(start).toHaveBeenCalledWith('svc-pve1', undefined, { dryRun: true }));
    expect(await screen.findByRole('dialog', { name: 'Start impact' })).toHaveTextContent('15 seconds');
  });

  it('shows an indefinite downtime label for stop (estimatedDowntimeSeconds: 0)', async () => {
    stop.mockResolvedValue({
      targetService: 'pve1',
      estimatedDowntimeSeconds: 0,
      dependentServices: [],
    });

    renderPage();
    fireEvent.click(await screen.findByRole('button', { name: 'Stop' }));

    await waitFor(() => expect(stop).toHaveBeenCalledWith('svc-pve1', undefined, { dryRun: true }));
    expect(await screen.findByRole('dialog', { name: 'Stop impact' })).toHaveTextContent(
      'Indefinite (until started again)'
    );
  });
});

describe('ServiceDetailPage config push', () => {
  it('pushes a field and shows a success toast', async () => {
    configFields.mockReturnValue({
      data: [{ key: 'memory', label: 'Memory (MB)', type: 'number', entityScope: true }],
    });
    configPush.mockResolvedValue({});

    renderPage();
    fireEvent.change(await screen.findByDisplayValue('This connector has no writable config fields.'), {
      target: { value: 'memory' },
    });
    fireEvent.change(screen.getByPlaceholderText('entityRef'), { target: { value: '100' } });
    fireEvent.change(screen.getByPlaceholderText('New value'), { target: { value: '4096' } });
    fireEvent.click(screen.getByRole('button', { name: 'Push' }));
    fireEvent.click(await screen.findByRole('button', { name: 'confirm-elevation' }));

    await waitFor(() =>
      expect(configPush).toHaveBeenCalledWith(
        'svc-pve1',
        expect.objectContaining({ entityRef: '100', fieldKey: 'memory', value: 4096 }),
        undefined
      )
    );
    expect(toastSuccess).toHaveBeenCalled();
  });

  it('shows a mismatch toast on a 409 response', async () => {
    configFields.mockReturnValue({
      data: [{ key: 'memory', label: 'Memory (MB)', type: 'number', entityScope: true }],
    });
    configPush.mockRejectedValue({ isAxiosError: true, response: { status: 409 } });

    renderPage();
    fireEvent.change(await screen.findByDisplayValue('This connector has no writable config fields.'), {
      target: { value: 'memory' },
    });
    fireEvent.change(screen.getByPlaceholderText('entityRef'), { target: { value: '100' } });
    fireEvent.change(screen.getByPlaceholderText('New value'), { target: { value: '4096' } });
    fireEvent.click(screen.getByRole('button', { name: 'Push' }));
    fireEvent.click(await screen.findByRole('button', { name: 'confirm-elevation' }));

    await waitFor(() => expect(toastWarning).toHaveBeenCalled());
  });
});
