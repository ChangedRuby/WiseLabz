/**
 * Markdown renderer for WiseLabz's generated docs, on react-markdown + remark-gfm
 * (tables, strikethrough, task lists, autolinks) instead of the old hand-rolled
 * line-scanner. `components` maps each element to the same Tailwind classes the
 * previous renderer used, so existing docs render pixel-identical; gfm adds
 * correctness (nested lists, escaping, ordered lists) the old parser didn't have.
 * Inline vs. block code is told apart with a CSS `:not(pre)` selector rather than
 * a JS heuristic, since react-markdown's `code` renderer no longer reports it.
 */
import { isValidElement } from 'react';
import ReactMarkdown, { type Components } from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Mermaid } from './Mermaid';

const headingClass: Record<'h1' | 'h2' | 'h3' | 'h4', string> = {
  h1: 'mt-1 mb-3 font-mono text-2xl font-semibold tracking-tight text-balance text-ink',
  h2: 'mb-2 mt-6 font-mono text-lg font-semibold tracking-tight text-ink',
  h3: 'mb-1.5 mt-5 font-mono text-xs font-semibold uppercase tracking-[0.16em] text-[var(--color-ink-muted)] text-ink',
  h4: 'mb-1.5 mt-5 font-mono text-xs font-semibold uppercase tracking-[0.16em] text-[var(--color-ink-muted)] text-ink',
};

const components: Components = {
  h1: ({ children }) => <h1 className={headingClass.h1}>{children}</h1>,
  h2: ({ children }) => <h2 className={headingClass.h2}>{children}</h2>,
  h3: ({ children }) => <h3 className={headingClass.h3}>{children}</h3>,
  h4: ({ children }) => <h4 className={headingClass.h4}>{children}</h4>,
  p: ({ children }) => <p className="my-2.5 text-sm leading-relaxed text-ink-muted text-pretty">{children}</p>,
  strong: ({ children }) => <strong className="font-semibold text-ink">{children}</strong>,
  blockquote: ({ children }) => (
    <blockquote className="my-3 flex gap-2.5 rounded-sm border border-line-soft bg-canvas-sunken px-4 py-2.5 text-sm text-ink-muted">
      <span className="mt-1.5 h-1.5 w-1.5 shrink-0 bg-accent-primary" />
      <div>{children}</div>
    </blockquote>
  ),
  ul: ({ children }) => <ul className="my-3 space-y-1.5 pl-1">{children}</ul>,
  ol: ({ children }) => <ol className="my-3 space-y-1.5 pl-5 text-sm text-ink-muted [list-style:decimal]">{children}</ol>,
  li: ({ children, className }) => {
    // Task-list items (remark-gfm) get a "task-list-item" className and render
    // their own checkbox — keep those bare, dot-prefix everything else.
    if (className?.includes('task-list-item')) {
      return <li className="flex items-center gap-2 text-sm text-ink-muted">{children}</li>;
    }
    return (
      <li className="flex gap-2.5 text-sm text-ink-muted">
        <span className="mt-2 h-1 w-1 shrink-0 rounded-full bg-accent-primary" />
        <span>{children}</span>
      </li>
    );
  },
  pre: ({ children }) => {
    // A fenced ```mermaid block renders through <Mermaid>, which owns its
    // own frame — skip the generic <pre> wrapper for it.
    const child = isValidElement<{ className?: string }>(children) ? children : null;
    if (child?.props.className?.includes('language-mermaid')) {
      return <>{children}</>;
    }
    return (
      <pre className="my-3 overflow-x-auto rounded-lg border border-line-soft bg-canvas-sunken p-3 font-mono text-xs leading-relaxed text-ink-muted">
        {children}
      </pre>
    );
  },
  code: ({ children, className }) => {
    if (className?.includes('language-mermaid')) {
      return <Mermaid chart={String(children).replace(/\n$/, '')} />;
    }
    return (
      <code className="rounded bg-canvas-sunken px-1.5 py-0.5 font-mono text-[0.85em] text-accent-primary-bright [pre_&]:rounded-none [pre_&]:bg-transparent [pre_&]:px-0 [pre_&]:py-0 [pre_&]:text-[1em] [pre_&]:text-inherit">
        {children}
      </code>
    );
  },
  table: ({ children }) => (
    <div className="my-3 overflow-x-auto rounded-lg border border-line-soft">
      <table className="w-full border-collapse text-sm">{children}</table>
    </div>
  ),
  tr: ({ children }) => <tr className="transition-colors hover:bg-surface-raised">{children}</tr>,
  th: ({ children }) => (
    <th className="border-b border-line-soft bg-canvas-sunken px-3 py-2 text-left text-xs font-semibold uppercase tracking-wide text-ink-muted">
      {children}
    </th>
  ),
  td: ({ children }) => (
    <td className="border-b border-line-soft px-3 py-2 text-ink last:border-0 [tr:last-child_&]:border-0">
      {children}
    </td>
  ),
};

export function Markdown({ source }: { source: string }) {
  return (
    <div className="max-w-[68ch]">
      <ReactMarkdown remarkPlugins={[remarkGfm]} components={components}>
        {source}
      </ReactMarkdown>
    </div>
  );
}
