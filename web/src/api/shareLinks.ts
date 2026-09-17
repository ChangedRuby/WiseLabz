/**
 * Hand-written client for the #240 PR2 read-only doc share-link endpoints.
 * TODO: fold into docs/openapi.yaml and regenerate via `npm run gen:api`
 * once the backend has landed and the spec is updated — see
 * src/api/permissions.ts for why these are hand-typed for now.
 *
 * The public `/share/{token}/...` endpoints are unauthenticated: they're
 * called with `customInstance` too, but since there's no access token set
 * for an anonymous visitor, no Authorization header goes out (see
 * axios-instance.ts) — no separate client needed.
 */
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { isAxiosError } from 'axios';
import { customInstance } from './axios-instance';

export interface ShareLink {
  id: string;
  docTreeRoot: string;
  createdBy: string;
  createdAt: string;
  expiresAt: string;
  revokedAt: string;
  lastAccessedAt: string;
}

export interface ShareLinkCreated extends ShareLink {
  /** Raw token, returned once at creation time — never shown again. */
  token: string;
}

export interface ShareTreeNode {
  docId: string;
  title: string;
  kind: string;
  children?: ShareTreeNode[];
}

export interface ShareDoc {
  docId: string;
  title: string;
  kind: string;
  serviceId?: string | null;
  content: string;
  currentVersion: number;
  createdAt: string;
  updatedAt: string;
}

const shareLinksQueryKey = ['/docs/share-links'] as const;

export function useListShareLinks() {
  return useQuery({
    queryKey: shareLinksQueryKey,
    queryFn: ({ signal }) => customInstance<ShareLink[]>({ url: '/docs/share-links', method: 'GET', signal }),
  });
}

export function useCreateShareLink() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: { docTreeRoot: string; expiresAt: string }) =>
      customInstance<ShareLinkCreated>({ url: '/docs/share-links', method: 'POST', data: body }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: shareLinksQueryKey }),
  });
}

export function useRevokeShareLink() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => customInstance<void>({ url: `/docs/share-links/${id}`, method: 'DELETE' }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: shareLinksQueryKey }),
  });
}

/** Distinguishes the dedicated expired/revoked states from a generic not-found. */
export type ShareLinkErrorCode = 'share_link_expired' | 'share_link_revoked' | 'not_found' | 'unknown';

export function shareLinkErrorCode(error: unknown): ShareLinkErrorCode {
  if (isAxiosError(error)) {
    const code = (error.response?.data as { code?: string } | undefined)?.code;
    if (code === 'share_link_expired' || code === 'share_link_revoked' || code === 'not_found') {
      return code;
    }
  }
  return 'unknown';
}

export function useShareLinkTree(token: string | undefined) {
  return useQuery({
    queryKey: ['/share', token, 'tree'] as const,
    queryFn: ({ signal }) => customInstance<ShareTreeNode>({ url: `/share/${token}/tree`, method: 'GET', signal }),
    enabled: !!token,
    retry: false,
  });
}

export function useShareLinkDoc(token: string | undefined, docId: string | undefined) {
  return useQuery({
    queryKey: ['/share', token, 'docs', docId] as const,
    queryFn: ({ signal }) =>
      customInstance<ShareDoc>({ url: `/share/${token}/docs/${docId}`, method: 'GET', signal }),
    enabled: !!token && !!docId,
    retry: false,
  });
}
