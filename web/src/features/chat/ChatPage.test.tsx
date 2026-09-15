import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import '../../i18n';
import { ChatPage } from './ChatPage';

const postChatConversations = vi.fn().mockResolvedValue({
  id: 'conv-1',
  userId: 'u1',
  scopeType: 'lab',
  scopeId: '',
  createdAt: '2026-01-01T00:00:00Z',
});
const postChatConversationsIdMessages = vi.fn().mockResolvedValue({
  id: 'm2',
  conversationId: 'conv-1',
  role: 'assistant',
  content: 'Your Proxmox host is pve1.',
  provider: 'openai',
  createdAt: '2026-01-01T00:00:01Z',
});

vi.mock('../../api/generated/chat/chat', () => ({
  useGetChatConversations: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
  useGetChatConversationsId: (id: string) =>
    id === 'conv-1'
      ? {
          data: {
            conversation: { id: 'conv-1', userId: 'u1', scopeType: 'lab', scopeId: '', createdAt: '2026-01-01T00:00:00Z' },
            messages: [
              { id: 'm1', conversationId: 'conv-1', role: 'user', content: 'Which host runs Proxmox?', createdAt: '2026-01-01T00:00:00Z' },
              { id: 'm2', conversationId: 'conv-1', role: 'assistant', content: 'Your Proxmox host is pve1.', provider: 'openai', createdAt: '2026-01-01T00:00:01Z' },
            ],
          },
          isLoading: false,
          isError: false,
          refetch: vi.fn(),
        }
      : { data: undefined, isLoading: false, isError: false, refetch: vi.fn() },
  postChatConversations: (...args: unknown[]) => postChatConversations(...args),
  postChatConversationsIdMessages: (...args: unknown[]) => postChatConversationsIdMessages(...args),
  getGetChatConversationsQueryKey: () => ['/chat/conversations'],
  getGetChatConversationsIdQueryKey: (id: string) => [`/chat/conversations/${id}`],
}));

vi.mock('../../api/generated/docs/docs', () => ({
  useGetDocs: () => ({
    data: { items: [{ docId: 'doc-pve1', title: 'pve1', kind: 'service', updatedAt: '2026-01-01T00:00:00Z' }], total: 1, page: 1, pageSize: 100 },
    isLoading: false,
    isError: false,
  }),
}));

afterEach(cleanup);

function renderChat() {
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter initialEntries={['/chat']}>
        <ChatPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('ChatPage', () => {
  it('asks a lab-scoped question and renders the answer with its history', async () => {
    renderChat();

    const input = screen.getByPlaceholderText('Ask a question…');
    fireEvent.change(input, { target: { value: 'Which host runs Proxmox?' } });
    fireEvent.click(screen.getByRole('button', { name: /Ask/ }));

    expect(await screen.findByText('Your Proxmox host is pve1.')).toBeInTheDocument();
    expect(screen.getByText('Which host runs Proxmox?')).toBeInTheDocument();
    expect(postChatConversations).toHaveBeenCalledWith({ scopeType: 'lab', scopeId: undefined });
    expect(postChatConversationsIdMessages).toHaveBeenCalledWith('conv-1', {
      content: 'Which host runs Proxmox?',
    });
  });

  it('requires a doc to be chosen before asking in doc scope', () => {
    renderChat();

    fireEvent.click(screen.getByRole('tab', { name: 'This doc' }));
    fireEvent.change(screen.getByPlaceholderText('Ask a question…'), {
      target: { value: 'What ports does it expose?' },
    });

    expect(screen.getByRole('button', { name: /Ask/ })).toBeDisabled();
  });
});
