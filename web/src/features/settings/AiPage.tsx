/**
 * Settings → AI (operator-only). Configure the opt-in, provider-agnostic AI doc
 * module and run a connection test. The API key is write-only: it's never returned
 * by the API, so the field stays blank on load and is only sent when re-entered.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetAiConfig,
  putAiConfig,
  postAiConfigTest,
  getGetAiConfigQueryKey,
  useGetAiConfigFallbackProviders,
  putAiConfigFallbackProviders,
  getGetAiConfigFallbackProvidersQueryKey,
} from '../../api/generated/settings/settings';
import { AiConfigMode, AiConfigProvider, AiFallbackProviderProvider } from '../../api/model';
import type { AiConfig, AiFallbackProvider, TestResult } from '../../api/model';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState } from '../../components/ui/states';
import { ToneTag } from '../../components/ui/ToneTag';
import { toast } from '../../lib/toast';
import { PlusIcon, XIcon, ArrowRightIcon } from '../../components/icons';
import { SubHeader, Section, Field, TextInput, Select, ToggleRow } from './parts';

export function AiPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetAiConfig();

  const [form, setForm] = useState<AiConfig | null>(null);
  const [apiKey, setApiKey] = useState('');
  const [testResult, setTestResult] = useState<TestResult | null>(null);
  // Seed the editable form from server data, re-seeding whenever the query
  // returns a fresh reference (e.g. after an invalidate). Adjusting state during
  // render is the React-blessed alternative to a syncing effect.
  const [seeded, setSeeded] = useState<AiConfig | null>(null);
  if (data && data !== seeded) {
    setSeeded(data);
    setForm({ ...data });
  }

  const save = useMutation({
    mutationFn: (body: AiConfig) => putAiConfig(apiKey ? { ...body, apiKey } : body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetAiConfigQueryKey() });
      toast.success(t('settings.ai.saved'));
      setApiKey('');
    },
    onError: () => toast.error(t('settings.ai.saveError')),
  });

  const test = useMutation({
    mutationFn: () => postAiConfigTest(),
    onSuccess: (result) => {
      setTestResult(result);
      if (result.ok) toast.success(t('settings.ai.testOk'));
      else toast.error(result.message ?? t('settings.ai.testFail'));
    },
    onError: () => toast.error(t('settings.ai.testFail')),
  });

  if (isLoading) return <Loading />;
  if (isError || !data || !form)
    return (
      <div>
        <SubHeader title={t('settings.ai.title')} />
        <ErrorState description={t('settings.ai.loadError')} onRetry={() => refetch()} />
      </div>
    );

  const set = <K extends keyof AiConfig>(key: K, value: AiConfig[K]) =>
    setForm((f) => (f ? { ...f, [key]: value } : f));
  const isOllama = form.provider === AiConfigProvider.ollama;

  return (
    <div>
      <SubHeader title={t('settings.ai.title')} description={t('settings.ai.subtitle')} />

      <Section title={t('settings.ai.moduleTitle')}>
        <ToggleRow
          title={t('settings.ai.enableTitle')}
          description={t('settings.ai.enableDesc')}
          checked={form.enabled}
          onChange={(enabled) => set('enabled', enabled)}
        />
      </Section>

      <Section title={t('settings.ai.providerTitle')}>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label={t('settings.ai.provider')} htmlFor="ai-provider">
            <Select
              id="ai-provider"
              value={form.provider ?? ''}
              disabled={!form.enabled}
              onChange={(e) => set('provider', (e.target.value || null) as AiConfig['provider'])}
            >
              <option value="">{t('common.none')}</option>
              <option value={AiConfigProvider.anthropic}>Anthropic</option>
              <option value={AiConfigProvider.openai}>OpenAI</option>
              <option value={AiConfigProvider.ollama}>Ollama</option>
            </Select>
          </Field>
          <Field label={t('settings.ai.model')} htmlFor="ai-model">
            <TextInput
              id="ai-model"
              value={form.model ?? ''}
              disabled={!form.enabled}
              placeholder="claude-sonnet-4-5"
              onChange={(e) => set('model', e.target.value)}
            />
          </Field>
          {isOllama ? (
            <Field
              label={t('settings.ai.baseUrl')}
              htmlFor="ai-baseurl"
              hint={t('settings.ai.baseUrlHint')}
            >
              <TextInput
                id="ai-baseurl"
                value={form.baseUrl ?? ''}
                disabled={!form.enabled}
                placeholder="http://localhost:11434"
                onChange={(e) => set('baseUrl', e.target.value)}
              />
            </Field>
          ) : (
            <Field
              label={t('settings.ai.apiKey')}
              htmlFor="ai-apikey"
              hint={t('settings.ai.apiKeyHint')}
            >
              <TextInput
                id="ai-apikey"
                type="password"
                autoComplete="off"
                value={apiKey}
                disabled={!form.enabled}
                placeholder="••••••••••••"
                onChange={(e) => setApiKey(e.target.value)}
              />
            </Field>
          )}
          <Field label={t('settings.ai.mode')} htmlFor="ai-mode" hint={t('settings.ai.modeHint')}>
            <Select
              id="ai-mode"
              value={form.mode}
              disabled={!form.enabled}
              onChange={(e) => set('mode', e.target.value as AiConfig['mode'])}
            >
              <option value={AiConfigMode.suggest_only}>{t('settings.ai.modeSuggest')}</option>
              <option value={AiConfigMode.auto_update}>{t('settings.ai.modeAuto')}</option>
            </Select>
          </Field>
        </div>

        <div className="mt-4 flex items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              disabled={!form.enabled || test.isPending}
              onClick={() => test.mutate()}
            >
              {test.isPending ? t('settings.ai.testing') : t('settings.ai.test')}
            </Button>
            {testResult && (
              <span className="flex items-center gap-2 text-2xs text-ink-muted">
                <ToneTag
                  tone={testResult.ok ? 'ok' : 'err'}
                  label={testResult.ok ? t('settings.ai.ok') : t('settings.ai.failed')}
                />
                {testResult.latencyMs != null && (
                  <span className="font-mono">{testResult.latencyMs}ms</span>
                )}
              </span>
            )}
          </div>
          <Button
            variant="primary"
            size="sm"
            disabled={save.isPending}
            onClick={() => save.mutate(form)}
          >
            {t('common.save')}
          </Button>
        </div>
      </Section>

      <FallbackProvidersSection />
    </div>
  );
}

// Ordered fallback chain (issue #238, piece 3/3): tried in list order when the
// provider before it errors with 429, 5xx, or a timeout. Priority 1 is the
// primary provider configured above; this list starts at priority 2.
function FallbackProvidersSection() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetAiConfigFallbackProviders();

  const [rows, setRows] = useState<AiFallbackProvider[] | null>(null);
  const [seeded, setSeeded] = useState<AiFallbackProvider[] | null>(null);
  if (data && data !== seeded) {
    setSeeded(data);
    setRows(data.map((p) => ({ ...p })));
  }

  const save = useMutation({
    mutationFn: (body: AiFallbackProvider[]) => putAiConfigFallbackProviders(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetAiConfigFallbackProvidersQueryKey() });
      toast.success(t('settings.ai.fallbackSaved'));
    },
    onError: () => toast.error(t('settings.ai.fallbackSaveError')),
  });

  if (isLoading || !rows) return null;
  if (isError)
    return (
      <Section title={t('settings.ai.fallbackTitle')}>
        <ErrorState description={t('settings.ai.fallbackLoadError')} onRetry={() => refetch()} />
      </Section>
    );

  const update = (i: number, patch: Partial<AiFallbackProvider>) =>
    setRows((r) => r && r.map((row, idx) => (idx === i ? { ...row, ...patch } : row)));
  const move = (i: number, dir: -1 | 1) =>
    setRows((r) => {
      if (!r) return r;
      const j = i + dir;
      if (j < 0 || j >= r.length) return r;
      const next = [...r];
      [next[i], next[j]] = [next[j], next[i]];
      return next;
    });

  return (
    <Section title={t('settings.ai.fallbackTitle')} description={t('settings.ai.fallbackDesc')}>
      <div className="flex flex-col gap-3">
        {rows.map((row, i) => (
          <div key={i} className="grid items-end gap-3 sm:grid-cols-[auto_1fr_1fr_1fr_1fr_auto]">
            <div className="flex gap-1">
              <Button
                variant="ghost"
                size="sm"
                disabled={i === 0}
                onClick={() => move(i, -1)}
                aria-label={t('settings.ai.fallbackMoveUp')}
              >
                <ArrowRightIcon size={13} className="-rotate-90" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                disabled={i === rows.length - 1}
                onClick={() => move(i, 1)}
                aria-label={t('settings.ai.fallbackMoveDown')}
              >
                <ArrowRightIcon size={13} className="rotate-90" />
              </Button>
            </div>
            <Field label={t('settings.ai.provider')}>
              <Select
                value={row.provider}
                onChange={(e) => update(i, { provider: e.target.value as AiFallbackProvider['provider'] })}
              >
                <option value={AiFallbackProviderProvider.anthropic}>Anthropic</option>
                <option value={AiFallbackProviderProvider.openai}>OpenAI</option>
                <option value={AiFallbackProviderProvider.ollama}>Ollama</option>
              </Select>
            </Field>
            <Field label={t('settings.ai.model')}>
              <TextInput value={row.model ?? ''} onChange={(e) => update(i, { model: e.target.value })} />
            </Field>
            {row.provider === AiFallbackProviderProvider.ollama ? (
              <Field label={t('settings.ai.baseUrl')}>
                <TextInput
                  value={row.baseUrl ?? ''}
                  placeholder="http://localhost:11434"
                  onChange={(e) => update(i, { baseUrl: e.target.value })}
                />
              </Field>
            ) : (
              <Field label={t('settings.ai.apiKey')}>
                <TextInput
                  type="password"
                  autoComplete="off"
                  placeholder="••••••••••••"
                  onChange={(e) => update(i, { apiKey: e.target.value })}
                />
              </Field>
            )}
            <div />
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setRows((r) => r && r.filter((_, idx) => idx !== i))}
              aria-label={t('settings.ai.fallbackRemove')}
            >
              <XIcon size={13} />
            </Button>
          </div>
        ))}

        <div className="flex items-center justify-between gap-3">
          <Button
            variant="secondary"
            size="sm"
            onClick={() =>
              setRows((r) => [...(r ?? []), { provider: AiFallbackProviderProvider.openai, model: '' }])
            }
          >
            <PlusIcon size={13} />
            {t('settings.ai.fallbackAdd')}
          </Button>
          <Button variant="primary" size="sm" disabled={save.isPending} onClick={() => save.mutate(rows)}>
            {t('common.save')}
          </Button>
        </div>
      </div>
    </Section>
  );
}

function Loading() {
  return (
    <div>
      <SubHeader title="AI" />
      <SkeletonRows rows={6} />
    </div>
  );
}
