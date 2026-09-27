import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import '../../../i18n';
import { SnapshotsPage } from './SnapshotsPage';

const { get, download } = vi.hoisted(() => ({ get: vi.fn(), download: vi.fn() }));
vi.mock('../../../api/axios-instance', () => ({ AXIOS_INSTANCE: { get } }));
vi.mock('../../../lib/download', () => ({ downloadBlob: download, filenameFromContentDisposition: (_: string | undefined, fallback: string) => fallback }));
vi.mock('../../../hooks/useRole', () => ({ useConnectorRole: () => 'operator', useIsInstanceAdmin: () => false }));
vi.mock('../../../api/generated/connectors/connectors', () => ({
  useGetConnectorsConnectorId: () => ({ data: { name: 'server' } }),
  useGetConnectorsConnectorIdGoldenSnapshot: () => ({ data: { snapshotId: 'old' } }),
  getGetConnectorsConnectorIdGoldenSnapshotQueryKey: () => ['golden'],
  useGetConnectorsConnectorIdSnapshotsSnapshotId: () => ({ data: undefined, isLoading: false }),
  useGetConnectorsConnectorIdSnapshotsDiff: (_id: string, params: { from: string; to: string }, options: { query: { enabled: boolean } }) => options.query.enabled ? ({
    data: {
      provenance: { from: { id: params.from }, to: { id: params.to } },
      summary: { sectionsAdded: 0, sectionsRemoved: 0, sectionsModified: 1, entitiesAdded: 0, entitiesRemoved: 0, entitiesModified: 1, dependenciesAdded: 0, dependenciesRemoved: 0 },
      sections: [{ type: 'modified', summary: 'CPU changed', patches: [{ section: 'State', old: 'one', new: 'two' }] }],
      entities: [{ kind: 'host', key: 'h1', name: 'host1', field: 'ip', change: 'modified', old: '1', new: '2' }],
      dependencies: [],
    }, isLoading: false,
  }) : ({ data: undefined, isLoading: false }),
  postConnectorsConnectorIdGoldenSnapshot: vi.fn(),
  deleteConnectorsConnectorIdGoldenSnapshot: vi.fn(),
}));
vi.mock('../../../components/diff/DiffViewer', () => ({ DocDiff: ({ before, after }: { before: string; after: string }) => <div>{before} → {after}</div> }));

function mount(initial = '/services/svc/snapshots') {
  return render(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}><MemoryRouter initialEntries={[initial]}><Routes><Route path="/services/:id/snapshots" element={<SnapshotsPage />} /></Routes></MemoryRouter></QueryClientProvider>);
}

beforeEach(() => {
  get.mockReset();
  download.mockReset();
  get.mockImplementation((_url: string, options?: { responseType?: string }) => options?.responseType === 'blob'
    ? Promise.resolve({ data: new Blob(['report']), headers: {} })
    : Promise.resolve({ data: [
      { id: 'new', fetchedAt: '2026-09-26T10:00:00Z', sizeBytes: 1024, golden: false },
      { id: 'old', fetchedAt: '2026-09-25T10:00:00Z', sizeBytes: 512, golden: true },
    ], headers: {} }));
});

describe('SnapshotsPage', () => {
  it('compares two selected snapshots and renders section and entity changes', async () => {
    mount();
    fireEvent.click(await screen.findByLabelText('Select snapshot old'));
    fireEvent.click(screen.getByLabelText('Select snapshot new'));
    fireEvent.click(screen.getByText('Compare selected'));
    expect(await screen.findByText('CPU changed')).toBeInTheDocument();
    expect(screen.getByText('one → two')).toBeInTheDocument();
    expect(screen.getByText('host1')).toBeInTheDocument();
  });

  it('compares a row with its predecessor and exports a blob', async () => {
    mount();
    await screen.findByLabelText('Select snapshot new');
    fireEvent.click(screen.getAllByText('Vs previous')[0]);
    expect(await screen.findByText('CPU changed')).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('Export'), { target: { value: 'csv' } });
    await waitFor(() => expect(download).toHaveBeenCalledOnce());
    expect(get).toHaveBeenCalledWith('/connectors/svc/snapshots/diff', expect.objectContaining({ params: { from: 'old', to: 'new', format: 'csv' }, responseType: 'blob' }));
  });

  it('finds the predecessor when a linked sync snapshot is on a later page', async () => {
    get.mockImplementation((_url: string, options?: { params?: { cursor?: string } }) => options?.params?.cursor
      ? Promise.resolve({ data: [
        { id: 'target', fetchedAt: '2026-09-25T10:00:00Z', sizeBytes: 512, golden: false },
        { id: 'older', fetchedAt: '2026-09-24T10:00:00Z', sizeBytes: 256, golden: false },
      ], headers: {} })
      : Promise.resolve({ data: [
        { id: 'newer', fetchedAt: '2026-09-26T10:00:00Z', sizeBytes: 1024, golden: false },
      ], headers: { 'x-next-cursor': 'page-2' } }));

    mount('/services/svc/snapshots?b=target');
    expect(await screen.findByText('older → target')).toBeInTheDocument();
    expect(get).toHaveBeenCalledTimes(2);
    expect(get).toHaveBeenCalledWith('/connectors/svc/snapshots', expect.objectContaining({ params: { limit: 30, cursor: 'page-2' } }));
  });
});
