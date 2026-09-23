import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import '../../i18n';
import type { NotificationConfig } from '../../api/model';
import { EventRoutingTable } from './EventRoutingTable';

vi.mock('../../api/generated/connectors/connectors', () => ({
  useGetConnectors: () => ({ data: [] }),
}));

const channels: NotificationConfig['channels'] = [
  { type: 'in_app', enabled: true },
  { type: 'webhook', enabled: true, config: { url: 'https://hooks.example.com' } },
];

describe('EventRoutingTable', () => {
  afterEach(() => {
    cleanup();
  });

  it('shows known event types that have no saved route yet', () => {
    const config: NotificationConfig = {
      channels,
      routing: [{ eventType: 'alert.created', channel: 'webhook', enabled: true }],
    };
    render(<EventRoutingTable config={config} onChange={vi.fn()} />);
    expect(screen.getByText('alert created')).toBeInTheDocument();
    expect(screen.getByText('finding created')).toBeInTheDocument();
    expect(screen.getByText('system job failed')).toBeInTheDocument();

    const cell = screen.getByRole('switch', { name: 'system job failed → Webhook' });
    expect(cell).not.toBeDisabled();
    expect(cell).toHaveAttribute('aria-checked', 'false');
  });

  it('keeps custom event types from the saved routing', () => {
    const config: NotificationConfig = {
      channels,
      routing: [{ eventType: 'sync.failed', channel: 'webhook', enabled: true }],
    };
    render(<EventRoutingTable config={config} onChange={vi.fn()} />);
    expect(screen.getByText('sync failed')).toBeInTheDocument();
  });

  it('creates routes when enabling an unrouted event', () => {
    const onChange = vi.fn();
    const config: NotificationConfig = {
      channels,
      routing: [{ eventType: 'alert.created', channel: 'webhook', enabled: true }],
    };
    render(<EventRoutingTable config={config} onChange={onChange} />);
    fireEvent.click(screen.getByRole('switch', { name: 'system job failed → Webhook' }));

    const next = onChange.mock.calls[0][0] as NotificationConfig;
    expect(next.routing).toContainEqual({
      eventType: 'alert.created',
      channel: 'webhook',
      enabled: true,
    });
    expect(next.routing).toContainEqual(
      expect.objectContaining({ eventType: 'system.job_failed', channel: 'webhook', enabled: true })
    );
    expect(next.routing).toContainEqual(
      expect.objectContaining({ eventType: 'system.job_failed', channel: 'in_app', enabled: false })
    );
  });

  it('treats an empty routing as deliver-everything and keeps that on the first edit', () => {
    const onChange = vi.fn();
    render(<EventRoutingTable config={{ channels, routing: [] }} onChange={onChange} />);
    const cell = screen.getByRole('switch', { name: 'system job failed → Webhook' });
    expect(cell).toHaveAttribute('aria-checked', 'true');

    fireEvent.click(cell);
    const next = onChange.mock.calls[0][0] as NotificationConfig;
    const find = (eventType: string, channel: string) =>
      next.routing.find((r) => r.eventType === eventType && r.channel === channel);
    expect(find('system.job_failed', 'webhook')?.enabled).toBe(false);
    expect(find('alert.created', 'webhook')?.enabled).toBe(true);
    expect(find('finding.created', 'webhook')?.enabled).toBe(true);
  });
});
