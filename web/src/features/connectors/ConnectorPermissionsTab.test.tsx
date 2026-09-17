import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { ConnectorPermissionsTab } from './ConnectorPermissionsTab';

let isInstanceAdmin = true;
vi.mock('../../hooks/useRole', () => ({
  useIsInstanceAdmin: () => isInstanceAdmin,
}));

const { upsertMock, deleteMock } = vi.hoisted(() => ({
  upsertMock: vi.fn(),
  deleteMock: vi.fn(),
}));

const grants = [
  { id: 'g1', userId: 'u1', connectorId: 'c1', role: 'operator', createdAt: '', updatedAt: '' },
];

vi.mock('../../api/permissions', () => ({
  useGetConnectorPermissions: () => ({ data: grants, isLoading: false, isError: false, refetch: vi.fn() }),
  useUpsertConnectorPermission: () => ({ mutate: upsertMock, isPending: false }),
  useDeleteConnectorPermission: () => ({ mutate: deleteMock, isPending: false }),
}));

vi.mock('../../api/generated/users/users', () => ({
  useGetUsers: () => ({
    data: [
      { id: 'u1', username: 'alice', displayName: 'Alice' },
      { id: 'u2', username: 'bob', displayName: 'Bob' },
    ],
  }),
}));

function renderTab() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ConnectorPermissionsTab connectorId="c1" />
    </QueryClientProvider>
  );
}

describe('ConnectorPermissionsTab (#240 PR1)', () => {
  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it('renders nothing for a non-instance-admin', () => {
    isInstanceAdmin = false;
    const { container } = renderTab();
    expect(container).toBeEmptyDOMElement();
    isInstanceAdmin = true;
  });

  it('lists existing grants and offers only ungranted users to add', () => {
    renderTab();
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByRole('option', { name: 'Bob' })).toBeInTheDocument();
    expect(screen.queryByRole('option', { name: 'Alice' })).not.toBeInTheDocument();
  });

  it('adding a user calls upsert with the default viewer role', async () => {
    renderTab();
    fireEvent.change(screen.getByDisplayValue(/add a user/i), { target: { value: 'u2' } });
    fireEvent.click(screen.getByRole('button', { name: /add/i }));
    await waitFor(() => expect(upsertMock).toHaveBeenCalledWith({ userId: 'u2', role: 'viewer' }, expect.anything()));
  });

  it('changing a grant role calls upsert with the new role', () => {
    renderTab();
    fireEvent.change(screen.getByDisplayValue('Operator'), { target: { value: 'viewer' } });
    expect(upsertMock).toHaveBeenCalledWith({ userId: 'u1', role: 'viewer' });
  });

  it('revoking a grant calls delete', () => {
    renderTab();
    fireEvent.click(screen.getByRole('button', { name: /revoke/i }));
    expect(deleteMock).toHaveBeenCalledWith('u1');
  });
});
