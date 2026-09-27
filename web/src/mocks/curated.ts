/**
 * Hand-authored MSW handlers for the surfaces we showcase. Registered BEFORE the
 * generated faker handlers in handlers.ts, so MSW's first-match wins and these
 * deterministic fixtures render instead of random faker data. Anything not
 * covered here still falls through to the generated mocks.
 */
import { http, HttpResponse, delay } from 'msw';
import {
  alertPage,
  changeDetail,
  changePage,
  connectorSchemas,
  connectors,
  dashboardOverview,
  docTree,
  docVersionContent,
  docVersionMeta,
  docs,
  findingPage,
  findings,
  removalImpact,
  resolveFinding,
  serviceSnapshot,
  syncRunsFor,
  user,
} from '../data/fixtures';
import type { QualityCheckType, QualityFindingStatus } from '../api/model';

// Small artificial latency so loading skeletons are actually exercised on first paint.
const LATENCY = 280;
const goldenSnapshots = new Map<string, string>();

function snapshotHistory(connectorId: string) {
  const current = serviceSnapshot(connectorId);
  if (!current) return [];
  return [0, 1, 2, 3].map((index) => ({
    id: `${connectorId}-snapshot-${4 - index}`,
    connectorId,
    serviceName: current.serviceName,
    type: current.type,
    sections: current.sections.map(({ title, content }) => ({
      title,
      content: index === 0 ? content : `${content}\n\nRevision ${4 - index}`,
    })),
    entities: [{ kind: 'host', name: current.serviceName, ip: `192.168.1.${10 + index}`, attributes: { revision: 4 - index } }],
    dependencies: [{ kind: 'service', name: 'gateway', ref: 'gateway-1' }],
    metadata: { source: 'curated' },
    fetchedAt: new Date(Date.now() - index * 3_600_000).toISOString(),
  }));
}

// Session-mutable auth config so the Settings step-up toggle persists across reads.
const authConfig = {
  localEnabled: true,
  accessTokenTtl: 900,
  refreshTokenTtl: 2592000,
  stepUpForDestructive: true,
  oidcProviders: [] as unknown[],
};

// Mock session state. Login/oidc set it true; refresh succeeds only when set (i.e.
// a refresh "cookie" exists); logout clears it. Resets on full page reload, which
// matches the in-memory-token model (no persistent client session).
let hasSession = false;
const newSession = () => ({
  accessToken: `acc-${Math.random().toString(36).slice(2)}`,
  expiresIn: authConfig.accessTokenTtl,
  user,
});

export const curatedHandlers = [
  http.get('*/auth/providers', async () => {
    await delay(LATENCY);
    return HttpResponse.json({
      localEnabled: authConfig.localEnabled,
      oidc: [
        {
          id: 'authentik',
          displayName: 'Authentik',
          authUrl: '/auth/callback?providerId=authentik&code=mock-oidc-code&state=mock-state',
        },
      ],
    });
  }),

  http.post('*/auth/login', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as { username?: string; password?: string };
    if (!body.username || !body.password) {
      return HttpResponse.json({ code: 'invalid-credentials', message: 'Invalid credentials' }, { status: 401 });
    }
    hasSession = true;
    return HttpResponse.json(newSession());
  }),

  http.post('*/auth/oidc/callback', async () => {
    await delay(LATENCY);
    hasSession = true;
    return HttpResponse.json(newSession());
  }),

  http.post('*/auth/refresh', async () => {
    await delay(LATENCY);
    if (!hasSession) {
      return HttpResponse.json({ code: 'no-session', message: 'No active session' }, { status: 401 });
    }
    return HttpResponse.json(newSession());
  }),

  http.post('*/auth/logout', async () => {
    await delay(LATENCY);
    hasSession = false;
    return new HttpResponse(null, { status: 204 });
  }),

  http.post('*/ws/ticket', () => HttpResponse.json({ ticket: 'mock-ticket' })),

  http.get('*/auth/config', async () => {
    await delay(LATENCY);
    return HttpResponse.json(authConfig);
  }),

  http.put('*/auth/config', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as Partial<typeof authConfig>;
    Object.assign(authConfig, body);
    return HttpResponse.json(authConfig);
  }),

  // Manager mutations — operate on the fixture connectors in place so toggles and
  // removals visibly persist for the session.
  http.put('*/connectors/:connectorId/enabled', async ({ params, request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as { enabled?: boolean };
    const c = connectors.find((x) => x.id === params.connectorId);
    if (!c) return new HttpResponse(null, { status: 404 });
    c.enabled = body.enabled ?? c.enabled;
    return HttpResponse.json(c);
  }),

  http.post('*/connectors/:connectorId/test', async () => {
    await delay(LATENCY);
    return HttpResponse.json({ ok: true, message: 'Connection succeeded', latencyMs: 42 });
  }),

  http.post('*/connectors/:connectorId/sync', async ({ params }) => {
    await delay(LATENCY);
    return HttpResponse.json({ jobId: `job-${Date.now().toString(36)}`, serviceId: params.connectorId });
  }),

  http.post('*/sync', async () => {
    await delay(LATENCY);
    return HttpResponse.json({ jobId: `job-${Date.now().toString(36)}`, serviceId: null });
  }),

  // Bulk connector actions — deterministic against the submitted ids and the
  // fixture connector set, instead of the generated faker mock's random 1-4
  // fake ids/statuses unrelated to what was actually selected.
  http.post('*/connectors/bulk-sync', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as { ids?: string[] };
    const results = (body.ids ?? []).map((id) =>
      connectors.some((c) => c.id === id)
        ? { id, status: 'success' as const, jobId: `job-${Date.now().toString(36)}` }
        : { id, status: 'error' as const, reason: 'not_found' }
    );
    return HttpResponse.json({ results });
  }),

  http.post('*/connectors/bulk-reauth', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as { ids?: string[] };
    const results = (body.ids ?? []).map((id) =>
      connectors.some((c) => c.id === id)
        ? { id, status: 'success' as const }
        : { id, status: 'error' as const, reason: 'not_found' }
    );
    return HttpResponse.json({ results });
  }),

  http.post('*/connectors/bulk-restart', async ({ request }) => {
    await delay(LATENCY);
    if (authConfig.stepUpForDestructive && !request.headers.get('X-Elevation-Token')) {
      return HttpResponse.json(
        { code: 'elevation-required', message: 'Step-up required' },
        { status: 400 },
      );
    }
    const body = (await request.json().catch(() => ({}))) as { ids?: string[] };
    const results = (body.ids ?? []).map((id) =>
      connectors.some((c) => c.id === id)
        ? { id, status: 'success' as const }
        : { id, status: 'error' as const, reason: 'not_found' }
    );
    return HttpResponse.json({ results });
  }),

  http.post('*/connectors', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
    const created = {
      id: `svc-${Date.now().toString(36)}`,
      name: String(body.name ?? 'new-connector'),
      category: body.category,
      type: String(body.type ?? 'custom'),
      enabled: true,
      status: 'online',
      url: String(body.url ?? ''),
      verifyTls: Boolean(body.verifyTls),
      lastSyncAt: '',
    };
    connectors.push(created as (typeof connectors)[number]);
    return HttpResponse.json(created, { status: 201 });
  }),

  // Destructive: honour step-up — reject when the toggle is on and the elevation
  // token header is absent, otherwise drop the connector.
  http.delete('*/connectors/:connectorId', async ({ params, request }) => {
    await delay(LATENCY);
    if (authConfig.stepUpForDestructive && !request.headers.get('X-Elevation-Token')) {
      return HttpResponse.json(
        { code: 'elevation-required', message: 'Step-up required' },
        { status: 403 },
      );
    }
    const idx = connectors.findIndex((x) => x.id === params.connectorId);
    if (idx >= 0) connectors.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  http.get('*/me', async () => {
    await delay(LATENCY);
    return HttpResponse.json(user);
  }),

  http.get('*/dashboard/overview', async () => {
    await delay(LATENCY);
    return HttpResponse.json(dashboardOverview());
  }),

  http.get('*/connectors/schema', async () => {
    await delay(LATENCY);
    return HttpResponse.json(connectorSchemas);
  }),

  http.get('*/connectors/:connectorId/removal-impact', async ({ params }) => {
    await delay(LATENCY);
    return HttpResponse.json(removalImpact(params.connectorId as string));
  }),

  http.post('*/connectors/:connectorId/restart', async ({ params }) => {
    await delay(LATENCY);
    const connector = connectors.find((item) => item.id === params.connectorId);
    if (!connector) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json({
      targetService: connector.name,
      estimatedDowntimeSeconds: 30,
      dependentServices: [],
    });
  }),

  // Live raw-state snapshot for the service detail page.
  http.get('*/connectors/:connectorId/data', async ({ params }) => {
    await delay(LATENCY);
    const snap = serviceSnapshot(params.connectorId as string);
    if (!snap) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(snap);
  }),

  http.get('*/connectors/:connectorId/snapshots/diff', async ({ params, request }) => {
    await delay(LATENCY);
    const history = snapshotHistory(params.connectorId as string);
    const query = new URL(request.url).searchParams;
    const from = history.find((snapshot) => snapshot.id === query.get('from'));
    const to = history.find((snapshot) => snapshot.id === query.get('to'));
    if (!from || !to) return new HttpResponse(null, { status: 404 });
    const report = {
      provenance: {
        connectorId: params.connectorId,
        connectorName: to.serviceName,
        from: { id: from.id, fetchedAt: from.fetchedAt, sha256: 'mock-from-hash' },
        to: { id: to.id, fetchedAt: to.fetchedAt, sha256: 'mock-to-hash' },
        generatedAt: new Date().toISOString(),
        generatedBy: user.id,
      },
      summary: { sectionsAdded: 0, sectionsRemoved: 0, sectionsModified: 1, entitiesAdded: 0, entitiesRemoved: 0, entitiesModified: 1, dependenciesAdded: 0, dependenciesRemoved: 0 },
      sections: [{ type: 'modified', severity: 'info', summary: 'Revision changed', detail: '', patches: [{ section: from.sections[0]?.title ?? 'State', old: from.sections[0]?.content ?? '', new: to.sections[0]?.content ?? '' }] }],
      entities: [{ kind: 'host', key: to.serviceName, name: to.serviceName, field: 'ip', change: 'modified', old: from.entities[0].ip, new: to.entities[0].ip }],
      dependencies: [],
    };
    const format = query.get('format');
    if (!format) return HttpResponse.json(report);
    const body = format === 'json' ? JSON.stringify(report, null, 2) : format === 'csv' ? 'kind,key,field,change,old,new\n' : format === 'html' ? '<h1>Snapshot diff</h1>' : '# Snapshot diff\n';
    return new HttpResponse(body, { headers: { 'Content-Disposition': `attachment; filename="wiselabz-snapshot-diff-${params.connectorId}.${format}"`, 'Content-Type': format === 'html' ? 'text/html' : 'text/plain' } });
  }),

  http.get('*/connectors/:connectorId/snapshots/:snapshotId', async ({ params }) => {
    await delay(LATENCY);
    const snapshot = snapshotHistory(params.connectorId as string).find((item) => item.id === params.snapshotId);
    return snapshot ? HttpResponse.json(snapshot) : new HttpResponse(null, { status: 404 });
  }),

  http.get('*/connectors/:connectorId/snapshots', async ({ params, request }) => {
    await delay(LATENCY);
    const history = snapshotHistory(params.connectorId as string);
    if (history.length === 0) return new HttpResponse(null, { status: 404 });
    const query = new URL(request.url).searchParams;
    const offset = Number(query.get('cursor') ?? 0);
    const limit = Math.min(Number(query.get('limit') ?? 30), 100);
    const page = history.slice(offset, offset + limit).map((snapshot) => ({ id: snapshot.id, fetchedAt: snapshot.fetchedAt, sizeBytes: JSON.stringify(snapshot).length, golden: goldenSnapshots.get(params.connectorId as string) === snapshot.id }));
    const next = offset + limit < history.length ? String(offset + limit) : '';
    return HttpResponse.json(page, { headers: next ? { 'X-Next-Cursor': next } : {} });
  }),

  http.post('*/connectors/:connectorId/golden-snapshot', async ({ params, request }) => {
    const body = (await request.json()) as { snapshotId?: string };
    const snapshotId = body.snapshotId ?? snapshotHistory(params.connectorId as string)[0]?.id;
    if (!snapshotId) return new HttpResponse(null, { status: 404 });
    goldenSnapshots.set(params.connectorId as string, snapshotId);
    return HttpResponse.json({ connectorId: params.connectorId, snapshotId, pinnedBy: user.id, pinnedAt: new Date().toISOString() });
  }),

  http.get('*/connectors/:connectorId/golden-snapshot', ({ params }) => {
    const snapshotId = goldenSnapshots.get(params.connectorId as string);
    return snapshotId ? HttpResponse.json({ connectorId: params.connectorId, snapshotId, pinnedBy: user.id, pinnedAt: new Date().toISOString() }) : new HttpResponse(null, { status: 404 });
  }),

  http.delete('*/connectors/:connectorId/golden-snapshot', ({ params }) => {
    goldenSnapshots.delete(params.connectorId as string);
    return HttpResponse.json({});
  }),

  http.put('*/connectors/:connectorId', async ({ params, request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
    const c = connectors.find((x) => x.id === params.connectorId);
    if (!c) return new HttpResponse(null, { status: 404 });
    if (typeof body.name === 'string') c.name = body.name;
    if (typeof body.url === 'string') c.url = body.url;
    if (typeof body.verifyTls === 'boolean') c.verifyTls = body.verifyTls;
    if ('scheduleSeconds' in body) c.scheduleSeconds = body.scheduleSeconds as number | null;
    return HttpResponse.json(c);
  }),

  // Recent sync history for the service detail page's sync-history panel.
  http.get('*/connectors/:connectorId/syncs', async ({ params, request }) => {
    await delay(LATENCY);
    const limit = Number(new URL(request.url).searchParams.get('limit')) || undefined;
    const snapshots = snapshotHistory(params.connectorId as string);
    const runs = syncRunsFor(params.connectorId as string).map((run, index) => ({
      ...run,
      snapshotId: run.status === 'success' ? snapshots[index]?.id ?? null : null,
    }));
    return HttpResponse.json(limit ? runs.slice(0, limit) : runs);
  }),

  http.get('*/connectors', async () => {
    await delay(LATENCY);
    return HttpResponse.json(connectors);
  }),

  // Single connector — registered AFTER /schema, /removal-impact, /data so those win.
  http.get('*/connectors/:connectorId', async ({ params }) => {
    await delay(LATENCY);
    const c = connectors.find((x) => x.id === params.connectorId);
    if (!c) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(c);
  }),

  // Step-up: mint a short-lived elevation token for a single destructive action.
  http.post('*/auth/elevate', async () => {
    await delay(LATENCY);
    return HttpResponse.json({
      token: `elev-${Math.random().toString(36).slice(2)}`,
      expiresAt: new Date(Date.now() + 2 * 60_000).toISOString(),
    });
  }),

  http.get('*/changes', async ({ request }) => {
    await delay(LATENCY);
    const serviceId = new URL(request.url).searchParams.get('serviceId');
    const page = changePage();
    if (!serviceId) return HttpResponse.json(page);
    const items = page.items.filter((c) => c.serviceId === serviceId);
    return HttpResponse.json({ ...page, items, total: items.length, pageSize: items.length });
  }),

  http.get('*/changes/:changeId', async ({ params }) => {
    await delay(LATENCY);
    const detail = changeDetail(params.changeId as string);
    if (!detail) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(detail);
  }),

  http.get('*/alerts', async () => {
    await delay(LATENCY);
    return HttpResponse.json(alertPage());
  }),

  http.get('*/findings/:findingId', async ({ params }) => {
    await delay(LATENCY);
    const finding = findings.find((f) => f.id === params.findingId);
    if (!finding) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(finding);
  }),

  http.post('*/findings/:findingId/resolve', async ({ params }) => {
    await delay(LATENCY);
    const finding = resolveFinding(params.findingId as string);
    if (!finding) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(finding);
  }),

  http.get('*/findings', async ({ request }) => {
    await delay(LATENCY);
    const url = new URL(request.url);
    const status = (url.searchParams.get('status') as QualityFindingStatus | null) ?? undefined;
    const checkType = (url.searchParams.get('checkType') as QualityCheckType | null) ?? undefined;
    const connectorId = url.searchParams.get('connectorId') ?? undefined;
    const page = Number(url.searchParams.get('page')) || undefined;
    const pageSize = Number(url.searchParams.get('pageSize')) || undefined;
    return HttpResponse.json(findingPage({ status, checkType, connectorId, page, pageSize }));
  }),

  http.get('*/docs/tree', async () => {
    await delay(LATENCY);
    return HttpResponse.json(docTree);
  }),

  http.get('*/docs/:docId/versions/:rev', async ({ params }) => {
    await delay(LATENCY);
    const byRev = docVersionContent[params.docId as string];
    const v = byRev?.[Number(params.rev)];
    if (!v) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(v);
  }),

  http.get('*/docs/:docId/versions', async ({ params }) => {
    await delay(LATENCY);
    return HttpResponse.json(docVersionMeta[params.docId as string] ?? []);
  }),

  // Save an edited doc → new version, rejecting a stale editor base version.
  http.put('*/docs/:docId', async ({ params, request }) => {
    await delay(LATENCY);
    const id = params.docId as string;
    const doc = docs[id];
    if (!doc) return new HttpResponse(null, { status: 404 });
    const body = (await request.json().catch(() => ({}))) as { content?: string; baseVersion?: number };
    if (body.baseVersion !== undefined && body.baseVersion !== doc.currentVersion) {
      return HttpResponse.json(doc, { status: 409 });
    }
    const rev = doc.currentVersion + 1;
    const now = new Date().toISOString();
    doc.content = body.content ?? doc.content;
    doc.currentVersion = rev;
    doc.updatedAt = now;
    const meta = { rev, createdAt: now, author: 'you', trigger: 'manual' as const };
    (docVersionMeta[id] ??= []).unshift(meta);
    (docVersionContent[id] ??= {})[rev] = { ...meta, content: doc.content };
    return HttpResponse.json(doc);
  }),

  // Batched AI suggestion — returns a ref; the editor renders the proposed revision
  // as a review-diff (no streaming in v1, per the plan).
  http.post('*/docs/:docId/ai-suggest', async () => {
    await delay(LATENCY * 2);
    return HttpResponse.json({ requestId: `ai-${Math.random().toString(36).slice(2)}` });
  }),

  // Restore a past revision → creates a new head version with that content.
  http.post('*/docs/:docId/versions/:rev/restore', async ({ params }) => {
    await delay(LATENCY);
    const id = params.docId as string;
    const doc = docs[id];
    if (!doc) return new HttpResponse(null, { status: 404 });
    const from = docVersionContent[id]?.[Number(params.rev)];
    const rev = doc.currentVersion + 1;
    const now = new Date().toISOString();
    doc.content = from?.content ?? doc.content;
    doc.currentVersion = rev;
    doc.updatedAt = now;
    const meta = { rev, createdAt: now, author: 'you', trigger: 'manual' as const };
    (docVersionMeta[id] ??= []).unshift(meta);
    (docVersionContent[id] ??= {})[rev] = { ...meta, content: doc.content };
    return HttpResponse.json(doc);
  }),

  http.get('*/docs/:docId', async ({ params }) => {
    await delay(LATENCY);
    const doc = docs[params.docId as string];
    if (!doc) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(doc);
  }),

  // Flat searchable doc list — deterministic (filtered from the fixture set)
  // instead of the generated faker mock, so TopologyPage's "does a Lab
  // Topology doc already exist" lookup behaves like the real backend
  // instead of matching random titles.
  http.get('*/docs', async ({ request }) => {
    await delay(LATENCY);
    const search = new URL(request.url).searchParams.get('search')?.toLowerCase() ?? '';
    const items = Object.values(docs).filter((d) => d.title.toLowerCase().includes(search));
    return HttpResponse.json({ items, total: items.length, page: 1, pageSize: items.length || 20 });
  }),

  // Generate (or refresh) the lab-wide topology doc, deterministically —
  // the generated faker mock returns a random docId that no fixture knows
  // about, so the page TopologyPage redirects into 404s every time.
  http.post('*/docs/topology', async () => {
    await delay(LATENCY);
    const existing = Object.values(docs).find((d) => d.title === 'Lab Topology');
    const id = existing?.docId ?? 'doc-topology';
    const now = new Date().toISOString();
    const content = '# Lab Topology\n\nGenerated from the current connector set.\n';
    docs[id] = {
      docId: id,
      title: 'Lab Topology',
      kind: 'lab',
      content,
      currentVersion: (existing?.currentVersion ?? 0) + 1,
      updatedAt: now,
    };
    return HttpResponse.json({ docId: id, title: 'Lab Topology', content });
  }),
];
