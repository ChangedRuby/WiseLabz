/**
 * Freehand + AI-assisted doc editor (`/docs/:docId/edit`). CodeMirror 6 on the
 * left, live Markdown preview on the right; save writes a new version
 * (last-write-wins). If the doc is regenerated elsewhere while editing, a
 * non-destructive "newer version available" banner appears — saving still creates
 * a new version over the top.
 *
 * AI assist is batched (no streaming, per the plan): "Suggest update" calls the
 * mock once, then renders the proposed revision as a review-diff. Accept applies it
 * to the editor and marks the draft as AI-drafted (provenance); Reject discards it.
 * Operator-gated (the route guards, and the save button respects role too).
 */
import { useEffect, useRef, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { isAxiosError } from 'axios';
import CodeMirror from '@uiw/react-codemirror';
import { markdown } from '@codemirror/lang-markdown';
import { EditorView } from '@codemirror/view';
import {
  useGetDocsDocId,
  putDocsDocId,
  postDocsDocIdAiSuggest,
  postDocsDocIdLock,
  postDocsDocIdLockRelease,
  getGetDocsDocIdQueryKey,
  getGetDocsDocIdVersionsQueryKey,
  getGetDocsTreeQueryKey,
} from '../../api/generated/docs/docs';
import { useConnectorRole, useIsInstanceAdmin } from '../../hooks/useRole';
import { useAuth } from '../../store/auth';
import { useLive } from '../../store/live';
import { Button } from '../../components/ui/Button';
import { Panel } from '../../components/ui/Panel';
import { Skeleton, SkeletonRows, ErrorState } from '../../components/ui/states';
import { Markdown } from '../../components/docs/Markdown';
import { DocDiff } from '../../components/diff/DiffViewer';
import { toast } from '../../lib/toast';
import {
  ArrowRightIcon,
  SparklesIcon,
  CheckIcon,
  XIcon,
  FileTextIcon,
} from '../../components/icons';

// Compact CodeMirror theme that sits on the app's dark chrome instead of the
// library's default white surface.
const cmTheme = EditorView.theme(
  {
    '&': { backgroundColor: 'transparent', color: 'var(--color-ink)', fontSize: '13px' },
    '.cm-content': { fontFamily: 'var(--font-mono, monospace)', caretColor: 'var(--color-accent-primary)' },
    '.cm-gutters': {
      backgroundColor: 'transparent',
      color: 'var(--color-ink-faint)',
      border: 'none',
    },
    '.cm-activeLine': { backgroundColor: 'var(--color-surface-raised)' },
    '.cm-activeLineGutter': { backgroundColor: 'var(--color-surface-raised)' },
    '&.cm-focused': { outline: 'none' },
    '.cm-selectionBackground, ::selection': { backgroundColor: 'var(--color-accent-primary-tint)' },
  },
  { dark: true }
);

type Provenance = 'manual' | 'ai-draft';

/** Deterministic stand-in for an LLM doc revision (mock-only). */
function synthesizeSuggestion(draft: string): string {
  const note =
    '\n\n> _AI draft — reviewed the latest synced state and tightened the wording above._\n';
  if (!/^##?\s+Summary/im.test(draft)) {
    const withSummary = draft.replace(
      /^(#.*\n)/,
      '$1\n## Summary\n\nAuto-generated overview of this service based on its most recent sync.\n'
    );
    return withSummary + note;
  }
  return draft + note;
}

export function DocEditorPage() {
  const { t } = useTranslation();
  const { docId = '' } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isInstanceAdmin = useIsInstanceAdmin();

  const userId = useAuth((s) => s.user?.id);
  const docLock = useLive((s) => s.docLocks[docId]);

  const doc = useGetDocsDocId(docId);
  // Service docs are gated on that connector's operator role; the lab
  // overview doc has no owning connector, so it falls back to instance-admin.
  const connectorRole = useConnectorRole(doc.data?.serviceId ?? undefined);
  const canMutate = doc.data?.serviceId ? connectorRole === 'operator' : isInstanceAdmin;

  const [draft, setDraft] = useState<string | null>(null);
  const [provenance, setProvenance] = useState<Provenance>('manual');
  const [suggestion, setSuggestion] = useState<string | null>(null);
  // State (not refs) so `newerAvailable` can derive from them during render.
  const [baseVersion, setBaseVersion] = useState<number | null>(null);
  const [justSaved, setJustSaved] = useState(false);
  const lockAcquiredRef = useRef(false);

  // Seed the draft once the doc loads; capture the version the edit is based on.
  // Adjusting state during render is the React-blessed alternative to an effect.
  const [seeded, setSeeded] = useState(false);
  if (doc.data && !seeded) {
    setSeeded(true);
    setDraft(doc.data.content);
    setBaseVersion(doc.data.currentVersion);
  }

  const dirty = draft !== null && doc.data != null && draft !== doc.data.content;

  // Warn on tab close with unsaved edits; release lock on unload.
  useEffect(() => {
    const onBeforeUnload = (e: BeforeUnloadEvent) => {
      if (dirty) {
        e.preventDefault();
      }
      if (lockAcquiredRef.current) {
        postDocsDocIdLockRelease(docId).catch(() => {});
      }
    };
    window.addEventListener('beforeunload', onBeforeUnload);
    return () => {
      window.removeEventListener('beforeunload', onBeforeUnload);
      if (lockAcquiredRef.current) {
        postDocsDocIdLockRelease(docId).catch(() => {});
      }
    };
  }, [dirty, docId]);

  const save = useMutation({
    mutationFn: () =>
      putDocsDocId(docId, {
        content: draft ?? '',
        baseVersion: baseVersion ?? undefined,
        trigger: provenance === 'ai-draft' ? 'ai' : 'manual',
      }),
    onSuccess: (updated) => {
      setJustSaved(true);
      setBaseVersion(updated.currentVersion);
      if (lockAcquiredRef.current) {
        postDocsDocIdLockRelease(docId).catch(() => {});
        lockAcquiredRef.current = false;
      }
      queryClient.invalidateQueries({ queryKey: getGetDocsDocIdQueryKey(docId) });
      queryClient.invalidateQueries({ queryKey: getGetDocsDocIdVersionsQueryKey(docId) });
      queryClient.invalidateQueries({ queryKey: getGetDocsTreeQueryKey() });
      toast.success(t('docs.editor.saved', { version: updated.currentVersion }));
      navigate(`/docs/${docId}`);
    },
    onError: (error) => {
      // Stale baseVersion: server rejected with 409, current doc is untouched.
      // Refetch so `newerAvailable` flips true and the existing "load latest"
      // banner (below) becomes the recovery path — the draft stays intact.
      if (isAxiosError(error) && error.response?.status === 409) {
        queryClient.invalidateQueries({ queryKey: getGetDocsDocIdQueryKey(docId) });
        toast.error(t('docs.editor.saveConflict'));
        return;
      }
      toast.error(t('docs.editor.saveError'));
    },
  });

  const suggest = useMutation({
    mutationFn: () => postDocsDocIdAiSuggest(docId, { prompt: 'improve' }),
    onSuccess: () => setSuggestion(synthesizeSuggestion(draft ?? '')),
    onError: () => toast.error(t('docs.editor.aiError')),
  });

  // A newer version landed (e.g. a regen) while editing and it isn't our own save.
  const newerAvailable =
    doc.data != null && baseVersion != null && doc.data.currentVersion > baseVersion && !justSaved;

  // Acquire and maintain lock while editing.
  useEffect(() => {
    if (!dirty) {
      lockAcquiredRef.current = false;
      return;
    }

    let heartbeat: ReturnType<typeof setInterval> | undefined;

    // Acquire lock on first edit
    postDocsDocIdLock(docId).catch(() => {});
    lockAcquiredRef.current = true;

    // Renew lock every 30 seconds
    heartbeat = setInterval(() => {
      postDocsDocIdLock(docId).catch(() => {});
    }, 30000);

    return () => {
      if (heartbeat) clearInterval(heartbeat);
    };
  }, [dirty, docId]);

  if (doc.isLoading) {
    return (
      <div className="mx-auto max-w-330 px-6 py-6">
        <Panel className="p-6">
          <Skeleton className="mb-4 h-7 w-1/3" />
          <SkeletonRows rows={8} />
        </Panel>
      </div>
    );
  }
  if (doc.isError || !doc.data) {
    return (
      <div className="mx-auto max-w-330 px-6 py-6">
        <Panel className="min-h-[50vh]">
          <ErrorState description={t('docs.docLoadError')} onRetry={() => doc.refetch()} />
        </Panel>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-330 px-6 py-6">
      <button
        onClick={() => navigate(`/docs/${docId}`)}
        className="mb-4 inline-flex items-center gap-1.5 text-xs text-ink-muted transition-colors hover:text-ink"
      >
        <ArrowRightIcon size={13} className="rotate-180" />
        {t('docs.editor.back')}
      </button>

      {/* Toolbar */}
      <header className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2.5">
          <FileTextIcon size={18} className="text-accent-primary" />
          <h1 className="text-lg font-semibold tracking-tight text-ink">{doc.data.title}</h1>
          {provenance === 'ai-draft' && (
            <span className="inline-flex items-center gap-1 rounded bg-accent-primary-tint px-1.5 py-0.5 text-2xs font-medium text-accent-primary">
              <SparklesIcon size={11} /> {t('docs.editor.aiDrafted')}
            </span>
          )}
          {dirty && (
            <span className="inline-flex items-center gap-1 text-2xs text-ink-faint">
              <span className="h-1.5 w-1.5 rounded-sm bg-warn" aria-hidden />{' '}
              {t('docs.editor.unsaved')}
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            size="sm"
            onClick={() => suggest.mutate()}
            disabled={suggest.isPending || !canMutate}
          >
            <SparklesIcon size={14} />
            {suggest.isPending ? t('docs.editor.aiThinking') : t('docs.editor.aiSuggest')}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              setDraft(doc.data!.content);
              setProvenance('manual');
            }}
            disabled={!dirty || save.isPending}
          >
            {t('docs.editor.discard')}
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => save.mutate()}
            disabled={!dirty || save.isPending || !canMutate}
          >
            <CheckIcon size={14} />
            {save.isPending ? t('docs.editor.saving') : t('docs.editor.save')}
          </Button>
        </div>
      </header>

      {newerAvailable && (
        <div className="mb-4 flex items-center justify-between gap-3 rounded-md border border-warn bg-warn-tint px-3 py-2 text-xs text-warn">
          <span>{t('docs.editor.newerBanner', { version: doc.data.currentVersion })}</span>
          <button
            onClick={() => {
              setDraft(doc.data!.content);
              setBaseVersion(doc.data!.currentVersion);
              setProvenance('manual');
            }}
            className="rounded-sm font-medium underline-offset-2 hover:underline"
          >
            {t('docs.editor.loadLatest')}
          </button>
        </div>
      )}

      {docLock && docLock.userId !== userId && (
        <div className="mb-4 rounded-md border border-info bg-info-tint px-3 py-2 text-xs text-info">
          <span>{t('docs.editor.lockBanner', { userId: docLock.userId })}</span>
        </div>
      )}

      {/* AI review-diff panel */}
      {suggestion !== null && (
        <Panel className="mb-4 p-4">
          <div className="mb-3 flex items-center justify-between">
            <span className="flex items-center gap-2 text-sm font-semibold text-ink">
              <SparklesIcon size={15} className="text-accent-primary" />
              {t('docs.editor.aiReviewTitle')}
            </span>
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" onClick={() => setSuggestion(null)}>
                <XIcon size={14} /> {t('docs.editor.aiReject')}
              </Button>
              <Button
                variant="primary"
                size="sm"
                onClick={() => {
                  setDraft(suggestion);
                  setProvenance('ai-draft');
                  setSuggestion(null);
                  toast.success(t('docs.editor.aiAccepted'));
                }}
              >
                <CheckIcon size={14} /> {t('docs.editor.aiAccept')}
              </Button>
            </div>
          </div>
          <DocDiff
            before={draft ?? ''}
            after={suggestion}
            baseLabel={t('docs.editor.aiDiffCurrent')}
            headLabel={t('docs.editor.aiDiffProposed')}
            headTrigger="ai"
            label={t('docs.editor.aiDiffLabel')}
          />
        </Panel>
      )}

      {/* Editor + preview */}
      <div className="grid gap-4 lg:grid-cols-2">
        <Panel className="overflow-hidden">
          <div className="border-b border-line-soft px-3 py-2 text-2xs text-ink-faint">
            {t('docs.editor.markdownLabel')}
          </div>
          <CodeMirror
            value={draft ?? ''}
            onChange={(v) => setDraft(v)}
            extensions={[markdown(), cmTheme, EditorView.lineWrapping]}
            basicSetup={{ lineNumbers: true, foldGutter: false, highlightActiveLine: true }}
            editable={canMutate}
            className="min-h-[60vh] text-sm"
          />
        </Panel>
        <Panel className="overflow-hidden">
          <div className="border-b border-line-soft px-3 py-2 text-2xs text-ink-faint">
            {t('docs.editor.previewLabel')}
          </div>
          <div className="max-h-[70vh] overflow-y-auto px-5 py-4">
            <Markdown source={draft ?? ''} />
          </div>
        </Panel>
      </div>
    </div>
  );
}
