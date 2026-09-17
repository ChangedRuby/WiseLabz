import { cleanup, render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { AxiosError } from 'axios';
import '../../i18n';
import { ShareLinkPage } from './ShareLinkPage';

let treeResult: { data?: unknown; isLoading: boolean; isError: boolean; error?: unknown } = {
  data: undefined,
  isLoading: true,
  isError: false,
};
let docResult: { data?: unknown; isLoading: boolean; isError: boolean } = {
  data: undefined,
  isLoading: false,
  isError: false,
};

vi.mock('../../api/shareLinks', async () => {
  const actual = await vi.importActual<typeof import('../../api/shareLinks')>('../../api/shareLinks');
  return {
    ...actual,
    useShareLinkTree: () => treeResult,
    useShareLinkDoc: () => docResult,
  };
});

function renderPage(path = '/share/tok123') {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="/share/:token" element={<ShareLinkPage />} />
      </Routes>
    </MemoryRouter>
  );
}

describe('ShareLinkPage (#240 PR2)', () => {
  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it('renders a loading skeleton while the tree loads', () => {
    treeResult = { data: undefined, isLoading: true, isError: false };
    const { container } = renderPage();
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('renders the tree and the first doc once loaded', () => {
    treeResult = {
      isLoading: false,
      isError: false,
      data: {
        docId: 'root',
        title: 'Shared Documentation',
        kind: 'lab',
        children: [
          {
            docId: 'c1',
            title: 'Connector One',
            kind: 'service',
            children: [{ docId: 'd1', title: 'Runbook', kind: 'doc' }],
          },
        ],
      },
    };
    docResult = {
      isLoading: false,
      isError: false,
      data: { docId: 'd1', title: 'Runbook', kind: 'doc', content: '# Hello', currentVersion: 1, createdAt: '', updatedAt: '' },
    };
    renderPage();
    expect(screen.getByText('Connector One')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Runbook' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Runbook' })).toBeInTheDocument();
    expect(screen.getByText('Hello')).toBeInTheDocument();
  });

  it('shows a distinct message for an expired link', () => {
    const err = new AxiosError('fail');
    err.response = { data: { code: 'share_link_expired' }, status: 410 } as never;
    treeResult = { data: undefined, isLoading: false, isError: true, error: err };
    renderPage();
    expect(screen.getByText(/this link has expired/i)).toBeInTheDocument();
  });

  it('shows a distinct message for a revoked link', () => {
    const err = new AxiosError('fail');
    err.response = { data: { code: 'share_link_revoked' }, status: 410 } as never;
    treeResult = { data: undefined, isLoading: false, isError: true, error: err };
    renderPage();
    expect(screen.getByText(/this link has been revoked/i)).toBeInTheDocument();
  });

  it('falls back to a generic error message for an unrecognized failure', () => {
    treeResult = { data: undefined, isLoading: false, isError: true, error: new Error('boom') };
    renderPage();
    expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
  });
});
