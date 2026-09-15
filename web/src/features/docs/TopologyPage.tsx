/**
 * Lab-wide topology: finds the single "Lab Topology" doc (created/refreshed
 * by POST /docs/topology) and redirects into DocsPage, which already knows
 * how to render a doc (Markdown, incl. the Mermaid diagram inside it) and
 * its version history — nothing to duplicate here. Operators additionally
 * get a button to generate it on first visit / refresh it on demand.
 */
import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { useGetDocs, postDocsTopology } from '../../api/generated/docs/docs';
import { Panel } from '../../components/ui/Panel';
import { Button } from '../../components/ui/Button';
import { RoleGate } from '../../components/ui/RoleGate';
import { SkeletonRows, ErrorState, EmptyState } from '../../components/ui/states';
import { NetworkIcon } from '../../components/icons';

const TOPOLOGY_TITLE = 'Lab Topology';

export function TopologyPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { data, isLoading, isError, refetch } = useGetDocs({ search: TOPOLOGY_TITLE, pageSize: 5 });

  const generate = useMutation({
    mutationFn: () => postDocsTopology(),
    onSuccess: (result) => navigate(`/docs/${result.docId}`, { replace: true }),
  });

  const existing = data?.items.find((d) => d.title === TOPOLOGY_TITLE && d.kind === 'lab');

  useEffect(() => {
    if (existing) navigate(`/docs/${existing.docId}`, { replace: true });
  }, [existing, navigate]);

  if (existing) return null;

  return (
    <div className="mx-auto max-w-205 px-6 py-6">
      <Panel className="min-h-[50vh]">
        {isLoading ? (
          <div className="p-6">
            <SkeletonRows rows={6} />
          </div>
        ) : isError ? (
          <ErrorState description={t('docs.topology.loadError')} onRetry={() => refetch()} />
        ) : (
          <EmptyState
            icon={<NetworkIcon size={20} />}
            title={t('docs.topology.emptyTitle')}
            description={
              generate.isPending
                ? t('docs.topology.generating')
                : t('docs.topology.emptyDesc')
            }
            action={
              <RoleGate fallback={null}>
                <Button size="sm" onClick={() => generate.mutate()} disabled={generate.isPending}>
                  {t('docs.topology.generateAction')}
                </Button>
              </RoleGate>
            }
          />
        )}
      </Panel>
    </div>
  );
}
