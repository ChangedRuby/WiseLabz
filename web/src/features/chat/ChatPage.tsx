/**
 * "Ask your lab" chat — a dedicated page for asking questions answered by
 * retrieving relevant doc sections (backend piece 1/3 of issue #238). A new
 * conversation is scoped to either one doc or the whole lab at creation time;
 * the scope can't change mid-conversation, so the selector only shows while
 * composing the first message.
 */
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  useGetChatConversations,
  useGetChatConversationsId,
  postChatConversations,
  postChatConversationsIdMessages,
  getGetChatConversationsQueryKey,
  getGetChatConversationsIdQueryKey,
} from '../../api/generated/chat/chat';
import { useGetDocs } from '../../api/generated/docs/docs';
import { Panel, PanelHeader } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { Markdown } from '../../components/docs/Markdown';
import { cn } from '../../lib/cn';
import { relativeTime } from '../../lib/time';
import { toast } from '../../lib/toast';
import { ChatIcon, FileTextIcon, LayersIcon, PlusIcon, ArrowRightIcon } from '../../components/icons';
import type { ChatConversationCreateScopeType } from '../../api/model';

export function ChatPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const conversations = useGetChatConversations();
  const docs = useGetDocs({ pageSize: 100 });
  const docTitle = (docId: string) =>
    docs.data?.items.find((d) => d.docId === docId)?.title ?? docId;

  const [activeId, setActiveId] = useState<string | null>(null);
  const [scopeType, setScopeType] = useState<ChatConversationCreateScopeType>('lab');
  const [scopeDocId, setScopeDocId] = useState('');
  const [question, setQuestion] = useState('');

  const active = useGetChatConversationsId(activeId ?? '', { query: { enabled: !!activeId } });

  const ask = useMutation({
    mutationFn: async (content: string) => {
      let id = activeId;
      if (!id) {
        const conversation = await postChatConversations({
          scopeType,
          scopeId: scopeType === 'doc' ? scopeDocId : undefined,
        });
        id = conversation.id;
      }
      await postChatConversationsIdMessages(id, { content });
      return id;
    },
    onSuccess: (id) => {
      setActiveId(id);
      setQuestion('');
      queryClient.invalidateQueries({ queryKey: getGetChatConversationsQueryKey() });
      queryClient.invalidateQueries({ queryKey: getGetChatConversationsIdQueryKey(id) });
    },
    onError: () => toast.error(t('chat.askError')),
  });

  const startNewChat = () => {
    setActiveId(null);
    setScopeType('lab');
    setScopeDocId('');
    setQuestion('');
  };

  const canAsk = question.trim() !== '' && (scopeType === 'lab' || scopeDocId !== '' || !!activeId);

  return (
    <div className="mx-auto flex h-[calc(100vh-4rem)] max-w-330 gap-6 px-6 py-6">
      {/* Conversation history */}
      <aside className="hidden w-64 shrink-0 lg:flex lg:flex-col">
        <Panel className="flex-1 overflow-hidden">
          <PanelHeader
            title={t('chat.history')}
            icon={<ChatIcon size={16} />}
            action={
              <Button variant="ghost" size="sm" onClick={startNewChat}>
                <PlusIcon size={13} />
                {t('chat.new')}
              </Button>
            }
          />
          <div className="flex-1 overflow-y-auto">
            {conversations.isLoading ? (
              <SkeletonRows rows={5} />
            ) : conversations.isError || !conversations.data ? (
              <ErrorState
                description={t('chat.historyLoadError')}
                onRetry={() => conversations.refetch()}
              />
            ) : conversations.data.length === 0 ? (
              <EmptyState title={t('chat.historyEmptyTitle')} description={t('chat.historyEmptyDesc')} />
            ) : (
              <ul>
                {conversations.data.map((c) => (
                  <li key={c.id} className="border-b border-line-soft last:border-0">
                    <button
                      onClick={() => setActiveId(c.id)}
                      aria-current={activeId === c.id}
                      className={cn(
                        'flex w-full flex-col gap-0.5 px-3 py-2.5 text-left transition-colors hover:bg-surface-raised',
                        activeId === c.id && 'bg-surface-raised'
                      )}
                    >
                      <span className="flex items-center gap-1.5 truncate text-xs font-medium text-ink">
                        {c.scopeType === 'doc' ? (
                          <FileTextIcon size={12} className="shrink-0 text-ink-faint" />
                        ) : (
                          <LayersIcon size={12} className="shrink-0 text-ink-faint" />
                        )}
                        <span className="truncate">
                          {c.scopeType === 'doc' ? docTitle(c.scopeId) : t('chat.scopeLab')}
                        </span>
                      </span>
                      <span className="font-mono text-2xs text-ink-faint">
                        {t('common.ago', { time: relativeTime(c.createdAt) })}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </Panel>
      </aside>

      {/* Conversation */}
      <section className="flex min-w-0 flex-1 flex-col">
        <Panel className="flex flex-1 flex-col overflow-hidden">
          <PanelHeader
            title={
              activeId && active.data
                ? active.data.conversation.scopeType === 'doc'
                  ? docTitle(active.data.conversation.scopeId)
                  : t('chat.scopeLab')
                : t('chat.title')
            }
            icon={<ChatIcon size={16} />}
          />

          {!activeId ? (
            <div className="border-b border-line-soft px-4 py-3">
              <p className="mb-2 text-xs text-ink-muted">{t('chat.scopeLabel')}</p>
              <div className="flex flex-wrap items-center gap-2">
                <div
                  role="tablist"
                  className="flex items-center gap-1 rounded-lg border border-line-soft bg-canvas-sunken p-0.5"
                >
                  <button
                    role="tab"
                    aria-selected={scopeType === 'lab'}
                    onClick={() => setScopeType('lab')}
                    className={cn(
                      'rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors',
                      scopeType === 'lab' ? 'bg-surface-raised text-ink' : 'text-ink-muted hover:text-ink'
                    )}
                  >
                    {t('chat.scopeLab')}
                  </button>
                  <button
                    role="tab"
                    aria-selected={scopeType === 'doc'}
                    onClick={() => setScopeType('doc')}
                    className={cn(
                      'rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors',
                      scopeType === 'doc' ? 'bg-surface-raised text-ink' : 'text-ink-muted hover:text-ink'
                    )}
                  >
                    {t('chat.scopeDoc')}
                  </button>
                </div>
                {scopeType === 'doc' && (
                  <select
                    value={scopeDocId}
                    onChange={(e) => setScopeDocId(e.target.value)}
                    className="h-9 rounded-sm border border-line bg-surface px-2.5 text-sm text-ink outline-none focus-visible:border-accent-primary-soft"
                  >
                    <option value="">{t('chat.scopeDocPlaceholder')}</option>
                    {docs.data?.items.map((d) => (
                      <option key={d.docId} value={d.docId}>
                        {d.title}
                      </option>
                    ))}
                  </select>
                )}
              </div>
            </div>
          ) : null}

          <div className="flex-1 overflow-y-auto px-4 py-4">
            {activeId && active.isLoading ? (
              <SkeletonRows rows={4} />
            ) : activeId && (active.isError || !active.data) ? (
              <ErrorState description={t('chat.conversationLoadError')} onRetry={() => active.refetch()} />
            ) : activeId && active.data ? (
              <MessageList messages={active.data.messages} pending={ask.isPending} />
            ) : (
              <EmptyState
                icon={<ChatIcon size={20} />}
                title={t('chat.emptyTitle')}
                description={t('chat.emptyDesc')}
              />
            )}
          </div>

          <form
            onSubmit={(e) => {
              e.preventDefault();
              if (canAsk) ask.mutate(question);
            }}
            className="flex items-center gap-2 border-t border-line-soft px-4 py-3"
          >
            <input
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder={t('chat.askPlaceholder')}
              disabled={ask.isPending}
              className="h-9 flex-1 rounded-sm border border-line bg-surface px-2.5 text-sm text-ink outline-none placeholder:text-ink-faint focus-visible:border-accent-primary-soft"
            />
            <Button type="submit" variant="primary" size="sm" disabled={!canAsk || ask.isPending}>
              {t('chat.ask')}
              <ArrowRightIcon size={13} />
            </Button>
          </form>
        </Panel>
      </section>
    </div>
  );
}

function MessageList({
  messages,
  pending,
}: {
  messages: { id: string; role: string; content: string; provider?: string }[];
  pending: boolean;
}) {
  const { t } = useTranslation();
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView?.({ block: 'end' });
  }, [messages.length, pending]);

  return (
    <div className="flex flex-col gap-3">
      {messages.map((m) => (
        <div
          key={m.id}
          className={cn('max-w-[85%] rounded-lg px-3 py-2 text-sm', {
            'ml-auto bg-accent-primary-tint text-ink': m.role === 'user',
            'bg-surface-raised text-ink': m.role === 'assistant',
          })}
        >
          {m.role === 'assistant' ? (
            <Markdown source={m.content} />
          ) : (
            <p className="whitespace-pre-wrap">{m.content}</p>
          )}
          {m.provider && (
            <p className="mt-1 font-mono text-2xs text-ink-faint">{t('chat.answeredBy', { provider: m.provider })}</p>
          )}
        </div>
      ))}
      {pending && (
        <div className="max-w-[85%] rounded-lg bg-surface-raised px-3 py-2 text-sm text-ink-faint">
          {t('chat.thinking')}
        </div>
      )}
      <div ref={bottomRef} />
    </div>
  );
}
