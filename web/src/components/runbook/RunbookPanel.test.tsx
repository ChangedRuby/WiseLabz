import { afterAll, afterEach, beforeAll, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import '../../i18n';
import type { RunbookPage } from '../../api/model';
import { RunbookPanel } from './RunbookPanel';

vi.mock('../ui/Dialog', () => ({
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
vi.mock('../manager/ElevationConfirm', () => ({
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
        <button onClick={() => onConfirm('tok-123')}>confirm-elevation</button>
      </div>
    ) : null,
}));

const server = setupServer();

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

function renderPanel(response: RunbookPage) {
  server.use(http.get('/api/runbooks', () => HttpResponse.json(response)));
  return render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <RunbookPanel changeType="vm.created" />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('RunbookPanel', () => {
  it('renders nothing when no runbook matches the target', async () => {
    const { container } = renderPanel({ items: [], total: 0, page: 1, pageSize: 20 });

    // Let the query settle before asserting the still-empty DOM.
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(container).toBeEmptyDOMElement();
  });

  it('renders the title and body when a runbook is found', async () => {
    renderPanel({
      items: [
        {
          id: 'rb-1',
          title: 'Restart the hung agent',
          body: 'Step 1: SSH in.\nStep 2: restart the service.',
          targetType: 'change_type',
          targetValue: 'vm.created',
          steps: [],
          createdAt: '2026-01-01T00:00:00Z',
          updatedAt: '2026-01-01T00:00:00Z',
        },
      ],
      total: 1,
      page: 1,
      pageSize: 20,
    });

    expect(await screen.findByText('Restart the hung agent')).toBeInTheDocument();
    expect(screen.getByText('Step 1: SSH in.', { exact: false })).toBeInTheDocument();
  });
});

describe('RunbookPanel steps', () => {
  const stepsRunbook: RunbookPage = {
    items: [
      {
        id: 'rb-1',
        title: 'Restart the hung agent',
        body: 'Step 1: SSH in.',
        targetType: 'change_type',
        targetValue: 'vm.created',
        steps: [
          {
            id: 'step-1',
            position: 0,
            title: 'Restart the sync worker',
            connectorId: 'conn-1',
            connectorName: 'pve1',
            verb: 'restart',
            entityRef: '',
            canExecute: true,
            executeBlockedReason: '',
          },
          {
            id: 'step-2',
            position: 1,
            title: 'Stop the agent',
            connectorId: 'conn-2',
            connectorName: 'docker1',
            verb: 'stop',
            entityRef: '100',
            canExecute: false,
            executeBlockedReason: 'no_operator_grant',
          },
        ],
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
      },
    ],
    total: 1,
    page: 1,
    pageSize: 20,
  };

  it('renders steps with verb, connector and entityRef', async () => {
    renderPanel(stepsRunbook);
    expect(await screen.findByText('Restart the sync worker')).toBeInTheDocument();
    expect(screen.getByText('restart · pve1')).toBeInTheDocument();
    expect(screen.getByText('stop · docker1 · Entity: 100')).toBeInTheDocument();
  });

  it('disables execute with a visible reason when canExecute is false', async () => {
    renderPanel(stepsRunbook);
    await screen.findByText('Stop the agent');
    const buttons = screen.getAllByRole('button', { name: 'Execute' });
    expect(buttons[1]).toBeDisabled();
    expect(screen.getByText('You need operator access to docker1.')).toBeInTheDocument();
  });

  it('runs the dry-run preview then executes with the elevation token', async () => {
    let capturedToken: string | null = null;
    server.use(
      http.post('/api/runbooks/rb-1/steps/step-1/execute', ({ request }) => {
        const url = new URL(request.url);
        if (url.searchParams.get('dryRun') === 'true') {
          return HttpResponse.json({
            targetService: 'pve1',
            estimatedDowntimeSeconds: 20,
            dependentServices: [],
          });
        }
        capturedToken = request.headers.get('X-Elevation-Token');
        return HttpResponse.json({ status: 'ok' });
      })
    );

    renderPanel(stepsRunbook);
    const executeButtons = await screen.findAllByRole('button', { name: 'Execute' });
    fireEvent.click(executeButtons[0]);

    const previewDialog = await screen.findByRole('dialog', { name: 'Execute impact' });
    expect(previewDialog).toHaveTextContent('20 seconds');

    fireEvent.click(within(previewDialog).getByRole('button', { name: 'Execute' }));
    fireEvent.click(await screen.findByRole('button', { name: 'confirm-elevation' }));

    await waitFor(() => expect(capturedToken).toBe('tok-123'));
  });
});
