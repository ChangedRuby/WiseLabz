import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { RunbooksPage } from './RunbooksPage';

const { postRunbooks, putRunbooksRunbookId, deleteRunbooksRunbookId } = vi.hoisted(() => ({
  postRunbooks: vi.fn().mockResolvedValue({}),
  putRunbooksRunbookId: vi.fn().mockResolvedValue({}),
  deleteRunbooksRunbookId: vi.fn().mockResolvedValue({}),
}));

let runbooks: unknown[] = [];

vi.mock('../../api/generated/runbooks/runbooks', () => ({
  useGetRunbooks: () => ({ data: { items: runbooks }, isLoading: false, isError: false, refetch: vi.fn() }),
  getGetRunbooksQueryKey: () => ['runbooks'],
  postRunbooks,
  putRunbooksRunbookId,
  deleteRunbooksRunbookId,
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

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectors: () => ({
    data: [{ id: 'conn-1', name: 'pve1', type: 'proxmox' }],
  }),
  useGetConnectorsSchema: () => ({
    data: [{ type: 'proxmox', category: 'hypervisor', displayName: 'Proxmox', fields: [], lifecycleVerbs: ['restart', 'stop'] }],
  }),
  useGetConnectorsConnectorIdSnapshots: () => ({ data: [{ id: 'snap-1' }] }),
  useGetConnectorsConnectorIdSnapshotsSnapshotId: () => ({
    data: { entities: [{ name: 'vm-100', externalId: '100' }] },
  }),
}));

function renderPage() {
  return render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <RunbooksPage />
    </QueryClientProvider>
  );
}

describe('RunbooksPage steps editor', () => {
  afterEach(() => {
    runbooks = [];
    vi.clearAllMocks();
  });

  it('creates a runbook with a step', async () => {
    renderPage();

    fireEvent.click(screen.getAllByRole('button', { name: 'New runbook' })[0]);
    fireEvent.change(screen.getByLabelText('Title'), { target: { value: 'Restart worker' } });
    fireEvent.change(screen.getByLabelText('Target', { exact: false }), {
      target: { value: 'vm.created' },
    });

    fireEvent.click(screen.getByRole('button', { name: 'Add step' }));
    fireEvent.change(screen.getByPlaceholderText('e.g. Restart the sync worker'), {
      target: { value: 'Restart the worker' },
    });

    fireEvent.change(screen.getByLabelText('Connector'), { target: { value: 'conn-1' } });
    fireEvent.change(screen.getByLabelText('Action'), { target: { value: 'stop' } });
    fireEvent.change(screen.getByLabelText('Entity'), { target: { value: '100' } });

    fireEvent.click(screen.getByRole('button', { name: 'Save' }));

    await waitFor(() =>
      expect(postRunbooks).toHaveBeenCalledWith(
        expect.objectContaining({
          steps: [
            expect.objectContaining({
              title: 'Restart the worker',
              connectorId: 'conn-1',
              verb: 'stop',
              entityRef: '100',
            }),
          ],
        })
      )
    );
  });

  it('surfaces a field error for an invalid step', async () => {
    postRunbooks.mockRejectedValueOnce({
      isAxiosError: true,
      response: { status: 400, data: { code: 'invalid_request', message: 'bad', details: [{ field: 'steps[0].verb', msg: 'unsupported verb' }] } },
    });

    renderPage();
    fireEvent.click(screen.getAllByRole('button', { name: 'New runbook' })[0]);
    fireEvent.change(screen.getByLabelText('Title'), { target: { value: 'Restart worker' } });
    fireEvent.change(screen.getByLabelText('Target', { exact: false }), {
      target: { value: 'vm.created' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Add step' }));
    fireEvent.change(screen.getByLabelText('Connector'), { target: { value: 'conn-1' } });

    fireEvent.click(screen.getByRole('button', { name: 'Save' }));

    expect(await screen.findByText('unsupported verb')).toBeInTheDocument();
  });
});
