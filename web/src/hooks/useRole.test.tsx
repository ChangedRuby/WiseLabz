import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { ReactNode } from 'react';
import { useCanMutate, useConnectorRole, useIsInstanceAdmin, useOperatorConnectorIds } from './useRole';

let meData: unknown = { id: 'u1', instanceAdminRole: 'user' };
let connectorsData: unknown = [
  { id: 'c1', myRole: 'operator' },
  { id: 'c2', myRole: 'viewer' },
  { id: 'c3', myRole: '' },
];

vi.mock('../api/generated/me/me', () => ({
  useGetMe: () => ({ data: meData }),
}));

vi.mock('../api/generated/connectors/connectors', () => ({
  useGetConnectors: () => ({ data: connectorsData }),
}));

function wrapper({ children }: { children: ReactNode }) {
  const client = new QueryClient();
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

describe('useRole hooks (#240 PR1)', () => {
  it('useIsInstanceAdmin reads instanceAdminRole off /me', () => {
    meData = { instanceAdminRole: 'admin' };
    expect(renderHook(() => useIsInstanceAdmin(), { wrapper }).result.current).toBe(true);

    meData = { instanceAdminRole: 'user' };
    expect(renderHook(() => useIsInstanceAdmin(), { wrapper }).result.current).toBe(false);
  });

  it('useConnectorRole returns the connector-specific role, or undefined without a grant', () => {
    connectorsData = [
      { id: 'c1', myRole: 'operator' },
      { id: 'c2', myRole: '' },
    ];
    expect(renderHook(() => useConnectorRole('c1'), { wrapper }).result.current).toBe('operator');
    expect(renderHook(() => useConnectorRole('c2'), { wrapper }).result.current).toBeUndefined();
    expect(renderHook(() => useConnectorRole('missing'), { wrapper }).result.current).toBeUndefined();
  });

  it('useOperatorConnectorIds collects only connectors with an operator grant', () => {
    connectorsData = [
      { id: 'c1', myRole: 'operator' },
      { id: 'c2', myRole: 'viewer' },
      { id: 'c3', myRole: 'operator' },
    ];
    const { result } = renderHook(() => useOperatorConnectorIds(), { wrapper });
    expect(result.current).toEqual(new Set(['c1', 'c3']));
  });

  it('useCanMutate is true for an instance admin with zero connector grants', () => {
    meData = { instanceAdminRole: 'admin' };
    connectorsData = [];
    expect(renderHook(() => useCanMutate(), { wrapper }).result.current).toBe(true);
  });

  it('useCanMutate is true for a non-admin with at least one operator grant', () => {
    meData = { instanceAdminRole: 'user' };
    connectorsData = [{ id: 'c1', myRole: 'operator' }];
    expect(renderHook(() => useCanMutate(), { wrapper }).result.current).toBe(true);
  });

  it('useCanMutate is false for a non-admin viewer with no operator grants', () => {
    meData = { instanceAdminRole: 'user' };
    connectorsData = [{ id: 'c1', myRole: 'viewer' }];
    expect(renderHook(() => useCanMutate(), { wrapper }).result.current).toBe(false);
  });
});
