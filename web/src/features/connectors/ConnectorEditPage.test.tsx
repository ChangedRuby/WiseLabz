import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { ConnectorEditPage } from './ConnectorEditPage';

const { putConnectorsConnectorId, testMock } = vi.hoisted(() => ({
  putConnectorsConnectorId: vi.fn().mockResolvedValue({}),
  testMock: vi.fn(),
}));

let connectorData: Record<string, unknown> = {
  id: 'c1',
  name: 'pve1',
  type: 'proxmox',
  url: 'https://pve1.example',
  verifyTls: true,
  enabled: true,
  status: 'online',
  secretRotatedAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
  userExpiresAt: '',
  rotationMaxAgeDays: null,
};

let schemas: Array<Record<string, unknown>> = [
  { type: 'proxmox', category: 'virtualization', displayName: 'Proxmox', fields: [], isCredentialRefresher: false },
];

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectorsConnectorId: () => ({ data: connectorData, isLoading: false, isError: false, refetch: vi.fn() }),
  useGetConnectorsSchema: () => ({ data: schemas }),
  putConnectorsConnectorId: (...args: unknown[]) => putConnectorsConnectorId(...args),
  postConnectorsConnectorIdTest: (...args: unknown[]) => testMock(...args),
  getGetConnectorsQueryKey: () => [],
}));

vi.mock('../../hooks/useRole', () => ({ useConnectorRole: () => 'operator' }));

vi.mock('./ConnectorPermissionsTab', () => ({ ConnectorPermissionsTab: () => null }));

function renderPage() {
  return render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter initialEntries={['/connectors/c1/edit']}>
        <Routes>
          <Route path="/connectors/:id/edit" element={<ConnectorEditPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('ConnectorEditPage rotation fields (#239 PR1)', () => {
  it('shows the "last rotated" text and rotation inputs for a non-refresher type', () => {
    renderPage();
    expect(screen.getByText(/secret last rotated 5 days ago/i)).toBeInTheDocument();
    expect(screen.getByText(/credential expires on/i)).toBeInTheDocument();
    expect(screen.getByText(/rotation reminder/i)).toBeInTheDocument();
  });

  it('hides rotation fields for a refresher connector type', () => {
    schemas = [
      { type: 'proxmox', category: 'virtualization', displayName: 'Proxmox', fields: [], isCredentialRefresher: true },
    ];
    renderPage();
    expect(screen.queryByText(/last rotated/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/credential expires on/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/rotation reminder/i)).not.toBeInTheDocument();
    schemas = [
      { type: 'proxmox', category: 'virtualization', displayName: 'Proxmox', fields: [], isCredentialRefresher: false },
    ];
  });

  it('submits userExpiresAt and rotationMaxAgeDays when set', async () => {
    renderPage();
    fireEvent.change(screen.getByLabelText(/credential expires on/i), { target: { value: '2027-01-01' } });
    fireEvent.change(screen.getByLabelText(/rotation reminder/i), { target: { value: '30' } });
    fireEvent.click(screen.getByRole('button', { name: /save/i }));
    await waitFor(() =>
      expect(putConnectorsConnectorId).toHaveBeenCalledWith(
        'c1',
        expect.objectContaining({ userExpiresAt: '2027-01-01T00:00:00Z', rotationMaxAgeDays: 30 }),
      ),
    );
  });
});
