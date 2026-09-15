/** Single source of truth for primary navigation, shared by every shell variant. */
import {
  GaugeIcon,
  LayersIcon,
  FileTextIcon,
  DiffIcon,
  BellIcon,
  AlertTriangleIcon,
  SparklesIcon,
  ChatIcon,
  NetworkIcon,
} from '../icons';

export interface NavItem {
  to: string;
  label: string;
  Icon: React.ComponentType<{ size?: number; className?: string }>;
  /** show the live pending-alerts badge on this item */
  badge?: boolean;
  badgeSource?: 'alerts' | 'findings' | 'attention';
}

export const NAV: NavItem[] = [
  { to: '/dashboard', label: 'Dashboard', Icon: GaugeIcon },
  { to: '/services', label: 'Services', Icon: LayersIcon },
  { to: '/docs', label: 'Docs', Icon: FileTextIcon },
  { to: '/topology', label: 'Topology', Icon: NetworkIcon },
  { to: '/chat', label: 'Chat', Icon: ChatIcon },
  { to: '/changes', label: 'Changes', Icon: DiffIcon },
  { to: '/attention', label: 'Attention', Icon: SparklesIcon, badge: true, badgeSource: 'attention' },
  { to: '/alerts', label: 'Alerts', Icon: BellIcon, badge: true },
  { to: '/findings', label: 'Findings', Icon: AlertTriangleIcon, badge: true, badgeSource: 'findings' },
];
