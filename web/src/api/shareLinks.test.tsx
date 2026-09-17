import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { AxiosError } from 'axios';
import type { ReactNode } from 'react';
import {
  useListShareLinks,
  useCreateShareLink,
  useRevokeShareLink,
  useShareLinkTree,
  useShareLinkDoc,
  shareLinkErrorCode,
} from './shareLinks';

const customInstanceMock = vi.fn();
vi.mock('./axios-instance', () => ({
  customInstance: (config: unknown) => customInstanceMock(config),
}));

function wrapper({ children }: { children: ReactNode }) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

describe('shareLinks API client (#240 PR2)', () => {
  it('useListShareLinks GETs /docs/share-links', async () => {
    customInstanceMock.mockResolvedValueOnce([{ id: 's1' }]);
    const { result } = renderHook(() => useListShareLinks(), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(customInstanceMock).toHaveBeenCalledWith(
      expect.objectContaining({ url: '/docs/share-links', method: 'GET' })
    );
  });

  it('useCreateShareLink POSTs docTreeRoot + expiresAt', async () => {
    customInstanceMock.mockResolvedValueOnce({ id: 's1', token: 'wlz_share_abc' });
    const { result } = renderHook(() => useCreateShareLink(), { wrapper });
    result.current.mutate({ docTreeRoot: 'root', expiresAt: '2099-01-01T00:00:00Z' });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(customInstanceMock).toHaveBeenCalledWith(
      expect.objectContaining({
        url: '/docs/share-links',
        method: 'POST',
        data: { docTreeRoot: 'root', expiresAt: '2099-01-01T00:00:00Z' },
      })
    );
    expect(result.current.data?.token).toBe('wlz_share_abc');
  });

  it('useRevokeShareLink DELETEs by id', async () => {
    customInstanceMock.mockResolvedValueOnce(undefined);
    const { result } = renderHook(() => useRevokeShareLink(), { wrapper });
    result.current.mutate('s1');
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(customInstanceMock).toHaveBeenCalledWith(
      expect.objectContaining({ url: '/docs/share-links/s1', method: 'DELETE' })
    );
  });

  it('useShareLinkTree GETs /share/{token}/tree and is disabled without a token', async () => {
    customInstanceMock.mockResolvedValueOnce({ docId: 'root', title: 'Shared Documentation', kind: 'lab' });
    const { result } = renderHook(() => useShareLinkTree('tok123'), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(customInstanceMock).toHaveBeenCalledWith(
      expect.objectContaining({ url: '/share/tok123/tree', method: 'GET' })
    );

    customInstanceMock.mockClear();
    const disabled = renderHook(() => useShareLinkTree(undefined), { wrapper });
    expect(disabled.result.current.fetchStatus).toBe('idle');
    expect(customInstanceMock).not.toHaveBeenCalled();
  });

  it('useShareLinkDoc GETs /share/{token}/docs/{docId}', async () => {
    customInstanceMock.mockResolvedValueOnce({ docId: 'd1', title: 'Doc', content: 'hi' });
    const { result } = renderHook(() => useShareLinkDoc('tok123', 'd1'), { wrapper });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(customInstanceMock).toHaveBeenCalledWith(
      expect.objectContaining({ url: '/share/tok123/docs/d1', method: 'GET' })
    );
  });

  describe('shareLinkErrorCode', () => {
    function axiosErrorWithCode(code: string) {
      const err = new AxiosError('fail');
      err.response = { data: { code }, status: 410 } as never;
      return err;
    }

    it('recognizes share_link_expired', () => {
      expect(shareLinkErrorCode(axiosErrorWithCode('share_link_expired'))).toBe('share_link_expired');
    });

    it('recognizes share_link_revoked', () => {
      expect(shareLinkErrorCode(axiosErrorWithCode('share_link_revoked'))).toBe('share_link_revoked');
    });

    it('recognizes not_found', () => {
      expect(shareLinkErrorCode(axiosErrorWithCode('not_found'))).toBe('not_found');
    });

    it('falls back to unknown for an unrecognized code or non-axios error', () => {
      expect(shareLinkErrorCode(axiosErrorWithCode('internal_error'))).toBe('unknown');
      expect(shareLinkErrorCode(new Error('boom'))).toBe('unknown');
      expect(shareLinkErrorCode(undefined)).toBe('unknown');
    });
  });
});
