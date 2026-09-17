import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { RequireAuth, RequireInstanceAdmin } from './guards';
import { useAuth } from '../../store/auth';

let isInstanceAdmin: boolean;

vi.mock('../../hooks/useRole', () => ({
  useIsInstanceAdmin: () => isInstanceAdmin,
}));

function Location() {
  const location = useLocation();
  return <output>{location.pathname}</output>;
}

function renderAuthGuard() {
  return render(
    <MemoryRouter initialEntries={['/docs/guide']}>
      <Routes>
        <Route path="/docs/guide" element={<RequireAuth><p>protected</p></RequireAuth>} />
        <Route path="/login" element={<Location />} />
      </Routes>
    </MemoryRouter>,
  );
}

function renderRoleGuard() {
  return render(
    <MemoryRouter initialEntries={['/settings/users']}>
      <Routes>
        <Route path="/settings/users" element={<RequireInstanceAdmin><p>operator controls</p></RequireInstanceAdmin>} />
        <Route path="/forbidden" element={<Location />} />
      </Routes>
    </MemoryRouter>,
  );
}

describe('route guards', () => {
  afterEach(cleanup);

  beforeEach(() => {
    isInstanceAdmin = false;
    useAuth.setState({ status: 'unknown', user: null });
  });

  it('keeps the splash up until the session resolves, then sends anonymous users to login', () => {
    const { rerender } = renderAuthGuard();

    expect(screen.queryByText('protected')).not.toBeInTheDocument();
    expect(screen.queryByText('/login')).not.toBeInTheDocument();

    useAuth.setState({ status: 'anonymous', user: null });
    rerender(
      <MemoryRouter initialEntries={['/docs/guide']}>
        <Routes>
          <Route path="/docs/guide" element={<RequireAuth><p>protected</p></RequireAuth>} />
          <Route path="/login" element={<Location />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByText('/login')).toBeInTheDocument();
  });

  it('blocks non-admins from instance-admin routes while allowing admins through', () => {
    isInstanceAdmin = false;
    const { unmount } = renderRoleGuard();
    expect(screen.getByText('/forbidden')).toBeInTheDocument();

    unmount();
    isInstanceAdmin = true;
    renderRoleGuard();

    expect(screen.getByText('operator controls')).toBeInTheDocument();
  });
});
