/**
 * Renders its children only for users who may mutate; everyone else sees
 * `fallback` (or nothing). With `connectorId`, gates on that connector's
 * per-user role (#240 PR1); without it, gates on the flat instance-admin
 * role. This is a UI affordance only — the server still enforces the real
 * boundary on every mutating endpoint (ARCHITECTURE.md).
 */
import type { ReactNode } from 'react';
import { useConnectorRole, useIsInstanceAdmin } from '../../hooks/useRole';

interface RoleGateProps {
  children: ReactNode;
  fallback?: ReactNode;
  /** Gate on this connector's role instead of the flat instance-admin role. */
  connectorId?: string;
  /** Minimum connector role required when `connectorId` is set. Default: operator. */
  minRole?: 'viewer' | 'operator';
}

export function RoleGate({ children, fallback, connectorId, minRole = 'operator' }: RoleGateProps) {
  const isInstanceAdmin = useIsInstanceAdmin();
  const connectorRole = useConnectorRole(connectorId);

  const allowed = connectorId
    ? connectorRole === 'operator' || (minRole === 'viewer' && connectorRole === 'viewer')
    : isInstanceAdmin;

  return <>{allowed ? children : (fallback ?? null)}</>;
}
