/**
 * Small panel showing the runbook attached to a change type or alert severity,
 * if one exists. The backend enforces a unique target per runbook, so at most
 * one can match — silent (renders nothing) while loading or when none exists,
 * since most change types/severities won't have one.
 *
 * Each step links a connector restart/start/stop; execute follows the same
 * dry-run-preview / elevation-confirm flow as ServiceDetailPage (#282). Linking
 * a step grants nothing by itself — `canExecute` reflects whether the caller
 * currently holds an operator grant on the step's connector, and a disabled
 * button surfaces why. There is no "run all".
 */
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useGetRunbooks, useGetRunbooksRunbookId } from '../../api/generated/runbooks/runbooks';
import { executeRunbookStep } from '../../api/generated/runbooks/runbooks';
import type { RestartPreview, RunbookStep } from '../../api/model';
import type { Severity } from '../../api/model/severity';
import { FileTextIcon, PlayIcon } from '../icons';
import { useMutatingOp } from '../manager/useMutatingOp';
import { MutatingOpDialogs, type MutatingOpMessages } from '../manager/LifecycleOp';
import { toast } from '../../lib/toast';

type RunbookPanelProps =
  | { changeType: string; alertSeverity?: never; runbookId?: never }
  | { alertSeverity: Severity; changeType?: never; runbookId?: never }
  | { runbookId: string; changeType?: never; alertSeverity?: never };

export function RunbookPanel(props: RunbookPanelProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const byTarget = useGetRunbooks(
    'changeType' in props && props.changeType
      ? { changeType: props.changeType }
      : 'alertSeverity' in props && props.alertSeverity
        ? { alertSeverity: props.alertSeverity }
        : undefined,
    { query: { enabled: !('runbookId' in props && props.runbookId) } }
  );
  const byId = useGetRunbooksRunbookId('runbookId' in props ? (props.runbookId ?? '') : '', {
    query: { enabled: Boolean('runbookId' in props && props.runbookId) },
  });

  const isLoading = 'runbookId' in props && props.runbookId ? byId.isLoading : byTarget.isLoading;
  const runbook = 'runbookId' in props && props.runbookId ? byId.data : byTarget.data?.items[0];

  if (isLoading || !runbook) return null;

  return (
    <section className="mt-4 rounded-lg border border-line-soft bg-canvas-sunken p-3" aria-labelledby="runbook-heading">
      <h2 id="runbook-heading" className="text-xs font-medium text-ink">
        {t('runbooks.heading')}
      </h2>
      <p className="mt-1.5 text-sm font-medium text-ink">{runbook.title}</p>
      <p className="mt-1 whitespace-pre-wrap text-xs leading-relaxed text-ink-muted">{runbook.body}</p>
      <div className="mt-2 flex flex-wrap items-center gap-2">
        {runbook.docId && (
          <button
            onClick={() => navigate(`/docs/${runbook.docId}`)}
            className="inline-flex items-center gap-1.5 rounded-md border border-line-soft bg-canvas-sunken px-2 py-1 text-xs text-ink-muted transition-colors hover:border-accent-secondary-soft hover:text-accent-secondary-bright"
          >
            <FileTextIcon size={13} />
            {runbook.docId}
          </button>
        )}
        {runbook.snapshotId && (
          <span className="text-2xs text-ink-faint">
            {t('runbooks.snapshotRef', { snapshotId: runbook.snapshotId })}
          </span>
        )}
      </div>

      {runbook.steps.length > 0 && (
        <div className="mt-3 border-t border-line-soft pt-3">
          <h3 className="text-2xs font-medium text-ink-faint">{t('runbooks.steps.heading')}</h3>
          <ul className="mt-1.5 space-y-1.5">
            {runbook.steps.map((step) => (
              <RunbookStepRow key={step.id} runbookId={runbook.id} step={step} />
            ))}
          </ul>
        </div>
      )}
    </section>
  );
}

function RunbookStepRow({ runbookId, step }: { runbookId: string; step: RunbookStep }) {
  const { t } = useTranslation();

  const op = useMutatingOp({
    previewFn: () =>
      executeRunbookStep(runbookId, step.id, { dryRun: true }) as Promise<RestartPreview>,
    executeFn: (token) =>
      executeRunbookStep(
        runbookId,
        step.id,
        { dryRun: false },
        token ? { headers: { 'X-Elevation-Token': token } } : undefined
      ),
    onSuccess: () => toast.success(t('runbooks.steps.toastSuccess')),
  });

  const messages: MutatingOpMessages = {
    previewTitle: t('runbooks.steps.previewTitle'),
    previewNotice: t('runbooks.steps.previewNotice'),
    previewError: t('runbooks.steps.previewError'),
    targetLabel: t('services.detail.opTarget'),
    downtimeLabel: t('services.detail.opDowntime'),
    downtimeSeconds: (count) => t('services.detail.opSeconds', { count }),
    downtimeIndefinite: t('services.detail.opIndefinite'),
    dependenciesLabel: t('services.detail.opDependencies'),
    noDependencies: t('services.detail.opNoDependencies'),
    failed: t('runbooks.steps.failed'),
    retry: t('common.retry'),
    confirmTitle: t('runbooks.steps.confirmTitle', { title: step.title }),
    confirmDescription: t('runbooks.steps.confirmDescription', {
      verb: step.verb,
      connector: step.connectorName,
    }),
    confirmLabel: t('runbooks.steps.execute'),
  };

  const blockedReason =
    step.executeBlockedReason === 'no_operator_grant'
      ? t('runbooks.steps.blockedNoOperatorGrant', { connector: step.connectorName })
      : undefined;

  return (
    <li className="rounded-md border border-line-soft bg-surface px-2.5 py-1.5">
      <div className="flex items-center justify-between gap-2">
        <div className="min-w-0 flex-1">
          <p className="truncate text-xs font-medium text-ink">{step.title}</p>
          <p className="truncate font-mono text-2xs text-ink-faint">
            {step.verb} · {step.connectorName}
            {step.entityRef && ` · ${t('runbooks.steps.entityRef', { entityRef: step.entityRef })}`}
          </p>
        </div>
        <button
          type="button"
          onClick={op.open}
          disabled={!step.canExecute || op.preview.isPending}
          aria-describedby={blockedReason ? `runbook-step-${step.id}-reason` : undefined}
          className="inline-flex shrink-0 items-center gap-1.5 rounded-sm border border-line-strong px-2 py-1 font-mono text-2xs text-ink transition-colors hover:border-accent-primary-soft disabled:cursor-not-allowed disabled:opacity-40"
        >
          <PlayIcon size={12} />
          {t('runbooks.steps.execute')}
        </button>
      </div>
      {blockedReason && (
        <p id={`runbook-step-${step.id}-reason`} className="mt-1 text-2xs text-ink-faint">
          {blockedReason}
        </p>
      )}

      <MutatingOpDialogs
        op={op}
        action={`connector.${step.verb}`}
        resourceName={step.connectorName}
        messages={messages}
      />
    </li>
  );
}
