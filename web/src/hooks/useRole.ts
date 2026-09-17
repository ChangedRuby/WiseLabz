/**
 * Role helpers (#240 PR1). Access is now per-connector: a user's role on a
 * connector comes from `myRole` on the fetched connector, not a flat
 * instance-wide role. A separate `instanceAdminRole` on `/me` covers
 * non-connector actions (user management, API keys, granting permissions).
 * The UI uses these only to hide controls the user can't action — the server
 * enforces the real boundary on every mutating endpoint (ARCHITECTURE.md).
 */
import { useGetMe } from '../api/generated/me/me';
import { useGetConnectors } from '../api/generated/connectors/connectors';
import type { Connector, User } from '../api/model';

// TODO: fold into docs/openapi.yaml once the backend PR1 spec update lands —
// see src/api/permissions.ts for why these are hand-typed for now.
type UserWithInstanceAdmin = User & { instanceAdminRole?: 'admin' | 'user' };
type ConnectorWithRole = Connector & { myRole?: 'viewer' | 'operator' | '' };

/** True for the flat, non-connector-scoped instance-admin role. */
export function useIsInstanceAdmin(): boolean {
  const { data } = useGetMe();
  return (data as UserWithInstanceAdmin | undefined)?.instanceAdminRole === 'admin';
}

/** The current user's role on one connector, or undefined if no grant. */
export function useConnectorRole(connectorId: string | undefined): 'viewer' | 'operator' | undefined {
  const { data } = useGetConnectors();
  const connectors = (data as ConnectorWithRole[] | undefined) ?? [];
  const connector = connectors.find((c) => c.id === connectorId);
  return connector?.myRole || undefined;
}

/** Every connector id the current user holds at least operator on. */
export function useOperatorConnectorIds(): Set<string> {
  const { data } = useGetConnectors();
  const connectors = (data as ConnectorWithRole[] | undefined) ?? [];
  return new Set(connectors.filter((c) => c.myRole === 'operator').map((c) => c.id));
}

/**
 * True when the current user may perform *some* mutating action: an instance
 * admin, or an operator on at least one connector. Used by page chrome
 * (command palette, shortcuts, topbar, settings nav) that isn't scoped to a
 * specific connector — those surfaces show/hide broadly, the actual mutation
 * is still gated per-connector where it happens.
 */
export function useCanMutate(): boolean {
  const isInstanceAdmin = useIsInstanceAdmin();
  const operatorIds = useOperatorConnectorIds();
  return isInstanceAdmin || operatorIds.size > 0;
}
