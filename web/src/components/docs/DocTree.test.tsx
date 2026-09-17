import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import '../../i18n';
import { DocTree } from './DocTree';
import type { DocNode } from '../../api/model';

let isInstanceAdmin = false;
let connectorRole: 'viewer' | 'operator' | undefined;
vi.mock('../../hooks/useRole', () => ({
  useIsInstanceAdmin: () => isInstanceAdmin,
  useConnectorRole: () => connectorRole,
}));

const tree: DocNode = {
  docId: 'root',
  title: 'Lab Documentation',
  kind: 'lab',
  children: [{ docId: 'c1', title: 'Connector One', kind: 'service' }],
};

function renderTree(onShare?: (n: { docId: string; title: string }) => void) {
  return render(
    <MemoryRouter>
      <DocTree tree={tree} onShare={onShare} />
    </MemoryRouter>
  );
}

describe('DocTree Share affordance (#240 PR2)', () => {
  afterEach(() => {
    cleanup();
    isInstanceAdmin = false;
    connectorRole = undefined;
    vi.clearAllMocks();
  });

  it('renders no Share buttons when onShare is not provided', () => {
    renderTree();
    expect(screen.queryByRole('button', { name: /share/i })).not.toBeInTheDocument();
  });

  it('hides the connector Share button for a viewer', () => {
    connectorRole = 'viewer';
    renderTree(vi.fn());
    // Only the root's Share button (instance-admin gated, also hidden here) should be absent too.
    expect(screen.queryAllByRole('button', { name: /share/i })).toHaveLength(0);
  });

  it('shows the connector Share button for an operator on that connector', () => {
    connectorRole = 'operator';
    renderTree(vi.fn());
    expect(screen.getAllByRole('button', { name: /share/i })).toHaveLength(1);
  });

  it('shows the root Share button for an instance admin', () => {
    isInstanceAdmin = true;
    renderTree(vi.fn());
    expect(screen.getAllByRole('button', { name: /share/i }).length).toBeGreaterThanOrEqual(1);
  });

  it('clicking a connector Share button calls onShare with that node', () => {
    connectorRole = 'operator';
    const onShare = vi.fn();
    renderTree(onShare);
    fireEvent.click(screen.getByRole('button', { name: /share/i }));
    expect(onShare).toHaveBeenCalledWith({ docId: 'c1', title: 'Connector One', connectorId: 'c1' });
  });
});
