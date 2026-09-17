import { lazy, Suspense, useEffect } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { QueryClientProvider } from '@tanstack/react-query';
import { Link, RouterProvider, createBrowserRouter, Navigate } from 'react-router-dom';
import { queryClient } from './app/queryClient';
import { MotionProvider } from './components/MotionProvider';
import { WebSocketProvider } from './ws/WebSocketProvider';
import { useAuth } from './store/auth';
import { PublicLayout } from './features/auth/PublicLayout';
import { LoginPage } from './features/auth/LoginPage';
import {
  RequireAuth,
  RequireInstanceAdmin,
  RequireOnboarded,
  ForbiddenPage,
  Splash,
} from './features/auth/guards';
import { EmptyState } from './components/ui/states';

const AppShell = lazy(() =>
  import('./components/shell/AppShell').then(({ AppShell }) => ({ default: AppShell }))
);
const AuthCallbackPage = lazy(() =>
  import('./features/auth/AuthCallbackPage').then(({ AuthCallbackPage }) => ({
    default: AuthCallbackPage,
  }))
);
const OnboardingPage = lazy(() =>
  import('./features/onboarding/OnboardingPage').then(({ OnboardingPage }) => ({
    default: OnboardingPage,
  }))
);
const DashboardPage = lazy(() =>
  import('./features/dashboard/DashboardPage').then(({ DashboardPage }) => ({
    default: DashboardPage,
  }))
);
const ServicesPage = lazy(() =>
  import('./features/services/ServicesPage').then(({ ServicesPage }) => ({ default: ServicesPage }))
);
const ServiceDetailPage = lazy(() =>
  import('./features/services/ServiceDetailPage').then(({ ServiceDetailPage }) => ({
    default: ServiceDetailPage,
  }))
);
const AddConnectorPage = lazy(() =>
  import('./features/connectors/AddConnectorPage').then(({ AddConnectorPage }) => ({
    default: AddConnectorPage,
  }))
);
const ConnectorEditPage = lazy(() =>
  import('./features/connectors/ConnectorEditPage').then(({ ConnectorEditPage }) => ({
    default: ConnectorEditPage,
  }))
);
const DocsPage = lazy(() =>
  import('./features/docs/DocsPage').then(({ DocsPage }) => ({ default: DocsPage }))
);
const AllDocsPage = lazy(() =>
  import('./features/docs/AllDocsPage').then(({ AllDocsPage }) => ({ default: AllDocsPage }))
);
const TopologyPage = lazy(() =>
  import('./features/docs/TopologyPage').then(({ TopologyPage }) => ({ default: TopologyPage }))
);
const DocEditorPage = lazy(() =>
  import('./features/docs/DocEditorPage').then(({ DocEditorPage }) => ({ default: DocEditorPage }))
);
const ChatPage = lazy(() =>
  import('./features/chat/ChatPage').then(({ ChatPage }) => ({ default: ChatPage }))
);
const ChangesPage = lazy(() =>
  import('./features/changes/ChangesPage').then(({ ChangesPage }) => ({ default: ChangesPage }))
);
const ChangeDetailPage = lazy(() =>
  import('./features/changes/ChangeDetailPage').then(({ ChangeDetailPage }) => ({
    default: ChangeDetailPage,
  }))
);
const AttentionPage = lazy(() =>
  import('./features/attention/AttentionPage').then(({ AttentionPage }) => ({ default: AttentionPage }))
);
const AlertsPage = lazy(() =>
  import('./features/alerts/AlertsPage').then(({ AlertsPage }) => ({ default: AlertsPage }))
);
const FindingsPage = lazy(() =>
  import('./features/findings/FindingsPage').then(({ FindingsPage }) => ({ default: FindingsPage }))
);
const TemplatesPage = lazy(() =>
  import('./features/templates/TemplatesPage').then(({ TemplatesPage }) => ({
    default: TemplatesPage,
  }))
);
const TemplateEditorPage = lazy(() =>
  import('./features/templates/TemplateEditorPage').then(({ TemplateEditorPage }) => ({
    default: TemplateEditorPage,
  }))
);
const SettingsLayout = lazy(() =>
  import('./features/settings').then(({ SettingsLayout }) => ({ default: SettingsLayout }))
);
const ProfilePage = lazy(() =>
  import('./features/settings').then(({ ProfilePage }) => ({ default: ProfilePage }))
);
const UsersPage = lazy(() =>
  import('./features/settings').then(({ UsersPage }) => ({ default: UsersPage }))
);
const AuthPage = lazy(() =>
  import('./features/settings').then(({ AuthPage }) => ({ default: AuthPage }))
);
const AiPage = lazy(() =>
  import('./features/settings').then(({ AiPage }) => ({ default: AiPage }))
);
const NotificationsPage = lazy(() =>
  import('./features/settings').then(({ NotificationsPage }) => ({ default: NotificationsPage }))
);
const SystemPage = lazy(() =>
  import('./features/settings').then(({ SystemPage }) => ({ default: SystemPage }))
);
const RetentionPage = lazy(() =>
  import('./features/settings').then(({ RetentionPage }) => ({ default: RetentionPage }))
);
const RunbooksPage = lazy(() =>
  import('./features/settings').then(({ RunbooksPage }) => ({ default: RunbooksPage }))
);
const AuditPage = lazy(() =>
  import('./features/settings').then(({ AuditPage }) => ({ default: AuditPage }))
);
const AppearancePage = lazy(() =>
  import('./features/settings').then(({ AppearancePage }) => ({ default: AppearancePage }))
);

function NotFound() {
  return (
    <div className="flex h-full items-center justify-center">
      <EmptyState
        title="Page not found"
        description="That route doesn't exist."
        action={
          <Link
            to="/dashboard"
            className="inline-flex h-7 items-center rounded-sm border border-line-strong px-2.5 font-mono text-xs text-ink hover:bg-surface"
          >
            Back to dashboard
          </Link>
        }
      />
    </div>
  );
}

const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <PublicLayout>
        <LoginPage />
      </PublicLayout>
    ),
  },
  { path: '/auth/callback', element: <AuthCallbackPage /> },
  { path: '/forbidden', element: <ForbiddenPage /> },
  {
    // Onboarding sits under auth only (NOT RequireOnboarded) so it never loops.
    path: '/onboarding',
    element: (
      <RequireAuth>
        <OnboardingPage />
      </RequireAuth>
    ),
  },
  {
    path: '/',
    element: (
      <RequireAuth>
        <RequireOnboarded>
          <AppShell />
        </RequireOnboarded>
      </RequireAuth>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <DashboardPage /> },
      { path: 'services', element: <ServicesPage /> },
      { path: 'services/new', element: <AddConnectorPage /> },
      { path: 'services/:id', element: <ServiceDetailPage /> },
      { path: 'connectors/:id/edit', element: <ConnectorEditPage /> },
      { path: 'docs', element: <DocsPage /> },
      { path: 'docs/all', element: <AllDocsPage /> },
      { path: 'topology', element: <TopologyPage /> },
      { path: 'docs/:docId', element: <DocsPage /> },
      { path: 'docs/:docId/edit', element: <DocEditorPage /> },
      { path: 'docs/:docId/history', element: <DocsPage /> },
      { path: 'chat', element: <ChatPage /> },
      { path: 'changes', element: <ChangesPage /> },
      { path: 'changes/:changeId', element: <ChangeDetailPage /> },
      { path: 'attention', element: <AttentionPage /> },
      { path: 'alerts', element: <AlertsPage /> },
      { path: 'findings', element: <FindingsPage /> },
      {
        path: 'templates',
        element: (
          <RequireInstanceAdmin>
            <TemplatesPage />
          </RequireInstanceAdmin>
        ),
      },
      {
        path: 'templates/:id',
        element: (
          <RequireInstanceAdmin>
            <TemplateEditorPage />
          </RequireInstanceAdmin>
        ),
      },
      {
        path: 'settings',
        element: <SettingsLayout />,
        children: [
          { index: true, element: <Navigate to="/settings/profile" replace /> },
          { path: 'profile', element: <ProfilePage /> },
          {
            path: 'users',
            element: (
              <RequireInstanceAdmin>
                <UsersPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'auth',
            element: (
              <RequireInstanceAdmin>
                <AuthPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'ai',
            element: (
              <RequireInstanceAdmin>
                <AiPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'notifications',
            element: (
              <RequireInstanceAdmin>
                <NotificationsPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'system',
            element: (
              <RequireInstanceAdmin>
                <SystemPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'runbooks',
            element: (
              <RequireInstanceAdmin>
                <RunbooksPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'retention',
            element: (
              <RequireInstanceAdmin>
                <RetentionPage />
              </RequireInstanceAdmin>
            ),
          },
          {
            path: 'audit',
            element: (
              <RequireInstanceAdmin>
                <AuditPage />
              </RequireInstanceAdmin>
            ),
          },
          { path: 'appearance', element: <AppearancePage /> },
        ],
      },
      { path: '*', element: <NotFound /> },
    ],
  },
]);

function Root() {
  const status = useAuth((s) => s.status);
  const bootstrap = useAuth((s) => s.bootstrap);
  useEffect(() => {
    void bootstrap();
  }, [bootstrap]);
  if (status === 'unknown') return <Splash />;
  return (
    <ErrorBoundary
      fallback={
        <div role="alert" className="p-6 text-sm text-ink">
          Unable to load this page.
        </div>
      }
    >
      <Suspense
        fallback={
          <div role="status" aria-live="polite" className="p-6 text-sm text-ink-muted">
            Loading page…
          </div>
        }
      >
        <RouterProvider router={router} />
      </Suspense>
    </ErrorBoundary>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <WebSocketProvider>
        <MotionProvider>
          <Root />
        </MotionProvider>
      </WebSocketProvider>
    </QueryClientProvider>
  );
}
