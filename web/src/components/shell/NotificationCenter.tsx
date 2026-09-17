/**
 * Notification bell + dropdown: persisted in-app notifications (unread state,
 * mark-read, deep link to the affected resource). Lives in the Topbar's bell
 * slot; live updates ride the existing WS `alert.created`/`alert.resolved`
 * invalidation of the notifications query in WebSocketProvider.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import { AnimatePresence, motion } from 'motion/react';
import {
  getGetNotificationsQueryKey,
  useGetNotifications,
  postNotificationsNotificationIdRead,
  postNotificationsReadAll,
} from '../../api/generated/notifications/notifications';
import type { Notification } from '../../api/model';
import { IconButton } from '../ui/Button';
import { TimeAgo } from '../ui/TimeAgo';
import { BellIcon, CheckIcon } from '../icons';
import { navigateTo } from '../../lib/navigation';
import { cn } from '../../lib/cn';

export function NotificationCenter() {
  const { t } = useTranslation();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);

  // ponytail: flat recent-N list, no infinite scroll — add pagination if the
  // panel needs more than the latest page.
  const list = useGetNotifications({ pageSize: 20 });
  const unread = useGetNotifications({ unread: true, pageSize: 1 });
  const unreadCount = unread.data?.total ?? 0;

  const invalidate = () => {
    qc.invalidateQueries({ queryKey: getGetNotificationsQueryKey() });
  };

  const markRead = async (n: Notification) => {
    if (!n.read) {
      await postNotificationsNotificationIdRead(n.id);
      invalidate();
    }
    if (n.alertId) {
      navigateTo('/alerts');
      setOpen(false);
    } else if (n.eventType?.startsWith('finding.')) {
      navigateTo('/findings');
      setOpen(false);
    }
  };

  const markAllRead = async () => {
    await postNotificationsReadAll();
    invalidate();
  };

  const items = list.data?.items ?? [];

  return (
    <div className="relative">
      <IconButton
        label={t('notifications.title')}
        onClick={() => setOpen((v) => !v)}
        className="relative"
      >
        <BellIcon size={18} />
        <AnimatePresence>
          {unreadCount > 0 && (
            <motion.span
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              exit={{ scale: 0 }}
              transition={{ type: 'spring', stiffness: 600, damping: 22 }}
              className="nums absolute -right-0.5 -top-0.5 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-err px-1 text-[10px] font-bold text-accent-primary-ink"
            >
              {unreadCount}
            </motion.span>
          )}
        </AnimatePresence>
      </IconButton>

      <AnimatePresence>
        {open && (
          <>
            <div className="fixed inset-0 z-(--z-dropdown)" onClick={() => setOpen(false)} />
            <motion.div
              initial={{ opacity: 0, y: -6, scale: 0.98 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -6, scale: 0.98 }}
              transition={{ duration: 0.16, ease: [0.16, 1, 0.3, 1] }}
              className="absolute right-0 top-[calc(100%+8px)] z-(--z-dropdown) w-80 overflow-hidden rounded-sm border border-line bg-surface-overlay shadow-(--shadow-pop)"
            >
              <div className="flex items-center justify-between border-b border-line-soft px-3 py-2.5">
                <p className="text-sm font-medium text-ink">{t('notifications.title')}</p>
                {unreadCount > 0 && (
                  <button
                    onClick={() => void markAllRead()}
                    className="flex items-center gap-1 text-xs text-ink-faint transition-colors hover:text-ink"
                  >
                    <CheckIcon size={13} />
                    {t('notifications.markAllRead')}
                  </button>
                )}
              </div>

              <div className="max-h-96 overflow-y-auto p-1">
                {items.length === 0 && (
                  <p className="px-2.5 py-6 text-center text-sm text-ink-faint">
                    {t('notifications.empty')}
                  </p>
                )}
                {items.map((n) => {
                  const willNavigate = n.alertId || n.eventType?.startsWith('finding.');
                  return (
                  <button
                    key={n.id}
                    onClick={() => void markRead(n)}
                    className={cn(
                      'flex w-full flex-col items-start gap-0.5 rounded-md px-2.5 py-2 text-left transition-colors',
                      willNavigate && 'hover:bg-surface-raised',
                      !n.read && 'bg-accent-primary-tint/40',
                    )}
                  >
                    <div className="flex w-full items-center gap-1.5">
                      {!n.read && <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-accent-primary" />}
                      <span
                        className={cn(
                          'flex-1 truncate text-sm',
                          n.read ? 'text-ink-muted' : 'font-medium text-ink',
                        )}
                      >
                        {n.title}
                      </span>
                      <TimeAgo at={n.createdAt} className="shrink-0 text-2xs text-ink-faint" />
                    </div>
                    {n.message && (
                      <p className="line-clamp-2 text-xs text-ink-faint">{n.message}</p>
                    )}
                  </button>
                  );
                })}
              </div>
            </motion.div>
          </>
        )}
      </AnimatePresence>
    </div>
  );
}
