import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { fireEvent, render, screen, cleanup } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import '../../i18n';
import { RulesPage } from './RulesPage';

const rule = { id: 'r1', name: 'No privileged containers', connectorType: 'docker', entityKind: 'container', conditions: [{ attribute: 'privileged', op: 'eq', value: true }], severity: 'critical', title: 'Privileged container', remediationLink: '', enabled: false, createdAt: '', updatedAt: '' };
const server = setupServer(
  http.get('/api/compliance/rules', () => HttpResponse.json({ items: [rule] })),
  http.get('/api/compliance/schema', () => HttpResponse.json({ docker: { container: [{ name: 'privileged', type: 'boolean', description: '' }] } })),
);
beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => { server.resetHandlers(); cleanup(); });
afterAll(() => server.close());

function renderPage() { return render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><RulesPage /></QueryClientProvider>); }

describe('RulesPage', () => {
  it('renders a rule and toggles it through the API', async () => {
    let enabled = false;
    server.use(http.put('/api/compliance/rules/r1', async ({ request }) => { enabled = (await request.json() as typeof rule).enabled; return HttpResponse.json({ ...rule, enabled }); }));
    renderPage();
    expect(await screen.findByText('No privileged containers')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('switch', { name: 'Enable No privileged containers' }));
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(enabled).toBe(true);
  });
});
