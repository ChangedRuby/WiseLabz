/** Global change feed — every detected infra change, newest first. */
import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { motion } from 'motion/react';
import {
  useGetChanges,
  postChangesBulkResolve,
  getGetChangesQueryKey,
} from '../../api/generated/changes/changes';
import { SeverityTag } from '../../components/ui/StatusDot';
import { Panel } from '../../components/ui/Panel';
import { Pagination } from '../../components/ui/Pagination';
import { SavedViewsMenu } from '../../components/views/SavedViewsMenu';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { relativeTime } from '../../lib/time';
import { toast } from '../../lib/toast';
import { useOperatorConnectorIds } from '../../hooks/useRole';
import { ArrowRightIcon, CheckIcon, XIcon } from '../../components/icons';
import type { Severity } from '../../api/model';

const FILTERS: { key: string; value: Severity | 'all' }[] = [
  { key: 'changes.filterAll', value: 'all' },
  { key: 'changes.filterCritical', value: 'critical' },
  { key: 'changes.filterWarning', value: 'warning' },
  { key: 'changes.filterInfo', value: 'info' },
];

// Low-risk = server's bulk-resolve boundary (severity != critical). Kept in
// sync by hand with backend/internal/api/changes/handlers.go's BulkResolve —
// the server re-checks this independently either way, so a mismatch here is
// a UX nuisance, never a safety hole.
const isLowRisk = (severity: Severity) => severity !== 'critical';

export function ChangesPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const operatorIds = useOperatorConnectorIds();
  const [searchParams, setSearchParams] = useSearchParams();
  const pageSize = 20;
  const page = Math.max(1, Number(searchParams.get('page')) || 1);
  const filter = (searchParams.get('severity') as Severity | 'all' | null) ?? 'all';
  const { data, isLoading, isError, refetch } = useGetChanges({ page, pageSize });
  const [selected, setSelected] = useState<Set<string>>(new Set());

  const setFilter = (next: Severity | 'all') =>
    setSearchParams((prev) => {
      const params = new URLSearchParams(prev);
      if (next === 'all') params.delete('severity');
      else params.set('severity', next);
      params.delete('page');
      return params;
    });

  const setPage = (next: number) =>
    setSearchParams((prev) => {
      const params = new URLSearchParams(prev);
      if (next <= 1) params.delete('page');
      else params.set('page', String(next));
      return params;
    });

  const items = (data?.items ?? []).filter((c) => filter === 'all' || c.severity === filter);
  const pageCount = data ? Math.max(1, Math.ceil(data.total / data.pageSize)) : 1;

  const toggleSelected = (id: string) =>
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });

  const bulkResolve = useMutation({
    mutationFn: (status: 'acknowledged' | 'dismissed') =>
      postChangesBulkResolve({ ids: Array.from(selected), status }),
    onSuccess: (res) => {
      const succeeded = res.results.filter((r) => r.status === 'success').length;
      const failed = res.results.filter((r) => r.status === 'error');
      if (failed.length === 0) {
        toast.success(t('changes.bulkAllSucceeded', { count: succeeded }));
      } else {
        toast.warning(
          t('changes.bulkPartial', { succeeded, failed: failed.length, reason: failed[0].reason })
        );
      }
      setSelected(new Set());
      queryClient.invalidateQueries({ queryKey: getGetChangesQueryKey() });
    },
    onError: () => toast.error(t('changes.bulkError')),
  });

  return (
    <div className="mx-auto max-w-225 px-6 py-6">
      <header className="mb-5 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight text-ink">{t('changes.title')}</h1>
          <p className="text-sm text-ink-muted">
            {data ? t('changes.countDetected', { count: data.total }) : t('changes.subtitle')}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div
            role="tablist"
            className="flex items-center gap-1 rounded-lg border border-line-soft bg-canvas-sunken p-0.5"
          >
            {FILTERS.map((f) => (
              <button
                key={f.value}
                role="tab"
                aria-selected={filter === f.value}
                onClick={() => setFilter(f.value)}
                className="relative rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors"
                style={{ color: filter === f.value ? 'var(--color-ink)' : 'var(--color-ink-muted)' }}
              >
                {filter === f.value && (
                  <motion.span
                    layoutId="chg-filter"
                    className="absolute inset-0 -z-10 rounded-md bg-surface-raised"
                    transition={{ type: 'spring', stiffness: 500, damping: 36 }}
                  />
                )}
                {t(f.key)}
              </button>
            ))}
          </div>
          <SavedViewsMenu
            surface="changes"
            filters={{ severity: filter }}
            onApply={(f) => setFilter(f.severity ?? 'all')}
          />
        </div>
      </header>

      {operatorIds.size > 0 && selected.size > 0 && (
        <div className="mb-3 flex items-center justify-between gap-3 rounded-lg border border-line-soft bg-canvas-sunken px-4 py-2.5">
          <span className="text-xs text-ink-muted">
            {t('changes.bulkSelectedCount', { count: selected.size })}
          </span>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              disabled={bulkResolve.isPending}
              onClick={() => bulkResolve.mutate('dismissed')}
            >
              <XIcon size={13} /> {t('common.dismiss')}
            </Button>
            <Button
              variant="primary"
              size="sm"
              disabled={bulkResolve.isPending}
              onClick={() => bulkResolve.mutate('acknowledged')}
            >
              <CheckIcon size={13} /> {t('common.acknowledge')}
            </Button>
          </div>
        </div>
      )}

      <Panel>
        {isLoading ? (
          <SkeletonRows rows={6} />
        ) : isError || !data ? (
          <ErrorState description={t('changes.loadError')} onRetry={() => refetch()} />
        ) : items.length === 0 ? (
          <EmptyState title={t('changes.emptyTitle')} description={t('changes.emptyDesc')} />
        ) : (
          items.map((c, idx) => (
            <motion.div
              key={c.id}
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: idx * 0.03, duration: 0.25 }}
              className="group flex w-full items-center gap-4 border-b border-line-soft px-4 py-3 transition-colors last:border-0 hover:bg-surface-raised"
            >
              {operatorIds.has(c.serviceId) && (
                <input
                  type="checkbox"
                  aria-label={t('changes.bulkSelectLabel', { summary: c.summary })}
                  checked={selected.has(c.id)}
                  disabled={!isLowRisk(c.severity)}
                  onChange={() => toggleSelected(c.id)}
                  title={isLowRisk(c.severity) ? undefined : t('changes.bulkNotLowRisk')}
                  className="size-3.5 shrink-0 accent-[var(--color-accent-primary)] disabled:opacity-30"
                />
              )}
              <button
                onClick={() => navigate(`/changes/${c.id}`)}
                className="flex min-w-0 flex-1 items-center gap-4 text-left"
              >
                <SeverityTag severity={c.severity} />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm text-ink">{c.summary}</p>
                  <p className="flex items-center gap-1.5 font-mono text-2xs text-ink-faint">
                    <span className="text-accent-secondary-bright">{c.serviceName}</span>
                    <span>·</span>
                    <span>{c.changeType}</span>
                    <span>·</span>
                    <span>{t('common.ago', { time: relativeTime(c.detectedAt) })}</span>
                  </p>
                </div>
                {c.willTriggerAi && (
                  <span className="hidden rounded bg-accent-primary-tint px-1.5 py-0.5 text-2xs font-medium text-accent-primary sm:block">
                    {t('changes.aiUpdate')}
                  </span>
                )}
                <ArrowRightIcon
                  size={15}
                  className="shrink-0 text-line-strong transition-colors group-hover:text-ink-muted"
                />
              </button>
            </motion.div>
          ))
        )}
      </Panel>
      {data && pageCount > 1 && (
        <Pagination page={page} pageCount={pageCount} onPage={setPage} className="mt-4 justify-center" />
      )}
    </div>
  );
}
