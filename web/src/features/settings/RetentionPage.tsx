/**
 * Settings → Retention (operator-only). Data-retention cleanup: how long
 * snapshots, doc revisions, alerts, sync run history, and the audit log are
 * kept before the scheduled cleanup job purges them, plus the cron schedule
 * that job runs on. DB-backed (retention_settings, id='default'), mirroring
 * AuthPage's save/dirty/toast pattern.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetSystemSettingsRetention,
  putSystemSettingsRetention,
  getGetSystemSettingsRetentionQueryKey,
} from '../../api/generated/system/system';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState } from '../../components/ui/states';
import { toast } from '../../lib/toast';
import { SubHeader, Section, Field, TextInput } from './parts';

interface FormState {
  snapshotDays: string;
  docVersionDays: string;
  alertDays: string;
  syncRunDays: string;
  auditDays: string;
  healthCheckDays: string;
  cronExpr: string;
}

function toForm(data: {
  snapshotDays: number;
  docVersionDays: number;
  alertDays: number;
  syncRunDays: number;
  auditDays: number;
  healthCheckDays: number;
  cronExpr: string;
}): FormState {
  return {
    snapshotDays: String(data.snapshotDays),
    docVersionDays: String(data.docVersionDays),
    alertDays: String(data.alertDays),
    syncRunDays: String(data.syncRunDays),
    auditDays: String(data.auditDays),
    healthCheckDays: String(data.healthCheckDays),
    cronExpr: data.cronExpr,
  };
}

export function RetentionPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetSystemSettingsRetention();

  const [form, setForm] = useState<FormState | null>(null);
  // Adjust state during render: re-seed whenever the query yields a fresh
  // reference (initial load, or after an invalidate following a save).
  const [seeded, setSeeded] = useState<typeof data | null>(null);
  if (data && data !== seeded) {
    setSeeded(data);
    setForm(toForm(data));
  }

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: getGetSystemSettingsRetentionQueryKey() });

  const save = useMutation({
    mutationFn: (body: Parameters<typeof putSystemSettingsRetention>[0]) =>
      putSystemSettingsRetention(body),
    onSuccess: () => {
      invalidate();
      toast.success(t('settings.retention.saved'));
    },
    onError: () => toast.error(t('settings.retention.saveError')),
  });

  if (isLoading || !form) return <Loading />;
  if (isError || !data)
    return (
      <div>
        <SubHeader title={t('settings.retention.title')} />
        <ErrorState description={t('settings.retention.loadError')} onRetry={() => refetch()} />
      </div>
    );

  const dirty =
    form.snapshotDays !== String(data.snapshotDays) ||
    form.docVersionDays !== String(data.docVersionDays) ||
    form.alertDays !== String(data.alertDays) ||
    form.syncRunDays !== String(data.syncRunDays) ||
    form.auditDays !== String(data.auditDays) ||
    form.healthCheckDays !== String(data.healthCheckDays) ||
    form.cronExpr !== data.cronExpr;

  // Empty or non-numeric day fields must not silently coerce to 0 (which
  // would disable that retention category) when saved.
  const isValidDays = (v: string) => v.trim() !== '' && Number.isInteger(Number(v)) && Number(v) >= 0;
  const daysValid =
    isValidDays(form.snapshotDays) &&
    isValidDays(form.docVersionDays) &&
    isValidDays(form.alertDays) &&
    isValidDays(form.syncRunDays) &&
    isValidDays(form.auditDays) &&
    isValidDays(form.healthCheckDays);

  const set = (field: keyof FormState) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => (f ? { ...f, [field]: e.target.value } : f));

  return (
    <div>
      <SubHeader
        title={t('settings.retention.title')}
        description={t('settings.retention.subtitle')}
      />

      <Section
        title={t('settings.retention.daysTitle')}
        description={t('settings.retention.daysDesc')}
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label={t('settings.retention.snapshotDays')} htmlFor="ret-snapshot-days">
            <TextInput
              id="ret-snapshot-days"
              type="number"
              min={0}
              value={form.snapshotDays}
              onChange={set('snapshotDays')}
            />
          </Field>
          <Field label={t('settings.retention.docVersionDays')} htmlFor="ret-doc-version-days">
            <TextInput
              id="ret-doc-version-days"
              type="number"
              min={0}
              value={form.docVersionDays}
              onChange={set('docVersionDays')}
            />
          </Field>
          <Field label={t('settings.retention.alertDays')} htmlFor="ret-alert-days">
            <TextInput
              id="ret-alert-days"
              type="number"
              min={0}
              value={form.alertDays}
              onChange={set('alertDays')}
            />
          </Field>
          <Field label={t('settings.retention.syncRunDays')} htmlFor="ret-sync-run-days">
            <TextInput
              id="ret-sync-run-days"
              type="number"
              min={0}
              value={form.syncRunDays}
              onChange={set('syncRunDays')}
            />
          </Field>
          <Field label={t('settings.retention.auditDays')} htmlFor="ret-audit-days">
            <TextInput
              id="ret-audit-days"
              type="number"
              min={0}
              value={form.auditDays}
              onChange={set('auditDays')}
            />
          </Field>
          <Field label={t('settings.retention.healthCheckDays')} htmlFor="ret-health-check-days">
            <TextInput
              id="ret-health-check-days"
              type="number"
              min={0}
              value={form.healthCheckDays}
              onChange={set('healthCheckDays')}
            />
          </Field>
          <Field
            label={t('settings.retention.cronExpr')}
            htmlFor="ret-cron-expr"
            hint={t('settings.retention.cronHint')}
          >
            <TextInput
              id="ret-cron-expr"
              value={form.cronExpr}
              onChange={set('cronExpr')}
              className="font-mono"
            />
          </Field>
        </div>
        <div className="mt-4 flex justify-end">
          <Button
            variant="primary"
            size="sm"
            disabled={!dirty || !daysValid || save.isPending}
            onClick={() =>
              save.mutate({
                snapshotDays: Number(form.snapshotDays),
                docVersionDays: Number(form.docVersionDays),
                alertDays: Number(form.alertDays),
                syncRunDays: Number(form.syncRunDays),
                auditDays: Number(form.auditDays),
                healthCheckDays: Number(form.healthCheckDays),
                cronExpr: form.cronExpr,
              })
            }
          >
            {t('common.save')}
          </Button>
        </div>
      </Section>
    </div>
  );
}

function Loading() {
  return (
    <div>
      <SubHeader title="Retention" />
      <SkeletonRows rows={6} />
    </div>
  );
}
