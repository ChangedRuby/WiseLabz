/**
 * Authenticated app frame. Navigation is the locked bottom dock (PRODUCT.md);
 * the routed content (error boundary + Suspense + the active page) and the
 * command palette (⌘K) live here.
 */
import { Suspense, useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { ErrorBoundary } from 'react-error-boundary';
import { Toaster } from 'sonner';
import { ShellDock } from './ShellDock';
import { CommandPalette } from '../command/CommandPalette';
import { ShortcutsModal } from '../command/ShortcutsModal';
import { ErrorState, SkeletonRows } from '../ui/states';
import { setAppNavigator } from '../../lib/navigation';

/** Registers react-router's navigate with the navigator singleton so non-React
 *  code (the WS handler) can drive SPA navigation, e.g. click-to-jump toasts. */
function NavigatorBridge() {
  const navigate = useNavigate();
  useEffect(() => {
    setAppNavigator((to) => navigate(to));
    return () => setAppNavigator(null);
  }, [navigate]);
  return null;
}

function Content() {
  const location = useLocation();
  return (
    <ErrorBoundary
      // Reset whenever the route changes — otherwise a crash on one page
      // keeps showing this fallback for every page navigated to afterward,
      // since Outlet just swaps children under the same mounted boundary.
      resetKeys={[location.pathname]}
      fallbackRender={({ resetErrorBoundary }) => (
        <ErrorState
          title="This page hit an error"
          description="The view failed to render. Retry, or head back to the dashboard."
          onRetry={resetErrorBoundary}
        />
      )}
    >
      <Suspense fallback={<SkeletonRows rows={6} className="m-6 max-w-2xl" />}>
        <Outlet />
      </Suspense>
    </ErrorBoundary>
  );
}

export function AppShell() {
  return (
    <>
      <NavigatorBridge />
      <ShellDock>
        <Content />
      </ShellDock>
      <CommandPalette />
      <ShortcutsModal />
      <Toaster
        position="top-right"
        gap={8}
        toastOptions={{
          style: {
            background: 'var(--color-surface-raised)',
            color: 'var(--color-ink)',
            border: '1px solid var(--color-line)',
            borderRadius: 'var(--radius-md)',
            boxShadow: 'var(--shadow-pop)',
            fontFamily: 'var(--font-sans)',
            fontSize: 'var(--text-sm)',
          },
        }}
      />
    </>
  );
}
