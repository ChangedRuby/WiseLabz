/**
 * Hand-written client for the #240 PR1 per-connector permissions endpoints.
 * TODO: fold into docs/openapi.yaml and regenerate via `npm run gen:api`
 * once the backend has landed and the spec is updated — written by hand for
 * now since this fork ran in parallel with the backend implementation and
 * touching the shared generated client/spec here would risk clobbering it.
 */
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { customInstance } from './axios-instance';

export type ConnectorRole = 'viewer' | 'operator';

export interface ConnectorGrant {
  id: string;
  userId: string;
  connectorId: string;
  role: ConnectorRole;
  createdAt: string;
  updatedAt: string;
}

const getConnectorPermissionsQueryKey = (connectorId: string) =>
  [`/connectors/${connectorId}/permissions`] as const;

export function useGetConnectorPermissions(connectorId: string, enabled = true) {
  return useQuery({
    queryKey: getConnectorPermissionsQueryKey(connectorId),
    queryFn: ({ signal }) =>
      customInstance<ConnectorGrant[]>({
        url: `/connectors/${connectorId}/permissions`,
        method: 'GET',
        signal,
      }),
    enabled: enabled && !!connectorId,
  });
}

export function useUpsertConnectorPermission(connectorId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: ConnectorRole }) =>
      customInstance<ConnectorGrant>({
        url: `/connectors/${connectorId}/permissions/${userId}`,
        method: 'PUT',
        data: { role },
      }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: getConnectorPermissionsQueryKey(connectorId) }),
  });
}

export function useDeleteConnectorPermission(connectorId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) =>
      customInstance<void>({ url: `/connectors/${connectorId}/permissions/${userId}`, method: 'DELETE' }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: getConnectorPermissionsQueryKey(connectorId) }),
  });
}

export type InstanceAdminRole = 'admin' | 'user';

/**
 * Raw calls for the `instanceAdminRole` field replacing the old flat `role`
 * on the users endpoints (#240 PR1). Bypasses the generated
 * postUsers/patchUsersUserId, which are still typed against the old
 * UserCreate/UserUpdate `role` field until the OpenAPI spec is updated.
 */
export const patchUserInstanceAdmin = (
  userId: string,
  body: { instanceAdminRole?: InstanceAdminRole; disabled?: boolean; canManageDashboardDefaults?: boolean }
) => customInstance<unknown>({ url: `/users/${userId}`, method: 'PATCH', data: body });

export const postUserInstanceAdmin = (body: {
  username: string;
  email?: string;
  instanceAdminRole: InstanceAdminRole;
  canManageDashboardDefaults?: boolean;
}) => customInstance<unknown>({ url: `/users`, method: 'POST', data: body });
