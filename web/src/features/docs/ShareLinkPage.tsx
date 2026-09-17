/**
 * Public, unauthenticated read-only viewer for a doc share link (#240 PR2).
 * Mounted at /share/:token, outside RequireAuth. Fetches the scoped tree and
 * (optionally) one doc from it; expired/revoked links get a dedicated
 * message, not a generic 404 — `share_link_expired` / `share_link_revoked`
 * are distinct backend error codes for exactly this.
 */
import { useState, type ReactNode } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  useShareLinkTree,
  useShareLinkDoc,
  shareLinkErrorCode,
  type ShareTreeNode,
} from '../../api/shareLinks';
import { Panel } from '../../components/ui/Panel';
import { Markdown } from '../../components/docs/Markdown';
import { Skeleton, SkeletonRows, EmptyState } from '../../components/ui/states';
import { FileTextIcon, LayersIcon, ClockIcon } from '../../components/icons';
import { cn } from '../../lib/cn';

function findFirstDoc(node: ShareTreeNode): string | undefined {
  if (!node.children || node.children.length === 0) return node.kind !== 'lab' && node.kind !== 'service' ? node.docId : undefined;
  for (const child of node.children) {
    if (child.kind !== 'lab' && child.kind !== 'service') return child.docId;
    const nested = findFirstDoc(child);
    if (nested) return nested;
  }
  return undefined;
}

export function ShareLinkPage() {
  const { t } = useTranslation();
  const { token, docId: routeDocId } = useParams<{ token: string; docId?: string }>();
  const tree = useShareLinkTree(token);
  const [selectedDocId, setSelectedDocId] = useState<string | undefined>(routeDocId);

  const activeDocId = selectedDocId ?? (tree.data ? findFirstDoc(tree.data) : undefined);
  const doc = useShareLinkDoc(token, activeDocId);

  if (tree.isLoading) {
    return (
      <ShareShell>
        <SkeletonRows rows={8} />
      </ShareShell>
    );
  }

  if (tree.isError) {
    return <ShareErrorState error={tree.error} />;
  }

  if (!tree.data) return null;

  return (
    <ShareShell>
      <div className="flex gap-6">
        <aside className="hidden w-56 shrink-0 md:block">
          <ShareTree node={tree.data} activeDocId={activeDocId} onSelect={setSelectedDocId} />
        </aside>
        <section className="min-w-0 flex-1">
          {doc.isLoading ? (
            <Panel className="p-8">
              <Skeleton className="mb-4 h-7 w-1/3" />
              <SkeletonRows rows={8} />
            </Panel>
          ) : doc.isError || !doc.data ? (
            <Panel className="min-h-[40vh]">
              <EmptyState
                icon={<FileTextIcon size={20} />}
                title={t('docs.share.viewer.selectDoc', { defaultValue: 'Select a document' })}
              />
            </Panel>
          ) : (
            <Panel className="p-6">
              <h1 className="mb-4 font-mono text-lg font-semibold text-ink">{doc.data.title}</h1>
              <Markdown source={doc.data.content} />
            </Panel>
          )}
        </section>
      </div>
    </ShareShell>
  );
}

function ShareTree({
  node,
  activeDocId,
  onSelect,
}: {
  node: ShareTreeNode;
  activeDocId: string | undefined;
  onSelect: (docId: string) => void;
}) {
  return (
    <Panel className="p-2">
      <nav className="flex flex-col gap-0.5">
        {(node.children ?? []).map((connector) => (
          <div key={connector.docId} className="mb-1">
            <div className="flex items-center gap-2 px-2.5 py-1.5 text-xs font-medium text-ink-muted">
              <LayersIcon size={14} />
              {connector.title}
            </div>
            <div className="ml-3 flex flex-col gap-0.5 border-l border-line-soft pl-2">
              {(connector.children ?? []).map((docNode) => (
                <button
                  key={docNode.docId}
                  type="button"
                  onClick={() => onSelect(docNode.docId)}
                  className={cn(
                    'flex items-center gap-2 rounded-md px-2.5 py-1.5 text-left text-sm transition-colors',
                    activeDocId === docNode.docId
                      ? 'bg-accent-primary-tint text-ink'
                      : 'text-ink-muted hover:bg-surface-raised hover:text-ink'
                  )}
                >
                  <FileTextIcon size={14} className="shrink-0 opacity-70" />
                  <span className="truncate">{docNode.title}</span>
                </button>
              ))}
            </div>
          </div>
        ))}
      </nav>
    </Panel>
  );
}

function ShareShell({ children }: { children: ReactNode }) {
  const { t } = useTranslation();
  return (
    <div className="min-h-dvh bg-canvas">
      <header className="border-b border-line-soft px-6 py-4">
        <p className="font-mono text-xs uppercase tracking-[0.16em] text-ink-faint">
          {t('docs.share.viewer.badge', { defaultValue: 'Shared documentation — read only' })}
        </p>
      </header>
      <div className="mx-auto max-w-330 px-6 py-6">{children}</div>
    </div>
  );
}

function ShareErrorState({ error }: { error: unknown }) {
  const { t } = useTranslation();
  const code = shareLinkErrorCode(error);

  const copy: Record<string, { title: string; description: string }> = {
    share_link_expired: {
      title: t('docs.share.viewer.expiredTitle', { defaultValue: 'This link has expired' }),
      description: t('docs.share.viewer.expiredDesc', {
        defaultValue: 'Ask whoever shared it with you to create a new link.',
      }),
    },
    share_link_revoked: {
      title: t('docs.share.viewer.revokedTitle', { defaultValue: 'This link has been revoked' }),
      description: t('docs.share.viewer.revokedDesc', {
        defaultValue: 'The person who shared it has turned off access.',
      }),
    },
    not_found: {
      title: t('docs.share.viewer.notFoundTitle', { defaultValue: 'Link not found' }),
      description: t('docs.share.viewer.notFoundDesc', {
        defaultValue: "This share link doesn't exist or was never valid.",
      }),
    },
    unknown: {
      title: t('docs.share.viewer.errorTitle', { defaultValue: 'Something went wrong' }),
      description: t('docs.share.viewer.errorDesc', { defaultValue: 'Could not load this share link.' }),
    },
  };
  const { title, description } = copy[code];

  return (
    <div className="flex min-h-dvh items-center justify-center bg-canvas p-6">
      <EmptyState icon={<ClockIcon size={22} />} title={title} description={description} />
    </div>
  );
}
