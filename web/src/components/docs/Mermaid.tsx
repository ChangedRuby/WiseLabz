/**
 * Renders a Mermaid diagram source string as SVG. Themed off the app's live
 * CSS custom properties (set by the palette system in ../../store/theme.ts)
 * rather than a light/dark boolean, since this app has no such flag — just
 * re-reads the current --color-* values, so it re-subscribes to useTheme to
 * re-render whenever the palette changes.
 */
import { useEffect, useId, useRef } from 'react';
import { useTheme } from '../../store/theme';

function cssVar(name: string, fallback: string): string {
  if (typeof window === 'undefined') return fallback;
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  return value || fallback;
}

// Resolves a CSS custom property to a color format mermaid's own color
// parser (khroma) understands. The app's palette tokens are oklch(...),
// which khroma can't parse, so raw getPropertyValue text isn't usable
// directly — instead let the browser's own CSS engine resolve it (any
// valid CSS color resolves through `color` to an rgb()/rgba() string).
function resolveColor(varName: string, fallback: string): string {
  if (typeof document === 'undefined') return fallback;
  const probe = document.createElement('span');
  probe.style.color = `var(${varName})`;
  document.body.appendChild(probe);
  const resolved = getComputedStyle(probe).color;
  document.body.removeChild(probe);
  // A CSS engine that can't resolve var()/oklch() (e.g. jsdom in tests)
  // echoes the property back unparsed rather than erroring — fall back
  // rather than hand mermaid's color parser something it can't read.
  return resolved && !resolved.includes('var(') ? resolved : fallback;
}

export function Mermaid({ chart }: { chart: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const id = useId().replace(/[^a-zA-Z0-9]/g, '');
  const mode = useTheme((s) => s.mode);
  const preset = useTheme((s) => s.preset);
  const custom = useTheme((s) => s.custom);

  useEffect(() => {
    let cancelled = false;
    const container = containerRef.current;
    if (!container) return;

    import('mermaid')
      .then(({ default: mermaid }) => {
        mermaid.initialize({
          startOnLoad: false,
          theme: 'base',
          themeVariables: {
            background: resolveColor('--color-canvas', '#0b0c0f'),
            primaryColor: resolveColor('--color-canvas-sunken', '#15171b'),
            primaryTextColor: resolveColor('--color-ink', '#e8e8ea'),
            primaryBorderColor: resolveColor('--color-line-soft', '#2a2d33'),
            lineColor: resolveColor('--color-accent-primary', '#5b8cff'),
            fontFamily: cssVar('--font-mono', 'ui-monospace, monospace'),
          },
        });
        return mermaid.render(`mermaid-${id}`, chart);
      })
      .then(({ svg }) => {
        if (!cancelled && containerRef.current) containerRef.current.innerHTML = svg;
      })
      .catch((err: unknown) => {
        if (!cancelled && containerRef.current) {
          containerRef.current.textContent = `Diagram error: ${err instanceof Error ? err.message : String(err)}`;
        }
      });

    return () => {
      cancelled = true;
    };
  }, [chart, id, mode, preset, custom]);

  return (
    <div
      ref={containerRef}
      className="my-3 overflow-x-auto rounded-lg border border-line-soft bg-canvas-sunken p-3 text-xs text-ink-muted"
    />
  );
}
