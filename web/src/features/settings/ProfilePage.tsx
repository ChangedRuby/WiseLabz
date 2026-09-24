/**
 * Settings → Profile. Edit display name / email, change password (local accounts),
 * and review + revoke active sessions. Available to every authenticated user.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  getGetMeQueryKey,
  getGetMeSessionsQueryKey,
  patchMe,
  postMePassword,
  getMeSessions,
  deleteMeSessionsSessionId,
} from '../../api/generated/me/me';
import { useGetMe } from '../../api/generated/me/me';
import {
  useGetAuthApiKeys,
  postAuthApiKeys,
  deleteAuthApiKeysId,
  getGetAuthApiKeysQueryKey,
} from '../../api/generated/auth/auth';
import { useGetConnectors } from '../../api/generated/connectors/connectors';
import type { Session, ApiKey } from '../../api/model';
import { ApiKeyScope } from '../../api/model/apiKeyScope';
import { UserDigestCadence } from '../../api/model/userDigestCadence';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { ConfirmDialog } from '../../components/ui/ConfirmDialog';
import { Dialog } from '../../components/ui/Dialog';
import { TimeAgo } from '../../components/ui/TimeAgo';
import { ToneTag } from '../../components/ui/ToneTag';
import { toast } from '../../lib/toast';
import { SubHeader, Section, Field, TextInput, Select } from './parts';
import { UserIcon, KeyIcon, CopyIcon } from '../../components/icons';

export function ProfilePage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data: me } = useGetMe();

  const [displayName, setDisplayName] = useState('');
  const [email, setEmail] = useState('');
  const [digestCadence, setDigestCadence] = useState<typeof UserDigestCadence[keyof typeof UserDigestCadence]>('off');
  const [digestTimezone, setDigestTimezone] = useState('');
  // Adjust state during render (React-blessed alternative to a syncing effect):
  // re-seed whenever the query yields a fresh reference, e.g. after an invalidate.
  const [seeded, setSeeded] = useState<typeof me | null>(null);
  if (me && me !== seeded) {
    setSeeded(me);
    setDisplayName(me.displayName ?? '');
    setEmail(me.email ?? '');
    setDigestCadence(me.digestCadence ?? 'off');
    setDigestTimezone(me.digestTimezone || Intl.DateTimeFormat().resolvedOptions().timeZone);
  }

  const saveProfile = useMutation({
    mutationFn: () => patchMe({ displayName, email, digestCadence, digestTimezone }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetMeQueryKey() });
      toast.success(t('settings.profile.saved'));
    },
    onError: () => toast.error(t('settings.profile.saveError')),
  });

  const isLocal = me?.authSource === 'local';
  const dirty = me
    ? displayName !== (me.displayName ?? '') ||
      email !== (me.email ?? '') ||
      digestCadence !== (me.digestCadence ?? 'off') ||
      digestTimezone !== (me.digestTimezone ?? '')
    : false;

  return (
    <div>
      <SubHeader title={t('settings.profile.title')} description={t('settings.profile.subtitle')} />

      <Section
        title={t('settings.profile.detailsTitle')}
        description={me?.authSource === 'oidc' ? t('settings.profile.oidcNote') : undefined}
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label={t('settings.profile.username')}>
            <TextInput value={me?.username ?? ''} readOnly disabled />
          </Field>
          <Field label={t('settings.profile.role')}>
            <div className="flex h-9.5 items-center">
              <ToneTag tone={me?.role === 'operator' ? 'signal' : 'idle'} label={me?.role ?? '—'} />
            </div>
          </Field>
          <Field label={t('settings.profile.displayName')} htmlFor="profile-name">
            <TextInput
              id="profile-name"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
            />
          </Field>
          <Field label={t('settings.profile.email')} htmlFor="profile-email">
            <TextInput
              id="profile-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </Field>
          <Field label={t('settings.profile.digestCadence')} htmlFor="profile-digest-cadence">
            <Select
              id="profile-digest-cadence"
              value={digestCadence}
              onChange={(e) => setDigestCadence(e.target.value as typeof UserDigestCadence[keyof typeof UserDigestCadence])}
            >
              <option value={UserDigestCadence.off}>{t('settings.profile.digestCadenceOff')}</option>
              <option value={UserDigestCadence.daily}>{t('settings.profile.digestCadenceDaily')}</option>
              <option value={UserDigestCadence.weekly}>{t('settings.profile.digestCadenceWeekly')}</option>
            </Select>
          </Field>
          <Field
            label={t('settings.profile.digestTimezone')}
            htmlFor="profile-digest-timezone"
            hint={t('settings.profile.digestTimezoneHint')}
          >
            <TextInput
              id="profile-digest-timezone"
              value={digestTimezone}
              onChange={(e) => setDigestTimezone(e.target.value)}
            />
          </Field>
        </div>
        <div className="mt-4 flex justify-end">
          <Button
            variant="primary"
            size="sm"
            disabled={!dirty || saveProfile.isPending}
            onClick={() => saveProfile.mutate()}
          >
            {t('common.save')}
          </Button>
        </div>
      </Section>

      {isLocal && <ChangePassword />}

      <ApiKeysSection />

      <SessionsSection />
    </div>
  );
}

function ApiKeysSection() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetAuthApiKeys();

  const [name, setName] = useState('');
  const [expiresAt, setExpiresAt] = useState('');
  const [scope, setScope] = useState<ApiKeyScope>(ApiKeyScope.full);
  const [connectorIds, setConnectorIds] = useState<string[]>([]);
  const [created, setCreated] = useState<string | null>(null);
  const [revoke, setRevoke] = useState<ApiKey | null>(null);
  const { data: connectors } = useGetConnectors();
  const connectorName = (id: string) => connectors?.find((c) => c.id === id)?.name ?? id;

  const create = useMutation({
    mutationFn: () =>
      postAuthApiKeys({
        name,
        expiresAt: expiresAt ? new Date(`${expiresAt}T23:59:59.999Z`).toISOString() : undefined,
        scope,
        connectorIds: connectorIds.length > 0 ? connectorIds : undefined,
      }),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: getGetAuthApiKeysQueryKey() });
      setCreated(result.token);
      setName('');
      setExpiresAt('');
      setScope(ApiKeyScope.full);
      setConnectorIds([]);
    },
    onError: () => toast.error(t('settings.profile.apiKeys.createError')),
  });

  const doRevoke = useMutation({
    mutationFn: (id: string) => deleteAuthApiKeysId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetAuthApiKeysQueryKey() });
      toast.success(t('settings.profile.apiKeys.revoked'));
      setRevoke(null);
    },
    onError: () => toast.error(t('settings.profile.apiKeys.revokeError')),
  });

  return (
    <Section
      title={t('settings.profile.apiKeys.title')}
      description={t('settings.profile.apiKeys.subtitle')}
    >
      <form
        className="grid gap-4 sm:grid-cols-3"
        onSubmit={(e) => {
          e.preventDefault();
          if (name.trim() && !create.isPending) create.mutate();
        }}
      >
        <Field label={t('settings.profile.apiKeys.name')} htmlFor="apikey-name">
          <TextInput
            id="apikey-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t('settings.profile.apiKeys.namePlaceholder')}
          />
        </Field>
        <Field label={t('settings.profile.apiKeys.expiresAt')} htmlFor="apikey-expires">
          <TextInput
            id="apikey-expires"
            type="date"
            value={expiresAt}
            onChange={(e) => setExpiresAt(e.target.value)}
          />
        </Field>
        <Field
          label={t('settings.profile.apiKeys.scope')}
          htmlFor="apikey-scope"
          hint={t('settings.profile.apiKeys.scopeHint')}
        >
          <Select
            id="apikey-scope"
            value={scope}
            onChange={(e) => setScope(e.target.value as ApiKeyScope)}
          >
            <option value={ApiKeyScope.full}>{t('settings.profile.apiKeys.scopeFull')}</option>
            <option value={ApiKeyScope.read}>{t('settings.profile.apiKeys.scopeRead')}</option>
          </Select>
        </Field>
        {connectors && connectors.length > 0 && (
          <fieldset className="sm:col-span-3">
            <legend className="mb-1.5 block font-mono text-2xs text-ink-faint">
              {t('settings.profile.apiKeys.connectors')}
            </legend>
            <div className="flex flex-wrap gap-x-4 gap-y-1.5">
              {connectors.map((c) => (
                <label key={c.id} className="flex items-center gap-1.5 text-xs text-ink">
                  <input
                    type="checkbox"
                    checked={connectorIds.includes(c.id)}
                    onChange={(e) =>
                      setConnectorIds((ids) =>
                        e.target.checked ? [...ids, c.id] : ids.filter((id) => id !== c.id),
                      )
                    }
                  />
                  {c.name}
                </label>
              ))}
            </div>
            <span className="mt-1 block text-2xs leading-relaxed text-ink-faint">
              {t('settings.profile.apiKeys.connectorsHint')}
            </span>
          </fieldset>
        )}
        <div className="flex items-end">
          <Button
            type="submit"
            variant="primary"
            size="sm"
            disabled={!name.trim() || create.isPending}
          >
            {t('settings.profile.apiKeys.create')}
          </Button>
        </div>
      </form>

      <div className="mt-5 border-t border-line-soft pt-4">
        {isLoading ? (
          <SkeletonRows rows={2} className="p-0" />
        ) : isError || !data ? (
          <ErrorState
            description={t('settings.profile.apiKeys.loadError')}
            onRetry={() => refetch()}
          />
        ) : data.length === 0 ? (
          <EmptyState icon={<KeyIcon size={18} />} title={t('settings.profile.apiKeys.empty')} />
        ) : (
          <ul className="divide-y divide-line-soft">
            {data.map((k) => (
              <li key={k.id} className="flex items-center gap-4 py-3 first:pt-0 last:pb-0">
                <div className="min-w-0 flex-1">
                  <p className="flex items-center gap-2 text-sm text-ink">
                    <span className="truncate">{k.name}</span>
                    {k.scope === ApiKeyScope.read && (
                      <ToneTag tone="idle" label={t('settings.profile.apiKeys.readOnlyLabel')} />
                    )}
                    {k.revokedAt && (
                      <ToneTag tone="err" label={t('settings.profile.apiKeys.revokedLabel')} />
                    )}
                  </p>
                  {k.connectorIds.length > 0 && (
                    <p className="mt-0.5 truncate text-2xs text-ink-faint">
                      {t('settings.profile.apiKeys.limitedTo', {
                        connectors: k.connectorIds.map(connectorName).join(', '),
                      })}
                    </p>
                  )}
                  <p className="mt-0.5 flex items-center gap-1.5 font-mono text-2xs text-ink-faint">
                    <span>
                      {t('settings.profile.apiKeys.created')} <TimeAgo at={k.createdAt} />
                    </span>
                    <span>·</span>
                    <span>
                      {k.lastUsedAt ? (
                        <>
                          {t('settings.profile.apiKeys.lastUsed')} <TimeAgo at={k.lastUsedAt} />
                        </>
                      ) : (
                        t('settings.profile.apiKeys.neverUsed')
                      )}
                    </span>
                    {k.expiresAt && (
                      <>
                        <span>·</span>
                        <span>
                          {t('settings.profile.apiKeys.expires')} <TimeAgo at={k.expiresAt} />
                        </span>
                      </>
                    )}
                  </p>
                </div>
                {!k.revokedAt && (
                  <Button
                    variant="danger"
                    size="sm"
                    disabled={doRevoke.isPending}
                    onClick={() => setRevoke(k)}
                  >
                    {t('settings.profile.apiKeys.revoke')}
                  </Button>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>

      <ConfirmDialog
        open={revoke !== null}
        onClose={() => setRevoke(null)}
        onConfirm={() => revoke && doRevoke.mutate(revoke.id)}
        tone="danger"
        title={t('settings.profile.apiKeys.revokeTitle')}
        description={t('settings.profile.apiKeys.revokeConfirm')}
        confirmLabel={t('settings.profile.apiKeys.revoke')}
        cancelLabel={t('common.cancel')}
        confirmDisabled={doRevoke.isPending}
      />

      <NewApiKeyDialog token={created} onClose={() => setCreated(null)} />
    </Section>
  );
}

function NewApiKeyDialog({ token, onClose }: { token: string | null; onClose: () => void }) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    if (!token) return;
    await navigator.clipboard.writeText(token);
    setCopied(true);
  };

  return (
    <Dialog
      open={token !== null}
      onClose={() => {
        setCopied(false);
        onClose();
      }}
      title={t('settings.profile.apiKeys.createdTitle')}
      size="sm"
    >
      <p className="text-sm leading-relaxed text-ink-muted">
        {t('settings.profile.apiKeys.createdDesc')}
      </p>
      <div className="mt-3 flex items-center gap-2 rounded-md border border-line-soft bg-canvas-sunken p-2">
        <code className="flex-1 overflow-x-auto whitespace-nowrap font-mono text-xs text-ink">
          {token}
        </code>
        <Button variant="ghost" size="sm" onClick={() => void copy()}>
          <CopyIcon size={14} />
          {copied ? t('settings.profile.apiKeys.copied') : t('settings.profile.apiKeys.copy')}
        </Button>
      </div>
      <div className="mt-5 flex justify-end">
        <Button
          variant="primary"
          size="sm"
          onClick={() => {
            setCopied(false);
            onClose();
          }}
        >
          {t('settings.profile.apiKeys.done')}
        </Button>
      </div>
    </Dialog>
  );
}

function ChangePassword() {
  const { t } = useTranslation();
  const [current, setCurrent] = useState('');
  const [next, setNext] = useState('');
  const [confirm, setConfirm] = useState('');

  const change = useMutation({
    mutationFn: () => postMePassword({ currentPassword: current, newPassword: next }),
    onSuccess: () => {
      toast.success(t('settings.profile.passwordChanged'));
      setCurrent('');
      setNext('');
      setConfirm('');
    },
    onError: () => toast.error(t('settings.profile.passwordError')),
  });

  const mismatch = next.length > 0 && confirm.length > 0 && next !== confirm;
  const canSubmit = current.length > 0 && next.length >= 8 && next === confirm && !change.isPending;

  return (
    <Section
      title={t('settings.profile.passwordTitle')}
      description={t('settings.profile.passwordDesc')}
    >
      <form
        className="grid gap-4 sm:grid-cols-3"
        onSubmit={(e) => {
          e.preventDefault();
          if (canSubmit) change.mutate();
        }}
      >
        <Field label={t('settings.profile.currentPassword')} htmlFor="pw-current">
          <TextInput
            id="pw-current"
            type="password"
            autoComplete="current-password"
            value={current}
            onChange={(e) => setCurrent(e.target.value)}
          />
        </Field>
        <Field label={t('settings.profile.newPassword')} htmlFor="pw-new">
          <TextInput
            id="pw-new"
            type="password"
            autoComplete="new-password"
            value={next}
            onChange={(e) => setNext(e.target.value)}
          />
        </Field>
        <Field
          label={t('settings.profile.confirmPassword')}
          htmlFor="pw-confirm"
          hint={mismatch ? t('settings.profile.passwordMismatch') : undefined}
        >
          <TextInput
            id="pw-confirm"
            type="password"
            autoComplete="new-password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
          />
        </Field>
        <div className="sm:col-span-3 flex justify-end">
          <Button type="submit" variant="primary" size="sm" disabled={!canSubmit}>
            {t('settings.profile.updatePassword')}
          </Button>
        </div>
      </form>
    </Section>
  );
}

function SessionsSection() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: getGetMeSessionsQueryKey(),
    queryFn: () => getMeSessions(),
  });
  const [revoke, setRevoke] = useState<Session | null>(null);

  const doRevoke = useMutation({
    mutationFn: (id: string) => deleteMeSessionsSessionId(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetMeSessionsQueryKey() });
      toast.success(t('settings.profile.sessionRevoked'));
      setRevoke(null);
    },
    onError: () => toast.error(t('settings.profile.sessionRevokeError')),
  });

  return (
    <Section
      title={t('settings.profile.sessionsTitle')}
      description={t('settings.profile.sessionsDesc')}
    >
      {isLoading ? (
        <SkeletonRows rows={3} className="p-0" />
      ) : isError || !data ? (
        <ErrorState description={t('settings.profile.sessionsError')} onRetry={() => refetch()} />
      ) : data.length === 0 ? (
        <EmptyState icon={<UserIcon size={18} />} title={t('settings.profile.noSessions')} />
      ) : (
        <ul className="divide-y divide-line-soft">
          {data.map((s) => (
            <li key={s.id} className="flex items-center gap-4 py-3 first:pt-0 last:pb-0">
              <div className="min-w-0 flex-1">
                <p className="flex items-center gap-2 text-sm text-ink">
                  <span className="truncate">
                    {s.userAgent ?? t('settings.profile.unknownDevice')}
                  </span>
                  {s.current && <ToneTag tone="ok" label={t('common.current')} />}
                </p>
                <p className="mt-0.5 flex items-center gap-1.5 font-mono text-2xs text-ink-faint">
                  <span>{s.ip ?? '—'}</span>
                  <span>·</span>
                  <span>
                    {t('settings.profile.lastSeen')} <TimeAgo at={s.lastSeenAt} />
                  </span>
                </p>
              </div>
              {!s.current && (
                <Button
                  variant="danger"
                  size="sm"
                  disabled={doRevoke.isPending}
                  onClick={() => setRevoke(s)}
                >
                  {t('settings.profile.revoke')}
                </Button>
              )}
            </li>
          ))}
        </ul>
      )}

      <ConfirmDialog
        open={revoke !== null}
        onClose={() => setRevoke(null)}
        onConfirm={() => revoke && doRevoke.mutate(revoke.id)}
        tone="danger"
        title={t('settings.profile.revokeTitle')}
        description={t('settings.profile.revokeConfirm')}
        confirmLabel={t('settings.profile.revoke')}
        cancelLabel={t('common.cancel')}
        confirmDisabled={doRevoke.isPending}
      />
    </Section>
  );
}
