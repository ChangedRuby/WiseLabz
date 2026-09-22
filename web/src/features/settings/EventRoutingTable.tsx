/**
 * The notification routing matrix: event types (rows) × channels (columns). Each
 * cell is an on/off toggle deciding whether that event reaches that channel; each
 * row carries a minimum-severity threshold applied to all its channels. Channel
 * columns are disabled when the channel itself is off. Controlled — the parent
 * owns the NotificationConfig and persists it.
 */
import { useTranslation } from 'react-i18next';
import type { NotificationConfig, NotificationChannelType, Severity } from '../../api/model';
import { Severity as SeverityEnum, ConnectorCategory } from '../../api/model';
import { useGetConnectors } from '../../api/generated/connectors/connectors';
import { Toggle, Select } from './parts';
import { cn } from '../../lib/cn';

const CHANNEL_LABELS: Record<NotificationChannelType, string> = {
  in_app: 'In-app',
  smtp: 'Email',
  webhook: 'Webhook',
  discord: 'Discord',
  slack: 'Slack',
  ntfy: 'ntfy',
  telegram: 'Telegram',
};

function eventLabel(eventType: string): string {
  return eventType.replace(/[._]/g, ' ');
}

export function EventRoutingTable({
  config,
  onChange,
  disabled,
}: {
  config: NotificationConfig;
  onChange: (next: NotificationConfig) => void;
  disabled?: boolean;
}) {
  const { t } = useTranslation();
  const channels = config.channels;
  const eventTypes = Array.from(new Set(config.routing.map((r) => r.eventType)));
  const { data: connectors = [] } = useGetConnectors();

  const findRoute = (eventType: string, channel: NotificationChannelType) =>
    config.routing.find((r) => r.eventType === eventType && r.channel === channel);

  const setCell = (eventType: string, channel: NotificationChannelType, enabled: boolean) => {
    onChange({
      ...config,
      routing: config.routing.map((r) =>
        r.eventType === eventType && r.channel === channel ? { ...r, enabled } : r
      ),
    });
  };

  const setRowSeverity = (eventType: string, minSeverity: Severity) => {
    onChange({
      ...config,
      routing: config.routing.map((r) => (r.eventType === eventType ? { ...r, minSeverity } : r)),
    });
  };

  const rowSeverity = (eventType: string): Severity =>
    findRoute(eventType, channels[0]?.type)?.minSeverity ?? SeverityEnum.info;

  const setRowConnectorCategory = (eventType: string, connectorCategory: string) => {
    onChange({
      ...config,
      routing: config.routing.map((r) =>
        r.eventType === eventType
          ? {
              ...r,
              connectorCategory: connectorCategory === '' ? undefined : connectorCategory,
            }
          : r
      ),
    });
  };

  const rowConnectorCategory = (eventType: string): string =>
    findRoute(eventType, channels[0]?.type)?.connectorCategory ?? '';

  const setRowConnectorId = (eventType: string, connectorId: string) => {
    onChange({
      ...config,
      routing: config.routing.map((r) =>
        r.eventType === eventType
          ? { ...r, connectorId: connectorId === '' ? undefined : connectorId }
          : r
      ),
    });
  };

  const rowConnectorId = (eventType: string): string =>
    findRoute(eventType, channels[0]?.type)?.connectorId ?? '';

  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-sm">
        <thead>
          <tr className="border-b border-line-soft">
            <th className="py-2 pr-4 text-left font-mono text-2xs text-ink-faint">
              {t('settings.notifications.event')}
            </th>
            {channels.map((c) => (
              <th key={c.type} className="px-3 py-2 text-center">
                <span
                  className={cn(
                    'font-mono text-2xs',
                    c.enabled ? 'text-ink-muted' : 'text-ink-faint line-through'
                  )}
                >
                  {t(`settings.notifications.channel.${c.type}`, {
                    defaultValue: CHANNEL_LABELS[c.type],
                  })}
                </span>
              </th>
            ))}
            <th className="px-3 py-2 text-right font-mono text-2xs text-ink-faint">
              {t('settings.notifications.minSeverity')}
            </th>
            <th className="px-3 py-2 text-right font-mono text-2xs text-ink-faint">
              {t('settings.notifications.connectorCategory')}
            </th>
            <th className="px-3 py-2 text-right font-mono text-2xs text-ink-faint">
              {t('settings.notifications.connectorId')}
            </th>
          </tr>
        </thead>
        <tbody>
          {eventTypes.map((eventType) => (
            <tr key={eventType} className="border-b border-line-soft last:border-0">
              <td className="py-2.5 pr-4">
                <span className="font-mono text-xs text-ink">{eventLabel(eventType)}</span>
              </td>
              {channels.map((c) => {
                const route = findRoute(eventType, c.type);
                const cellDisabled = disabled || !c.enabled || !route;
                return (
                  <td key={c.type} className="px-3 py-2.5">
                    <div className="flex justify-center">
                      <Toggle
                        size="sm"
                        checked={Boolean(route?.enabled) && c.enabled}
                        disabled={cellDisabled}
                        onChange={(enabled) => setCell(eventType, c.type, enabled)}
                        label={`${eventLabel(eventType)} → ${CHANNEL_LABELS[c.type]}`}
                      />
                    </div>
                  </td>
                );
              })}
              <td className="px-3 py-2.5 text-right">
                <Select
                  aria-label={t('settings.notifications.minSeverityFor', {
                    event: eventLabel(eventType),
                  })}
                  value={rowSeverity(eventType)}
                  disabled={disabled}
                  onChange={(e) => setRowSeverity(eventType, e.target.value as Severity)}
                  className="w-28 py-1.5 text-xs"
                >
                  <option value={SeverityEnum.info}>{t('status.severity.info')}</option>
                  <option value={SeverityEnum.warning}>{t('status.severity.warning')}</option>
                  <option value={SeverityEnum.critical}>{t('status.severity.critical')}</option>
                </Select>
              </td>
              <td className="px-3 py-2.5 text-right">
                <Select
                  aria-label={`${t('settings.notifications.connectorCategory')} for ${eventLabel(eventType)}`}
                  value={rowConnectorCategory(eventType)}
                  disabled={disabled}
                  onChange={(e) => setRowConnectorCategory(eventType, e.target.value)}
                  className="w-28 py-1.5 text-xs"
                >
                  <option value="">{t('settings.notifications.anyConnectorCategory')}</option>
                  {Object.values(ConnectorCategory).map((cat) => (
                    <option key={cat} value={cat}>
                      {cat}
                    </option>
                  ))}
                </Select>
              </td>
              <td className="px-3 py-2.5 text-right">
                <Select
                  aria-label={`${t('settings.notifications.connectorId')} for ${eventLabel(eventType)}`}
                  value={rowConnectorId(eventType)}
                  disabled={disabled}
                  onChange={(e) => setRowConnectorId(eventType, e.target.value)}
                  className="w-28 py-1.5 text-xs"
                >
                  <option value="">{t('settings.notifications.anyConnector')}</option>
                  {connectors.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </Select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
