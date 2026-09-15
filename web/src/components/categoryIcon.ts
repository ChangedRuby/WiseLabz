/** Connector category → icon. Kept out of icons.tsx so that file stays
 *  component-only (react-refresh / fast-refresh friendliness). */
import { ServerIcon, BoxIcon, NetworkIcon, GlobeIcon } from './icons';

export const categoryIcon = {
  virtualization: ServerIcon,
  containers_paas: BoxIcon,
  networking: NetworkIcon,
  dns: GlobeIcon,
} as const;
