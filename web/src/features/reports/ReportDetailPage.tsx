import { useNavigate, useParams } from 'react-router-dom';
import { getReportsReportIdDownload, useGetReportsReportId } from '../../api/generated/reports/reports';
import { Button } from '../../components/ui/Button';
import { Markdown } from '../../components/docs/Markdown';
import { ErrorState, SkeletonRows } from '../../components/ui/states';
import { Panel } from '../../components/ui/Panel';
import { DownloadIcon } from '../../components/icons';
import { toast } from '../../lib/toast';

export function ReportDetailPage() {
  const { id = '' } = useParams();
  const navigate = useNavigate();
  const report = useGetReportsReportId(id, { query: { enabled: !!id } });
  const download = async (format: 'md' | 'html') => {
    try {
      const body = await getReportsReportIdDownload(id, { format });
      const url = URL.createObjectURL(new Blob([body], { type: format === 'md' ? 'text/markdown' : 'text/html' }));
      const link = document.createElement('a'); link.href = url; link.download = `report-${id}.${format}`; link.click(); URL.revokeObjectURL(url);
    } catch { toast.error('Could not download report.'); }
  };
  if (report.isLoading) return <div className="mx-auto max-w-225 px-6 py-6"><SkeletonRows rows={6} /></div>;
  if (report.isError || !report.data) return <div className="mx-auto max-w-225 px-6 py-6"><ErrorState description="Could not load this report." onRetry={() => report.refetch()} /></div>;
  const data = report.data;
  return <div className="mx-auto max-w-225 px-6 py-6"><header className="mb-5 flex flex-wrap items-start justify-between gap-3"><div><Button size="sm" onClick={() => navigate('/reports')}>Back to reports</Button><h1 className="mt-3 text-xl font-semibold tracking-tight text-ink">{data.definitionName}</h1><p className="text-sm text-ink-muted">{new Date(data.periodStart).toLocaleString()} – {new Date(data.periodEnd).toLocaleString()}</p></div><div className="flex gap-2"><Button size="sm" onClick={() => void download('md')}><DownloadIcon size={14} />Markdown</Button><Button size="sm" onClick={() => void download('html')}><DownloadIcon size={14} />HTML</Button></div></header>{(data.status === 'partial' || data.truncated) && <div className="mb-4 rounded-sm border border-warn/40 bg-warn-tint px-3 py-2 text-sm text-ink">{data.status === 'partial' ? 'This report is partial because one or more sections were unavailable.' : 'This report window was truncated to the maximum period.'}</div>}<Panel className="p-5"><Markdown source={data.markdown} /></Panel></div>;
}
