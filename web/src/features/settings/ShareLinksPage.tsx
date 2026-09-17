/**
 * "My Share Links" (#240 PR2) — every read-only doc share link the current
 * user created, with revoke. Row layout mirrors ConnectorPermissionsTab.
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useListShareLinks, useRevokeShareLink } from '../../api/shareLinks';
import { Panel } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { toast } from '../../lib/toast';
import { fullDate } from '../../lib/time';

export function ShareLinksPage() {
  const { t } = useTranslation();
  const links = useListShareLinks();
  const revoke = useRevokeShareLink();
  // Lazy init: read "now" once at mount, not on every render
  // (react-hooks/purity forbids calling Date.now() during render).
  const [nowIso] = useState(() => new Date().toISOString());

  const onRevoke = (id: string) => {
    revoke.mutate(id, {
      onError: () => toast.error(t('settings.shareLinks.revokeError', { defaultValue: 'Could not revoke link' })),
    });
  };

  return (
    <Panel className="p-5">
      <h2 className="mb-1 text-sm font-semibold text-ink">
        {t('settings.shareLinks.title', { defaultValue: 'Share links' })}
      </h2>
      <p className="mb-4 text-xs text-ink-muted">
        {t('settings.shareLinks.subtitle', {
          defaultValue: 'Read-only links you created into the doc tree. Revoke one to cut off access immediately.',
        })}
      </p>

      {links.isLoading ? (
        <SkeletonRows rows={3} />
      ) : links.isError ? (
        <ErrorState
          description={t('settings.shareLinks.loadError', { defaultValue: 'Could not load share links' })}
          onRetry={() => links.refetch()}
        />
      ) : (links.data ?? []).length === 0 ? (
        <EmptyState title={t('settings.shareLinks.empty', { defaultValue: 'No share links yet' })} />
      ) : (
        <ul className="divide-y divide-line-soft">
          {(links.data ?? []).map((link) => {
            const revoked = !!link.revokedAt;
            const expired = !revoked && link.expiresAt <= nowIso;
            return (
              <li key={link.id} className="flex items-center justify-between gap-3 py-2.5">
                <div className="min-w-0">
                  <p className="truncate font-mono text-sm text-ink">{link.docTreeRoot}</p>
                  <p className="truncate text-2xs text-ink-faint">
                    {t('settings.shareLinks.created', { defaultValue: 'Created' })} {fullDate(link.createdAt)}
                    {' · '}
                    {revoked
                      ? t('settings.shareLinks.revoked', { defaultValue: 'Revoked' })
                      : expired
                        ? t('settings.shareLinks.expired', { defaultValue: 'Expired' })
                        : `${t('settings.shareLinks.expires', { defaultValue: 'Expires' })} ${fullDate(link.expiresAt)}`}
                  </p>
                </div>
                {!revoked && !expired && (
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={revoke.isPending}
                    onClick={() => onRevoke(link.id)}
                    className="shrink-0"
                  >
                    {t('settings.shareLinks.revoke', { defaultValue: 'Revoke' })}
                  </Button>
                )}
              </li>
            );
          })}
        </ul>
      )}
    </Panel>
  );
}
