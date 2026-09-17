import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getGetTemplatesQueryKey,
  getGetTemplatesTemplateIdQueryKey,
  getGetTemplatesTemplateIdVersionsQueryKey,
  postTemplatesTemplateIdVersionsRevRestore,
  useGetTemplatesTemplateIdVersions,
  useGetTemplatesTemplateIdVersionsRev,
} from '../../api/generated/templates/templates';
import type { TemplateVersion, TemplateVersionMetaTrigger } from '../../api/model';
import { DocDiff } from '../../components/diff/DiffViewer';
import { Button } from '../../components/ui/Button';
import { EmptyState, ErrorState, Skeleton, SkeletonRows } from '../../components/ui/states';
import { useIsInstanceAdmin } from '../../hooks/useRole';
import { toast } from '../../lib/toast';
import { cn } from '../../lib/cn';
import { relativeTime } from '../../lib/time';
import { CheckIcon, HistoryIcon, HistoryIcon as RestoreIcon } from '../../components/icons';

const TRIGGER: Record<
  TemplateVersionMetaTrigger,
  { Icon: React.ComponentType<{ size?: number }>; tone: string }
> = {
  save: { Icon: CheckIcon, tone: 'var(--color-accent-primary)' },
  restore: { Icon: RestoreIcon, tone: 'var(--color-accent-secondary-bright)' },
};

function versionContent(version?: TemplateVersion): string {
  if (!version) return '';
  const { name, description, appliesTo, sections } = version;
  return JSON.stringify({ name, description, appliesTo, sections }, null, 2);
}

export function TemplateHistory({
  templateId,
  currentVersion,
}: {
  templateId: string;
  currentVersion: number;
}) {
  const { t } = useTranslation();
  const canMutate = useIsInstanceAdmin();
  const queryClient = useQueryClient();
  const versions = useGetTemplatesTemplateIdVersions(templateId);
  const [selected, setSelected] = useState(currentVersion);

  const restore = useMutation({
    mutationFn: () => postTemplatesTemplateIdVersionsRevRestore(templateId, selected),
    onSuccess: (updated) => {
      setSelected(updated.currentVersion);
      queryClient.setQueryData(getGetTemplatesTemplateIdQueryKey(templateId), updated);
      queryClient.invalidateQueries({ queryKey: getGetTemplatesTemplateIdQueryKey(templateId) });
      queryClient.invalidateQueries({
        queryKey: getGetTemplatesTemplateIdVersionsQueryKey(templateId),
      });
      queryClient.invalidateQueries({ queryKey: getGetTemplatesQueryKey() });
      toast.success(t('templates.restore.done', { version: selected }));
    },
    onError: () => toast.error(t('templates.restore.error')),
  });

  const prevRev = selected - 1;
  const current = useGetTemplatesTemplateIdVersionsRev(templateId, selected, {
    query: { enabled: selected > 0 },
  });
  const previous = useGetTemplatesTemplateIdVersionsRev(templateId, prevRev, {
    query: { enabled: prevRev > 0, retry: false },
  });
  const metas = versions.data ?? [];
  const beforeText = versionContent(previous.data);
  const afterText = versionContent(current.data);
  const ready = !current.isLoading && (prevRev <= 0 || !previous.isFetching);
  const hasContent = useMemo(() => afterText !== '' || beforeText !== '', [afterText, beforeText]);

  if (versions.isLoading) return <SkeletonRows rows={4} />;
  if (versions.isError)
    return (
      <ErrorState
        description={t('templates.historyLoadError')}
        onRetry={() => versions.refetch()}
      />
    );
  if (metas.length === 0)
    return (
      <EmptyState
        icon={<HistoryIcon size={20} />}
        title={t('templates.noRevisionsTitle')}
        description={t('templates.noRevisionsDesc')}
      />
    );

  return (
    <div className="grid gap-5 lg:grid-cols-[200px_1fr]">
      <ol className="flex flex-col gap-1">
        {metas.map((version) => {
          const trigger = TRIGGER[version.trigger];
          const active = version.rev === selected;
          return (
            <li key={version.rev}>
              <button
                onClick={() => setSelected(version.rev)}
                className={cn(
                  'flex w-full items-center gap-2.5 rounded-md border px-2.5 py-2 text-left transition-colors',
                  active
                    ? 'border-accent-primary-soft bg-accent-primary-tint'
                    : 'border-transparent hover:bg-surface-raised'
                )}
              >
                <span
                  className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md"
                  style={{ color: trigger.tone, backgroundColor: 'var(--color-canvas-sunken)' }}
                >
                  <trigger.Icon size={13} />
                </span>
                <span className="min-w-0 flex-1">
                  <span className="nums block text-sm font-medium text-ink">
                    v{version.rev}
                    {version.rev === currentVersion && (
                      <span className="ml-1.5 text-2xs font-normal text-ink-faint">
                        {t('common.current')}
                      </span>
                    )}
                  </span>
                  <span className="block text-2xs text-ink-faint">
                    {t(`templates.trigger.${version.trigger}`)} ·{' '}
                    {t('common.ago', { time: relativeTime(version.createdAt) })}
                  </span>
                </span>
              </button>
            </li>
          );
        })}
      </ol>

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
                ? t('templates.restore.pending')
                : t('templates.restore.action', { version: selected })}
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
            title={t('templates.revisionUnavailableTitle')}
            description={t('templates.revisionUnavailableDesc')}
          />
        ) : !hasContent ? (
          <EmptyState
            title={t('templates.noTextualTitle')}
            description={t('templates.noTextualDesc')}
          />
        ) : (
          <DocDiff
            before={beforeText}
            after={afterText}
            label={`v${prevRev > 0 ? prevRev : 0} → v${selected}`}
          />
        )}
      </div>
    </div>
  );
}
