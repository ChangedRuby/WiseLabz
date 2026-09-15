/** Single change — metadata, the DiffViewer, and the accept/reject resolution loop. */
import { useState } from 'react';
import { AnimatePresence, motion } from 'motion/react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetChangesChangeId,
  postChangesChangeIdAck,
  postChangesChangeIdDismiss,
  postChangesChangeIdAiUpdate,
  postChangesChangeIdExplain,
  getGetChangesChangeIdQueryKey,
} from '../../api/generated/changes/changes';
import { getGetChangesQueryKey } from '../../api/generated/changes/changes';
import { DiffViewer } from '../../components/diff/DiffViewer';
import { RunbookPanel } from '../../components/runbook/RunbookPanel';
import { SeverityTag } from '../../components/ui/StatusDot';
import { Button } from '../../components/ui/Button';
import { Panel } from '../../components/ui/Panel';
import { Skeleton, SkeletonRows, ErrorState } from '../../components/ui/states';
import { fullDate } from '../../lib/time';
import {
  ArrowRightIcon,
  CheckIcon,
  XIcon,
  SparklesIcon,
  FileTextIcon,
  ChatIcon,
} from '../../components/icons';

export function ChangeDetailPage() {
  const { t } = useTranslation();
  const { changeId } = useParams<{ changeId: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetChangesChangeId(changeId ?? '');

  // The accept/reject loop. On success the panel plays the resolve motion (PRODUCT
  // motion moment #4) and routes back to the feed; the caches are invalidated so
  // the resolved change drops out of the list.
  const [resolved, setResolved] = useState<null | 'acknowledged' | 'dismissed'>(null);
  const resolve = useMutation({
    mutationFn: (action: 'ack' | 'dismiss') =>
      action === 'ack'
        ? postChangesChangeIdAck(changeId ?? '')
        : postChangesChangeIdDismiss(changeId ?? ''),
    onSuccess: (_res, action) => {
      setResolved(action === 'ack' ? 'acknowledged' : 'dismissed');
      queryClient.invalidateQueries({ queryKey: getGetChangesQueryKey() });
      if (changeId) {
        queryClient.invalidateQueries({ queryKey: getGetChangesChangeIdQueryKey(changeId) });
      }
    },
  });
  const pending = resolve.isPending;

  // Fire-and-forget AI update request; the WS doc.ai_suggestion event and the
  // backend's willTriggerAi flag are wired separately — this just queues it.
  const [aiRequested, setAiRequested] = useState(false);
  const aiUpdate = useMutation({
    mutationFn: () => postChangesChangeIdAiUpdate(changeId ?? ''),
    onSuccess: () => {
      setAiRequested(true);
      if (changeId) {
        queryClient.invalidateQueries({ queryKey: getGetChangesChangeIdQueryKey(changeId) });
      }
    },
  });
  const aiQueued = aiRequested || Boolean(data?.willTriggerAi);

  // On-demand narration (issue #238, piece 2/3): generated once by the
  // backend and cached on the change record — refetching the change reuses
  // the cached text instead of re-invoking the AI provider.
  const explain = useMutation({
    mutationFn: () => postChangesChangeIdExplain(changeId ?? ''),
    onSuccess: () => {
      if (changeId) {
        queryClient.invalidateQueries({ queryKey: getGetChangesChangeIdQueryKey(changeId) });
      }
    },
  });

  return (
    <div className="mx-auto max-w-210 px-6 py-6">
      <button
        onClick={() => navigate('/changes')}
        className="mb-4 inline-flex items-center gap-1.5 text-xs text-ink-muted transition-colors hover:text-ink"
      >
        <ArrowRightIcon size={13} className="rotate-180" />
        {t('changes.back')}
      </button>

      {isLoading ? (
        <Panel className="p-6">
          <Skeleton className="mb-3 h-6 w-2/3" />
          <SkeletonRows rows={5} />
        </Panel>
      ) : isError || !data ? (
        <Panel className="min-h-[40vh]">
          <ErrorState description={t('changes.detailLoadError')} onRetry={() => refetch()} />
        </Panel>
      ) : (
        <AnimatePresence onExitComplete={() => navigate('/changes')}>
          {!resolved && (
            <motion.div
              key="panel"
              exit={{ opacity: 0, y: -8, scale: 0.985 }}
              transition={{ duration: 0.22, ease: [0.22, 1, 0.36, 1] }}
            >
              <Panel>
                <div className="border-b border-line-soft px-6 py-5">
                  <div className="mb-2 flex flex-wrap items-center gap-2">
                    <SeverityTag severity={data.severity} />
                    <span className="font-mono text-2xs text-ink-faint">{data.changeType}</span>
                    {data.willTriggerAi && (
                      <span className="inline-flex items-center gap-1 rounded bg-accent-primary-tint px-1.5 py-0.5 text-2xs font-medium text-accent-primary">
                        <SparklesIcon size={11} /> {t('changes.aiUpdateQueued')}
                      </span>
                    )}
                  </div>
                  <h1 className="text-lg font-semibold tracking-tight text-balance text-ink">
                    {data.summary}
                  </h1>
                  <p className="mt-1 font-mono text-2xs text-ink-faint">
                    <span className="text-accent-secondary-bright">{data.serviceName}</span> ·{' '}
                    {t('changes.detected', { date: fullDate(data.detectedAt) })}
                  </p>
                </div>

                <div className="px-6 py-5">
                  {data.narration ? (
                    <section className="mb-4 rounded-lg border border-line-soft bg-canvas-sunken p-3" aria-labelledby="narration-heading">
                      <h2 id="narration-heading" className="flex items-center gap-1.5 text-xs font-medium text-ink">
                        <ChatIcon size={13} /> {t('changes.explainHeading')}
                      </h2>
                      <p className="mt-1.5 text-xs text-ink-muted">{data.narration}</p>
                    </section>
                  ) : (
                    <div className="mb-4 flex items-center justify-between rounded-lg border border-dashed border-line-soft p-3">
                      <span className="text-xs text-ink-faint">{t('changes.explainPrompt')}</span>
                      <Button
                        variant="ghost"
                        size="sm"
                        disabled={explain.isPending}
                        onClick={() => explain.mutate()}
                      >
                        <ChatIcon size={13} />
                        {explain.isPending ? t('changes.explainLoading') : t('changes.explain')}
                      </Button>
                    </div>
                  )}
                  {explain.isError && (
                    <p className="mb-4 text-2xs text-[var(--color-err)]">{t('changes.explainError')}</p>
                  )}

                  <DiffViewer diff={data.diff} />

                  {data.diff.format === 'doc' && data.provenance && data.provenance.length > 0 && (
                    <section className="mt-4 rounded-lg border border-line-soft bg-canvas-sunken p-3" aria-labelledby="provenance-heading">
                      <h2 id="provenance-heading" className="text-xs font-medium text-ink">
                        Change provenance
                      </h2>
                      <ul className="mt-2 space-y-2">
                        {data.provenance.map((entry, index) => (
                          <li key={`${entry.snapshotPath}-${entry.templateSection}-${index}`} className="text-xs text-ink-muted">
                            <code className="text-accent-secondary-bright">{entry.snapshotPath}</code>
                            <span> → {entry.templateSection} → </span>
                            {entry.diffLines.map((line, lineIndex) => (
                              <span key={line}>
                                {lineIndex > 0 && ', '}
                                <a className="text-accent-secondary-bright underline underline-offset-2" href={`#diff-line-${line}`}>
                                  line {line}
                                </a>
                              </span>
                            ))}
                          </li>
                        ))}
                      </ul>
                    </section>
                  )}

                  {data.affectedDocIds && data.affectedDocIds.length > 0 && (
                    <div className="mt-4 flex flex-wrap items-center gap-2">
                      <span className="text-2xs text-ink-faint">
                        {t('changes.affectedDocs')}
                      </span>
                      {data.affectedDocIds.map((id) => (
                        <button
                          key={id}
                          onClick={() => navigate(`/docs/${id}`)}
                          className="inline-flex items-center gap-1.5 rounded-md border border-line-soft bg-canvas-sunken px-2 py-1 text-xs text-ink-muted transition-colors hover:border-accent-secondary-soft hover:text-accent-secondary-bright"
                        >
                          <FileTextIcon size={13} />
                          {id}
                        </button>
                      ))}
                    </div>
                  )}

                  <RunbookPanel changeType={data.changeType} />
                </div>

                <div className="flex items-center justify-end gap-2 border-t border-line-soft px-6 py-4">
                  {aiQueued && (
                    <span className="text-2xs text-ink-faint">{t('changes.aiUpdateQueued')}</span>
                  )}
                  <Button
                    variant="ghost"
                    size="md"
                    disabled={aiUpdate.isPending || aiQueued}
                    onClick={() => aiUpdate.mutate()}
                  >
                    <SparklesIcon size={15} /> {t('changes.aiUpdate')}
                  </Button>
                  <Button
                    variant="ghost"
                    size="md"
                    disabled={pending}
                    onClick={() => resolve.mutate('dismiss')}
                  >
                    <XIcon size={15} /> {t('common.dismiss')}
                  </Button>
                  <Button
                    variant="primary"
                    size="md"
                    disabled={pending}
                    onClick={() => resolve.mutate('ack')}
                  >
                    <CheckIcon size={15} /> {t('common.acknowledge')}
                  </Button>
                </div>
              </Panel>
            </motion.div>
          )}
        </AnimatePresence>
      )}
    </div>
  );
}
