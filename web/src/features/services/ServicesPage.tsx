/** Services list — every connector with live status, type, endpoint, last sync,
 *  and the operator manager actions: sync, enable/disable, and remove. Mutating
 *  controls are hidden for viewers (server still enforces the boundary). */
import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion } from 'motion/react';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { matchSorter } from 'match-sorter';
import {
  useGetConnectors,
  putConnectorsConnectorIdEnabled,
  getGetConnectorsQueryKey,
  useGetConnectorsMaintenanceWindows,
} from '../../api/generated/connectors/connectors';
import { useLive } from '../../store/live';
import { useCanMutate } from '../../hooks/useRole';
import { runSync } from '../../lib/runSync';
import { StatusPill } from '../../components/ui/StatusDot';
import { Button, IconButton } from '../../components/ui/Button';
import { SavedViewsMenu } from '../../components/views/SavedViewsMenu';
import { Panel } from '../../components/ui/Panel';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { ConfirmDestructive } from '../../components/manager/ConfirmDestructive';
import { MaintenanceWindowMenu } from '../../components/manager/MaintenanceWindowMenu';
import { relativeTime } from '../../lib/time';
import { SearchIcon, SyncIcon, PlusIcon, XIcon } from '../../components/icons';
import { categoryIcon } from '../../components/categoryIcon';
import type { Connector, ServiceStatus } from '../../api/model';

export function ServicesPage() {
  const { t } = useTranslation();
  const { data, isLoading, isError, refetch } = useGetConnectors();
  const { data: maintenanceWindows } = useGetConnectorsMaintenanceWindows();
  const activeMaintenanceIds = useMemo(
    () => new Set((maintenanceWindows ?? []).map((w) => w.connectorId)),
    [maintenanceWindows]
  );
  const overrides = useLive((s) => s.statusOverrides);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const canMutate = useCanMutate();
  const [q, setQ] = useState('');
  const [removing, setRemoving] = useState<Connector | null>(null);

  const toggleEnabled = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) =>
      putConnectorsConnectorIdEnabled(id, { enabled }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() }),
  });

  const rows = useMemo(() => {
    if (!data) return [];
    const withStatus = data.map((c) => ({ ...c, status: overrides[c.id] ?? c.status }) as Connector);
    const query = q.trim();
    return query ? matchSorter(withStatus, query, { keys: ['name', 'type'] }) : withStatus;
  }, [data, overrides, q]);

  return (
    <div className="mx-auto max-w-275 px-6 py-6">
      <header className="mb-5 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight text-ink">{t('services.title')}</h1>
          <p className="text-sm text-ink-muted">
            {data ? t('services.countConnectors', { count: data.length }) : t('services.subtitle')}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex h-9 items-center gap-2 rounded-md border border-line-soft bg-surface px-2.5">
            <SearchIcon size={15} className="text-ink-faint" />
            <label htmlFor="services-filter" className="sr-only">
              {t('services.filterPlaceholder')}
            </label>
            <input
              id="services-filter"
              name="services-filter"
              autoComplete="off"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder={t('services.filterPlaceholder')}
              className="w-36 bg-transparent text-sm text-ink placeholder:text-ink-faint sm:w-44"
            />
          </div>
          <SavedViewsMenu surface="services" filters={{ q }} onApply={(f) => setQ(f.q ?? '')} />
          {canMutate && (
            <Button variant="primary" size="md" onClick={() => navigate('/services/new')}>
              <PlusIcon size={15} /> {t('services.addConnector')}
            </Button>
          )}
        </div>
      </header>

      <Panel>
        {isLoading ? (
          <SkeletonRows rows={6} />
        ) : isError || !data ? (
          <ErrorState description={t('services.loadError')} onRetry={() => refetch()} />
        ) : rows.length === 0 ? (
          <EmptyState
            title={t('services.noMatchTitle')}
            description={t('services.noMatchDesc', { query: q })}
          />
        ) : (
          <div className="overflow-x-auto">
          <table className="w-full min-w-140 text-sm">
            <thead>
              <tr className="border-b border-line-soft text-left text-2xs text-ink-faint">
                <th className="px-4 py-2.5 font-semibold">{t('services.col.service')}</th>
                <th className="hidden px-4 py-2.5 font-semibold sm:table-cell">
                  {t('services.col.category')}
                </th>
                <th className="hidden px-4 py-2.5 font-semibold md:table-cell">
                  {t('services.col.endpoint')}
                </th>
                <th className="px-4 py-2.5 font-semibold">{t('services.col.status')}</th>
                <th className="px-4 py-2.5 text-right font-semibold">
                  {t('services.col.lastSync')}
                </th>
                <th className="px-4 py-2.5" />
              </tr>
            </thead>
            <tbody>
              {rows.map((c, idx) => {
                const Icon = categoryIcon[c.category];
                return (
                  <motion.tr
                    key={c.id}
                    initial={{ opacity: 0, y: 6 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.03, duration: 0.25 }}
                    className="group border-b border-line-soft transition-colors last:border-0 hover:bg-surface-raised"
                  >
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2.5">
                        <span className="flex h-8 w-8 items-center justify-center rounded-md bg-canvas-sunken text-ink-faint">
                          <Icon size={16} />
                        </span>
                        <div>
                          <p className="flex items-center gap-2 font-medium text-ink">
                            <button
                              onClick={() => navigate(`/services/${c.id}`)}
                              className="rounded-sm text-left transition-colors hover:text-accent-secondary focus-visible:outline focus-visible:outline-accent-primary-soft"
                            >
                              {c.name}
                            </button>
                            {!c.enabled && (
                              <span className="rounded bg-idle-tint px-1.5 py-0.5 text-2xs font-medium text-ink-faint">
                                {t('services.disabledTag')}
                              </span>
                            )}
                            {activeMaintenanceIds.has(c.id) && (
                              <span className="rounded bg-idle-tint px-1.5 py-0.5 text-2xs font-medium text-ink-faint">
                                {t('services.maintenance.badge')}
                              </span>
                            )}
                          </p>
                          <p className="font-mono text-2xs text-ink-faint">{c.type}</p>
                        </div>
                      </div>
                    </td>
                    <td className="hidden px-4 py-3 text-ink-muted sm:table-cell">
                      {t(`services.category.${c.category}`)}
                    </td>
                    <td className="hidden px-4 py-3 font-mono text-2xs text-ink-faint md:table-cell">
                      {c.url?.replace(/^https?:\/\//, '')}
                    </td>
                    <td className="px-4 py-3">
                      <StatusPill status={c.status as ServiceStatus} />
                      {c.statusMessage && (
                        <p className="mt-0.5 max-w-50 truncate text-2xs text-ink-faint">
                          {c.statusMessage}
                        </p>
                      )}
                    </td>
                    <td className="nums px-4 py-3 text-right font-mono text-2xs text-ink-faint">
                      {t('common.ago', { time: relativeTime(c.lastSyncAt) })}
                    </td>
                    <td className="px-4 py-3 text-right">
                      {canMutate && (
                        <div className="flex items-center justify-end gap-1 opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100 sm:focus-within:opacity-100">
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => runSync(c.id)}
                            disabled={!c.enabled}
                          >
                            <SyncIcon size={14} /> {t('common.sync')}
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => toggleEnabled.mutate({ id: c.id, enabled: !c.enabled })}
                            disabled={toggleEnabled.isPending}
                          >
                            {c.enabled ? t('common.disable') : t('common.enable')}
                          </Button>
                          <MaintenanceWindowMenu connectorId={c.id} active={activeMaintenanceIds.has(c.id)} />
                          <IconButton
                            label={t('services.removeLabel', { name: c.name })}
                            onClick={() => setRemoving(c)}
                            className="hover:bg-err-tint hover:text-err"
                          >
                            <XIcon size={15} />
                          </IconButton>
                        </div>
                      )}
                    </td>
                  </motion.tr>
                );
              })}
            </tbody>
          </table>
          </div>
        )}
      </Panel>

      <ConfirmDestructive
        open={!!removing}
        connectorId={removing?.id ?? ''}
        connectorName={removing?.name ?? ''}
        onClose={() => setRemoving(null)}
        onConfirmed={() => setRemoving(null)}
      />
    </div>
  );
}
