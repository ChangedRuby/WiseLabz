/**
 * Curated MSW handlers for the Settings hub (Phase 8). Only endpoints NOT already
 * served by `mocks/curated.ts` live here — `/auth/config` (GET/PUT) and
 * `/auth/providers` are reused from curated, not duplicated.
 *
 * State is module-level and session-mutable so edits (create a user, revoke a
 * session, flip a routing cell, save AI config) visibly persist across reads for
 * the lifetime of the page. Register by spreading `settingsHandlers` in
 * `mocks/handlers.ts` (see this file's PR note) BEFORE the generated faker mocks
 * so these deterministic fixtures win.
 */
import { http, HttpResponse, delay } from 'msw';
import type {
  AiConfig,
  AiFallbackProvider,
  NotificationConfig,
  ProfileUpdate,
  Role,
  User,
  UserCreate,
  UserUpdate,
} from '../api/model';
import { user as currentUser } from '../data/fixtures';
import {
  settingsAiConfig,
  settingsHealth,
  settingsNotificationConfig,
  settingsSessions,
  settingsSystemInfo,
  settingsUsers,
} from '../data/settings.fixtures';

const LATENCY = 280;

// Deep-ish clones so mutating mock state never edits the source fixtures.
const users: User[] = settingsUsers.map((u) => ({ ...u }));
let sessions = settingsSessions.map((s) => ({ ...s }));
const aiConfig: AiConfig = { ...settingsAiConfig };
let aiFallbackProviders: AiFallbackProvider[] = [];
const notificationConfig: NotificationConfig = {
  channels: settingsNotificationConfig.channels.map((c) => ({ ...c })),
  routing: settingsNotificationConfig.routing.map((r) => ({ ...r })),
};

export const settingsHandlers = [
  // ── Profile (current user) ────────────────────────────────────────────────
  http.patch('*/me', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as ProfileUpdate;
    if (body.displayName !== undefined) currentUser.displayName = body.displayName;
    if (body.email !== undefined) currentUser.email = body.email;
    return HttpResponse.json(currentUser);
  }),

  http.post('*/me/password', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as {
      currentPassword?: string;
      newPassword?: string;
    };
    if (!body.currentPassword || !body.newPassword) {
      return HttpResponse.json(
        { code: 'invalid-password', message: 'Both current and new password are required.' },
        { status: 400 },
      );
    }
    return new HttpResponse(null, { status: 204 });
  }),

  http.get('*/me/sessions', async () => {
    await delay(LATENCY);
    return HttpResponse.json(sessions);
  }),

  http.delete('*/me/sessions/:sessionId', async ({ params }) => {
    await delay(LATENCY);
    sessions = sessions.filter((s) => s.id !== params.sessionId || s.current);
    return new HttpResponse(null, { status: 204 });
  }),

  // ── Users (operator) ──────────────────────────────────────────────────────
  http.get('*/users', async () => {
    await delay(LATENCY);
    return HttpResponse.json(users);
  }),

  http.post('*/users', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as UserCreate;
    const created: User = {
      id: `usr-${Math.random().toString(36).slice(2, 8)}`,
      username: body.username,
      displayName: body.username,
      email: body.email,
      role: (body.role ?? 'viewer') as Role,
      authSource: 'local',
      disabled: false,
      createdAt: new Date().toISOString(),
    };
    users.push(created);
    return HttpResponse.json(created, { status: 201 });
  }),

  http.patch('*/users/:userId', async ({ params, request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as UserUpdate;
    const u = users.find((x) => x.id === params.userId);
    if (!u) return new HttpResponse(null, { status: 404 });
    if (body.role !== undefined) u.role = body.role;
    if (body.disabled !== undefined) u.disabled = body.disabled;
    if (body.displayName !== undefined) u.displayName = body.displayName;
    return HttpResponse.json(u);
  }),

  http.delete('*/users/:userId', async ({ params }) => {
    await delay(LATENCY);
    const idx = users.findIndex((x) => x.id === params.userId);
    if (idx >= 0) users.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  http.post('*/users/:userId/reset-password', async ({ params }) => {
    await delay(LATENCY);
    const exists = users.some((x) => x.id === params.userId);
    if (!exists) return new HttpResponse(null, { status: 404 });
    return new HttpResponse(null, { status: 204 });
  }),

  // ── AI module config (operator) ───────────────────────────────────────────
  http.get('*/ai/config', async () => {
    await delay(LATENCY);
    return HttpResponse.json(aiConfig);
  }),

  http.put('*/ai/config', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as Partial<AiConfig>;
    // apiKey is write-only: accept it, never echo it back.
    const { apiKey: _apiKey, ...rest } = body;
    Object.assign(aiConfig, rest);
    return HttpResponse.json(aiConfig);
  }),

  http.post('*/ai/config/test', async () => {
    await delay(700);
    if (!aiConfig.enabled || !aiConfig.provider) {
      return HttpResponse.json({ ok: false, message: 'AI module is disabled or no provider selected.' });
    }
    return HttpResponse.json({ ok: true, message: 'Provider reachable.', latencyMs: 412 });
  }),

  http.get('*/ai/config/fallback-providers', async () => {
    await delay(LATENCY);
    return HttpResponse.json(aiFallbackProviders);
  }),

  http.put('*/ai/config/fallback-providers', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => [])) as AiFallbackProvider[];
    // apiKey is write-only: accept it, never echo it back.
    aiFallbackProviders = body.map(({ apiKey: _apiKey, ...rest }) => rest);
    return HttpResponse.json(aiFallbackProviders);
  }),

  // ── Notifications config (operator) ───────────────────────────────────────
  http.get('*/notifications/config', async () => {
    await delay(LATENCY);
    return HttpResponse.json(notificationConfig);
  }),

  http.put('*/notifications/config', async ({ request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as NotificationConfig;
    if (body.channels) notificationConfig.channels = body.channels;
    if (body.routing) notificationConfig.routing = body.routing;
    return HttpResponse.json(notificationConfig);
  }),

  http.post('*/notifications/config/test', async ({ request }) => {
    await delay(600);
    const body = (await request.json().catch(() => ({}))) as { channel?: string };
    const ch = notificationConfig.channels.find((c) => c.type === body.channel);
    if (!ch?.enabled) {
      return HttpResponse.json({ ok: false, message: 'Channel is disabled.' });
    }
    return HttpResponse.json({ ok: true, message: `Test sent on ${body.channel}.`, latencyMs: 120 });
  }),

  // ── System / instance (operator) ──────────────────────────────────────────
  http.get('*/system/info', async () => {
    await delay(LATENCY);
    return HttpResponse.json(settingsSystemInfo);
  }),

  http.get('*/health', async () => {
    await delay(LATENCY);
    return HttpResponse.json(settingsHealth);
  }),

  // ── OIDC provider enable/disable (operator) — the one missing auth mutation ─
  // Providers themselves are file/env-defined and surfaced read-only via the
  // reused `/auth/config`; this only flips the persisted enabled flag.
  http.put('*/auth/providers/:providerId/enabled', async ({ params, request }) => {
    await delay(LATENCY);
    const body = (await request.json().catch(() => ({}))) as { enabled?: boolean };
    return HttpResponse.json({
      id: String(params.providerId),
      displayName: String(params.providerId),
      issuerUrl: '',
      clientId: '',
      secretConfigured: true,
      enabled: body.enabled ?? true,
      source: 'file',
    });
  }),
];
