/**
 * Curated, deterministic homelab dataset for the mock layer. The generated faker
 * mocks produce random noise that reads as slop; these hand-authored fixtures let
 * the real screens (wired through the generated React Query hooks) render a
 * believable lab. Shapes are the orval-generated contract types — no drift.
 */
import type {
  Alert,
  AlertPage,
  ChangeDetail,
  ChangePage,
  ChangeSummary,
  Connector,
  ConnectorTypeSchema,
  DashboardOverview,
  Doc,
  DocNode,
  DocVersion,
  DocVersionMeta,
  QualityCheckType,
  QualityFinding,
  QualityFindingPage,
  QualityFindingStatus,
  RemovalImpact,
  ServiceSnapshot,
  SyncRun,
  User,
} from '../api/model';

const minsAgo = (m: number) => new Date(Date.now() - m * 60_000).toISOString();
const hrsAgo = (h: number) => new Date(Date.now() - h * 3_600_000).toISOString();
const daysAgo = (d: number) => new Date(Date.now() - d * 86_400_000).toISOString();
const minsFromNow = (m: number) => new Date(Date.now() + m * 60_000).toISOString();

// TODO: fold instanceAdminRole into the generated User type once the OpenAPI
// spec is updated for #240 PR1 — see src/hooks/useRole.ts.
export const user = {
  id: 'usr-1',
  username: 'ops',
  displayName: 'Ada',
  email: 'ada@homelab.lan',
  role: 'operator',
  instanceAdminRole: 'admin',
  authSource: 'local',
  createdAt: daysAgo(212),
} satisfies User & { instanceAdminRole: 'admin' | 'user' };

// TODO: fold myRole into the generated Connector type once the OpenAPI spec
// is updated for #240 PR1 — see src/hooks/useRole.ts. Every fixture connector
// grants the fixture user (an instance admin) operator access, matching
// today's single-admin-homelab default.
const connectorsBase: Connector[] = [
  {
    id: 'svc-pve1',
    name: 'pve1',
    category: 'virtualization',
    type: 'proxmox',
    enabled: true,
    status: 'online',
    url: 'https://10.0.0.11:8006',
    verifyTls: false,
    lastSyncAt: minsAgo(4),
    scheduleSeconds: 900,
    nextRunAt: minsFromNow(11),
    lastSyncDurationMs: 1240,
    retryCount: 0,
  },
  {
    id: 'svc-pve2',
    name: 'pve2',
    category: 'virtualization',
    type: 'proxmox',
    enabled: true,
    status: 'online',
    url: 'https://10.0.0.12:8006',
    verifyTls: false,
    lastSyncAt: minsAgo(6),
    scheduleSeconds: null,
    nextRunAt: '',
    lastSyncDurationMs: 980,
    retryCount: 0,
  },
  {
    id: 'svc-opnsense',
    name: 'opnsense',
    category: 'networking',
    type: 'opnsense',
    enabled: true,
    status: 'online',
    url: 'https://10.0.0.1',
    verifyTls: true,
    lastSyncAt: minsAgo(4),
    scheduleSeconds: 3600,
    nextRunAt: minsFromNow(56),
    lastSyncDurationMs: 610,
    retryCount: 0,
  },
  {
    id: 'svc-portainer',
    name: 'portainer',
    category: 'containers_paas',
    type: 'portainer',
    enabled: true,
    status: 'degraded',
    url: 'https://10.0.0.20:9443',
    verifyTls: false,
    lastSyncAt: minsAgo(11),
    statusMessage: '1 stack unhealthy — media (sonarr exited)',
    scheduleSeconds: 1800,
    nextRunAt: minsFromNow(19),
    lastSyncDurationMs: 5230,
    lastSyncError: 'stack health check timed out after 5s',
    retryCount: 2,
  },
  {
    id: 'svc-pbs',
    name: 'pbs',
    category: 'virtualization',
    type: 'proxmox-backup',
    enabled: true,
    status: 'online',
    url: 'https://10.0.0.13:8007',
    verifyTls: false,
    lastSyncAt: hrsAgo(1),
    scheduleSeconds: 86400,
    nextRunAt: minsFromNow(1380),
    lastSyncDurationMs: 2100,
    retryCount: 0,
  },
  {
    id: 'svc-docker-nas',
    name: 'docker-nas',
    category: 'containers_paas',
    type: 'docker',
    enabled: true,
    status: 'offline',
    url: 'tcp://10.0.0.30:2375',
    verifyTls: false,
    lastSyncAt: hrsAgo(2),
    statusMessage: 'connection refused — host unreachable',
    scheduleSeconds: 900,
    nextRunAt: minsFromNow(3),
    lastSyncDurationMs: null,
    lastSyncError: 'connection refused — host unreachable',
    retryCount: 4,
  },
];

export const connectors = connectorsBase.map((c) => ({ ...c, myRole: 'operator' as const }));

export const changes: ChangeSummary[] = [
  {
    id: 'chg-1',
    serviceId: 'svc-pve1',
    serviceName: 'pve1',
    changeType: 'vm.created',
    severity: 'warning',
    summary: 'New VM "k8s-worker-03" detected on node pve1',
    willTriggerAi: true,
    detectedAt: minsAgo(4),
  },
  {
    id: 'chg-2',
    serviceId: 'svc-opnsense',
    serviceName: 'opnsense',
    changeType: 'firewall.rule.modified',
    severity: 'critical',
    summary: 'WAN rule "allow any → any" added to the firewall',
    willTriggerAi: true,
    detectedAt: minsAgo(12),
  },
  {
    id: 'chg-3',
    serviceId: 'svc-portainer',
    serviceName: 'portainer',
    changeType: 'container.exited',
    severity: 'warning',
    summary: 'Container "sonarr" in stack media exited (code 1)',
    willTriggerAi: false,
    detectedAt: minsAgo(13),
  },
  {
    id: 'chg-4',
    serviceId: 'svc-pve2',
    serviceName: 'pve2',
    changeType: 'vm.deleted',
    severity: 'warning',
    summary: 'VM "test-sandbox" removed from node pve2',
    willTriggerAi: true,
    detectedAt: hrsAgo(2),
  },
  {
    id: 'chg-5',
    serviceId: 'svc-portainer',
    serviceName: 'portainer',
    changeType: 'image.updated',
    severity: 'info',
    summary: 'jellyfin image linuxserver/jellyfin:10.9.3 → 10.9.4',
    willTriggerAi: false,
    detectedAt: hrsAgo(5),
  },
  {
    id: 'chg-6',
    serviceId: 'svc-pve1',
    serviceName: 'pve1',
    changeType: 'storage.threshold',
    severity: 'info',
    summary: 'local-zfs pool usage crossed 70% (was 64%)',
    willTriggerAi: false,
    detectedAt: daysAgo(1),
  },
  {
    id: 'chg-7',
    serviceId: 'svc-pve1',
    serviceName: 'pve1',
    changeType: 'doc.section.updated',
    severity: 'info',
    summary: 'pve1 doc: k8s-worker-03 row + autoscaler note drafted (rev 3 → 4)',
    willTriggerAi: true,
    detectedAt: minsAgo(4),
  },
];

export const changeDetails: Record<string, ChangeDetail> = {
  'chg-2': {
    ...changes[1],
    status: 'new',
    affectedDocIds: ['doc-opnsense'],
    diff: {
      format: 'infra',
      hunks: [
        {
          path: 'firewall.rules[wan].in[12]',
          before: null,
          after:
            'pass in quick on wan proto any from any to any  # allow-any-any',
        },
        {
          path: 'firewall.rules[wan].count',
          before: '11',
          after: '12',
        },
        {
          path: 'firewall.lastEditedBy',
          before: 'root@local',
          after: 'admin@10.0.0.99',
        },
      ],
    },
  },
  'chg-1': {
    ...changes[0],
    status: 'acknowledged',
    affectedDocIds: ['doc-pve1'],
    diff: {
      format: 'infra',
      hunks: [
        { path: 'nodes.pve1.vms[112].name', before: null, after: 'k8s-worker-03' },
        { path: 'nodes.pve1.vms[112].cores', before: null, after: '4' },
        { path: 'nodes.pve1.vms[112].memoryMB', before: null, after: '8192' },
        { path: 'nodes.pve1.vmCount', before: '11', after: '12' },
      ],
    },
  },
};

/** Authored detail when present, otherwise synthesized from the summary so any
 *  change in the feed opens to a real detail view. */
export function changeDetail(id: string): ChangeDetail | null {
  if (changeDetails[id]) return changeDetails[id];
  // Doc-format change: full before/after revisions drive the JetBrains-style
  // diff. Built lazily here (not in the static map) so it can reference the doc
  // revision constants declared further down without a temporal-dead-zone error.
  if (id === 'chg-7') {
    return {
      ...changes[6],
      status: 'new',
      affectedDocIds: ['doc-pve1'],
      diff: {
        format: 'doc',
        language: 'md',
        baseLabel: 'rev 3',
        headLabel: 'rev 4',
        baseText: PVE1_DOC_V3,
        headText: PVE1_DOC,
      },
      provenance: [
        {
          snapshotPath: 'nodes.pve1.vms[112]',
          templateSection: 'Virtual machines',
          diffLines: [11, 12, 13],
        },
      ],
    };
  }
  const summary = changes.find((c) => c.id === id);
  if (!summary) return null;
  return {
    ...summary,
    status: 'new',
    affectedDocIds: [],
    diff: {
      format: 'infra',
      hunks: [
        {
          path: `${summary.serviceName}.${summary.changeType}`,
          before: null,
          after: summary.summary,
        },
      ],
    },
  };
}

export const alerts: Alert[] = [
  {
    id: 'alt-1',
    changeId: 'chg-2',
    serviceId: 'svc-opnsense',
    serviceName: 'opnsense',
    severity: 'critical',
    title: 'Permissive WAN firewall rule added',
    description:
      'A rule allowing any source to any destination was added to the WAN interface. This exposes internal services to the internet.',
    status: 'pending',
    createdAt: minsAgo(12),
  },
  {
    id: 'alt-2',
    changeId: 'chg-1',
    serviceId: 'svc-pve1',
    serviceName: 'pve1',
    severity: 'warning',
    title: 'Undocumented VM created',
    description: 'k8s-worker-03 has no matching documentation section.',
    status: 'pending',
    createdAt: minsAgo(4),
  },
  {
    id: 'alt-3',
    changeId: 'chg-3',
    serviceId: 'svc-portainer',
    serviceName: 'portainer',
    severity: 'warning',
    title: 'Container repeatedly exiting',
    description: 'sonarr in stack media has exited 3 times in the last hour.',
    status: 'pending',
    createdAt: minsAgo(13),
  },
];

// ─── Docs ───────────────────────────────────────────────────────────────────

export const docTree: DocNode = {
  docId: 'doc-lab',
  title: 'Home Lab',
  kind: 'lab',
  children: [
    { docId: 'doc-pve1', title: 'pve1 — Proxmox VE', kind: 'service', serviceId: 'svc-pve1' },
    { docId: 'doc-pve2', title: 'pve2 — Proxmox VE', kind: 'service', serviceId: 'svc-pve2' },
    { docId: 'doc-opnsense', title: 'opnsense — Firewall', kind: 'service', serviceId: 'svc-opnsense' },
    { docId: 'doc-portainer', title: 'portainer — Containers', kind: 'service', serviceId: 'svc-portainer' },
    { docId: 'doc-pbs', title: 'pbs — Backup Server', kind: 'service', serviceId: 'svc-pbs' },
  ],
};

const LAB_DOC = `# Home Lab

A two-node Proxmox cluster fronted by an OPNsense firewall, with container
workloads on Portainer and backups handled by Proxmox Backup Server. This page is
generated and kept current by WiseLabz, last reconciled against live state on
each sync.

## Topology

| Layer | Service | Host | Status |
| --- | --- | --- | --- |
| Edge | opnsense | 10.0.0.1 | online |
| Compute | pve1 | 10.0.0.11 | online |
| Compute | pve2 | 10.0.0.12 | online |
| Containers | portainer | 10.0.0.20 | degraded |
| Backup | pbs | 10.0.0.13 | online |

## Conventions

- VLAN **10** carries management traffic; **20** is workloads; **30** is storage.
- Every VM is tagged with its owning stack so drift maps back to a service.
- Secrets live in the host \`config.yaml\`, never in this document.

> Generated by WiseLabz. Edits made here are preserved across syncs unless the
> source of truth changes.
`;

const PVE1_DOC = `# pve1 — Proxmox VE

Primary compute node. Hosts the Kubernetes control plane and the bulk of
long-running virtual machines.

## Node

| Field | Value |
| --- | --- |
| Version | 8.2.4 |
| CPU | AMD Ryzen 9 5950X (16c / 32t) |
| Memory | 128 GB ECC |
| Storage | local-zfs (4 TB NVMe mirror) |

## Virtual machines

| VMID | Name | Cores | Memory | State |
| --- | --- | --- | --- | --- |
| 101 | k8s-control-01 | 4 | 8 GB | running |
| 102 | k8s-worker-01 | 4 | 8 GB | running |
| 103 | k8s-worker-02 | 4 | 8 GB | running |
| 112 | k8s-worker-03 | 4 | 8 GB | running |

## Notes

\`k8s-worker-03\` was added automatically by the cluster autoscaler. WiseLabz
detected it on the last sync and drafted this row for review.
`;

const svcDoc = (name: string, body: string): string => `# ${name}\n\n${body}\n`;

export const docs: Record<string, Doc> = {
  'doc-lab': {
    docId: 'doc-lab',
    title: 'Home Lab',
    kind: 'lab',
    content: LAB_DOC,
    currentVersion: 9,
    updatedAt: minsAgo(4),
  },
  'doc-pve1': {
    docId: 'doc-pve1',
    title: 'pve1 — Proxmox VE',
    kind: 'service',
    serviceId: 'svc-pve1',
    content: PVE1_DOC,
    currentVersion: 4,
    updatedAt: minsAgo(4),
  },
  'doc-pve2': {
    docId: 'doc-pve2',
    title: 'pve2 — Proxmox VE',
    kind: 'service',
    serviceId: 'svc-pve2',
    content: svcDoc(
      'pve2 — Proxmox VE',
      'Secondary compute node used for staging and migration targets.\n\n## Virtual machines\n\n| VMID | Name | State |\n| --- | --- | --- |\n| 201 | staging-db | running |\n| 202 | migration-target | stopped |',
    ),
    currentVersion: 2,
    updatedAt: hrsAgo(2),
  },
  'doc-opnsense': {
    docId: 'doc-opnsense',
    title: 'opnsense — Firewall',
    kind: 'service',
    serviceId: 'svc-opnsense',
    content: svcDoc(
      'opnsense — Firewall',
      'Edge firewall and router. Terminates WAN and segments the lab into management, workload, and storage VLANs.\n\n## Interfaces\n\n| Interface | Network | Role |\n| --- | --- | --- |\n| WAN | dhcp | uplink |\n| LAN | 10.0.0.0/24 | management |\n| WORK | 10.0.20.0/24 | workloads |',
    ),
    currentVersion: 6,
    updatedAt: minsAgo(12),
  },
  'doc-portainer': {
    docId: 'doc-portainer',
    title: 'portainer — Containers',
    kind: 'service',
    serviceId: 'svc-portainer',
    content: svcDoc(
      'portainer — Containers',
      'Container management for the lab. Stacks are defined as compose files under version control.\n\n## Stacks\n\n| Stack | Containers | Status |\n| --- | --- | --- |\n| media | 5 | degraded |\n| monitoring | 3 | healthy |\n| proxy | 2 | healthy |',
    ),
    currentVersion: 5,
    updatedAt: minsAgo(11),
  },
  'doc-pbs': {
    docId: 'doc-pbs',
    title: 'pbs — Backup Server',
    kind: 'service',
    serviceId: 'svc-pbs',
    content: svcDoc(
      'pbs — Backup Server',
      'Proxmox Backup Server. Holds deduplicated backups of every VM and container volume.\n\n## Datastores\n\n| Name | Used | Total |\n| --- | --- | --- |\n| main | 1.8 TB | 4 TB |',
    ),
    currentVersion: 3,
    updatedAt: hrsAgo(1),
  },
};

// Version history for doc-pve1 — drives the DiffViewer documentation view.
export const docVersionMeta: Record<string, DocVersionMeta[]> = {
  'doc-pve1': [
    { rev: 4, createdAt: minsAgo(4), author: null, trigger: 'ai' },
    { rev: 3, createdAt: daysAgo(2), author: 'ada', trigger: 'manual' },
    { rev: 2, createdAt: daysAgo(9), author: null, trigger: 'template' },
    { rev: 1, createdAt: daysAgo(40), author: 'ada', trigger: 'manual' },
  ],
};

const PVE1_DOC_V3 = `# pve1 — Proxmox VE

Primary compute node. Hosts the Kubernetes control plane and the bulk of
long-running virtual machines.

## Node

| Field | Value |
| --- | --- |
| Version | 8.2.4 |
| CPU | AMD Ryzen 9 5950X (16c / 32t) |
| Memory | 128 GB ECC |
| Storage | local-zfs (4 TB NVMe mirror) |

## Virtual machines

| VMID | Name | Cores | Memory | State |
| --- | --- | --- | --- | --- |
| 101 | k8s-control-01 | 4 | 8 GB | running |
| 102 | k8s-worker-01 | 4 | 8 GB | running |
| 103 | k8s-worker-02 | 4 | 8 GB | running |
`;

export const docVersionContent: Record<string, Record<number, DocVersion>> = {
  'doc-pve1': {
    4: { rev: 4, createdAt: minsAgo(4), author: null, trigger: 'ai', content: PVE1_DOC },
    3: { rev: 3, createdAt: daysAgo(2), author: 'ada', trigger: 'manual', content: PVE1_DOC_V3 },
  },
};

export function dashboardOverview(): DashboardOverview {
  const counts = connectors.reduce(
    (acc, c) => {
      acc[c.status] = (acc[c.status] ?? 0) + 1;
      return acc;
    },
    {} as Record<string, number>,
  );
  return {
    statusCounts: {
      online: counts.online ?? 0,
      degraded: counts.degraded ?? 0,
      offline: counts.offline ?? 0,
      unknown: counts.unknown ?? 0,
    },
    pendingAlerts: alerts.filter((a) => a.status === 'pending').length,
    recentChanges: changes.slice(0, 5),
    lastSyncAt: minsAgo(4),
  };
}

export function changePage(): ChangePage {
  return { items: changes, total: 24, page: 1, pageSize: changes.length };
}

// ─── Manager v1: blast radius + connector type schemas ───────────────────────

/** Blast radius for removing a connector — the concrete dependents destroyed. */
export function removalImpact(connectorId: string): RemovalImpact {
  const c = connectors.find((x) => x.id === connectorId);
  const name = c?.name ?? connectorId;
  return {
    trackedServices: 1,
    docSections: 4,
    snapshots: 36,
    items: [
      { type: 'service', name },
      { type: 'doc-section', name: `${name} — Node` },
      { type: 'doc-section', name: `${name} — Virtual machines` },
      { type: 'doc-section', name: `${name} — Notes` },
      { type: 'snapshot', name: `${name} · 36 retained snapshots` },
    ],
  };
}

/** Per-type config schemas for the dedicated add-connector flow. */
export const connectorSchemas: ConnectorTypeSchema[] = [
  {
    category: 'virtualization',
    type: 'proxmox',
    displayName: 'Proxmox VE',
    fields: [
      { name: 'url', label: 'API URL', kind: 'string', required: true, placeholder: 'https://10.0.0.11:8006' },
      { name: 'tokenId', label: 'API token ID', kind: 'string', required: true, placeholder: 'root@pam!wiselabz' },
      { name: 'tokenSecret', label: 'API token secret', kind: 'password', required: true, secret: true },
      { name: 'verifyTls', label: 'Verify TLS certificate', kind: 'boolean', required: false },
    ],
  },
  {
    category: 'containers_paas',
    type: 'portainer',
    displayName: 'Portainer',
    fields: [
      { name: 'url', label: 'API URL', kind: 'string', required: true, placeholder: 'https://10.0.0.20:9443' },
      { name: 'apiKey', label: 'API key', kind: 'password', required: true, secret: true },
      { name: 'verifyTls', label: 'Verify TLS certificate', kind: 'boolean', required: false },
    ],
  },
  {
    category: 'networking',
    type: 'opnsense',
    displayName: 'OPNsense',
    fields: [
      { name: 'url', label: 'API URL', kind: 'string', required: true, placeholder: 'https://10.0.0.1' },
      { name: 'apiKey', label: 'API key', kind: 'string', required: true },
      { name: 'apiSecret', label: 'API secret', kind: 'password', required: true, secret: true },
      { name: 'verifyTls', label: 'Verify TLS certificate', kind: 'boolean', required: false },
    ],
  },
];

export function alertPage(): AlertPage {
  return { items: alerts, total: alerts.length, page: 1, pageSize: alerts.length };
}

// ─── Quality findings ───────────────────────────────────────────────────────

export const findings: QualityFinding[] = [
  {
    id: 'find-1',
    connectorId: 'svc-pve1',
    connectorName: 'pve1',
    docId: 'doc-pve1',
    ruleId: null,
    checkType: 'stale',
    severity: 'info',
    title: 'Doc not reconciled after last sync',
    description: 'k8s-worker-03 was added 8 months ago but the pve1 doc still lists only two worker nodes.',
    remediationLink: '/docs/doc-pve1',
    status: 'open',
    detectedCount: 4,
    firstDetectedAt: daysAgo(240),
    lastSeenAt: daysAgo(240),
    resolvedAt: null,
  },
  {
    id: 'find-2',
    connectorId: 'svc-opnsense',
    connectorName: 'opnsense',
    docId: 'doc-opnsense',
    ruleId: null,
    checkType: 'ownership_incomplete',
    severity: 'critical',
    title: 'WAN firewall rule has no documented owner',
    description: 'The permissive WAN rule flagged in alerts has no linked runbook or owning team in the opnsense doc.',
    remediationLink: '/docs/doc-opnsense',
    status: 'open',
    detectedCount: 12,
    firstDetectedAt: daysAgo(180),
    lastSeenAt: daysAgo(180),
    resolvedAt: null,
  },
  {
    id: 'find-3',
    connectorId: 'svc-portainer',
    connectorName: 'portainer',
    docId: 'doc-portainer',
    ruleId: null,
    checkType: 'failing',
    severity: 'info',
    title: 'Container restart loop undocumented',
    description: 'sonarr in stack media has been restarting repeatedly with no incident note in the portainer doc.',
    remediationLink: '/docs/doc-portainer',
    status: 'open',
    detectedCount: 3,
    firstDetectedAt: daysAgo(90),
    lastSeenAt: daysAgo(90),
    resolvedAt: null,
  },
  {
    id: 'find-4',
    connectorId: 'svc-pbs',
    connectorName: 'pbs',
    docId: 'doc-pbs',
    ruleId: null,
    checkType: 'empty',
    severity: 'info',
    title: 'Backup retention policy section is empty',
    description: 'The pbs doc has a "Retention policy" heading with no content underneath it.',
    remediationLink: '/docs/doc-pbs',
    status: 'open',
    detectedCount: 6,
    firstDetectedAt: daysAgo(150),
    lastSeenAt: daysAgo(150),
    resolvedAt: null,
  },
];

export function findingPage(opts: {
  status?: QualityFindingStatus;
  checkType?: QualityCheckType;
  connectorId?: string;
  page?: number;
  pageSize?: number;
}): QualityFindingPage {
  const page = opts.page ?? 1;
  const pageSize = opts.pageSize ?? 20;
  const items = findings.filter(
    (f) =>
      (!opts.status || f.status === opts.status) &&
      (!opts.checkType || f.checkType === opts.checkType) &&
      (!opts.connectorId || f.connectorId === opts.connectorId),
  );
  const start = (page - 1) * pageSize;
  return { items: items.slice(start, start + pageSize), total: items.length, page, pageSize };
}

export function resolveFinding(findingId: string): QualityFinding | null {
  const finding = findings.find((f) => f.id === findingId);
  if (!finding) return null;
  finding.status = 'resolved';
  finding.resolvedAt = new Date().toISOString();
  return finding;
}

// Deterministic raw-state snapshots per connector category, rendered on the service
// detail page. Markdown bodies so the existing renderer shows real structure.
const snapshotSections: Record<string, { title: string; content: string; order: number }[]> = {
  virtualization: [
    {
      title: 'Nodes',
      order: 1,
      content: '| Node | CPU | Memory | Status |\n|---|---|---|---|\n| node-01 | 12% | 38/128 GB | `online` |\n| node-02 | 7% | 22/128 GB | `online` |',
    },
    {
      title: 'Virtual machines',
      order: 2,
      content: '- **101** `dns-01` — running, 2 vCPU, 2 GB\n- **102** `web-01` — running, 4 vCPU, 8 GB\n- **110** `backup` — stopped, 1 vCPU, 1 GB',
    },
    { title: 'Storage', order: 3, content: '`local-zfs` — 1.2 / 4.0 TB used (30%)' },
  ],
  containers_paas: [
    {
      title: 'Stacks',
      order: 1,
      content: '- **media** — 4 containers, all healthy\n- **monitoring** — 3 containers, all healthy',
    },
    {
      title: 'Containers',
      order: 2,
      content: '| Name | Image | State |\n|---|---|---|\n| jellyfin | jellyfin:10.9 | `running` |\n| grafana | grafana:11 | `running` |\n| prometheus | prom:v2.53 | `running` |',
    },
  ],
  networking: [
    {
      title: 'Interfaces',
      order: 1,
      content: '- **WAN** — up, 1 Gb/s\n- **LAN** — up, 2.5 Gb/s\n- **VLAN20 (iot)** — up',
    },
    {
      title: 'Firewall rules',
      order: 2,
      content: '`pass` LAN → any · `block` IOT → LAN · `pass` WAN → web-01:443',
    },
    { title: 'DHCP leases', order: 3, content: '42 active leases on LAN, 8 on IOT.' },
  ],
};

export function serviceSnapshot(connectorId: string): ServiceSnapshot | null {
  const c = connectors.find((x) => x.id === connectorId);
  if (!c) return null;
  return {
    serviceName: c.name,
    type: c.type,
    sections: snapshotSections[c.category] ?? [],
    fetchedAt: minsAgo(4),
  };
}

/** Recent sync run history per connector, newest first — for the service detail page. */
export function syncRunsFor(connectorId: string): SyncRun[] {
  const c = connectors.find((x) => x.id === connectorId);
  if (!c) return [];
  const runs: SyncRun[] = [];
  // Most recent run mirrors the connector's own lastSync* fields so the two stay
  // consistent; older runs behind it are synthesized as clean successes.
  runs.push({
    id: `sync-${c.id}-1`,
    connectorId: c.id,
    startedAt: c.lastSyncAt ?? minsAgo(4),
    finishedAt: c.lastSyncAt ?? minsAgo(4),
    durationMs: c.lastSyncDurationMs ?? null,
    status: c.lastSyncError ? 'error' : 'success',
    error: c.lastSyncError,
    attempt: (c.retryCount ?? 0) + 1,
    changesCount: c.lastSyncError ? 0 : 2,
    alertsCount: c.lastSyncError ? 1 : 0,
  });
  for (let i = 2; i <= 4; i++) {
    runs.push({
      id: `sync-${c.id}-${i}`,
      connectorId: c.id,
      startedAt: minsAgo(4 + i * (c.scheduleSeconds ? c.scheduleSeconds / 60 : 30)),
      finishedAt: minsAgo(4 + i * (c.scheduleSeconds ? c.scheduleSeconds / 60 : 30)),
      durationMs: 800 + i * 120,
      status: 'success',
      attempt: 1,
      changesCount: i === 2 ? 1 : 0,
      alertsCount: 0,
    });
  }
  return runs;
}
