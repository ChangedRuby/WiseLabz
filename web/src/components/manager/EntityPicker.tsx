/**
 * Entity picker for connector lifecycle ops (restart/start/stop) and runbook
 * steps (#282): lists the entities (VMs, containers, …) from a connector's
 * latest snapshot, keyed by SnapshotEntity.externalId, plus a "Whole
 * service" option (entityRef '') for connectors that manage a single
 * implicit service or when no specific entity should be targeted.
 */
import { useTranslation } from 'react-i18next';
import {
  useGetConnectorsConnectorIdSnapshots,
  useGetConnectorsConnectorIdSnapshotsSnapshotId,
} from '../../api/generated/connectors/connectors';
import { ChevronDownIcon } from '../icons';

export function EntityPicker({
  connectorId,
  value,
  onChange,
  id,
  label,
}: {
  connectorId: string;
  value: string;
  onChange: (entityRef: string) => void;
  id?: string;
  label?: string;
}) {
  const { t } = useTranslation();
  const latest = useGetConnectorsConnectorIdSnapshots(
    connectorId,
    { limit: 1 },
    { query: { enabled: !!connectorId } }
  );
  const latestId = latest.data?.[0]?.id ?? '';
  const full = useGetConnectorsConnectorIdSnapshotsSnapshotId(connectorId, latestId, {
    query: { enabled: !!connectorId && !!latestId },
  });
  const entities = (full.data?.entities ?? []).filter((e) => !!e.externalId);

  const pickerLabel = label ?? t('entityPicker.label');

  return (
    <label className="block">
      <span className="mb-1 block text-2xs text-ink-faint">{pickerLabel}</span>
      <div className="relative">
        <select
          id={id}
          aria-label={pickerLabel}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="h-8 w-full appearance-none rounded-sm border border-line bg-surface pl-2.5 pr-7 text-xs text-ink outline-none focus-visible:border-accent-primary-soft"
        >
          <option value="">{t('entityPicker.wholeService')}</option>
          {entities.map((entity) => (
            <option key={entity.externalId} value={entity.externalId}>
              {entity.name} ({entity.externalId})
            </option>
          ))}
        </select>
        <ChevronDownIcon
          size={12}
          className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-ink-faint"
        />
      </div>
    </label>
  );
}
