import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import '../../i18n';
import { ProfilePage } from './ProfilePage';

const {
  saveProfileMock,
  changePasswordMock,
  createApiKeyMock,
  revokeApiKeyMock,
  revokeSessionMock,
} = vi.hoisted(() => ({
  saveProfileMock: vi.fn(),
  changePasswordMock: vi.fn(),
  createApiKeyMock: vi.fn(),
  revokeApiKeyMock: vi.fn(),
  revokeSessionMock: vi.fn(),
}));

let me = {
  id: 'u1',
  username: 'testuser',
  displayName: 'Test User',
  email: 'test@example.com',
  role: 'operator' as const,
  authSource: 'local' as const,
  digestCadence: 'off' as 'off' | 'daily' | 'weekly',
  digestTimezone: '',
  createdAt: '2025-01-01T00:00:00Z',
  updatedAt: '2025-01-01T00:00:00Z',
};

vi.mock('../../api/generated/me/me', () => ({
  useGetMe: () => ({ data: me, isLoading: false, isError: false }),
  getGetMeQueryKey: () => ['getMe'],
  patchMe: saveProfileMock,
  postMePassword: changePasswordMock,
  getMeSessions: vi.fn(() => Promise.resolve([])),
  deleteMeSessionsSessionId: revokeSessionMock,
  getGetMeSessionsQueryKey: () => ['getMeSessions'],
}));

vi.mock('../../api/generated/auth/auth', () => ({
  useGetAuthApiKeys: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
  postAuthApiKeys: createApiKeyMock,
  deleteAuthApiKeysId: revokeApiKeyMock,
  getGetAuthApiKeysQueryKey: () => ['getAuthApiKeys'],
}));

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectors: () => ({
    data: [
      { id: 'c1', name: 'pve' },
      { id: 'c2', name: 'pfsense' },
    ],
  }),
}));

function renderProfilePage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <ProfilePage />
    </QueryClientProvider>
  );
}

describe('ProfilePage (#237 Phase 4 - Digest Settings)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    saveProfileMock.mockResolvedValue({ ...me });
  });

  afterEach(() => {
    cleanup();
  });

  it('renders digest cadence and timezone fields when me.digestTimezone is empty', () => {
    me = { ...me, digestTimezone: '' };
    renderProfilePage();

    const cadenceField = screen.getByLabelText(/digest cadence/i) as HTMLSelectElement;
    const timezoneField = screen.getByLabelText(/digest timezone/i) as HTMLInputElement;

    expect(cadenceField).toBeInTheDocument();
    expect(timezoneField).toBeInTheDocument();

    // When digestTimezone is empty, it should be auto-detected (filled by Intl.DateTimeFormat)
    expect(timezoneField.value).toBeTruthy();
    expect(timezoneField.value).toMatch(/^[A-Za-z_/]+$/); // IANA format
  });

  it('pre-fills digest cadence from me.digestCadence', () => {
    me = { ...me, digestCadence: 'daily' };
    renderProfilePage();

    const cadenceSelect = screen.getByDisplayValue('Daily') as HTMLSelectElement;
    expect(cadenceSelect.value).toBe('daily');
  });

  it('pre-fills digest timezone from me.digestTimezone when set', () => {
    me = { ...me, digestTimezone: 'America/New_York' };
    renderProfilePage();

    const timezoneField = screen.getByDisplayValue('America/New_York') as HTMLInputElement;
    expect(timezoneField.value).toBe('America/New_York');
  });

  it('save button is disabled when nothing has changed', () => {
    me = { ...me, digestCadence: 'daily', digestTimezone: 'Europe/London' };
    renderProfilePage();

    const saveButton = screen.getByRole('button', { name: /save/i });
    expect(saveButton).toBeDisabled();
  });

  it('save button becomes enabled when digest cadence is changed', () => {
    me = { ...me, digestCadence: 'off', digestTimezone: 'America/New_York' };
    renderProfilePage();

    const cadenceSelect = screen.getByLabelText(/digest cadence/i) as HTMLSelectElement;
    fireEvent.change(cadenceSelect, { target: { value: 'daily' } });

    const saveButton = screen.getByRole('button', { name: /save/i });
    expect(saveButton).not.toBeDisabled();
  });

  it('save button becomes enabled when digest timezone is changed', () => {
    me = { ...me, digestCadence: 'daily', digestTimezone: 'America/New_York' };
    renderProfilePage();

    const timezoneField = screen.getByLabelText(/digest timezone/i) as HTMLInputElement;
    fireEvent.change(timezoneField, { target: { value: 'Europe/London' } });

    const saveButton = screen.getByRole('button', { name: /save/i });
    expect(saveButton).not.toBeDisabled();
  });

  it('clicking save calls patchMe with digestCadence and digestTimezone', async () => {
    me = { ...me, digestCadence: 'off', digestTimezone: 'America/New_York' };
    renderProfilePage();

    const cadenceSelect = screen.getByLabelText(/digest cadence/i) as HTMLSelectElement;
    fireEvent.change(cadenceSelect, { target: { value: 'daily' } });

    const saveButton = screen.getByRole('button', { name: /save/i });
    fireEvent.click(saveButton);

    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(saveProfileMock).toHaveBeenCalledWith(
      expect.objectContaining({
        displayName: 'Test User',
        email: 'test@example.com',
        digestCadence: 'daily',
        digestTimezone: 'America/New_York',
      })
    );
  });

  it('shows options for all digest cadence values (off, daily, weekly)', () => {
    me = { ...me, digestCadence: 'off' };
    renderProfilePage();

    const cadenceSelect = screen.getByLabelText(/digest cadence/i) as HTMLSelectElement;
    const options = Array.from(cadenceSelect.options).map((o) => o.value);

    expect(options).toContain('off');
    expect(options).toContain('daily');
    expect(options).toContain('weekly');
  });
});

describe('ProfilePage API keys (#278 scopes)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    createApiKeyMock.mockResolvedValue({ token: 'wlz_x' });
  });

  afterEach(() => {
    cleanup();
  });

  it('sends the chosen scope and connector restriction', async () => {
    renderProfilePage();

    fireEvent.change(screen.getByLabelText(/^name$/i), { target: { value: 'ci' } });
    fireEvent.change(screen.getByLabelText(/^access/i), { target: { value: 'read' } });
    fireEvent.click(screen.getByLabelText('pfsense'));
    fireEvent.click(screen.getByRole('button', { name: /create key/i }));

    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(createApiKeyMock).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'ci', scope: 'read', connectorIds: ['c2'] })
    );
  });

  it('omits connectorIds when no connector is checked', async () => {
    renderProfilePage();

    fireEvent.change(screen.getByLabelText(/^name$/i), { target: { value: 'ci' } });
    fireEvent.click(screen.getByRole('button', { name: /create key/i }));

    await new Promise((resolve) => setTimeout(resolve, 0));

    expect(createApiKeyMock).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'ci', scope: 'full', connectorIds: undefined })
    );
  });
});
