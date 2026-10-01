import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AxiosError } from 'axios';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import '../../i18n';
import { ProfilePage } from './ProfilePage';

// jsdom doesn't implement <dialog>'s imperative API (used by ui/Dialog.tsx).
HTMLDialogElement.prototype.showModal = vi.fn(function (this: HTMLDialogElement) {
  this.open = true;
});
HTMLDialogElement.prototype.close = vi.fn(function (this: HTMLDialogElement) {
  this.open = false;
});

const { postMeMfaTotpMock, confirmMock, beginKeyMock, finishKeyMock, startRegistrationMock, elevateMock } = vi.hoisted(() => ({
  elevateMock: vi.fn(),
  postMeMfaTotpMock: vi.fn(),
  confirmMock: vi.fn(),
  beginKeyMock: vi.fn(),
  finishKeyMock: vi.fn(),
  startRegistrationMock: vi.fn(),
}));

vi.mock('@simplewebauthn/browser', () => ({ startRegistration: startRegistrationMock }));

let webauthnAvailable = false;

vi.mock('qrcode', () => ({
  toDataURL: vi.fn().mockResolvedValue('data:image/png;base64,fake'),
}));

const me = {
  id: 'u1',
  username: 'alice',
  displayName: 'Alice',
  email: 'alice@example.com',
  role: 'operator' as const,
  authSource: 'local' as const,
  digestCadence: 'off' as const,
  digestTimezone: 'UTC',
  createdAt: '2025-01-01T00:00:00Z',
};

vi.mock('../../api/generated/me/me', () => ({
  useGetMe: () => ({ data: me, isLoading: false, isError: false }),
  getGetMeQueryKey: () => ['getMe'],
  patchMe: vi.fn(),
  postMePassword: vi.fn(),
  getMeSessions: vi.fn(() => Promise.resolve([])),
  deleteMeSessionsSessionId: vi.fn(),
  getGetMeSessionsQueryKey: () => ['getMeSessions'],
  useGetMeMfa: () => ({
    data: { factors: [], recoveryCodesRemaining: 0, required: false, webauthnAvailable },
    isLoading: false,
    isError: false,
    refetch: vi.fn(),
  }),
  getGetMeMfaQueryKey: () => ['getMeMfa'],
  postMeMfaTotp: postMeMfaTotpMock,
  postMeMfaTotpFactorIdConfirm: confirmMock,
  postMeMfaWebauthnRegisterBegin: beginKeyMock,
  postMeMfaWebauthnRegisterFinish: finishKeyMock,
  postMeMfaRecoveryCodes: vi.fn(),
  deleteMeMfaFactorsFactorId: vi.fn(),
}));

vi.mock('../../api/generated/auth/auth', () => ({
  useGetAuthApiKeys: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
  postAuthApiKeys: vi.fn(),
  deleteAuthApiKeysId: vi.fn(),
  getGetAuthApiKeysQueryKey: () => ['getAuthApiKeys'],
  useGetAuthElevateMethods: () => ({ data: { methods: ['password'] }, isLoading: false }),
  postAuthElevate: elevateMock,
  postAuthElevateWebauthnBegin: vi.fn(),
}));

// What the server answers when a start-enrollment call lacks X-Elevation-Token.
function elevationRequired() {
  return new AxiosError('elevation required', '400', undefined, undefined, {
    status: 400,
    data: { code: 'elevation_required' },
  } as never);
}

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectors: () => ({ data: [] }),
}));

function renderProfilePage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ProfilePage />
    </QueryClientProvider>
  );
}

describe('Settings > Profile > Security (#279 enrollment happy path)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    webauthnAvailable = false;
    postMeMfaTotpMock.mockResolvedValue({
      factorId: 'f1',
      secret: 'JBSWY3DPEHPK3PXP',
      otpauthUrl: 'otpauth://totp/WiseLabz:alice?secret=JBSWY3DPEHPK3PXP&issuer=WiseLabz',
    });
    confirmMock.mockResolvedValue({
      factor: { id: 'f1', type: 'totp', name: 'Authenticator app', createdAt: '2025-01-01T00:00:00Z' },
      recoveryCodes: ['abcde-fghjk', 'klmno-pqrst'],
    });
  });

  afterEach(cleanup);

  it('walks enrollment from empty state to saved recovery codes', async () => {
    renderProfilePage();

    fireEvent.click(screen.getByRole('button', { name: /set up authenticator app/i }));

    await waitFor(() => expect(postMeMfaTotpMock).toHaveBeenCalled());
    await waitFor(() => expect(screen.getByText('JBSWY3DPEHPK3PXP')).toBeInTheDocument());

    fireEvent.change(screen.getByLabelText(/enter the 6-digit code/i), { target: { value: '123456' } });
    fireEvent.click(screen.getByRole('button', { name: /^confirm$/i }));

    await waitFor(() =>
      expect(confirmMock).toHaveBeenCalledWith('f1', { code: '123456' })
    );

    // Recovery codes are shown once, gated behind the "I saved these" checkbox.
    await waitFor(() => expect(screen.getByText('abcde-fghjk')).toBeInTheDocument());
    const doneButton = screen.getByRole('button', { name: /^done$/i });
    expect(doneButton).toBeDisabled();

    fireEvent.click(screen.getByRole('checkbox', { name: /saved these recovery codes/i }));
    expect(doneButton).not.toBeDisabled();

    fireEvent.click(doneButton);
    await waitFor(() => expect(screen.queryByText('abcde-fghjk')).not.toBeInTheDocument());
  });

  it('registers a security key and shows the first-factor recovery codes', async () => {
    webauthnAvailable = true;
    Object.defineProperty(window, 'isSecureContext', { configurable: true, value: true });
    beginKeyMock.mockResolvedValue({ publicKey: { challenge: 'challenge' } });
    startRegistrationMock.mockResolvedValue({ id: 'credential', response: {} });
    finishKeyMock.mockResolvedValue({
      factor: { id: 'k1', type: 'webauthn', name: 'Desk key', createdAt: '2025-01-01T00:00:00Z' },
      recoveryCodes: ['abcde-fghjk'],
    });
    renderProfilePage();

    fireEvent.click(screen.getByRole('button', { name: /add security key/i }));
    await screen.findByRole('button', { name: /register key/i });
    fireEvent.change(screen.getByPlaceholderText('Security key'), { target: { value: 'Desk key' } });
    fireEvent.click(screen.getByRole('button', { name: /register key/i }));

    await waitFor(() => expect(beginKeyMock).toHaveBeenCalledWith({ name: 'Desk key' }, undefined));
    await waitFor(() => expect(startRegistrationMock).toHaveBeenCalledWith({ optionsJSON: { challenge: 'challenge' } }));
    await waitFor(() => expect(finishKeyMock).toHaveBeenCalledWith({ id: 'credential', response: {} }));
    await waitFor(() => expect(screen.getByText('abcde-fghjk')).toBeInTheDocument());
  });

  it('asks for step-up when starting TOTP enrollment needs it, then retries with the token', async () => {
    elevateMock.mockResolvedValue({ token: 'elev-1' });
    postMeMfaTotpMock.mockRejectedValueOnce(elevationRequired());
    renderProfilePage();

    fireEvent.click(screen.getByRole('button', { name: /set up authenticator app/i }));
    fireEvent.change(await screen.findByLabelText(/confirm your password/i), { target: { value: 'hunter22' } });
    fireEvent.click(screen.getByRole('button', { name: /^verify$/i }));

    await waitFor(() => expect(elevateMock).toHaveBeenCalledWith({ password: 'hunter22', action: 'mfa.manage' }));
    await waitFor(() =>
      expect(postMeMfaTotpMock).toHaveBeenLastCalledWith({}, { headers: { 'X-Elevation-Token': 'elev-1' } })
    );
    await waitFor(() => expect(screen.getByText('JBSWY3DPEHPK3PXP')).toBeInTheDocument());
  });

  it('asks for step-up when registering a security key needs it, then retries with the token', async () => {
    webauthnAvailable = true;
    Object.defineProperty(window, 'isSecureContext', { configurable: true, value: true });
    elevateMock.mockResolvedValue({ token: 'elev-2' });
    beginKeyMock.mockRejectedValueOnce(elevationRequired()).mockResolvedValue({ publicKey: { challenge: 'challenge' } });
    startRegistrationMock.mockResolvedValue({ id: 'credential', response: {} });
    finishKeyMock.mockResolvedValue({ factor: { id: 'k1', type: 'webauthn', name: 'Desk key', createdAt: '2025-01-01T00:00:00Z' } });
    renderProfilePage();

    fireEvent.click(screen.getByRole('button', { name: /add security key/i }));
    fireEvent.change(await screen.findByPlaceholderText('Security key'), { target: { value: 'Desk key' } });
    fireEvent.click(screen.getByRole('button', { name: /register key/i }));
    fireEvent.change(await screen.findByLabelText(/confirm your password/i), { target: { value: 'hunter22' } });
    fireEvent.click(screen.getByRole('button', { name: /^verify$/i }));

    await waitFor(() =>
      expect(beginKeyMock).toHaveBeenLastCalledWith({ name: 'Desk key' }, { headers: { 'X-Elevation-Token': 'elev-2' } })
    );
    await waitFor(() => expect(finishKeyMock).toHaveBeenCalled());
  });
});
