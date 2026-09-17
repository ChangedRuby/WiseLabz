/**
 * Per-connector permissions (#240 PR1): who has viewer/operator access to
 * this one connector. Instance-admin only — the server gates the underlying
 * endpoints the same way, this just keeps non-admins from seeing a section
 * they can't use.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useGetUsers } from '../../api/generated/users/users';
import {
  useGetConnectorPermissions,
  useUpsertConnectorPermission,
  useDeleteConnectorPermission,
  type ConnectorRole,
} from '../../api/permissions';
import { useIsInstanceAdmin } from '../../hooks/useRole';
import { Panel } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { toast } from '../../lib/toast';

export function ConnectorPermissionsTab({ connectorId }: { connectorId: string }) {
  const { t } = useTranslation();
  const isInstanceAdmin = useIsInstanceAdmin();
  const grants = useGetConnectorPermissions(connectorId, isInstanceAdmin);
  const users = useGetUsers({ query: { enabled: isInstanceAdmin } });
  const upsert = useUpsertConnectorPermission(connectorId);
  const remove = useDeleteConnectorPermission(connectorId);
  const [addingUserId, setAddingUserId] = useState('');

  if (!isInstanceAdmin) return null;

  const userById = new Map((users.data ?? []).map((u) => [u.id, u]));
  const grantedUserIds = new Set((grants.data ?? []).map((g) => g.userId));
  const addableUsers = (users.data ?? []).filter((u) => !grantedUserIds.has(u.id));

  const onAdd = () => {
    if (!addingUserId) return;
    upsert.mutate(
      { userId: addingUserId, role: 'viewer' },
      {
        onSuccess: () => setAddingUserId(''),
        onError: () => toast.error(t('connectors.permissions.error')),
      }
    );
  };

  return (
    <Panel className="mt-4 p-5">
      <h2 className="mb-1 text-sm font-semibold text-ink">{t('connectors.permissions.title')}</h2>
      <p className="mb-4 text-xs text-ink-muted">{t('connectors.permissions.subtitle')}</p>

      {grants.isLoading ? (
        <SkeletonRows rows={3} />
      ) : grants.isError ? (
        <ErrorState description={t('connectors.permissions.loadError')} onRetry={() => grants.refetch()} />
      ) : (grants.data ?? []).length === 0 ? (
        <EmptyState title={t('connectors.permissions.emptyTitle')} />
      ) : (
        <ul className="mb-4 divide-y divide-line-soft">
          {(grants.data ?? []).map((g) => {
            const user = userById.get(g.userId);
            return (
              <li key={g.id} className="flex items-center justify-between gap-3 py-2.5">
                <div className="min-w-0">
                  <p className="truncate text-sm text-ink">{user?.displayName || user?.username || g.userId}</p>
                  {user?.username && <p className="truncate text-2xs text-ink-faint">{user.username}</p>}
                </div>
                <div className="flex shrink-0 items-center gap-2">
                  <select
                    value={g.role}
                    onChange={(e) =>
                      upsert.mutate({ userId: g.userId, role: e.target.value as ConnectorRole })
                    }
                    className="rounded-md border border-line-soft bg-canvas px-2 py-1 text-xs text-ink"
                  >
                    <option value="viewer">{t('connectors.permissions.viewer')}</option>
                    <option value="operator">{t('connectors.permissions.operator')}</option>
                  </select>
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={remove.isPending}
                    onClick={() => remove.mutate(g.userId)}
                  >
                    {t('connectors.permissions.revoke')}
                  </Button>
                </div>
              </li>
            );
          })}
        </ul>
      )}

      <div className="flex items-center gap-2 border-t border-line-soft pt-3">
        <select
          value={addingUserId}
          onChange={(e) => setAddingUserId(e.target.value)}
          className="min-w-0 flex-1 rounded-md border border-line-soft bg-canvas px-2 py-1.5 text-xs text-ink"
        >
          <option value="">{t('connectors.permissions.addUserPlaceholder')}</option>
          {addableUsers.map((u) => (
            <option key={u.id} value={u.id}>
              {u.displayName || u.username}
            </option>
          ))}
        </select>
        <Button variant="secondary" size="sm" disabled={!addingUserId || upsert.isPending} onClick={onAdd}>
          {t('connectors.permissions.addUser')}
        </Button>
      </div>
    </Panel>
  );
}
