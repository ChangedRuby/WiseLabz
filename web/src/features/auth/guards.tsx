/**
 * Route guards + the auth splash. RequireAuth gates the protected tree on an
 * authenticated session; RequireInstanceAdmin gates instance-wide surfaces
 * (user management, API keys, system settings — #240 PR1); RequireOnboarded
 * routes a fresh install (zero connectors) into onboarding. Connector-scoped
 * pages (a single connector or doc) don't get a route guard at all — the
 * server's default-deny list/get filtering already keeps unauthorized data
 * out, so the page handles "no access" as a normal empty state instead. The
 * server enforces the real boundary — these guards are navigation, not
 * security.
 */
import type { ReactNode } from 'react';
import { Link, Navigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../store/auth';
import { useIsInstanceAdmin } from '../../hooks/useRole';
import { useGetConnectors } from '../../api/generated/connectors/connectors';
import { EmptyState, SkeletonRows } from '../../components/ui/states';

/** Centered brand splash shown while the session resolves. */
export function Splash() {
  return (
    <div className="flex min-h-dvh w-screen items-center justify-center bg-canvas">
      <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-accent-primary shadow-(--shadow-raised)">
        <span className="font-mono text-lg font-bold text-accent-primary-ink">W</span>
      </div>
    </div>
  );
}

export function RequireAuth({ children }: { children: ReactNode }) {
  const status = useAuth((s) => s.status);
  const location = useLocation();
  if (status === 'unknown') return <Splash />;
  if (status === 'anonymous') {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }
  return <>{children}</>;
}

export function RequireInstanceAdmin({ children }: { children: ReactNode }) {
  const isInstanceAdmin = useIsInstanceAdmin();
  if (!isInstanceAdmin) {
    return <Navigate to="/forbidden" replace />;
  }
  return <>{children}</>;
}

export function RequireOnboarded({ children }: { children: ReactNode }) {
  const { data, isLoading } = useGetConnectors();
  if (isLoading) return <SkeletonRows rows={6} className="m-6 max-w-2xl" />;
  if (Array.isArray(data) && data.length === 0) {
    return <Navigate to="/onboarding" replace />;
  }
  return <>{children}</>;
}

export function ForbiddenPage() {
  const { t } = useTranslation();
  return (
    <div className="flex h-screen w-screen items-center justify-center bg-canvas p-6">
      <EmptyState
        title={t('auth.forbiddenTitle')}
        description={t('auth.forbiddenDesc')}
        action={
          <Link
            to="/dashboard"
            className="inline-flex h-7 items-center rounded-sm border border-line-strong px-2.5 font-mono text-xs text-ink hover:bg-surface"
          >
            {t('nav.dashboard')}
          </Link>
        }
      />
    </div>
  );
}
