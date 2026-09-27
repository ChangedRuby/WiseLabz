/**
 * Shared dry-run-preview + elevation-confirm dialogs for a connector
 * lifecycle mutation (restart/start/stop) — driven by {@link useMutatingOp}
 * (`useMutatingOp.ts`), reused by ServiceDetailPage and RunbookPanel (#282).
 */
import type { ReactNode } from 'react';
import { Button } from '../ui/Button';
import { Dialog } from '../ui/Dialog';
import { SkeletonRows } from '../ui/states';
import { ElevationConfirm } from './ElevationConfirm';
import type { MutatingOp } from './useMutatingOp';

/** Every string the dialogs render, pre-translated by the caller (each
 * consumer has its own i18n keys — ServiceDetailPage's verb-specific copy vs.
 * RunbookPanel's step copy). */
export interface MutatingOpMessages {
  previewTitle: string;
  previewNotice: string;
  previewError: string;
  targetLabel: string;
  downtimeLabel: string;
  downtimeSeconds: (count: number) => string;
  downtimeIndefinite: string;
  dependenciesLabel: string;
  noDependencies: string;
  failed: string;
  retry: string;
  confirmTitle: string;
  confirmDescription: string;
  confirmLabel: string;
}

export function MutatingOpDialogs({
  op,
  action,
  resourceName,
  messages,
  entityPicker,
}: {
  op: MutatingOp;
  /** Elevation action string, e.g. `connector.restart`. */
  action: string;
  resourceName: string;
  messages: MutatingOpMessages;
  /** Optional entity picker rendered above the preview, for connectors that
   * manage more than one entity (VMs, containers, …). */
  entityPicker?: ReactNode;
}) {
  return (
    <>
      <Dialog
        open={op.previewOpen}
        onClose={() => op.setPreviewOpen(false)}
        title={messages.previewTitle}
        size="sm"
      >
        <div className="space-y-4">
          {entityPicker}
          {op.preview.isPending ? (
            <SkeletonRows rows={3} />
          ) : op.preview.isError || !op.preview.data ? (
            <div className="space-y-3">
              <p className="text-sm text-err">{messages.previewError}</p>
              <Button size="sm" variant="secondary" onClick={() => op.rerunPreview()}>
                {messages.retry}
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              <p className="text-sm text-ink-muted">{messages.previewNotice}</p>
              <dl className="grid grid-cols-2 gap-x-4 gap-y-3 rounded-md border border-line-soft bg-canvas-sunken p-3 text-sm">
                <div>
                  <dt className="text-2xs text-ink-faint">{messages.targetLabel}</dt>
                  <dd className="mt-0.5 font-medium text-ink">{op.preview.data.targetService}</dd>
                </div>
                <div>
                  <dt className="text-2xs text-ink-faint">{messages.downtimeLabel}</dt>
                  <dd className="mt-0.5 font-medium text-ink">
                    {op.preview.data.estimatedDowntimeSeconds > 0
                      ? messages.downtimeSeconds(op.preview.data.estimatedDowntimeSeconds)
                      : messages.downtimeIndefinite}
                  </dd>
                </div>
              </dl>
              <div>
                <h3 className="text-sm font-semibold text-ink">{messages.dependenciesLabel}</h3>
                {op.preview.data.dependentServices.length === 0 ? (
                  <p className="mt-1 text-sm text-ink-muted">{messages.noDependencies}</p>
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
              {op.mutate.isError && <p className="text-2xs text-err">{messages.failed}</p>}
              <div className="flex justify-end">
                <Button size="sm" variant="danger" onClick={() => op.setConfirmOpen(true)}>
                  {messages.confirmLabel}
                </Button>
              </div>
            </div>
          )}
        </div>
      </Dialog>

      <ElevationConfirm
        open={op.confirmOpen}
        resourceName={resourceName}
        action={action}
        title={messages.confirmTitle}
        description={messages.confirmDescription}
        confirmLabel={messages.confirmLabel}
        isPending={op.mutate.isPending}
        onClose={() => op.setConfirmOpen(false)}
        onConfirm={(token) => op.mutate.mutate(token)}
      />
    </>
  );
}
