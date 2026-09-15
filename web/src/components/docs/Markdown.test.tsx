import { afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cleanup, render, screen, waitFor } from '@testing-library/react';
import { Markdown } from './Markdown';

describe('Markdown', () => {
  // jsdom has no layout engine, so SVGElement.getBBox() (used by mermaid to
  // measure node/edge label text) doesn't exist. Real browsers implement it;
  // this stub just lets a diagram finish laying out in a jsdom test.
  beforeAll(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- jsdom's SVGElement has no getBBox at all
    (SVGElement.prototype as any).getBBox = () => ({ x: 0, y: 0, width: 0, height: 0 });
  });

  afterEach(cleanup);

  it('renders a normal fenced code block as <pre><code>', () => {
    render(<Markdown source={'```json\n{"a":1}\n```'} />);
    expect(screen.getByText('{"a":1}').closest('pre')).toBeInTheDocument();
  });

  it('renders a ```mermaid fenced block via the Mermaid component instead of a plain code block', async () => {
    const { container } = render(<Markdown source={'```mermaid\ngraph LR\n  a --> b\n```'} />);

    // The Mermaid component owns its own container div and never falls back
    // to a <pre> wrapper for a mermaid block.
    expect(container.querySelector('pre')).not.toBeInTheDocument();

    await waitFor(() => {
      expect(container.querySelector('svg')).toBeInTheDocument();
    });
  });
});
