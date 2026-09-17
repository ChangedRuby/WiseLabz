/**
 * Settings hub section registry. The single source of truth for the section nav
 * (SettingsLayout) and for any command-palette "jump to settings section" entries.
 * `operatorOnly` sections are hidden for viewers (UI convenience only — the server
 * still enforces the boundary on every mutating endpoint).
 */
import type { ComponentType } from 'react';
import {
  UserIcon,
  ServerIcon,
  SettingsIcon,
  SparklesIcon,
  BellIcon,
  GaugeIcon,
  LayersIcon,
  ClockIcon,
  HistoryIcon,
  FileTextIcon,
  ShareIcon,
  AlertTriangleIcon,
} from '../../components/icons';

export interface SettingsSection {
  /** Path segment relative to /settings (also the route child path). */
  segment: string;
  /** i18n key for the label. */
  labelKey: string;
  /** Inline default for the label (en.ts is not editable in this phase). */
  labelDefault: string;
  Icon: ComponentType<{ size?: number }>;
  operatorOnly: boolean;
  instanceAdminOnly?: boolean;
}

export const SETTINGS_SECTIONS: SettingsSection[] = [
  {
    segment: 'profile',
    labelKey: 'settings.nav.profile',
    labelDefault: 'Profile',
    Icon: UserIcon,
    operatorOnly: false,
  },
  {
    segment: 'users',
    labelKey: 'settings.nav.users',
    labelDefault: 'Users',
    Icon: LayersIcon,
    operatorOnly: true,
  },
  {
    segment: 'auth',
    labelKey: 'settings.nav.auth',
    labelDefault: 'Authentication',
    Icon: ServerIcon,
    operatorOnly: true,
  },
  {
    segment: 'ai',
    labelKey: 'settings.nav.ai',
    labelDefault: 'AI',
    Icon: SparklesIcon,
    operatorOnly: true,
  },
  {
    segment: 'notifications',
    labelKey: 'settings.nav.notifications',
    labelDefault: 'Notifications',
    Icon: BellIcon,
    operatorOnly: true,
  },
  {
    segment: 'system',
    labelKey: 'settings.nav.system',
    labelDefault: 'System',
    Icon: GaugeIcon,
    operatorOnly: true,
  },
  {
    segment: 'runbooks',
    labelKey: 'settings.nav.runbooks',
    labelDefault: 'Runbooks',
    Icon: FileTextIcon,
    operatorOnly: true,
  },
  {
    segment: 'retention',
    labelKey: 'settings.nav.retention',
    labelDefault: 'Retention',
    Icon: ClockIcon,
    operatorOnly: true,
  },
  {
    segment: 'audit',
    labelKey: 'settings.nav.audit',
    labelDefault: 'Audit log',
    Icon: HistoryIcon,
    operatorOnly: true,
  },
  {
    segment: 'share-links',
    labelKey: 'settings.nav.shareLinks',
    labelDefault: 'Share links',
    Icon: ShareIcon,
    operatorOnly: true,
  },
  {
    segment: 'compliance',
    labelKey: 'settings.nav.compliance',
    labelDefault: 'Compliance rules',
    Icon: AlertTriangleIcon,
    operatorOnly: true,
    instanceAdminOnly: true,
  },
  {
    segment: 'appearance',
    labelKey: 'settings.nav.appearance',
    labelDefault: 'Appearance',
    Icon: SettingsIcon,
    operatorOnly: false,
  },
];
