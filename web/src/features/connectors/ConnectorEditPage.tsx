/**
 * Edit an existing connector (`/connectors/:id/edit`). The type is fixed; the
 * schema-driven fields are reused from the add flow (<Field/>), prefilled from the
 * connector read. Credentials are write-only over the API — secret fields come back
 * empty and are only sent when re-entered (credential rotation). Re-test and save
 * are operator actions, enforced server-side.
 */
import { useMemo, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetConnectorsConnectorId,
  useGetConnectorsSchema,
  putConnectorsConnectorId,
  postConnectorsConnectorIdTest,
  getGetConnectorsQueryKey,
} from '../../api/generated/connectors/connectors';
import { Field } from './ConnectorForm';
import { ConnectorPermissionsTab } from './ConnectorPermissionsTab';
import { Button } from '../../components/ui/Button';
import { Panel } from '../../components/ui/Panel';
import { SkeletonRows, ErrorState } from '../../components/ui/states';
import { toast } from '../../lib/toast';
import { ArrowRightIcon, CheckIcon } from '../../components/icons';
import { useConnectorRole } from '../../hooks/useRole';

type FormValues = Record<string, string | boolean>;

export function ConnectorEditPage() {
  const { t } = useTranslation();
  const { id = '' } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const connector = useGetConnectorsConnectorId(id);
  const { data: schemas } = useGetConnectorsSchema();
  const canEdit = useConnectorRole(id) === 'operator';

  const schema = useMemo(
    () => schemas?.find((s) => s.type === connector.data?.type) ?? null,
    [schemas, connector.data?.type],
  );

  const [name, setName] = useState<string | null>(null);
  const [owner, setOwner] = useState<string | null>(null);
  const [values, setValues] = useState<FormValues>({});
  // Prefilled values fall back to the connector read until the user edits a field.
  const nameValue = name ?? connector.data?.name ?? '';
  const ownerValue = owner ?? connector.data?.owner ?? '';

  const test = useMutation({
    mutationFn: () => postConnectorsConnectorIdTest(id),
    onSuccess: (r) => toast.success(t('connectors.edit.testOk', { ms: r.latencyMs ?? 0 })),
    onError: () => toast.error(t('connectors.edit.testFail')),
  });

  const save = useMutation({
    mutationFn: () => {
      const config: Record<string, unknown> = {};
      for (const f of schema?.fields ?? []) {
        if (f.name === 'url' || f.name === 'verifyTls') continue;
        // Only send secret/config fields the user actually re-entered.
        if (values[f.name] !== undefined && String(values[f.name]).length > 0) config[f.name] = values[f.name];
      }
      return putConnectorsConnectorId(id, {
        name: nameValue,
        owner: ownerValue,
        url: values.url !== undefined ? String(values.url) : connector.data?.url,
        verifyTls: values.verifyTls !== undefined ? Boolean(values.verifyTls) : connector.data?.verifyTls,
        config,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetConnectorsQueryKey() });
      toast.success(t('connectors.edit.saved'));
      navigate(`/services/${id}`);
    },
    onError: () => toast.error(t('connectors.edit.saveError')),
  });

  if (connector.isLoading) {
    return (
      <div className="mx-auto max-w-170 px-6 py-6">
        <Panel className="p-6">
          <SkeletonRows rows={5} />
        </Panel>
      </div>
    );
  }
  if (connector.isError || !connector.data) {
    return (
      <div className="mx-auto max-w-170 px-6 py-6">
        <Panel className="min-h-[30vh]">
          <ErrorState description={t('services.detail.notFound')} onRetry={() => connector.refetch()} />
        </Panel>
      </div>
    );
  }

  const c = connector.data;

  if (!canEdit) {
    return (
      <div className="mx-auto max-w-170 px-6 py-6">
        <Panel className="min-h-[30vh]">
          <ErrorState
            title={t('connectors.edit.noAccessTitle')}
            description={t('connectors.edit.noAccessDesc')}
          />
        </Panel>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-170 px-6 py-6">
      <button
        onClick={() => navigate(`/services/${id}`)}
        className="mb-4 inline-flex items-center gap-1.5 text-xs text-ink-muted transition-colors hover:text-ink"
      >
        <ArrowRightIcon size={13} className="rotate-180" />
        {t('connectors.edit.back')}
      </button>

      <header className="mb-5">
        <h1 className="text-xl font-semibold tracking-tight text-ink">{t('connectors.edit.title', { name: c.name })}</h1>
        <p className="text-sm text-ink-muted">{t('connectors.edit.subtitle')}</p>
      </header>

      <Panel className="p-5">
        <div className="space-y-3">
          <Field
            field={{ name: 'name', label: t('connectors.displayName'), kind: 'string', required: true }}
            value={nameValue}
            onChange={(v) => setName(String(v))}
          />
          <Field
            field={{ name: 'owner', label: t('connectors.owner'), kind: 'string', required: false }}
            value={ownerValue}
            onChange={(v) => setOwner(String(v))}
          />
          {schema?.fields.map((f) => (
            <Field
              key={f.name}
              field={f.secret ? { ...f, placeholder: t('connectors.edit.secretPlaceholder') } : f}
              value={
                values[f.name] ??
                (f.name === 'url'
                  ? (c.url ?? '')
                  : f.name === 'verifyTls'
                    ? Boolean(c.verifyTls)
                    : f.kind === 'boolean'
                      ? false
                      : '')
              }
              onChange={(v) => setValues((s) => ({ ...s, [f.name]: v }))}
            />
          ))}
        </div>

        <div className="mt-5 flex items-center justify-between gap-2">
          <Button
            variant="secondary"
            size="md"
            onClick={() => test.mutate()}
            disabled={test.isPending || schema?.stub}
            title={schema?.stub ? t('connectors.edit.testUnavailable') : undefined}
          >
            {test.isPending ? t('connectors.edit.testing') : t('connectors.edit.test')}
          </Button>
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="md" onClick={() => navigate(`/services/${id}`)}>
              {t('common.cancel')}
            </Button>
            <Button variant="primary" size="md" onClick={() => save.mutate()} disabled={!nameValue || save.isPending}>
              <CheckIcon size={15} />
              {save.isPending ? t('connectors.edit.saving') : t('common.save')}
            </Button>
          </div>
        </div>
      </Panel>

      <ConnectorPermissionsTab connectorId={id} />
    </div>
  );
}
