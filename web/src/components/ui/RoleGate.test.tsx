import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RoleGate } from './RoleGate';

let isInstanceAdmin = false;
let connectorRole: 'viewer' | 'operator' | undefined;

vi.mock('../../hooks/useRole', () => ({
  useIsInstanceAdmin: () => isInstanceAdmin,
  useConnectorRole: () => connectorRole,
}));

describe('RoleGate (#240 PR1)', () => {
  it('without connectorId, gates on instance-admin', () => {
    isInstanceAdmin = false;
    render(
      <RoleGate fallback={<p>hidden</p>}>
        <p>visible</p>
      </RoleGate>
    );
    expect(screen.getByText('hidden')).toBeInTheDocument();

    isInstanceAdmin = true;
    render(
      <RoleGate fallback={<p>hidden2</p>}>
        <p>visible2</p>
      </RoleGate>
    );
    expect(screen.getByText('visible2')).toBeInTheDocument();
  });

  it('with connectorId, gates on that connector role instead of instance-admin', () => {
    isInstanceAdmin = false;
    connectorRole = 'operator';
    render(
      <RoleGate connectorId="c1" fallback={<p>hidden</p>}>
        <p>operator content</p>
      </RoleGate>
    );
    expect(screen.getByText('operator content')).toBeInTheDocument();
  });

  it('with connectorId and minRole="viewer", a viewer grant is enough', () => {
    connectorRole = 'viewer';
    render(
      <RoleGate connectorId="c1" minRole="viewer" fallback={<p>hidden</p>}>
        <p>viewer content</p>
      </RoleGate>
    );
    expect(screen.getByText('viewer content')).toBeInTheDocument();
  });

  it('with connectorId, a viewer role does not clear the default operator minimum', () => {
    connectorRole = 'viewer';
    render(
      <RoleGate connectorId="c1" fallback={<p>hidden</p>}>
        <p>operator-only content</p>
      </RoleGate>
    );
    expect(screen.getByText('hidden')).toBeInTheDocument();
    expect(screen.queryByText('operator-only content')).not.toBeInTheDocument();
  });

  it('with connectorId and no grant, renders nothing without a fallback', () => {
    connectorRole = undefined;
    const { container } = render(
      <RoleGate connectorId="c1">
        <p>content</p>
      </RoleGate>
    );
    expect(container).toBeEmptyDOMElement();
  });
});
