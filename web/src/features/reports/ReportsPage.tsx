import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  deleteReportsDefinitionsReportDefinitionId,
  getGetReportsDefinitionsQueryKey,
  getGetReportsQueryKey,
  postReportsDefinitions,
  postReportsDefinitionsReportDefinitionIdRun,
  putReportsDefinitionsReportDefinitionId,
  useGetReports,
  useGetReportsDefinitions,
} from '../../api/generated/reports/reports';
import { ReportChannelType, ReportSection, type ReportDefinition, type ReportDefinitionInput } from '../../api/model';
import { Button } from '../../components/ui/Button';
import { Dialog } from '../../components/ui/Dialog';
import { EmptyState, ErrorState, SkeletonRows } from '../../components/ui/states';
import { Panel } from '../../components/ui/Panel';
import { PlusIcon } from '../../components/icons';
import { toast } from '../../lib/toast';

const sections = Object.values(ReportSection);
const channels = Object.values(ReportChannelType);
const emptyDefinition = (): ReportDefinitionInput => ({
  slug: '', name: '', enabled: true, cronExpr: '0 9 * * 1', timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC', sections: [...sections], connectorIds: [], channels: [],
});
const when = (value: string) => new Date(value).toLocaleString();

export function ReportsPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const definitions = useGetReportsDefinitions();
  const reports = useGetReports({ pageSize: 50 });
  const [editing, setEditing] = useState<ReportDefinition | null | undefined>(undefined);
  const [deleting, setDeleting] = useState<ReportDefinition | null>(null);
  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: getGetReportsDefinitionsQueryKey() });
    void queryClient.invalidateQueries({ queryKey: getGetReportsQueryKey() });
  };
  const save = useMutation({
    mutationFn: ({ existing, input }: { existing: ReportDefinition | null; input: ReportDefinitionInput }) =>
      existing ? putReportsDefinitionsReportDefinitionId(existing.id, input) : postReportsDefinitions(input),
    onSuccess: () => { invalidate(); setEditing(undefined); toast.success('Report definition saved.'); },
    onError: () => toast.error('Could not save report definition.'),
  });
  const run = useMutation({
    mutationFn: (id: string) => postReportsDefinitionsReportDefinitionIdRun(id),
    onSuccess: (report) => { invalidate(); toast.success('Report generated.'); navigate(`/reports/${report.id}`); },
    onError: () => toast.error('Could not generate report.'),
  });
  const remove = useMutation({
    mutationFn: (id: string) => deleteReportsDefinitionsReportDefinitionId(id),
    onSuccess: () => { invalidate(); setDeleting(null); toast.success('Report definition deleted.'); },
    onError: () => toast.error('Could not delete report definition.'),
  });

  return <div className="mx-auto max-w-300 px-6 py-6">
    <header className="mb-5 flex flex-wrap items-center justify-between gap-3">
      <div><h1 className="text-xl font-semibold tracking-tight text-ink">Reports</h1><p className="text-sm text-ink-muted">Scheduled snapshots of your lab, available to download and review.</p></div>
      <Button variant="primary" onClick={() => setEditing(null)}><PlusIcon size={15} />New definition</Button>
    </header>
    <section className="mb-7"><h2 className="mb-2 font-mono text-xs font-semibold uppercase tracking-[0.16em] text-ink-muted">Schedules</h2>
      <Panel>{definitions.isLoading ? <SkeletonRows rows={3} /> : definitions.isError || !definitions.data ? <ErrorState description="Could not load report definitions." onRetry={() => definitions.refetch()} /> : definitions.data.length === 0 ? <EmptyState title="No report definitions" description="Create a schedule to begin collecting reports." /> : definitions.data.map((definition) => <div key={definition.id} className="flex flex-wrap items-center gap-3 border-b border-line-soft px-4 py-3 last:border-0"><button className="min-w-45 flex-1 text-left" onClick={() => setEditing(definition)}><div className="font-medium text-ink">{definition.name}</div><div className="font-mono text-2xs text-ink-faint">{definition.cronExpr} · {definition.timezone} · {definition.enabled ? 'enabled' : 'disabled'}</div></button><div className="flex gap-2"><Button size="sm" onClick={() => run.mutate(definition.id)} disabled={run.isPending}>Run now</Button><Button size="sm" onClick={() => setEditing(definition)}>Edit</Button><Button size="sm" variant="danger" onClick={() => setDeleting(definition)}>Delete</Button></div></div>)}</Panel>
    </section>
    <section><h2 className="mb-2 font-mono text-xs font-semibold uppercase tracking-[0.16em] text-ink-muted">Generated reports</h2>
      <Panel>{reports.isLoading ? <SkeletonRows rows={5} /> : reports.isError || !reports.data ? <ErrorState description="Could not load reports." onRetry={() => reports.refetch()} /> : reports.data.items.length === 0 ? <EmptyState title="No reports yet" description="Run a definition now or wait for its schedule." /> : reports.data.items.map((report) => <button key={report.id} onClick={() => navigate(`/reports/${report.id}`)} className="flex w-full flex-wrap items-center gap-3 border-b border-line-soft px-4 py-3 text-left last:border-0 hover:bg-surface-raised"><span className="min-w-45 flex-1"><span className="block font-medium text-ink">{report.definitionName}</span><span className="font-mono text-2xs text-ink-faint">{when(report.periodStart)} – {when(report.periodEnd)}</span></span><span className={report.status === 'partial' ? 'font-mono text-xs text-warn' : 'font-mono text-xs text-ok'}>{report.status}</span><span className="font-mono text-2xs text-ink-faint">{report.trigger}</span></button>)}</Panel>
    </section>
    {editing !== undefined && <DefinitionDialog definition={editing} pending={save.isPending} onClose={() => setEditing(undefined)} onSave={(input) => save.mutate({ existing: editing, input })} />}
    <Dialog open={!!deleting} onClose={() => setDeleting(null)} title="Delete report definition" size="sm"><p className="text-sm text-ink-muted">Generated reports will be kept. This schedule will be removed.</p><div className="mt-5 flex justify-end gap-2"><Button size="sm" onClick={() => setDeleting(null)}>Cancel</Button><Button size="sm" variant="danger" disabled={remove.isPending} onClick={() => deleting && remove.mutate(deleting.id)}>Delete</Button></div></Dialog>
  </div>;
}

function DefinitionDialog({ definition, pending, onClose, onSave }: { definition: ReportDefinition | null; pending: boolean; onClose: () => void; onSave: (input: ReportDefinitionInput) => void }) {
  const [form, setForm] = useState<ReportDefinitionInput>(() => definition ? { slug: definition.slug, name: definition.name, enabled: definition.enabled, cronExpr: definition.cronExpr, timezone: definition.timezone, sections: definition.sections, connectorIds: definition.connectorIds, channels: definition.channels } : emptyDefinition());
  const toggleSection = (section: ReportSection) => setForm((current) => ({ ...current, sections: current.sections.includes(section) ? current.sections.filter((item) => item !== section) : [...current.sections, section] }));
  const toggleChannel = (channel: ReportChannelType) => setForm((current) => {
    const selected = current.channels ?? [];
    return { ...current, channels: selected.includes(channel) ? selected.filter((item) => item !== channel) : [...selected, channel] };
  });
  return <Dialog open onClose={onClose} title={definition ? 'Edit report definition' : 'New report definition'} size="lg"><form onSubmit={(event) => { event.preventDefault(); onSave(form); }} className="grid gap-4 sm:grid-cols-2">
    {(['name', 'slug', 'cronExpr', 'timezone'] as const).map((field) => <label key={field} className="block"><span className="mb-1 block text-xs text-ink-muted">{field === 'cronExpr' ? 'Schedule (cron)' : field}</span><input required={field !== 'cronExpr' || true} disabled={field === 'slug' && !!definition} value={form[field]} onChange={(event) => setForm({ ...form, [field]: event.target.value })} className="h-9 w-full rounded-sm border border-line bg-surface px-2.5 text-sm text-ink disabled:opacity-50" /></label>)}
    <label className="flex items-center gap-2 text-sm text-ink"><input type="checkbox" checked={form.enabled} onChange={(event) => setForm({ ...form, enabled: event.target.checked })} />Enabled</label>
    <fieldset className="sm:col-span-2"><legend className="mb-1 text-xs text-ink-muted">Sections</legend><div className="flex flex-wrap gap-3">{sections.map((section) => <label key={section} className="text-sm text-ink"><input type="checkbox" checked={form.sections.includes(section)} onChange={() => toggleSection(section)} /> {section}</label>)}</div></fieldset>
    <fieldset className="sm:col-span-2"><legend className="mb-1 text-xs text-ink-muted">Notification routing</legend><p className="mb-2 text-2xs text-ink-faint">Deliver directly to every enabled channel of the selected type.</p><div className="flex flex-wrap gap-3">{channels.map((channel) => <label key={channel} className="text-sm text-ink"><input type="checkbox" checked={form.channels?.includes(channel)} onChange={() => toggleChannel(channel)} /> {channel}</label>)}</div></fieldset>
    <div className="sm:col-span-2 flex justify-end gap-2"><Button type="button" size="sm" onClick={onClose}>Cancel</Button><Button type="submit" size="sm" variant="primary" disabled={pending || !form.name.trim() || !form.slug.trim() || !form.cronExpr.trim() || form.sections.length === 0}>{pending ? 'Saving…' : 'Save'}</Button></div>
  </form></Dialog>;
}
