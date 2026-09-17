/**
 * Documentation viewer + the doc-format DiffViewer. Left: hierarchical tree.
 * Center: the rendered doc, or its revision history with a version-to-version
 * diff. The diff is the payoff — it shows exactly what a sync (or an editor)
 * changed between revisions, the same view the changes feed links into.
 */
import { useEffect, useMemo, useRef, useState } from 'react';
import { Link, useParams, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { AnimatePresence, motion } from 'motion/react';
import { useGetDocsTree, useGetDocsDocId } from '../../api/generated/docs/docs';
import { Panel } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { IconButton } from '../../components/ui/Button';
import { RoleGate } from '../../components/ui/RoleGate';
import { Skeleton, SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { Markdown } from '../../components/docs/Markdown';
import { DocTree } from '../../components/docs/DocTree';
import { DocHistory } from './DocHistory';
import { ShareDialog } from './ShareDialog';
import { cn } from '../../lib/cn';
import { relativeTime, fullDate } from '../../lib/time';
import {
  FileTextIcon,
  HistoryIcon,
  SparklesIcon,
  EditIcon,
  SearchIcon,
  MenuIcon,
  XIcon,
} from '../../components/icons';

type Tab = 'read' | 'history';

export function DocsPage() {
  const { t } = useTranslation();
  const { docId } = useParams<{ docId: string }>();
  const tree = useGetDocsTree();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [shareNode, setShareNode] = useState<{ docId: string; title: string } | null>(null);
  const drawerTriggerRef = useRef<HTMLButtonElement>(null);
  const drawerRef = useRef<HTMLDivElement>(null);

  // Default to the lab root when no doc is selected in the URL.
  const activeId = docId ?? tree.data?.docId;

  useEffect(() => {
    if (!drawerOpen) return;
    const trigger = drawerTriggerRef.current;
    const drawer = drawerRef.current;
    if (!drawer) return;

    const focusable = () =>
      Array.from(
        drawer.querySelectorAll<HTMLElement>(
          'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
        )
      ).filter((element) => element.tabIndex >= 0);

    const focusFrame = requestAnimationFrame(() => (focusable()[0] ?? drawer).focus());
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setDrawerOpen(false);
        return;
      }
      if (event.key !== 'Tab') return;

      const elements = focusable();
      if (elements.length === 0) {
        event.preventDefault();
        drawer.focus();
        return;
      }

      const first = elements[0];
      const last = elements[elements.length - 1];
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    };
    document.addEventListener('keydown', onKeyDown);
    return () => {
      cancelAnimationFrame(focusFrame);
      document.removeEventListener('keydown', onKeyDown);
      if (trigger && document.contains(trigger)) trigger.focus();
    };
  }, [drawerOpen]);

  const treeContent = tree.isLoading ? (
    <SkeletonRows rows={6} />
  ) : tree.isError || !tree.data ? (
    <ErrorState description={t('docs.treeLoadError')} onRetry={() => tree.refetch()} />
  ) : (
    <DocTree tree={tree.data} onShare={setShareNode} />
  );

  return (
    <div className="mx-auto flex max-w-330 gap-6 px-6 py-6">
      {/* Tree — desktop */}
      <aside className="hidden w-60 shrink-0 lg:block">
        <Panel className="sticky top-6 p-2">
          <Link
            to="/docs/all"
            className="mb-1 flex items-center gap-2 rounded-md px-2.5 py-2 text-xs text-ink-muted transition-colors hover:bg-surface-raised hover:text-ink"
          >
            <SearchIcon size={14} />
            {t('docs.searchAllLink')}
          </Link>
          {treeContent}
        </Panel>
      </aside>

      {/* Tree — mobile drawer */}
      <button
        ref={drawerTriggerRef}
        onClick={() => setDrawerOpen(true)}
        aria-label={t('docs.browseDocs')}
        className="fixed bottom-30 right-4 z-(--z-sticky) flex h-11 w-11 items-center justify-center rounded-full border border-line bg-surface-raised text-ink shadow-(--shadow-pop) lg:hidden"
      >
        <MenuIcon size={18} />
      </button>
      <AnimatePresence>
        {drawerOpen && (
          <>
            <motion.div
              className="fixed inset-0 z-(--z-overlay) bg-canvas/70 lg:hidden"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setDrawerOpen(false)}
            />
            <motion.div
              ref={drawerRef}
              role="dialog"
              aria-modal="true"
              aria-label={t('docs.browseDocs')}
              tabIndex={-1}
              className="fixed inset-y-0 left-0 z-(--z-overlay) w-72 max-w-[85vw] overflow-y-auto bg-surface-overlay p-2 shadow-(--shadow-pop) lg:hidden"
              initial={{ x: '-100%' }}
              animate={{ x: 0 }}
              exit={{ x: '-100%' }}
              transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            >
              <div className="mb-1 flex items-center justify-between px-1 py-1">
                <Link
                  to="/docs/all"
                  onClick={() => setDrawerOpen(false)}
                  className="flex items-center gap-2 rounded-md px-2.5 py-2 text-xs text-ink-muted transition-colors hover:bg-surface-raised hover:text-ink"
                >
                  <SearchIcon size={14} />
                  {t('docs.searchAllLink')}
                </Link>
                <IconButton label={t('docs.closeDocTree')} onClick={() => setDrawerOpen(false)}>
                  <XIcon size={16} />
                </IconButton>
              </div>
              <div onClick={() => setDrawerOpen(false)}>{treeContent}</div>
            </motion.div>
          </>
        )}
      </AnimatePresence>

      {/* Content */}
      <section className="min-w-0 flex-1">
        {activeId ? (
          <DocReader key={activeId} docId={activeId} />
        ) : (
          <Panel className="min-h-[60vh]">
            <EmptyState
              icon={<FileTextIcon size={20} />}
              title={t('docs.selectTitle')}
              description={t('docs.selectDesc')}
            />
          </Panel>
        )}
      </section>
      <ShareDialog open={!!shareNode} onClose={() => setShareNode(null)} node={shareNode} />
    </div>
  );
}

function DocReader({ docId }: { docId: string }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { data, isLoading, isError, refetch } = useGetDocsDocId(docId);
  const tab: Tab = location.pathname.endsWith('/history') ? 'history' : 'read';

  const triggerLabel = useMemo(() => {
    if (!data) return null;
    return data.kind === 'lab' ? t('docs.labOverview') : t('docs.serviceDoc');
  }, [data, t]);

  if (isLoading)
    return (
      <Panel className="p-8">
        <Skeleton className="mb-4 h-7 w-1/3" />
        <SkeletonRows rows={8} />
      </Panel>
    );
  if (isError || !data)
    return (
      <Panel className="min-h-[50vh]">
        <ErrorState description={t('docs.docLoadError')} onRetry={() => refetch()} />
      </Panel>
    );

  return (
    <Panel>
      {/* Doc header */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-line-soft px-6 py-4">
        <div className="min-w-0">
          <div className="mb-1 flex items-center gap-2">
            <span className="rounded bg-accent-secondary-tint px-1.5 py-0.5 text-2xs font-medium text-accent-secondary-bright">
              {triggerLabel}
            </span>
            <span className="nums font-mono text-2xs text-ink-faint">
              {t('docs.versionMeta', {
                version: data.currentVersion,
                time: t('common.ago', { time: relativeTime(data.updatedAt) }),
              })}
            </span>
          </div>
          <h1 className="truncate text-lg font-semibold tracking-tight text-ink">{data.title}</h1>
        </div>

        <div className="flex items-center gap-2">
          <div
            role="tablist"
            className="flex items-center gap-1 rounded-lg border border-line-soft bg-canvas-sunken p-0.5"
          >
            <TabButton
              active={tab === 'read'}
              onClick={() => navigate(`/docs/${docId}`)}
              icon={<FileTextIcon size={14} />}
            >
              {t('docs.read')}
            </TabButton>
            <TabButton
              active={tab === 'history'}
              onClick={() => navigate(`/docs/${docId}/history`)}
              icon={<HistoryIcon size={14} />}
            >
              {t('docs.history')}
            </TabButton>
          </div>
          <RoleGate connectorId={data.serviceId ?? undefined}>
            <Button variant="secondary" size="sm" onClick={() => navigate(`/docs/${docId}/edit`)}>
              <EditIcon size={14} />
              {t('docs.editAction')}
            </Button>
          </RoleGate>
        </div>
      </div>

      <div className="px-6 py-5">
        <AnimatePresence mode="wait">
          {tab === 'read' ? (
            <motion.div
              key="read"
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            >
              {data.kind === 'lab' && (
                <div className="mb-4 flex items-center gap-2 rounded-md bg-accent-primary-tint px-3 py-2 text-xs text-accent-primary">
                  <SparklesIcon size={14} />
                  {t('docs.labBanner')}
                </div>
              )}
              <Markdown source={data.content} />
              <p className="mt-8 border-t border-line-soft pt-3 text-2xs text-ink-faint">
                {t('docs.reconciledFooter', { date: fullDate(data.updatedAt) })}
              </p>
            </motion.div>
          ) : (
            <motion.div
              key="history"
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
            >
              <DocHistory docId={docId} currentVersion={data.currentVersion} connectorId={data.serviceId ?? undefined} />
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </Panel>
  );
}

function TabButton({
  active,
  onClick,
  icon,
  children,
}: {
  active: boolean;
  onClick: () => void;
  icon: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <button
      role="tab"
      aria-selected={active}
      onClick={onClick}
      className={cn(
        'relative flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors',
        active ? 'text-ink' : 'text-ink-muted hover:text-ink'
      )}
    >
      {active && (
        <motion.span
          layoutId="doc-tab"
          className="absolute inset-0 -z-10 rounded-md bg-surface-raised"
          transition={{ type: 'spring', stiffness: 500, damping: 36 }}
        />
      )}
      {icon}
      {children}
    </button>
  );
}
