import { useEffect, useState } from 'react';
import { useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link, useParams, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  useGetConnectorsConnectorId,
  useGetConnectorsConnectorIdGoldenSnapshot,
  getGetConnectorsConnectorIdGoldenSnapshotQueryKey,
  useGetConnectorsConnectorIdSnapshotsSnapshotId,
  useGetConnectorsConnectorIdSnapshotsDiff,
  postConnectorsConnectorIdGoldenSnapshot,
  deleteConnectorsConnectorIdGoldenSnapshot,
} from '../../../api/generated/connectors/connectors';
import type { EntityChange, SnapshotDiff, SnapshotSummary } from '../../../api/model';
import { AXIOS_INSTANCE } from '../../../api/axios-instance';
import { Button } from '../../../components/ui/Button';
import { Panel } from '../../../components/ui/Panel';
import { EmptyState, ErrorState, SkeletonRows } from '../../../components/ui/states';
import { TimeAgo } from '../../../components/ui/TimeAgo';
import { ToneTag } from '../../../components/ui/ToneTag';
import { RoleGate } from '../../../components/ui/RoleGate';
import { DocDiff } from '../../../components/diff/DiffViewer';
import { Markdown } from '../../../components/docs/Markdown';
import { downloadBlob, filenameFromContentDisposition } from '../../../lib/download';
import { toast } from '../../../lib/toast';

const PAGE_SIZE = 30;

function valueText(value: unknown): string {
  return value == null ? '—' : typeof value === 'string' ? value : JSON.stringify(value);
}

export function SnapshotsPage() {
  const { id = '' } = useParams();
  const { t } = useTranslation();
  const [search, setSearch] = useSearchParams();
  const connector = useGetConnectorsConnectorId(id);
  const goldenSnapshot = useGetConnectorsConnectorIdGoldenSnapshot(id, { query: { retry: false } });
  const snapshotId = search.get('s') ?? '';
  const from = search.get('a') ?? '';
  const to = search.get('b') ?? '';
  const [selected, setSelected] = useState<string[]>([]);
  const [exporting, setExporting] = useState(false);
  const queryClient = useQueryClient();

  const timeline = useInfiniteQuery({
    queryKey: ['connector-snapshot-timeline', id],
    initialPageParam: '',
    queryFn: async ({ pageParam }) => {
      const response = await AXIOS_INSTANCE.get<SnapshotSummary[]>(`/connectors/${id}/snapshots`, {
        params: { limit: PAGE_SIZE, cursor: pageParam },
      });
      return { rows: response.data, next: (response.headers['x-next-cursor'] as string | undefined) ?? '' };
    },
    getNextPageParam: (page) => page.next || undefined,
    enabled: !!id,
  });
  const rows = timeline.data?.pages.flatMap((page) => page.rows) ?? [];
  const targetIndex = rows.findIndex((row) => row.id === to);
  const activeFrom = from || (to && targetIndex >= 0 ? rows[targetIndex + 1]?.id ?? '' : '');
  const view = useGetConnectorsConnectorIdSnapshotsSnapshotId(id, snapshotId, {
    query: { enabled: !!snapshotId },
  });
  const diff = useGetConnectorsConnectorIdSnapshotsDiff(id, { from: activeFrom, to }, {
    query: { enabled: !!activeFrom && !!to },
  });
  const pin = useMutation({
    mutationFn: (snapshotId: string) => postConnectorsConnectorIdGoldenSnapshot(id, { snapshotId }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['connector-snapshot-timeline', id] });
      void queryClient.invalidateQueries({ queryKey: getGetConnectorsConnectorIdGoldenSnapshotQueryKey(id) });
    },
    onError: () => toast.error(t('services.snapshots.pinError')),
  });
  const unpin = useMutation({
    mutationFn: () => deleteConnectorsConnectorIdGoldenSnapshot(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['connector-snapshot-timeline', id] });
      void queryClient.invalidateQueries({ queryKey: getGetConnectorsConnectorIdGoldenSnapshotQueryKey(id) });
    },
    onError: () => toast.error(t('services.snapshots.pinError')),
  });

  // A sync row may link here with only its snapshot ID. Follow cursors until its predecessor is known.
  useEffect(() => {
    if (to && !from && !activeFrom && timeline.hasNextPage && !timeline.isFetchingNextPage) {
      void timeline.fetchNextPage();
    }
  }, [to, from, activeFrom, timeline]);

  const navigateDiff = (a: string, b: string) => setSearch({ a, b });
  const exportDiff = async (format: 'json' | 'csv' | 'md' | 'html') => {
    setExporting(true);
    try {
      const response = await AXIOS_INSTANCE.get(`/connectors/${id}/snapshots/diff`, {
        params: { from: activeFrom, to, format },
        responseType: 'blob',
      });
      downloadBlob(response.data as Blob, filenameFromContentDisposition(
        response.headers['content-disposition'] as string | undefined,
        `wiselabz-snapshot-diff-${id}.${format}`,
      ));
    } catch {
      toast.error(t('services.snapshots.exportError'));
    } finally {
      setExporting(false);
    }
  };

  return (
    <div className="space-y-5">
      <div>
        <Link to={`/services/${id}`} className="font-mono text-xs text-ink-muted hover:text-ink">← {t('services.snapshots.back')}</Link>
        <h1 className="mt-2 text-xl font-semibold text-ink">{t('services.snapshots.title', { name: connector.data?.name ?? '' })}</h1>
      </div>

      {snapshotId ? (
        <Panel className="space-y-5 p-5">
          <Button size="sm" onClick={() => setSearch({})}>{t('services.snapshots.timeline')}</Button>
          {view.isLoading ? <SkeletonRows rows={4} /> : view.isError || !view.data ? (
            <ErrorState description={t('services.snapshots.loadError')} onRetry={() => void view.refetch()} />
          ) : (
            <>
              <h2 className="text-sm font-semibold text-ink">{t('services.snapshots.view')} · <TimeAgo at={view.data.fetchedAt} /></h2>
              <div className="space-y-3">{view.data.sections.map((section) => (
                <details key={section.title} open className="rounded-md border border-line-soft bg-canvas-sunken p-3">
                  <summary className="cursor-pointer text-xs font-semibold text-ink">{section.title}</summary>
                  <div className="mt-2 text-sm"><Markdown source={section.content} /></div>
                </details>
              ))}</div>
              <h3 className="text-sm font-semibold text-ink">{t('services.snapshots.entities')}</h3>
              <div className="overflow-x-auto"><table className="w-full text-left text-xs"><thead><tr><th>{t('services.snapshots.kind')}</th><th>{t('services.snapshots.name')}</th><th>IP</th><th>{t('services.snapshots.hostname')}</th><th>{t('services.snapshots.attributes')}</th></tr></thead><tbody>{view.data.entities.map((entity, index) => (
                <tr key={`${entity.kind}:${entity.externalId ?? entity.name}:${index}`} className="border-t border-line-soft"><td>{entity.kind}</td><td>{entity.name}</td><td>{entity.ip ?? '—'}</td><td>{entity.hostname ?? '—'}</td><td className="font-mono">{valueText(entity.attributes)}</td></tr>
              ))}</tbody></table></div>
              <h3 className="text-sm font-semibold text-ink">{t('services.snapshots.dependencies')}</h3>
              <ul className="text-xs text-ink-muted">{view.data.dependencies.map((dep, index) => <li key={`${dep.kind}:${dep.name}:${dep.ref}:${index}`}>{dep.kind} · {dep.name} · {dep.ref ?? '—'}</li>)}</ul>
            </>
          )}
        </Panel>
      ) : to ? (
        <Panel className="space-y-5 p-5">
          <Button size="sm" onClick={() => setSearch({})}>{t('services.snapshots.timeline')}</Button>
          {!activeFrom && timeline.hasNextPage ? <SkeletonRows rows={3} /> : !activeFrom ? (
            <EmptyState title={t('services.snapshots.noPrevious')} />
          ) : diff.isLoading ? <SkeletonRows rows={4} /> : diff.isError || !diff.data || typeof diff.data === 'string' ? (
            <ErrorState description={t('services.snapshots.loadError')} onRetry={() => void diff.refetch()} />
          ) : <DiffContent diff={diff.data as SnapshotDiff} exporting={exporting} onExport={exportDiff} />}
        </Panel>
      ) : (
        <Panel className="p-5">
          <div className="mb-4 flex items-center justify-between gap-3">
            <h2 className="text-sm font-semibold text-ink">{t('services.snapshots.timeline')}</h2>
            <Button size="sm" disabled={selected.length !== 2} onClick={() => {
              const chosen = rows.filter((row) => selected.includes(row.id));
              if (chosen.length === 2) navigateDiff(chosen[1].id, chosen[0].id);
            }}>{t('services.snapshots.compareSelected')}</Button>
          </div>
          {timeline.isLoading ? <SkeletonRows rows={4} /> : timeline.isError ? <ErrorState description={t('services.snapshots.loadError')} onRetry={() => void timeline.refetch()} /> : rows.length === 0 ? <EmptyState title={t('services.snapshots.empty')} /> : (
            <ul className="divide-y divide-line-soft">{rows.map((row, index) => {
              const previous = rows[index + 1];
              const goldenId = goldenSnapshot.data?.snapshotId ?? rows.find((candidate) => candidate.golden)?.id;
              const checked = selected.includes(row.id);
              return <li key={row.id} className="flex flex-wrap items-center gap-2 py-3">
                <label className="flex min-w-0 flex-1 items-center gap-3 text-xs text-ink">
                  <input type="checkbox" aria-label={t('services.snapshots.select', { id: row.id })} checked={checked} disabled={!checked && selected.length === 2} onChange={() => setSelected(checked ? selected.filter((item) => item !== row.id) : [...selected, row.id])} />
                  <TimeAgo at={row.fetchedAt} />
                  <span className="font-mono text-ink-muted">{(row.sizeBytes / 1024).toFixed(1)} KiB</span>
                  {row.golden && <ToneTag tone="ok" label={t('services.snapshots.golden')} />}
                </label>
                <Button size="sm" variant="ghost" onClick={() => setSearch({ s: row.id })}>{t('services.snapshots.view')}</Button>
                <Button size="sm" variant="ghost" disabled={!previous && !timeline.hasNextPage} onClick={() => previous ? navigateDiff(previous.id, row.id) : setSearch({ b: row.id })}>{t('services.snapshots.vsPrevious')}</Button>
                {goldenId && goldenId !== row.id && <Button size="sm" variant="ghost" onClick={() => navigateDiff(goldenId, row.id)}>{t('services.snapshots.vsGolden')}</Button>}
                <RoleGate connectorId={id}><Button size="sm" variant="ghost" disabled={pin.isPending || unpin.isPending} onClick={() => row.golden ? unpin.mutate() : pin.mutate(row.id)}>{t(row.golden ? 'services.snapshots.unpin' : 'services.snapshots.pin')}</Button></RoleGate>
              </li>;
            })}</ul>
          )}
          {timeline.hasNextPage && <Button className="mt-4" size="sm" disabled={timeline.isFetchingNextPage} onClick={() => void timeline.fetchNextPage()}>{t('services.snapshots.loadMore')}</Button>}
        </Panel>
      )}
    </div>
  );
}

function DiffContent({ diff, exporting, onExport }: { diff: SnapshotDiff; exporting: boolean; onExport: (format: 'json' | 'csv' | 'md' | 'html') => void }) {
  const { t } = useTranslation();
  const summary = diff.summary;
  const grouped = diff.entities.reduce<Record<string, EntityChange[]>>((groups, change) => {
    (groups[`${change.kind}:${change.key}`] ??= []).push(change);
    return groups;
  }, {});
  return <>
    <div className="flex flex-wrap items-center justify-between gap-3">
      <h2 className="text-sm font-semibold text-ink">{t('services.snapshots.diff')}</h2>
      <label className="font-mono text-xs text-ink-muted">{t('services.snapshots.export')}{' '}<select aria-label={t('services.snapshots.export')} className="rounded-sm border border-line bg-surface p-1 text-ink" disabled={exporting} defaultValue="" onChange={(event) => { if (event.target.value) onExport(event.target.value as 'json' | 'csv' | 'md' | 'html'); event.target.value = ''; }}><option value="" disabled>{t('services.snapshots.chooseFormat')}</option><option value="json">JSON</option><option value="csv">CSV</option><option value="md">Markdown</option><option value="html">HTML</option></select></label>
    </div>
    <p className="font-mono text-xs text-ink-muted">{diff.provenance.from.id} → {diff.provenance.to.id}</p>
    <p className="text-xs text-ink-muted">{t('services.snapshots.summary', { sections: summary.sectionsAdded + summary.sectionsRemoved + summary.sectionsModified, entities: summary.entitiesAdded + summary.entitiesRemoved + summary.entitiesModified, dependencies: summary.dependenciesAdded + summary.dependenciesRemoved })}</p>
    <h3 className="text-sm font-semibold text-ink">{t('services.snapshots.sections')}</h3>
    {diff.sections.map((section, index) => <div key={`${section.summary}:${index}`} className="space-y-2"><ToneTag tone={section.type === 'added' ? 'ok' : section.type === 'removed' ? 'err' : 'signal'} label={section.type} /><h4 className="text-xs font-semibold text-ink">{section.summary}</h4>{section.patches.map((patch, patchIndex) => <DocDiff key={`${patch.section}:${patchIndex}`} before={patch.old ?? ''} after={patch.new ?? ''} label={patch.section} />)}</div>)}
    <h3 className="text-sm font-semibold text-ink">{t('services.snapshots.entities')}</h3>
    {Object.entries(grouped).map(([key, changes]) => <div key={key} className="overflow-x-auto"><h4 className="mb-1 font-mono text-xs text-ink">{changes?.[0]?.name ?? key}</h4><table className="w-full text-left text-xs"><thead><tr><th>{t('services.snapshots.field')}</th><th>{t('services.snapshots.change')}</th><th>{t('services.snapshots.old')}</th><th>{t('services.snapshots.new')}</th></tr></thead><tbody>{changes?.map((change, index) => <tr key={`${change.field}:${index}`} className="border-t border-line-soft"><td>{change.field}</td><td>{change.change}</td><td>{valueText(change.old)}</td><td>{valueText(change.new)}</td></tr>)}</tbody></table></div>)}
    <h3 className="text-sm font-semibold text-ink">{t('services.snapshots.dependencies')}</h3>
    <ul className="text-xs text-ink-muted">{diff.dependencies.map((dep, index) => <li key={`${dep.kind}:${dep.name}:${dep.ref}:${index}`}>{dep.change} · {dep.kind} · {dep.name} · {dep.ref}</li>)}</ul>
  </>;
}
