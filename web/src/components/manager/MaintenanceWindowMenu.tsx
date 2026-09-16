/**
 * Maintenance window popover: preset (30/60/120m) + custom-minutes control to
 * open a time-boxed suppression on one connector. Operator-only, gated by
 * useCanMutate() alone — no ElevationConfirm, since a maintenance window is
 * reversible and time-boxed (see ADR context in issue #236 PR3).
 */
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { AnimatePresence, motion } from 'motion/react';
import * as Popover from '@radix-ui/react-popover';
import {
  postConnectorsConnectorIdMaintenanceWindow,
  deleteConnectorsConnectorIdMaintenanceWindow,
  getGetConnectorsMaintenanceWindowsQueryKey,
} from '../../api/generated/connectors/connectors';
import { Button } from '../ui/Button';
import { ClockIcon } from '../icons';

const presets = [30, 60, 120];

interface MaintenanceWindowMenuProps {
  connectorId: string;
  active: boolean;
}

export function MaintenanceWindowMenu({ connectorId, active }: MaintenanceWindowMenuProps) {
  const { t } = useTranslation();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [customMinutes, setCustomMinutes] = useState('');

  const invalidate = () =>
    qc.invalidateQueries({ queryKey: getGetConnectorsMaintenanceWindowsQueryKey() });

  const openWindow = useMutation({
    mutationFn: (durationMinutes: number) =>
      postConnectorsConnectorIdMaintenanceWindow(connectorId, { durationMinutes }),
    onSuccess: () => {
      setCustomMinutes('');
      setOpen(false);
      invalidate();
    },
  });

  const closeWindow = useMutation({
    mutationFn: () => deleteConnectorsConnectorIdMaintenanceWindow(connectorId),
    onSuccess: invalidate,
  });

  if (active) {
    return (
      <Button
        size="sm"
        variant="ghost"
        onClick={() => closeWindow.mutate()}
        disabled={closeWindow.isPending}
      >
        <ClockIcon size={14} /> {t('services.maintenance.close')}
      </Button>
    );
  }

  const customValue = Number(customMinutes);
  const customValid = customMinutes.trim() !== '' && customValue > 0;

  return (
    <Popover.Root open={open} onOpenChange={setOpen}>
      <Popover.Trigger asChild>
        <Button size="sm" variant="ghost">
          <ClockIcon size={14} /> {t('services.maintenance.action')}
        </Button>
      </Popover.Trigger>

      <AnimatePresence>
        {open && (
          <Popover.Portal forceMount>
            <Popover.Content
              asChild
              align="end"
              sideOffset={8}
              onOpenAutoFocus={(e) => e.preventDefault()}
            >
              <motion.div
                initial={{ opacity: 0, y: -6, scale: 0.98 }}
                animate={{ opacity: 1, y: 0, scale: 1 }}
                exit={{ opacity: 0, y: -6, scale: 0.98 }}
                transition={{ duration: 0.16, ease: [0.16, 1, 0.3, 1] }}
                className="z-(--z-dropdown) w-64 overflow-hidden rounded-sm border border-line bg-surface-overlay p-3 shadow-(--shadow-pop)"
              >
                <p className="text-xs font-medium text-ink">{t('services.maintenance.heading')}</p>
                <p className="mt-1 text-2xs text-ink-faint">{t('services.maintenance.description')}</p>

                <div className="mt-3 flex gap-1.5">
                  {presets.map((minutes) => (
                    <Button
                      key={minutes}
                      size="sm"
                      variant="secondary"
                      className="flex-1"
                      disabled={openWindow.isPending}
                      onClick={() => openWindow.mutate(minutes)}
                    >
                      {t(`services.maintenance.preset${minutes}` as const)}
                    </Button>
                  ))}
                </div>

                <div className="mt-2 flex items-center gap-1.5">
                  <label htmlFor={`maintenance-custom-${connectorId}`} className="sr-only">
                    {t('services.maintenance.customLabel')}
                  </label>
                  <input
                    id={`maintenance-custom-${connectorId}`}
                    type="number"
                    min={1}
                    value={customMinutes}
                    onChange={(e) => setCustomMinutes(e.target.value)}
                    placeholder={t('services.maintenance.customPlaceholder')}
                    className="h-8 min-w-0 flex-1 rounded-sm border border-line-soft bg-canvas-sunken px-2 text-sm text-ink outline-none placeholder:text-ink-faint"
                  />
                  <Button
                    size="sm"
                    variant="primary"
                    disabled={!customValid || openWindow.isPending}
                    onClick={() => openWindow.mutate(customValue)}
                  >
                    {t('services.maintenance.start')}
                  </Button>
                </div>
              </motion.div>
            </Popover.Content>
          </Popover.Portal>
        )}
      </AnimatePresence>
    </Popover.Root>
  );
}
