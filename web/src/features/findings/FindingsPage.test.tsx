import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cleanup, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import '../../i18n';
import { FindingsPage } from './FindingsPage';

const server = setupServer(http.get('/api/connectors', () => HttpResponse.json([])));

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => {
  server.resetHandlers(http.get('/api/connectors', () => HttpResponse.json([])));
  cleanup();
});
afterAll(() => server.close());

function renderFindings() {
  return render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <FindingsPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('FindingsPage credential_rotation (#239 PR1)', () => {
  it('lists credential_rotation in the check type filter', async () => {
    server.use(http.get('/api/findings', () => HttpResponse.json({ items: [], total: 0, page: 1, pageSize: 20 })));
    renderFindings();
    expect(await screen.findByRole('option', { name: 'Credential rotation due' })).toBeInTheDocument();
  });

  it('shows a credential_rotation finding card', async () => {
    server.use(
      http.get('/api/findings', () =>
        HttpResponse.json({
          items: [
            {
              id: 'f1',
              connectorId: 'c1',
              connectorName: 'pve1',
              checkType: 'credential_rotation',
              severity: 'critical',
              title: 'Credential rotation due',
              description: 'Secret last rotated 95 days ago; due 2026-01-01T00:00:00Z.',
              status: 'open',
              detectedCount: 1,
              firstDetectedAt: '2026-01-01T00:00:00Z',
              lastSeenAt: '2026-01-01T00:00:00Z',
              remediationLink: '/connectors/c1/edit',
            },
          ],
          total: 1,
          page: 1,
          pageSize: 20,
        }),
      ),
    );
    renderFindings();
    expect(await screen.findByText('Credential rotation due')).toBeInTheDocument();
    expect(await screen.findByText(/secret last rotated 95 days ago/i)).toBeInTheDocument();
  });

  it('shows compliance in the filter and links an admin to its rule', async () => {
    server.use(
      http.get('/api/findings', () => HttpResponse.json({ items: [{ id: 'f1', connectorId: 'c1', connectorName: 'docker', docId: null, ruleId: 'r1', checkType: 'compliance', severity: 'warning', title: 'Host network', description: 'container', remediationLink: '/services/c1', status: 'open', detectedCount: 1, firstDetectedAt: '2026-01-01T00:00:00Z', lastSeenAt: '2026-01-01T00:00:00Z', resolvedAt: null }], total: 1, page: 1, pageSize: 20 })),
      http.get('/api/me', () => HttpResponse.json({ id: 'u1', username: 'admin', role: 'admin', instanceAdminRole: 'admin' })),
    );
    renderFindings();
    expect(await screen.findByRole('option', { name: 'Compliance' })).toBeInTheDocument();
    expect(await screen.findByRole('link', { name: 'View rule' })).toHaveAttribute('href', '/settings/compliance');
  });
});
