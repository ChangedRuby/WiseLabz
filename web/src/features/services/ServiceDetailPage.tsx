/**
 * Service detail (`/services/:id`) — the "live picture + operate on it" surface for
 * one connector. Four panels: live raw-data snapshot (status + last sync), the
 * linked service doc, service-scoped change history, and a recent activity timeline
 * fed by the live store. Operator controls (sync / enable-disable / edit / remove)
 * are inline and role-gated; the destructive remove runs the full blast-radius +
 * step-up + type-to-confirm flow via <ConfirmDestructive/>.
 */
import { useMemo, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetConnectorsConnectorId,
  useGetConnectorsConnectorIdData,
  useGetConnectorsConnectorIdSyncs,
  useGetConnectorsConnectorIdConfigFields,
  useGetConnectorsSchema,
  postConnectorsConnectorIdRestart,
  postConnectorsConnectorIdStart,
  postConnectorsConnectorIdStop,
  postConnectorsConnectorIdConfigPush,
  postConnectorsConnectorIdHealth,
  putConnectorsConnectorIdEnabled,
  putConnectorsConnectorId,
  getGetConnectorsQueryKey,
} from '../../api/generated/connectors/connectors';
import type { RestartPreview, ConfigField } from '../../api/model';
import { isAxiosError } from 'axios';
import { toast } from '../../lib/toast';
import { useGetChanges } from '../../api/generated/changes/changes';
import {
  useGetDocsServiceConnectorId,
  getGetDocsServiceConnectorIdQueryKey,
  postDocsGenerate,
} from '../../api/generated/docs/docs';
import { useGetTemplates } from '../../api/generated/templates/templates';
import { useLive } from '../../store/live';
import { useCanMutate } from '../../hooks/useRole';
import { runSync } from '../../lib/runSync';
import { relativeTime, durationLabel } from '../../lib/time';
import { StatusPill, SeverityTag } from '../../components/ui/StatusDot';
import { toneColor, type Tone } from '../../components/ui/status';
import { Button } from '../../components/ui/Button';
import { Panel } from '../../components/ui/Panel';
import { Dialog } from '../../components/ui/Dialog';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { Markdown } from '../../components/docs/Markdown';
import { ConfirmDestructive } from '../../components/manager/ConfirmDestructive';
import { ElevationConfirm } from '../../components/manager/ElevationConfirm';
import {
  ArrowRightIcon,
  SyncIcon,
  FileTextIcon,
  DiffIcon,
  SparklesIcon,
  ChevronDownIcon,
  GaugeIcon,
} from '../../components/icons';
import type { ServiceStatus, SyncRunStatus } from '../../api/model';

/** Schedule cadence options (seconds); null = manual only. */
const SCHEDULE_OPTIONS: { value: number | null; labelKey: string }[] = [
  { value: null, labelKey: 'services.detail.scheduleOff' },
  { value: 900, labelKey: 'services.detail.schedule15m' },
  { value: 1800, labelKey: 'services.detail.schedule30m' },
  { value: 3600, labelKey: 'services.detail.schedule1h' },
  { value: 21600, labelKey: 'services.detail.schedule6h' },
  { value: 43200, labelKey: 'services.detail.schedule12h' },
  { value: 86400, labelKey: 'services.detail.schedule24h' },
];

export function ServiceDetailPage() {
  const { t } = useTranslation();
  const { id = '' } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const canMutate = useCanMutate();

  const connector = useGetConnectorsConnectorId(id);
  const { data: schemas } = useGetConnectorsSchema();
  const schema = useMemo(
    () => schemas?.find((s) => s.type === connector.data?.type) ?? null,
    [schemas, connector.data?.type],
  );
  const overrides = useLive((s) => s.statusOverrides);
  const activity = useLive((s) => s.activity);
  const [removing, setRemoving] = useState(false);

  const toggleEnabled = useMutation({
    mutationFn: (enabled: boolean) => putConnectorsConnectorIdEnabled(id, { enabled }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
      void connector.refetch();
    },
  });

  const updateSchedule = useMutation({
    mutationFn: (scheduleSeconds: number | null) => putConnectorsConnectorId(id, { scheduleSeconds }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
      void connector.refetch();
    },
  });

  const onOpSuccess = () => {
    queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
    void connector.refetch();
  };
  const restartOp = useMutatingOp(id, postConnectorsConnectorIdRestart, onOpSuccess);
  const startOp = useMutatingOp(id, postConnectorsConnectorIdStart, onOpSuccess);
  const stopOp = useMutatingOp(id, postConnectorsConnectorIdStop, onOpSuccess);

  const healthCheck = useMutation({
    mutationFn: () => postConnectorsConnectorIdHealth(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
      void connector.refetch();
    },
  });

  if (connector.isLoading) {
    return (
      <div className="mx-auto max-w-275 px-6 py-6">
        <Panel className="p-6">
          <SkeletonRows rows={6} />
        </Panel>
      </div>
    );
  }
  if (connector.isError || !connector.data) {
    return (
      <div className="mx-auto max-w-275 px-6 py-6">
        <Panel className="min-h-[40vh]">
          <ErrorState
            description={t('services.detail.notFound')}
            onRetry={() => connector.refetch()}
          />
        </Panel>
      </div>
    );
  }

  const c = connector.data;
  const status = (overrides[c.id] ?? c.status) as ServiceStatus;

  return (
    <div className="mx-auto max-w-275 px-6 py-6">
      <button
        onClick={() => navigate('/services')}
        className="mb-4 inline-flex items-center gap-1.5 text-xs text-ink-muted transition-colors hover:text-ink"
      >
        <ArrowRightIcon size={13} className="rotate-180" />
        {t('services.detail.back')}
      </button>

      {/* Header + operator controls */}
      <header className="mb-5 flex flex-wrap items-start justify-between gap-3">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-xl font-semibold tracking-tight text-ink">{c.name}</h1>
            <StatusPill status={status} />
          </div>
          <p className="mt-1 font-mono text-2xs text-ink-faint">
            {c.type} · {c.url?.replace(/^https?:\/\//, '')} ·{' '}
            {t('services.detail.lastSync', { time: relativeTime(c.lastSyncAt) })}
          </p>
          <p className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 font-mono text-2xs text-ink-faint">
            <span>
              {c.scheduleSeconds == null
                ? t('services.detail.scheduleManual')
                : t('services.detail.nextRun', { time: relativeTime(c.nextRunAt) })}
            </span>
            {c.lastSyncDurationMs != null && (
              <span>
                {t('services.detail.lastDuration', { duration: durationLabel(c.lastSyncDurationMs) })}
              </span>
            )}
            {(c.retryCount ?? 0) > 0 && (
              <span className="rounded-sm bg-warn-tint px-1.5 py-0.5 text-warn">
                {t('services.detail.retrying', { attempt: c.retryCount })}
              </span>
            )}
          </p>
          {c.lastSyncError && <p className="mt-1 text-2xs text-err">{c.lastSyncError}</p>}
        </div>
        {canMutate && (
          <div className="flex items-center gap-2">
            <label className="flex items-center gap-1.5 text-2xs text-ink-faint">
              {t('services.detail.scheduleLabel')}
              <div className="relative">
                <select
                  value={c.scheduleSeconds == null ? '' : String(c.scheduleSeconds)}
                  onChange={(e) =>
                    updateSchedule.mutate(e.target.value === '' ? null : Number(e.target.value))
                  }
                  disabled={updateSchedule.isPending}
                  className="h-8 appearance-none rounded-sm border border-line bg-surface pl-2.5 pr-7 text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
                >
                  {SCHEDULE_OPTIONS.map((opt) => (
                    <option key={opt.value ?? 'off'} value={opt.value ?? ''}>
                      {t(opt.labelKey)}
                    </option>
                  ))}
                </select>
                <ChevronDownIcon
                  size={12}
                  className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-ink-faint"
                />
              </div>
            </label>
            <Button
              size="sm"
              variant="secondary"
              onClick={() => runSync(c.id)}
              disabled={!c.enabled || schema?.stub}
              title={schema?.stub ? t('connectors.edit.testUnavailable') : undefined}
            >
              <SyncIcon size={14} /> {t('common.sync')}
            </Button>
            <Button size="sm" variant="ghost" onClick={startOp.open} disabled={startOp.preview.isPending}>
              {startOp.preview.isPending
                ? t('services.detail.startPreviewLoading')
                : t('services.detail.startPreview')}
            </Button>
            <Button size="sm" variant="ghost" onClick={stopOp.open} disabled={stopOp.preview.isPending}>
              {stopOp.preview.isPending
                ? t('services.detail.stopPreviewLoading')
                : t('services.detail.stopPreview')}
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={restartOp.open}
              disabled={restartOp.preview.isPending}
            >
              {restartOp.preview.isPending
                ? t('services.detail.restartPreviewLoading')
                : t('services.detail.restartPreview')}
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => healthCheck.mutate()}
              disabled={healthCheck.isPending}
              aria-describedby={healthCheck.isSuccess || healthCheck.isError ? 'health-check-result' : undefined}
            >
              <GaugeIcon size={14} />{' '}
              {healthCheck.isPending
                ? t('services.detail.healthCheckLoading')
                : t('services.detail.healthCheck')}
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => toggleEnabled.mutate(!c.enabled)}
              disabled={toggleEnabled.isPending}
            >
              {c.enabled ? t('common.disable') : t('common.enable')}
            </Button>
            <Button size="sm" variant="ghost" onClick={() => navigate(`/connectors/${c.id}/edit`)}>
              {t('common.edit')}
            </Button>
            <Button size="sm" variant="danger" onClick={() => setRemoving(true)}>
              {t('common.remove')}
            </Button>
          </div>
        )}
      </header>

      {(healthCheck.isSuccess || healthCheck.isError) && (
        <p
          id="health-check-result"
          role="status"
          className="mb-5 -mt-3 flex items-center gap-2 font-mono text-2xs text-ink-faint"
        >
          {healthCheck.isError ? (
            <span className="text-err">{t('services.detail.healthCheckError')}</span>
          ) : (
            <>
              <StatusPill status={healthCheck.data.status} />
              {healthCheck.data.message && <span>{healthCheck.data.message}</span>}
              {healthCheck.data.latencyMs != null && (
                <span>{t('services.detail.lastDuration', { duration: durationLabel(healthCheck.data.latencyMs) })}</span>
              )}
            </>
          )}
        </p>
      )}

      <div className="grid gap-4 lg:grid-cols-[1fr_320px]">
        <div className="space-y-4">
          <SnapshotPanel id={c.id} />
          <ServiceChangesPanel id={c.id} />
        </div>
        <div className="space-y-4">
          <LinkedDocPanel connectorId={c.id} connectorType={c.type} />
          <ConfigPushPanel id={c.id} connectorName={c.name} />
          <ActivityPanel activity={activity} />
          <SyncHistoryPanel id={c.id} />
        </div>
      </div>

      <ConfirmDestructive
        open={removing}
        connectorId={c.id}
        connectorName={c.name}
        onClose={() => setRemoving(false)}
        onConfirmed={() => {
          setRemoving(false);
          navigate('/services');
        }}
      />

      <MutatingOpDialogs op={restartOp} verb="restart" connectorName={c.name} />
      <MutatingOpDialogs op={startOp} verb="start" connectorName={c.name} />
      <MutatingOpDialogs op={stopOp} verb="stop" connectorName={c.name} />
    </div>
  );
}

/** Shared dry-run-preview + elevation-confirm state for one mutating verb
 * (restart/start/stop) — same shape ADR 0001/0002 give all three. */
function useMutatingOp(
  id: string,
  postFn: (
    id: string,
    body?: { entityRef?: string },
    params?: { dryRun?: boolean },
    options?: { headers?: Record<string, string> },
  ) => Promise<RestartPreview | { status: string }>,
  onSuccess: () => void,
) {
  const [previewOpen, setPreviewOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  const preview = useMutation({
    mutationFn: () => postFn(id, undefined, { dryRun: true }) as Promise<RestartPreview>,
  });

  const mutate = useMutation({
    mutationFn: (token: string | null) =>
      postFn(id, {}, { dryRun: false }, token ? { headers: { 'X-Elevation-Token': token } } : undefined),
    onSuccess: () => {
      setConfirmOpen(false);
      setPreviewOpen(false);
      onSuccess();
    },
  });

  const open = () => {
    preview.reset();
    setPreviewOpen(true);
    preview.mutate();
  };

  return { previewOpen, setPreviewOpen, confirmOpen, setConfirmOpen, preview, mutate, open };
}

type MutatingOp = ReturnType<typeof useMutatingOp>;

function MutatingOpDialogs({
  op,
  verb,
  connectorName,
}: {
  op: MutatingOp;
  verb: 'restart' | 'start' | 'stop';
  connectorName: string;
}) {
  const { t } = useTranslation();
  const k = (suffix: string) => `services.detail.${verb}${suffix}`;

  return (
    <>
      <Dialog
        open={op.previewOpen}
        onClose={() => op.setPreviewOpen(false)}
        title={t(k('PreviewTitle'))}
        size="sm"
      >
        {op.preview.isPending ? (
          <SkeletonRows rows={3} />
        ) : op.preview.isError || !op.preview.data ? (
          <div className="space-y-3">
            <p className="text-sm text-err">{t(k('PreviewError'))}</p>
            <Button size="sm" variant="secondary" onClick={() => op.preview.mutate()}>
              {t('common.retry')}
            </Button>
          </div>
        ) : (
          <div className="space-y-4">
            <p className="text-sm text-ink-muted">{t(k('PreviewNotice'))}</p>
            <dl className="grid grid-cols-2 gap-x-4 gap-y-3 rounded-md border border-line-soft bg-canvas-sunken p-3 text-sm">
              <div>
                <dt className="text-2xs text-ink-faint">{t('services.detail.opTarget')}</dt>
                <dd className="mt-0.5 font-medium text-ink">{op.preview.data.targetService}</dd>
              </div>
              <div>
                <dt className="text-2xs text-ink-faint">{t('services.detail.opDowntime')}</dt>
                <dd className="mt-0.5 font-medium text-ink">
                  {op.preview.data.estimatedDowntimeSeconds > 0
                    ? t('services.detail.opSeconds', {
                        count: op.preview.data.estimatedDowntimeSeconds,
                      })
                    : t('services.detail.opIndefinite')}
                </dd>
              </div>
            </dl>
            <div>
              <h3 className="text-sm font-semibold text-ink">{t('services.detail.opDependencies')}</h3>
              {op.preview.data.dependentServices.length === 0 ? (
                <p className="mt-1 text-sm text-ink-muted">{t('services.detail.opNoDependencies')}</p>
              ) : (
                <ul className="mt-2 divide-y divide-line-soft rounded-md border border-line-soft">
                  {op.preview.data.dependentServices.map((service) => (
                    <li key={`${service.kind}-${service.name}`} className="px-3 py-2 text-sm text-ink">
                      <span>{service.name}</span>
                      <span className="ml-2 font-mono text-2xs text-ink-faint">{service.kind}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>
            {op.mutate.isError && <p className="text-2xs text-err">{t(k('Failed'))}</p>}
            <div className="flex justify-end">
              <Button size="sm" variant="danger" onClick={() => op.setConfirmOpen(true)}>
                {t(k('Now'))}
              </Button>
            </div>
          </div>
        )}
      </Dialog>

      <ElevationConfirm
        open={op.confirmOpen}
        resourceName={connectorName}
        action={`connector.${verb}`}
        title={t(k('ConfirmTitle'), { name: connectorName })}
        description={t(k('ConfirmDescription'))}
        confirmLabel={t(k('Now'))}
        isPending={op.mutate.isPending}
        onClose={() => op.setConfirmOpen(false)}
        onConfirm={(token) => op.mutate.mutate(token)}
      />
    </>
  );
}

/** Field-level config-push form (ADR 0003): pick a whitelisted field, edit its
 * value, confirm with step-up. On a 409 (verify-diff mismatch), the backend
 * has already auto-reverted — this just surfaces that via toast. */
function ConfigPushPanel({ id, connectorName }: { id: string; connectorName: string }) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const fields = useGetConnectorsConnectorIdConfigFields(id);
  const [fieldKey, setFieldKey] = useState('');
  const [entityRef, setEntityRef] = useState('');
  const [value, setValue] = useState('');
  const [previousValue, setPreviousValue] = useState<string>('');
  const [confirmOpen, setConfirmOpen] = useState(false);

  const selected: ConfigField | undefined = fields.data?.find((f) => f.key === fieldKey);

  const push = useMutation({
    mutationFn: (token: string | null) =>
      postConnectorsConnectorIdConfigPush(
        id,
        {
          entityRef: selected?.entityScope ? entityRef : undefined,
          fieldKey,
          value: coerceFieldValue(selected?.type, value),
          previousValue: previousValue === '' ? undefined : coerceFieldValue(selected?.type, previousValue),
        },
        token ? { headers: { 'X-Elevation-Token': token } } : undefined,
      ),
    onSuccess: () => {
      setConfirmOpen(false);
      toast.success(t('services.detail.configPush'));
      queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
    },
    onError: (error) => {
      setConfirmOpen(false);
      if (isAxiosError(error) && error.response?.status === 409) {
        toast.warning(t('services.detail.configPushMismatchTitle'), {
          description: t('services.detail.configPushMismatchBody'),
        });
        return;
      }
      toast.error(t('services.detail.configPushFailed'));
    },
  });

  if (!fields.data || fields.data.length === 0) {
    return null;
  }

  return (
    <Panel className="p-5">
      <h2 className="mb-3 text-sm font-semibold text-ink">{t('services.detail.configPushTitle')}</h2>
      <div className="space-y-2.5">
        <label className="block">
          <span className="mb-1 block text-2xs text-ink-faint">{t('services.detail.configPushField')}</span>
          <select
            value={fieldKey}
            onChange={(e) => {
              setFieldKey(e.target.value);
              setValue('');
              setPreviousValue('');
            }}
            className="h-8 w-full appearance-none rounded-sm border border-line bg-surface px-2.5 text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
          >
            <option value="">{t('services.detail.configPushNone')}</option>
            {fields.data.map((f) => (
              <option key={f.key} value={f.key}>
                {f.label}
              </option>
            ))}
          </select>
        </label>
        {selected?.entityScope && (
          <input
            value={entityRef}
            onChange={(e) => setEntityRef(e.target.value)}
            placeholder="entityRef"
            className="h-8 w-full rounded-sm border border-line bg-surface px-2.5 font-mono text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
          />
        )}
        {selected && (
          <>
            <input
              value={previousValue}
              onChange={(e) => setPreviousValue(e.target.value)}
              placeholder="current value (for revert)"
              type={selected.type === 'number' ? 'number' : 'text'}
              className="h-8 w-full rounded-sm border border-line bg-surface px-2.5 text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
            />
            <input
              value={value}
              onChange={(e) => setValue(e.target.value)}
              placeholder={t('services.detail.configPushValue')}
              type={selected.type === 'number' ? 'number' : 'text'}
              className="h-8 w-full rounded-sm border border-line bg-surface px-2.5 text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
            />
          </>
        )}
        <div className="flex justify-end">
          <Button
            size="sm"
            variant="secondary"
            disabled={!selected || value === ''}
            onClick={() => setConfirmOpen(true)}
          >
            {t('services.detail.configPushSubmit')}
          </Button>
        </div>
      </div>

      <ElevationConfirm
        open={confirmOpen}
        resourceName={connectorName}
        action="connector.configPush"
        title={t('services.detail.configPushConfirmTitle', { field: selected?.label ?? fieldKey, name: connectorName })}
        description={t('services.detail.configPushConfirmDescription')}
        confirmLabel={t('services.detail.configPushSubmit')}
        isPending={push.isPending}
        onClose={() => setConfirmOpen(false)}
        onConfirm={(token) => push.mutate(token)}
      />
    </Panel>
  );
}

function coerceFieldValue(type: string | undefined, raw: string): unknown {
  if (type === 'number') return Number(raw);
  if (type === 'toggle') return raw === 'true';
  return raw;
}

function SnapshotPanel({ id }: { id: string }) {
  const { t } = useTranslation();
  const data = useGetConnectorsConnectorIdData(id);

  return (
    <Panel className="p-5">
      <div className="mb-3 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-ink">{t('services.detail.snapshot')}</h2>
        {data.data && (
          <span className="font-mono text-2xs text-ink-faint">
            {t('services.detail.fetchedAt', { time: relativeTime(data.data.fetchedAt) })}
          </span>
        )}
      </div>
      {data.isLoading ? (
        <SkeletonRows rows={4} />
      ) : data.isError || !data.data ? (
        <ErrorState
          description={t('services.detail.snapshotError')}
          onRetry={() => data.refetch()}
        />
      ) : data.data.sections.length === 0 ? (
        <EmptyState
          title={t('services.detail.neverSyncedTitle')}
          description={t('services.detail.neverSyncedDesc')}
        />
      ) : (
        <div className="space-y-3">
          {[...data.data.sections]
            .sort((a, b) => a.order - b.order)
            .map((s) => (
              <details
                key={s.title}
                open
                className="rounded-md border border-line-soft bg-canvas-sunken p-3"
              >
                <summary className="cursor-pointer text-xs font-semibold text-ink">
                  {s.title}
                </summary>
                <div className="mt-2 text-sm">
                  <Markdown source={s.content} />
                </div>
              </details>
            ))}
        </div>
      )}
    </Panel>
  );
}

function ServiceChangesPanel({ id }: { id: string }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const changes = useGetChanges({ serviceId: id, pageSize: 10 });

  return (
    <Panel className="p-5">
      <h2 className="mb-3 text-sm font-semibold text-ink">{t('services.detail.history')}</h2>
      {changes.isLoading ? (
        <SkeletonRows rows={3} />
      ) : changes.isError ? (
        <ErrorState
          description={t('services.detail.historyError')}
          onRetry={() => changes.refetch()}
        />
      ) : !changes.data || changes.data.items.length === 0 ? (
        <EmptyState
          title={t('services.detail.noChangesTitle')}
          description={t('services.detail.noChangesDesc')}
        />
      ) : (
        <ul className="divide-y divide-line-soft">
          {changes.data.items.map((ch) => (
            <li key={ch.id}>
              <button
                onClick={() => navigate(`/changes/${ch.id}`)}
                className="flex w-full items-start gap-3 py-2.5 text-left transition-colors hover:bg-surface-raised"
              >
                <SeverityTag severity={ch.severity} />
                <span className="flex-1">
                  <span className="block text-sm text-ink">{ch.summary}</span>
                  <span className="font-mono text-2xs text-ink-faint">
                    {ch.changeType} · {relativeTime(ch.detectedAt)}
                  </span>
                </span>
                <DiffIcon size={14} className="mt-1 shrink-0 text-ink-faint" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </Panel>
  );
}

function LinkedDocPanel({
  connectorId,
  connectorType,
}: {
  connectorId: string;
  connectorType: string;
}) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const canMutate = useCanMutate();
  const doc = useGetDocsServiceConnectorId(connectorId);
  const templates = useGetTemplates();
  const [picking, setPicking] = useState(false);
  const [templateId, setTemplateId] = useState('');

  // Templates with no appliesTo.type apply to any connector; others must match.
  const candidates = (templates.data ?? []).filter(
    (tpl) => !tpl.appliesTo?.type || tpl.appliesTo.type === connectorType
  );

  const generate = useMutation({
    mutationFn: () => postDocsGenerate({ templateId, connectorId }),
    onSuccess: (result) => {
      queryClient.invalidateQueries({
        queryKey: getGetDocsServiceConnectorIdQueryKey(connectorId),
      });
      setPicking(false);
      navigate(`/docs/${result.docId}`);
    },
  });

  const openPicker = () => {
    generate.reset();
    setTemplateId(candidates[0]?.id ?? '');
    setPicking(true);
  };

  return (
    <Panel className="p-5">
      <h2 className="mb-3 text-sm font-semibold text-ink">{t('services.detail.linkedDoc')}</h2>
      {doc.isError || !doc.data ? (
        <div className="space-y-3">
          <p className="text-sm text-ink-muted">{t('services.detail.noDoc')}</p>
          {canMutate && (
            <Button size="sm" variant="secondary" onClick={openPicker}>
              <SparklesIcon size={14} />
              {t('services.detail.generateDoc')}
            </Button>
          )}
        </div>
      ) : (
        <button
          onClick={() => navigate(`/docs/${doc.data.docId}`)}
          className="flex w-full items-center gap-2.5 rounded-md border border-line-soft p-3 text-left transition-colors hover:border-line-strong"
        >
          <FileTextIcon size={16} className="shrink-0 text-accent-secondary" />
          <span className="flex-1 text-sm text-ink">{doc.data.title}</span>
          <ArrowRightIcon size={14} className="text-ink-faint" />
        </button>
      )}

      <Dialog
        open={picking}
        onClose={() => setPicking(false)}
        title={t('services.detail.generateDocTitle')}
        size="sm"
      >
        <div className="space-y-3">
          <label className="block">
            <span className="mb-1 block text-2xs text-ink-faint">
              {t('services.detail.templateLabel')}
            </span>
            <div className="relative">
              <select
                value={templateId}
                onChange={(e) => setTemplateId(e.target.value)}
                disabled={candidates.length === 0}
                className="h-9 w-full appearance-none rounded-sm border border-line bg-surface px-2.5 pr-8 text-sm text-ink outline-none focus-visible:border-accent-primary-soft"
              >
                {candidates.length === 0 ? (
                  <option value="">{t('services.detail.noTemplates')}</option>
                ) : (
                  candidates.map((tpl) => (
                    <option key={tpl.id} value={tpl.id}>
                      {tpl.name}
                    </option>
                  ))
                )}
              </select>
              <ChevronDownIcon
                size={14}
                className="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 text-ink-faint"
              />
            </div>
          </label>

          {generate.isError && (
            <p className="text-xs text-err">{t('services.detail.generateDocError')}</p>
          )}

          <div className="flex items-center justify-end gap-2">
            <Button variant="ghost" size="sm" onClick={() => setPicking(false)}>
              {t('common.cancel')}
            </Button>
            <Button
              variant="primary"
              size="sm"
              disabled={!templateId || generate.isPending}
              onClick={() => generate.mutate()}
            >
              {generate.isPending ? t('services.detail.generating') : t('services.detail.generateDoc')}
            </Button>
          </div>
        </div>
      </Dialog>
    </Panel>
  );
}

function ActivityPanel({
  activity,
}: {
  activity: { id: string; label: string; detail?: string; at: string }[];
}) {
  const { t } = useTranslation();
  return (
    <Panel className="p-5">
      <h2 className="mb-3 text-sm font-semibold text-ink">{t('services.detail.activity')}</h2>
      {activity.length === 0 ? (
        <p className="text-sm text-ink-muted">{t('services.detail.noActivity')}</p>
      ) : (
        <ul className="space-y-2.5">
          {activity.slice(0, 8).map((a) => (
            <li key={a.id} className="flex items-start gap-2.5">
              <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-accent-primary" aria-hidden />
              <span className="flex-1">
                <span className="block text-xs text-ink">{a.label}</span>
                {a.detail && <span className="text-2xs text-ink-faint">{a.detail}</span>}
              </span>
              <span className="font-mono text-2xs text-ink-faint">{relativeTime(a.at)}</span>
            </li>
          ))}
        </ul>
      )}
    </Panel>
  );
}

const syncStatusMeta: Record<SyncRunStatus, { label: string; tone: Tone }> = {
  success: { label: 'success', tone: 'ok' },
  error: { label: 'error', tone: 'err' },
  skipped: { label: 'skipped', tone: 'idle' },
};

function SyncStatusTag({ status }: { status: SyncRunStatus }) {
  const { label, tone } = syncStatusMeta[status];
  const { fg, tint } = toneColor[tone];
  return (
    <span
      className="inline-flex items-center gap-1.5 rounded-sm px-1.5 py-0.5 font-mono text-2xs font-medium"
      style={{ color: fg, backgroundColor: tint }}
    >
      <span className="h-1.5 w-1.5" style={{ backgroundColor: fg }} />
      {label}
    </span>
  );
}

function SyncHistoryPanel({ id }: { id: string }) {
  const { t } = useTranslation();
  const syncs = useGetConnectorsConnectorIdSyncs(id, { limit: 10 });

  return (
    <Panel className="p-5">
      <h2 className="mb-3 text-sm font-semibold text-ink">{t('services.detail.syncHistory')}</h2>
      {syncs.isLoading ? (
        <SkeletonRows rows={3} />
      ) : syncs.isError ? (
        <ErrorState
          description={t('services.detail.syncHistoryError')}
          onRetry={() => syncs.refetch()}
        />
      ) : !syncs.data || syncs.data.length === 0 ? (
        <EmptyState
          title={t('services.detail.noSyncsTitle')}
          description={t('services.detail.noSyncsDesc')}
        />
      ) : (
        <ul className="divide-y divide-line-soft">
          {syncs.data.map((run) => (
            <li key={run.id} className="py-2.5">
              <div className="flex items-center gap-2.5">
                <SyncStatusTag status={run.status} />
                <span className="flex-1 font-mono text-2xs text-ink-faint">
                  {relativeTime(run.startedAt)} · {durationLabel(run.durationMs)} ·{' '}
                  {t('services.detail.attempt', { n: run.attempt })}
                </span>
              </div>
              <p className="mt-1 pl-5 font-mono text-2xs text-ink-faint">
                {t('services.detail.syncCounts', {
                  changes: run.changesCount ?? 0,
                  alerts: run.alertsCount ?? 0,
                })}
              </p>
              {run.status === 'error' && run.error && (
                <p className="mt-1 pl-5 text-2xs text-err">{run.error}</p>
              )}
            </li>
          ))}
        </ul>
      )}
    </Panel>
  );
}
