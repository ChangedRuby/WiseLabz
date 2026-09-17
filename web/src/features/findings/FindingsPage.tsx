import { useState } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'motion/react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getGetFindingsQueryKey,
  postFindingsFindingIdResolve,
  useGetFindings,
} from '../../api/generated/findings/findings';
import type { QualityCheckType, QualityFindingStatus } from '../../api/model';
import { SeverityTag } from '../../components/ui/StatusDot';
import { Button } from '../../components/ui/Button';
import { Panel } from '../../components/ui/Panel';
import { Pagination } from '../../components/ui/Pagination';
import { EmptyState, ErrorState, SkeletonRows } from '../../components/ui/states';
import { ArrowRightIcon, CheckIcon } from '../../components/icons';
import { useOperatorConnectorIds } from '../../hooks/useRole';
import { useIsInstanceAdmin } from '../../hooks/useRole';
import { relativeTime } from '../../lib/time';

const CHECK_TYPES: Array<{ value: QualityCheckType | 'all'; label: string }> = [
  { value: 'all', label: 'findings.allChecks' },
  { value: 'stale', label: 'findings.stale' },
  { value: 'empty', label: 'findings.empty' },
  { value: 'failing', label: 'findings.failing' },
  { value: 'ownership_incomplete', label: 'findings.ownershipIncomplete' },
  { value: 'credential_rotation', label: 'findings.credentialRotation' },
  { value: 'compliance', label: 'findings.compliance' },
];

export function FindingsPage() {
  const { t } = useTranslation();
  const operatorIds = useOperatorConnectorIds();
  const isInstanceAdmin = useIsInstanceAdmin();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [checkType, setCheckType] = useState<QualityCheckType | 'all'>('all');
  const [status, setStatus] = useState<QualityFindingStatus>('open');
  const pageSize = 20;
  const findings = useGetFindings({
    page,
    pageSize,
    status,
    ...(checkType === 'all' ? {} : { checkType }),
  });
  const pageCount = findings.data
    ? Math.max(1, Math.ceil(findings.data.total / findings.data.pageSize))
    : 1;

  const resolveFinding = useMutation({
    mutationFn: (findingId: string) => postFindingsFindingIdResolve(findingId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: getGetFindingsQueryKey() }),
  });

  return (
    <div className="mx-auto max-w-205 px-6 py-6">
      <header className="mb-5 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight text-ink">{t('findings.title')}</h1>
          <p className="text-sm text-ink-muted">{t('findings.subtitle')}</p>
        </div>
        <div className="flex flex-wrap items-end gap-2">
          <label className="text-xs text-ink-muted">
            {t('findings.checkType')}
            <select
              value={checkType}
              onChange={(event) => {
                setCheckType(event.target.value as QualityCheckType | 'all');
                setPage(1);
              }}
              className="mt-1 block h-9 rounded-sm border border-line bg-surface px-2.5 text-sm text-ink outline-none focus-visible:border-accent-primary-soft"
            >
              {CHECK_TYPES.map((option) => (
                <option key={option.value} value={option.value}>{t(option.label)}</option>
              ))}
            </select>
          </label>
          <label className="text-xs text-ink-muted">
            {t('findings.status')}
            <select
              value={status}
              onChange={(event) => {
                setStatus(event.target.value as QualityFindingStatus);
                setPage(1);
              }}
              className="mt-1 block h-9 rounded-sm border border-line bg-surface px-2.5 text-sm text-ink outline-none focus-visible:border-accent-primary-soft"
            >
              <option value="open">{t('findings.open')}</option>
              <option value="resolved">{t('findings.resolved')}</option>
            </select>
          </label>
        </div>
      </header>

      {findings.isLoading ? (
        <Panel><SkeletonRows rows={4} /></Panel>
      ) : findings.isError || !findings.data ? (
        <Panel className="min-h-[40vh]">
          <ErrorState description={t('findings.loadError')} onRetry={() => findings.refetch()} />
        </Panel>
      ) : findings.data.items.length === 0 ? (
        <Panel className="min-h-[40vh]">
          <EmptyState icon={<CheckIcon size={20} />} title={t('findings.emptyTitle')} description={t('findings.emptyDesc')} />
        </Panel>
      ) : (
        <div className="flex flex-col gap-3">
          {findings.data.items.map((finding, index) => (
            <motion.div
              key={finding.id}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.04, duration: 0.25 }}
            >
              <Panel className="p-4">
                <div className="flex items-start gap-3">
                  <SeverityTag severity={finding.severity} className="mt-0.5" />
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-ink">{finding.title}</p>
                    <p className="mt-1 text-sm leading-relaxed text-ink-muted">{finding.description}</p>
                    <p className="mt-1.5 font-mono text-2xs text-ink-faint">
                      <span className="text-accent-secondary-bright">{finding.connectorName}</span> ·{' '}
                      {finding.checkType === 'compliance' && <>{t('findings.compliance')} · </>}
                      {t('findings.detected', { count: finding.detectedCount })} ·{' '}
                      {t('common.ago', { time: relativeTime(finding.lastSeenAt) })}
                    </p>
                  </div>
                </div>
                <div className="mt-3 flex items-center justify-end gap-2 border-t border-line-soft pt-3">
                  {isInstanceAdmin && finding.ruleId && (
                    <Link to="/settings/compliance" className="font-mono text-xs text-accent-secondary hover:text-accent-secondary-bright">
                      {t('findings.viewRule')}
                    </Link>
                  )}
                  <Link
                    to={finding.remediationLink}
                    className="inline-flex h-7 items-center justify-center gap-2 rounded-sm px-2.5 font-mono text-xs font-medium text-accent-secondary transition-colors hover:bg-surface hover:text-accent-secondary-bright"
                  >
                    {t('findings.remediate')} <ArrowRightIcon size={14} />
                  </Link>
                  {operatorIds.has(finding.connectorId) && finding.status === 'open' && (
                    <Button
                      variant="primary"
                      size="sm"
                      disabled={resolveFinding.isPending}
                      onClick={() => resolveFinding.mutate(finding.id)}
                    >
                      <CheckIcon size={14} /> {t('common.resolve')}
                    </Button>
                  )}
                </div>
              </Panel>
            </motion.div>
          ))}
          {pageCount > 1 && (
            <Pagination page={page} pageCount={pageCount} onPage={setPage} className="justify-center pt-1" />
          )}
        </div>
      )}
    </div>
  );
}
