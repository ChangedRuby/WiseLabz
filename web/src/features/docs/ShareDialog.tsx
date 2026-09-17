/**
 * Creates a read-only doc share link for one subtree (#240 PR2). TTL preset
 * (server requires a future expiresAt, no open-ended links) then a one-time
 * reveal of the raw token — mirrors how a freshly-minted API key is shown
 * once and never again (see settings/AiPage.tsx's secret-reveal pattern).
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Dialog } from '../../components/ui/Dialog';
import { Button } from '../../components/ui/Button';
import { CopyIcon, CheckIcon } from '../../components/icons';
import { useCreateShareLink } from '../../api/shareLinks';
import { toast } from '../../lib/toast';
import { cn } from '../../lib/cn';

const TTL_PRESETS = [
  { labelKey: 'docs.share.ttl.day', labelDefault: '24 hours', hours: 24 },
  { labelKey: 'docs.share.ttl.week', labelDefault: '7 days', hours: 24 * 7 },
  { labelKey: 'docs.share.ttl.month', labelDefault: '30 days', hours: 24 * 30 },
] as const;

interface ShareDialogProps {
  open: boolean;
  onClose: () => void;
  node: { docId: string; title: string } | null;
}

export function ShareDialog({ open, onClose, node }: ShareDialogProps) {
  const { t } = useTranslation();
  const create = useCreateShareLink();
  const [ttlHours, setTtlHours] = useState<number>(TTL_PRESETS[1].hours);
  const [createdUrl, setCreatedUrl] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleClose = () => {
    setCreatedUrl(null);
    setCopied(false);
    create.reset();
    onClose();
  };

  const handleCreate = () => {
    if (!node) return;
    const expiresAt = new Date(Date.now() + ttlHours * 60 * 60 * 1000).toISOString();
    create.mutate(
      { docTreeRoot: node.docId, expiresAt },
      {
        onSuccess: (link) => {
          setCreatedUrl(`${window.location.origin}/share/${link.token}`);
        },
        onError: () => toast.error(t('docs.share.error', { defaultValue: 'Could not create share link' })),
      }
    );
  };

  const handleCopy = async () => {
    if (!createdUrl) return;
    try {
      await navigator.clipboard.writeText(createdUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error(t('docs.share.copyError', { defaultValue: 'Could not copy link' }));
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      title={t('docs.share.title', { defaultValue: 'Share' })}
    >
      {!node ? null : createdUrl ? (
        <div className="flex flex-col gap-3">
          <p className="text-xs text-ink-muted">
            {t('docs.share.onceNotice', {
              defaultValue: "Copy this link now — it won't be shown again.",
            })}
          </p>
          <div className="flex items-center gap-2">
            <input
              readOnly
              value={createdUrl}
              onFocus={(e) => e.currentTarget.select()}
              className="min-w-0 flex-1 rounded-md border border-line-soft bg-canvas px-2.5 py-1.5 font-mono text-xs text-ink"
            />
            <Button variant="secondary" size="sm" onClick={handleCopy}>
              {copied ? <CheckIcon size={14} /> : <CopyIcon size={14} />}
              {copied
                ? t('docs.share.copied', { defaultValue: 'Copied' })
                : t('docs.share.copy', { defaultValue: 'Copy' })}
            </Button>
          </div>
          <div className="flex justify-end">
            <Button variant="secondary" size="sm" onClick={handleClose}>
              {t('common.done', { defaultValue: 'Done' })}
            </Button>
          </div>
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          <p className="text-xs text-ink-muted">
            {t('docs.share.subtitle', {
              defaultValue: 'Anyone with the link can view "{{title}}" without an account, read-only.',
              title: node.title,
            })}
          </p>
          <div>
            <label className="mb-1.5 block text-xs font-medium text-ink-muted">
              {t('docs.share.expiresLabel', { defaultValue: 'Expires in' })}
            </label>
            <div className="flex gap-1.5">
              {TTL_PRESETS.map((preset) => (
                <button
                  key={preset.hours}
                  type="button"
                  onClick={() => setTtlHours(preset.hours)}
                  className={cn(
                    'flex-1 rounded-md border px-2.5 py-1.5 text-xs transition-colors',
                    ttlHours === preset.hours
                      ? 'border-accent-primary bg-accent-primary-tint text-ink'
                      : 'border-line-soft text-ink-muted hover:bg-surface-raised'
                  )}
                >
                  {t(preset.labelKey, { defaultValue: preset.labelDefault })}
                </button>
              ))}
            </div>
          </div>
          <div className="flex justify-end gap-2">
            <Button variant="secondary" size="sm" onClick={handleClose}>
              {t('common.cancel', { defaultValue: 'Cancel' })}
            </Button>
            <Button variant="primary" size="sm" disabled={create.isPending} onClick={handleCreate}>
              {t('docs.share.create', { defaultValue: 'Create link' })}
            </Button>
          </div>
        </div>
      )}
    </Dialog>
  );
}
