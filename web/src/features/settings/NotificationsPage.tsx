/**
 * Settings → Notifications (operator-only). Channel enablement + a per-channel
 * test, the full event×channel routing matrix (EventRoutingTable), and the
 * delivery history (pending/sent/failed attempts, filterable by status).
 * Local draft state; "Save routing" persists the whole config.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetNotificationsConfig,
  putNotificationsConfig,
  postNotificationsConfigTest,
  getGetNotificationsConfigQueryKey,
} from '../../api/generated/settings/settings';
import { useGetNotificationsDeliveries } from '../../api/generated/notifications/notifications';
import type {
  NotificationConfig,
  NotificationChannelType,
  GetNotificationsDeliveriesStatus,
} from '../../api/model';
import { Panel, PanelHeader } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { ToneTag } from '../../components/ui/ToneTag';
import { TimeAgo } from '../../components/ui/TimeAgo';
import { type Tone } from '../../components/ui/status';
import { toast } from '../../lib/toast';
import { SubHeader, Section, ToggleRow, Select } from './parts';
import { EventRoutingTable } from './EventRoutingTable';
import { HistoryIcon } from '../../components/icons';

const PAGE_SIZE = 25;

const deliveryTone: Record<GetNotificationsDeliveriesStatus, Tone> = {
  pending: 'idle',
  sent: 'ok',
  failed: 'err',
};

const CHANNEL_LABELS: Record<NotificationChannelType, string> = {
  in_app: 'In-app',
  smtp: 'Email (SMTP)',
  webhook: 'Webhook',
  discord: 'Discord',
  slack: 'Slack',
  ntfy: 'ntfy',
  telegram: 'Telegram',
};

const CHANNEL_DESC_KEYS: Record<NotificationChannelType, string> = {
  in_app: 'settings.notifications.inAppDesc',
  smtp: 'settings.notifications.smtpDesc',
  webhook: 'settings.notifications.webhookDesc',
  discord: 'settings.notifications.discordDesc',
  slack: 'settings.notifications.slackDesc',
  ntfy: 'settings.notifications.ntfyDesc',
  telegram: 'settings.notifications.telegramDesc',
};

export function NotificationsPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useGetNotificationsConfig();

  const [draft, setDraft] = useState<NotificationConfig | null>(null);
  // Adjust state during render (React-blessed alternative to a syncing effect):
  // re-seed whenever the query yields a fresh reference, e.g. after an invalidate.
  const [seeded, setSeeded] = useState<NotificationConfig | null>(null);
  if (data && data !== seeded) {
    setSeeded(data);
    setDraft({
      channels: data.channels.map((c) => ({ ...c })),
      routing: data.routing.map((r) => ({ ...r })),
    });
  }

  const save = useMutation({
    mutationFn: (body: NotificationConfig) => putNotificationsConfig(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getGetNotificationsConfigQueryKey() });
      toast.success(t('settings.notifications.saved'));
    },
    onError: () => toast.error(t('settings.notifications.saveError')),
  });

  const test = useMutation({
    mutationFn: (channel: NotificationChannelType) => postNotificationsConfigTest({ channel }),
    onSuccess: (result) => {
      if (result.ok) toast.success(result.message ?? t('settings.notifications.testOk'));
      else toast.error(result.message ?? t('settings.notifications.testFail'));
    },
    onError: () => toast.error(t('settings.notifications.testFail')),
  });

  if (isLoading) return <Loading />;
  if (isError || !data || !draft)
    return (
      <div>
        <SubHeader title={t('settings.notifications.title')} />
        <ErrorState description={t('settings.notifications.loadError')} onRetry={() => refetch()} />
      </div>
    );

  const setChannelEnabled = (type: NotificationChannelType, enabled: boolean) =>
    setDraft((d) =>
      d ? { ...d, channels: d.channels.map((c) => (c.type === type ? { ...c, enabled } : c)) } : d
    );

  const dirty = JSON.stringify(draft) !== JSON.stringify(data);

  return (
    <div>
      <SubHeader
        title={t('settings.notifications.title')}
        description={t('settings.notifications.subtitle')}
      />

      <Section title={t('settings.notifications.channelsTitle')}>
        <div className="space-y-3">
          {draft.channels.map((c) => (
            <ToggleRow
              key={c.type}
              title={t(`settings.notifications.channel.${c.type}`, {
                defaultValue: CHANNEL_LABELS[c.type],
              })}
              description={t(CHANNEL_DESC_KEYS[c.type])}
              checked={c.enabled}
              onChange={(enabled) => setChannelEnabled(c.type, enabled)}
            />
          ))}
        </div>
        <div className="mt-3 flex flex-wrap gap-2">
          {draft.channels.map((c) => (
            <Button
              key={c.type}
              variant="ghost"
              size="sm"
              disabled={!c.enabled || test.isPending}
              onClick={() => test.mutate(c.type)}
            >
              {t('settings.notifications.testChannel', {
                channel: t(`settings.notifications.channel.${c.type}`, {
                  defaultValue: CHANNEL_LABELS[c.type],
                }),
              })}
            </Button>
          ))}
        </div>
      </Section>

      <Section
        title={t('settings.notifications.routingTitle')}
        description={t('settings.notifications.routingDesc')}
        action={
          <Button
            variant="primary"
            size="sm"
            disabled={!dirty || save.isPending}
            onClick={() => save.mutate(draft)}
          >
            {t('settings.notifications.saveRouting')}
          </Button>
        }
      >
        <EventRoutingTable config={draft} onChange={setDraft} disabled={save.isPending} />
      </Section>

      <DeliveryHistorySection />
    </div>
  );
}

function DeliveryHistorySection() {
  const { t } = useTranslation();
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<GetNotificationsDeliveriesStatus | ''>('');

  const { data, isLoading, isError, refetch } = useGetNotificationsDeliveries({
    status: status || undefined,
    page,
    pageSize: PAGE_SIZE,
  });

  const totalPages = data ? Math.max(1, Math.ceil(data.total / data.pageSize)) : 1;

  return (
    <Panel>
      <PanelHeader
        title={t('settings.notifications.deliveries.title')}
        icon={<HistoryIcon size={14} />}
        count={data?.total}
        action={
          <Select
            aria-label={t('settings.notifications.deliveries.statusFilter')}
            value={status}
            onChange={(e) => {
              setStatus(e.target.value as GetNotificationsDeliveriesStatus | '');
              setPage(1);
            }}
            className="!w-auto"
          >
            <option value="">{t('settings.notifications.deliveries.anyStatus')}</option>
            <option value="pending">{t('settings.notifications.deliveries.pending')}</option>
            <option value="sent">{t('settings.notifications.deliveries.sent')}</option>
            <option value="failed">{t('settings.notifications.deliveries.failed')}</option>
          </Select>
        }
      />

      {isLoading ? (
        <SkeletonRows rows={6} />
      ) : isError || !data ? (
        <ErrorState
          description={t('settings.notifications.deliveries.loadError')}
          onRetry={() => refetch()}
        />
      ) : data.items.length === 0 ? (
        <EmptyState title={t('settings.notifications.deliveries.empty')} />
      ) : (
        <>
          <ul className="divide-y divide-line-soft">
            {data.items.map((d) => (
              <li key={d.id} className="flex items-center justify-between gap-4 px-4 py-2.5">
                <div className="min-w-0 flex-1">
                  <p className="flex items-center gap-2 text-sm text-ink">
                    <span className="font-mono text-xs text-ink-muted">{d.channel}</span>
                    <ToneTag tone={deliveryTone[d.status]} label={d.status} />
                  </p>
                  {d.status === 'failed' && d.lastError && (
                    <p className="mt-0.5 truncate text-2xs leading-relaxed text-err">
                      {d.lastError}
                    </p>
                  )}
                  <p className="mt-0.5 flex items-center gap-1.5 font-mono text-2xs text-ink-faint">
                    <span>
                      {t('settings.notifications.deliveries.attempts', { count: d.attempts })}
                    </span>
                    <span>·</span>
                    <TimeAgo at={d.createdAt} />
                  </p>
                </div>
              </li>
            ))}
          </ul>
          <div className="flex items-center justify-between gap-3 border-t border-line-soft px-4 py-2.5">
            <span className="font-mono text-2xs text-ink-faint">
              {page} / {totalPages}
            </span>
            <div className="flex gap-2">
              <Button
                variant="secondary"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                {t('common.previous')}
              </Button>
              <Button
                variant="secondary"
                size="sm"
                disabled={page >= totalPages}
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              >
                {t('common.next')}
              </Button>
            </div>
          </div>
        </>
      )}
    </Panel>
  );
}

function Loading() {
  return (
    <div>
      <SubHeader title="Notifications" />
      <SkeletonRows rows={6} />
    </div>
  );
}
