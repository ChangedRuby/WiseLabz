/**
 * Revision history + version-to-version diff for a document. Selecting a revision
 * diffs it against the one before it (or against empty for the first revision).
 * The trigger of each revision (ai / template / manual) is shown so the reader
 * knows whether a sync or a human last touched the doc.
 */
import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetDocsDocIdVersions,
  useGetDocsDocIdVersionsRev,
  postDocsDocIdVersionsRevRestore,
  getGetDocsDocIdQueryKey,
  getGetDocsDocIdVersionsQueryKey,
  getGetDocsTreeQueryKey,
} from '../../api/generated/docs/docs';
import { DocDiff } from '../../components/diff/DiffViewer';
import { Button } from '../../components/ui/Button';
import { Skeleton, SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { useCanMutate } from '../../hooks/useRole';
import { toast } from '../../lib/toast';
import { cn } from '../../lib/cn';
import { relativeTime } from '../../lib/time';
import {
  HistoryIcon,
  SparklesIcon,
  FileTextIcon,
  UserIcon,
  SyncIcon,
  HistoryIcon as RestoreIcon,
} from '../../components/icons';
import type { DocVersionMetaTrigger } from '../../api/model';

const TRIGGER: Record<
  DocVersionMetaTrigger,
  { Icon: React.ComponentType<{ size?: number }>; tone: string }
> = {
  ai: { Icon: SparklesIcon, tone: 'var(--color-accent-primary)' },
  template: { Icon: FileTextIcon, tone: 'var(--color-accent-primary-bright)' },
  manual: { Icon: UserIcon, tone: 'var(--color-ink-muted)' },
  sync: { Icon: SyncIcon, tone: 'var(--color-ink-muted)' },
};

export function DocHistory({ docId, currentVersion }: { docId: string; currentVersion: number }) {
  const { t } = useTranslation();
  const canMutate = useCanMutate();
  const queryClient = useQueryClient();
  const versions = useGetDocsDocIdVersions(docId);
  // Mounted fresh per doc (DocReader is keyed by docId), so initial = latest rev.
  const [selected, setSelected] = useState<number>(currentVersion);

  const restore = useMutation({
    mutationFn: () => postDocsDocIdVersionsRevRestore(docId, selected),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetDocsDocIdQueryKey(docId) });
      queryClient.invalidateQueries({ queryKey: getGetDocsDocIdVersionsQueryKey(docId) });
      queryClient.invalidateQueries({ queryKey: getGetDocsTreeQueryKey() });
      toast.success(t('docs.restore.done', { version: selected }));
    },
    onError: () => toast.error(t('docs.restore.error')),
  });

  const prevRev = selected - 1;
  const current = useGetDocsDocIdVersionsRev(docId, selected, {
    query: { enabled: selected > 0 },
  });
  const previous = useGetDocsDocIdVersionsRev(docId, prevRev, {
    query: { enabled: prevRev > 0, retry: false },
  });

  const metas = versions.data ?? [];
  const selectedTrigger = metas.find((v) => v.rev === selected)?.trigger;

  const beforeText = previous.data?.content ?? '';
  const afterText = current.data?.content ?? '';

  const ready = !current.isLoading && (prevRev <= 0 || !previous.isFetching);
  const hasContent = useMemo(() => afterText !== '' || beforeText !== '', [afterText, beforeText]);

  if (versions.isLoading) return <SkeletonRows rows={4} />;
  if (versions.isError)
    return (
      <ErrorState description={t('docs.historyLoadError')} onRetry={() => versions.refetch()} />
    );
  if (metas.length === 0)
    return (
      <EmptyState
        icon={<HistoryIcon size={20} />}
        title={t('docs.noRevisionsTitle')}
        description={t('docs.noRevisionsDesc')}
      />
    );

  return (
    <div className="grid gap-5 lg:grid-cols-[200px_1fr]">
      {/* Version rail */}
      <ol className="flex flex-col gap-1">
        {metas.map((v) => {
          const trig = TRIGGER[v.trigger];
          const active = v.rev === selected;
          return (
            <li key={v.rev}>
              <button
                onClick={() => setSelected(v.rev)}
                className={cn(
                  'flex w-full items-center gap-2.5 rounded-md border px-2.5 py-2 text-left transition-colors',
                  active
                    ? 'border-accent-primary-soft bg-accent-primary-tint'
                    : 'border-transparent hover:bg-surface-raised'
                )}
              >
                <span
                  className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md"
                  style={{ color: trig.tone, backgroundColor: 'var(--color-canvas-sunken)' }}
                >
                  <trig.Icon size={13} />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="nums block text-sm font-medium text-ink">
                    v{v.rev}
                    {v.rev === currentVersion && (
                      <span className="ml-1.5 text-2xs font-normal text-ink-faint">
                        {t('common.current')}
                      </span>
                    )}
                  </span>
                  <span className="block text-2xs text-ink-faint">
                    {t(`docs.trigger.${v.trigger}`)} ·{' '}
                    {t('common.ago', { time: relativeTime(v.createdAt) })}
                  </span>
                </span>
              </button>
            </li>
          );
        })}
      </ol>

      {/* Diff */}
      <div className="min-w-0">
        <div className="mb-2.5 flex items-center justify-between gap-3">
          <p className="nums font-mono text-2xs text-ink-faint">
            {prevRev > 0
              ? t('docs.comparing', { prev: prevRev, current: selected })
              : t('docs.firstRevision', { current: selected })}
          </p>
          {canMutate && selected !== currentVersion && (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => restore.mutate()}
              disabled={restore.isPending}
            >
              <RestoreIcon size={13} />
              {restore.isPending
                ? t('docs.restore.pending')
                : t('docs.restore.action', { version: selected })}
            </Button>
          )}
        </div>
        {!ready ? (
          <div className="space-y-2">
            <Skeleton className="h-4 w-1/2" />
            <Skeleton className="h-32 w-full" />
          </div>
        ) : current.isError ? (
          <ErrorState
            title={t('docs.revisionUnavailableTitle')}
            description={t('docs.revisionUnavailableDesc')}
          />
        ) : !hasContent ? (
          <EmptyState title={t('docs.noTextualTitle')} description={t('docs.noTextualDesc')} />
        ) : (
          <DocDiff
            before={beforeText}
            after={afterText}
            label={`v${prevRev > 0 ? prevRev : 0} → v${selected}`}
            headTrigger={selectedTrigger}
          />
        )}
      </div>
    </div>
  );
}
