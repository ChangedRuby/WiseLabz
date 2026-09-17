import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { ShareLinksPage } from './ShareLinksPage';

const { revokeMock } = vi.hoisted(() => ({ revokeMock: vi.fn() }));

let links: unknown[] = [];
vi.mock('../../api/shareLinks', () => ({
  useListShareLinks: () => ({ data: links, isLoading: false, isError: false, refetch: vi.fn() }),
  useRevokeShareLink: () => ({ mutate: revokeMock, isPending: false }),
}));

describe('ShareLinksPage (#240 PR2)', () => {
  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it('shows an empty state with no links', () => {
    links = [];
    render(<ShareLinksPage />);
    expect(screen.getByText(/no share links yet/i)).toBeInTheDocument();
  });

  it('lists an active link with its docTreeRoot and a revoke button', () => {
    links = [
      {
        id: 's1',
        docTreeRoot: 'c1',
        createdBy: 'u1',
        createdAt: '2025-01-01T00:00:00Z',
        expiresAt: '2099-01-01T00:00:00Z',
        revokedAt: '',
        lastAccessedAt: '',
      },
    ];
    render(<ShareLinksPage />);
    expect(screen.getByText('c1')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /revoke/i })).toBeInTheDocument();
  });

  it('marks an expired link as expired and hides its revoke button', () => {
    links = [
      {
        id: 's1',
        docTreeRoot: 'c1',
        createdBy: 'u1',
        createdAt: '2020-01-01T00:00:00Z',
        expiresAt: '2020-01-02T00:00:00Z',
        revokedAt: '',
        lastAccessedAt: '',
      },
    ];
    render(<ShareLinksPage />);
    expect(screen.getByText(/expired/i)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /revoke/i })).not.toBeInTheDocument();
  });

  it('marks a revoked link as revoked and hides its revoke button', () => {
    links = [
      {
        id: 's1',
        docTreeRoot: 'c1',
        createdBy: 'u1',
        createdAt: '2025-01-01T00:00:00Z',
        expiresAt: '2099-01-01T00:00:00Z',
        revokedAt: '2025-01-02T00:00:00Z',
        lastAccessedAt: '',
      },
    ];
    render(<ShareLinksPage />);
    expect(screen.getByText(/revoked/i)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /revoke/i })).not.toBeInTheDocument();
  });

  it('clicking revoke calls the mutation with the link id', () => {
    links = [
      {
        id: 's1',
        docTreeRoot: 'c1',
        createdBy: 'u1',
        createdAt: '2025-01-01T00:00:00Z',
        expiresAt: '2099-01-01T00:00:00Z',
        revokedAt: '',
        lastAccessedAt: '',
      },
    ];
    render(<ShareLinksPage />);
    fireEvent.click(screen.getByRole('button', { name: /revoke/i }));
    expect(revokeMock).toHaveBeenCalledWith('s1', expect.anything());
  });
});
