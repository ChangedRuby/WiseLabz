/**
 * Settings → System (operator-only). Instance health: backend + integration
 * health, the live WebSocket connection state, the sync schedule, the running
 * version, and backup operations (export, on-demand run, schedule, recent runs).
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  useGetSystemInfo,
  useGetHealth,
  useGetSystemJobs,
} from '../../api/generated/system/system';
import {
  getSystemBackupSchedule,
  putSystemBackupSchedule,
  getGetSystemBackupScheduleQueryKey,
  useGetSystemBackupRuns,
  postSystemBackupRun,
  getGetSystemBackupRunsQueryKey,
} from '../../api/generated/system/system';
import type { BackupSchedule, HealthComponentsItemStatus, JobInfo } from '../../api/model';
import { AXIOS_INSTANCE } from '../../api/axios-instance';
import { useLive, type WsState } from '../../store/live';
import { Panel, PanelHeader } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { ToneTag } from '../../components/ui/ToneTag';
import { TimeAgo } from '../../components/ui/TimeAgo';
import { type Tone } from '../../components/ui/status';
import { toast } from '../../lib/toast';
import { filenameFromContentDisposition, downloadBlob } from '../../lib/download';
import { SubHeader, Section, Field, TextInput } from './parts';
import {
  GaugeIcon,
  ServerIcon,
  ClockIcon,
  NetworkIcon,
  BoxIcon,
  DownloadIcon,
} from '../../components/icons';

const healthTone: Record<HealthComponentsItemStatus, Tone> = {
  ok: 'ok',
  degraded: 'warn',
  down: 'err',
};

const wsTone: Record<WsState, Tone> = {
  open: 'ok',
  connecting: 'warn',
  closed: 'err',
};

export function SystemPage() {
  const { t } = useTranslation();
  const ws = useLive((s) => s.ws);
  const info = useGetSystemInfo();
  const health = useGetHealth();
  const jobs = useGetSystemJobs();

  if (info.isLoading || health.isLoading) return <Loading />;

  return (
    <div>
      <SubHeader title={t('settings.system.title')} description={t('settings.system.subtitle')} />

      <div className="mb-4 grid gap-4 sm:grid-cols-3">
        <Stat
          icon={<NetworkIcon size={15} />}
          label={t('settings.system.wsStatus')}
          value={
            <ToneTag
              tone={wsTone[ws]}
              label={t(`settings.system.ws.${ws}`)}
            />
          }
        />
        <Stat
          icon={<ClockIcon size={15} />}
          label={t('settings.system.syncSchedule')}
          value={
            <span className="font-mono text-sm text-ink">{info.data?.syncSchedule ?? '—'}</span>
          }
        />
        <Stat
          icon={<GaugeIcon size={15} />}
          label={t('settings.system.version')}
          value={<span className="font-mono text-sm text-ink">v{info.data?.version ?? '—'}</span>}
        />
      </div>

      <Section
        title={t('settings.system.healthTitle')}
        description={
          health.data
            ? t('settings.system.overall', {
                status: health.data.status,
              })
            : undefined
        }
      >
        {health.isError || !health.data ? (
          <ErrorState
            description={t('settings.system.healthError')}
            onRetry={() => health.refetch()}
          />
        ) : (
          <ul className="divide-y divide-line-soft">
            {health.data.components.map((c) => (
              <li
                key={c.name}
                className="flex items-center justify-between gap-4 py-2.5 first:pt-0 last:pb-0"
              >
                <div className="min-w-0">
                  <p className="font-mono text-xs text-ink">{c.name}</p>
                  {c.detail && <p className="font-mono text-2xs text-ink-faint">{c.detail}</p>}
                </div>
                <ToneTag tone={healthTone[c.status]} label={c.status} />
              </li>
            ))}
          </ul>
        )}
      </Section>

      <Panel>
        <PanelHeader
          title={t('settings.system.integrationsTitle')}
          icon={<ServerIcon size={14} />}
          count={info.data?.integrations.length}
        />
        {info.isError || !info.data ? (
          <ErrorState description={t('settings.system.infoError')} onRetry={() => info.refetch()} />
        ) : (
          <ul className="divide-y divide-line-soft">
            {info.data.integrations.map((i, idx) => (
              <li
                key={i.name ?? idx}
                className="flex items-center justify-between gap-4 px-4 py-2.5"
              >
                <div className="min-w-0">
                  <p className="text-sm text-ink">{i.name}</p>
                  {i.detail && <p className="font-mono text-2xs text-ink-faint">{i.detail}</p>}
                </div>
                {i.status && <ToneTag tone={healthTone[i.status]} label={i.status} />}
              </li>
            ))}
          </ul>
        )}
      </Panel>

      <DiagnosticsSection />

      <BackupSection />

      <JobsSection jobs={jobs} />
    </div>
  );
}

function JobsSection({ jobs }: { jobs: ReturnType<typeof useGetSystemJobs> }) {
  const { t } = useTranslation();

  return (
    <Section
      title={t('settings.system.jobs.title')}
      description={t('settings.system.jobs.subtitle')}
    >
      {jobs.isLoading ? (
        <SkeletonRows rows={4} className="p-0" />
      ) : jobs.isError || !jobs.data ? (
        <ErrorState
          description={t('settings.system.jobs.loadError')}
          onRetry={() => jobs.refetch()}
        />
      ) : jobs.data.length === 0 ? (
        <EmptyState title={t('settings.system.jobs.empty')} />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-225 text-left text-xs">
            <thead className="border-b border-line-soft font-mono text-2xs text-ink-faint">
              <tr>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.name')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.status')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.lastRun')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.lastSuccess')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.lastFailure')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.nextRun')}</th>
                <th className="px-2 py-2 font-normal">{t('settings.system.jobs.lastError')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-line-soft">
              {jobs.data.map((job) => (
                <JobRow key={job.name} job={job} />
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Section>
  );
}

function JobRow({ job }: { job: JobInfo }) {
  const { t } = useTranslation();
  const status = job.lastStatus;
  const tone = status === 'failing' ? 'err' : status === 'ok' ? 'ok' : 'idle';

  return (
    <tr>
      <td className="px-2 py-2.5">
        <p className="font-mono text-xs text-ink">{job.name}</p>
        <p className="font-mono text-2xs text-ink-faint">{job.cronExpr}</p>
      </td>
      <td className="px-2 py-2.5">
        <ToneTag tone={tone} label={status ?? t('settings.system.jobs.neverRun')} />
      </td>
      <JobTime at={job.lastRunAt} />
      <JobTime at={job.lastSuccessAt} />
      <JobTime at={job.lastFailureAt} />
      <JobTime at={job.nextRunAt} />
      <td className="max-w-48 px-2 py-2.5 text-err">
        {job.lastError ? (
          <span className="block truncate" title={job.lastError}>
            {job.lastError}
          </span>
        ) : (
          '—'
        )}
      </td>
    </tr>
  );
}

function JobTime({ at }: { at?: string }) {
  return (
    <td className="whitespace-nowrap px-2 py-2.5 text-ink-muted">
      {at ? <TimeAgo at={at} /> : '—'}
    </td>
  );
}

function DiagnosticsSection() {
  const { t } = useTranslation();
  const [downloading, setDownloading] = useState(false);

  const handleDownload = async () => {
    setDownloading(true);
    try {
      const res = await AXIOS_INSTANCE.get('/system/diagnostics', { responseType: 'blob' });
      const filename = filenameFromContentDisposition(
        res.headers['content-disposition'] as string | undefined,
        'wiselabz-diagnostics.json'
      );
      downloadBlob(res.data as Blob, filename);
      toast.success(t('settings.system.diagnostics.downloaded'));
    } catch {
      toast.error(t('settings.system.diagnostics.downloadError'));
    } finally {
      setDownloading(false);
    }
  };

  return (
    <Section
      title={t('settings.system.diagnostics.title')}
      description={t('settings.system.diagnostics.subtitle')}
      action={
        <Button
          variant="secondary"
          size="sm"
          disabled={downloading}
          aria-label={t('settings.system.diagnostics.download')}
          onClick={() => void handleDownload()}
        >
          <DownloadIcon size={14} />
          {downloading ? t('settings.system.diagnostics.downloading') : t('settings.system.diagnostics.download')}
        </Button>
      }
    >
      <p className="text-xs leading-relaxed text-ink-muted">
        {t('settings.system.diagnostics.sanitizedNote')}
      </p>
    </Section>
  );
}

function BackupSection() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [exporting, setExporting] = useState(false);

  const schedule = useQuery({
    queryKey: getGetSystemBackupScheduleQueryKey(),
    queryFn: () => getSystemBackupSchedule(),
  });
  const runs = useGetSystemBackupRuns({ limit: 10 });

  const [draft, setDraft] = useState<BackupSchedule | null>(null);
  // Adjust state during render (React-blessed alternative to a syncing effect):
  // re-seed whenever the query yields a fresh reference, e.g. after an invalidate.
  const [seeded, setSeeded] = useState<BackupSchedule | null>(null);
  if (schedule.data && schedule.data !== seeded) {
    setSeeded(schedule.data);
    setDraft({ ...schedule.data });
  }

  const saveSchedule = useMutation({
    mutationFn: (body: BackupSchedule) => putSystemBackupSchedule(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetSystemBackupScheduleQueryKey() });
      toast.success(t('settings.system.backup.scheduleSaved'));
    },
    onError: () => toast.error(t('settings.system.backup.scheduleSaveError')),
  });

  const runBackup = useMutation({
    mutationFn: () => postSystemBackupRun(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetSystemBackupRunsQueryKey() });
      toast.success(t('settings.system.backup.runOk'));
    },
    onError: () => toast.error(t('settings.system.backup.runError')),
  });

  const handleExport = async () => {
    setExporting(true);
    try {
      const res = await AXIOS_INSTANCE.get('/system/backup/export', { responseType: 'blob' });
      const filename = filenameFromContentDisposition(
        res.headers['content-disposition'] as string | undefined,
        'wiselabz-backup.zip'
      );
      downloadBlob(res.data as Blob, filename);
    } catch {
      toast.error(t('settings.system.backup.exportError'));
    } finally {
      setExporting(false);
    }
  };

  const scheduleDirty = draft && schedule.data && JSON.stringify(draft) !== JSON.stringify(schedule.data);

  return (
    <Section
      title={t('settings.system.backup.title')}
      description={t('settings.system.backup.subtitle')}
      action={
        <div className="flex gap-2">
          <Button
            variant="secondary"
            size="sm"
            disabled={exporting}
            onClick={() => void handleExport()}
          >
            <DownloadIcon size={14} /> {t('settings.system.backup.export')}
          </Button>
          <Button
            variant="primary"
            size="sm"
            disabled={runBackup.isPending}
            onClick={() => runBackup.mutate()}
          >
            {t('settings.system.backup.runNow')}
          </Button>
        </div>
      }
    >
      {schedule.isLoading || !draft ? (
        <SkeletonRows rows={2} className="p-0" />
      ) : schedule.isError ? (
        <ErrorState
          description={t('settings.system.backup.scheduleError')}
          onRetry={() => schedule.refetch()}
        />
      ) : (
        <form
          className="grid gap-4 sm:grid-cols-3"
          onSubmit={(e) => {
            e.preventDefault();
            if (draft) saveSchedule.mutate(draft);
          }}
        >
          <Field label={t('settings.system.backup.cronExpr')} htmlFor="backup-cron">
            <TextInput
              id="backup-cron"
              value={draft.cronExpr}
              onChange={(e) => setDraft((d) => (d ? { ...d, cronExpr: e.target.value } : d))}
            />
          </Field>
          <Field label={t('settings.system.backup.maxBackups')} htmlFor="backup-max-backups">
            <TextInput
              id="backup-max-backups"
              type="number"
              value={draft.maxBackups}
              onChange={(e) =>
                setDraft((d) => (d ? { ...d, maxBackups: Number(e.target.value) } : d))
              }
            />
          </Field>
          <Field label={t('settings.system.backup.maxAgeHours')} htmlFor="backup-max-age">
            <TextInput
              id="backup-max-age"
              type="number"
              value={draft.maxAgeHours}
              onChange={(e) =>
                setDraft((d) => (d ? { ...d, maxAgeHours: Number(e.target.value) } : d))
              }
            />
          </Field>
          <div className="sm:col-span-3 flex justify-end">
            <Button
              type="submit"
              variant="secondary"
              size="sm"
              disabled={!scheduleDirty || saveSchedule.isPending}
            >
              {t('common.save')}
            </Button>
          </div>
        </form>
      )}

      <div className="mt-5 border-t border-line-soft pt-4">
        <h3 className="mb-3 font-mono text-2xs text-ink-faint">
          {t('settings.system.backup.recentRuns')}
        </h3>
        {runs.isLoading ? (
          <SkeletonRows rows={3} className="p-0" />
        ) : runs.isError || !runs.data ? (
          <ErrorState
            description={t('settings.system.backup.runsError')}
            onRetry={() => runs.refetch()}
          />
        ) : runs.data.runs.length === 0 ? (
          <EmptyState icon={<BoxIcon size={18} />} title={t('settings.system.backup.noRuns')} />
        ) : (
          <ul className="divide-y divide-line-soft">
            {runs.data.runs.map((r) => (
              <li key={r.id} className="flex items-center justify-between gap-4 py-2.5 first:pt-0 last:pb-0">
                <div className="min-w-0">
                  <p className="truncate font-mono text-xs text-ink">{r.filePath}</p>
                  <p className="mt-0.5 font-mono text-2xs text-ink-faint">
                    {(r.sizeBytes / (1024 * 1024)).toFixed(1)} MB · <TimeAgo at={r.createdAt} />
                  </p>
                </div>
                <ToneTag
                  tone={r.triggeredBy === 'manual' ? 'signal' : 'idle'}
                  label={t(`settings.system.backup.triggeredBy.${r.triggeredBy}`)}
                />
              </li>
            ))}
          </ul>
        )}
      </div>
    </Section>
  );
}

function Stat({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: React.ReactNode;
}) {
  return (
    <Panel className="p-4">
      <p className="mb-2 flex items-center gap-1.5 font-mono text-2xs text-ink-faint">
        <span className="text-accent-primary">{icon}</span>
        {label}
      </p>
      {value}
    </Panel>
  );
}

function Loading() {
  return (
    <div>
      <SubHeader title="System" />
      <SkeletonRows rows={6} />
    </div>
  );
}
