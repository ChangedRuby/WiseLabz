/**
 * Hierarchical doc navigation — lab root + per-service children. Current doc gets
 * the iris selection treatment; the lab root is always expanded (the tree is
 * shallow by design). Service nodes show their kind as a quiet mono tag.
 *
 * `onShare` (#240 PR2) surfaces a per-node "Share" affordance: the lab root
 * (docId "root") and each service child (docId is that connector's ID, per
 * the Tree() handler) can be the root of a read-only share link. Gated with
 * RoleGate so it only renders where the viewer could plausibly create one —
 * the server re-checks operator-on-every-covered-connector regardless.
 */
import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { cn } from '../../lib/cn';
import { FileTextIcon, LayersIcon, ShareIcon } from '../icons';
import { RoleGate } from '../ui/RoleGate';
import { IconButton } from '../ui/Button';
import type { DocNode } from '../../api/model';

interface DocTreeProps {
  tree: DocNode;
  onShare?: (node: { docId: string; title: string; connectorId?: string }) => void;
}

export function DocTree({ tree, onShare }: DocTreeProps) {
  const { t } = useTranslation();
  const shareLabel = t('docs.share.action', { defaultValue: 'Share' });

  return (
    <nav className="flex flex-col gap-0.5">
      <div className="group flex items-center gap-1">
        <NavLink
          to={`/docs/${tree.docId}`}
          end
          className={({ isActive }) =>
            cn(
              'flex flex-1 items-center gap-2.5 rounded-md px-2.5 py-2 text-sm font-medium transition-colors',
              isActive
                ? 'bg-accent-primary-tint text-ink'
                : 'text-ink-muted hover:bg-surface-raised hover:text-ink'
            )
          }
        >
          <LayersIcon size={16} className="text-accent-primary-bright" />
          {tree.title}
        </NavLink>
        {onShare && (
          <RoleGate>
            <IconButton
              label={shareLabel}
              onClick={() => onShare({ docId: tree.docId, title: tree.title })}
              className="shrink-0 opacity-0 group-hover:opacity-100"
            >
              <ShareIcon size={14} />
            </IconButton>
          </RoleGate>
        )}
      </div>

      <div className="ml-3 mt-0.5 flex flex-col gap-0.5 border-l border-line-soft pl-2">
        {(tree.children ?? []).map((node) => (
          <div key={node.docId} className="group flex items-center gap-1">
            <NavLink
              to={`/docs/${node.docId}`}
              className={({ isActive }) =>
                cn(
                  'flex flex-1 items-center gap-2.5 rounded-md px-2.5 py-1.5 text-sm transition-colors',
                  isActive
                    ? 'bg-accent-primary-tint text-ink'
                    : 'text-ink-muted hover:bg-surface-raised hover:text-ink'
                )
              }
            >
              <FileTextIcon size={15} className="shrink-0 opacity-70" />
              <span className="truncate">{node.title.split(' — ')[0]}</span>
              {node.title.includes(' — ') && (
                <span className="ml-auto truncate font-mono text-2xs text-ink-faint">
                  {node.title.split(' — ')[1]}
                </span>
              )}
            </NavLink>
            {onShare && (
              <RoleGate connectorId={node.docId} minRole="operator">
                <IconButton
                  label={shareLabel}
                  onClick={() => onShare({ docId: node.docId, title: node.title, connectorId: node.docId })}
                  className="shrink-0 opacity-0 group-hover:opacity-100"
                >
                  <ShareIcon size={14} />
                </IconButton>
              </RoleGate>
            )}
          </div>
        ))}
      </div>
    </nav>
  );
}
