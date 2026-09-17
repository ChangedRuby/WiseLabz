/**
 * Settings hub. Replaces the old single SettingsPage: a section nav down the left
 * (operator-only sections hidden for viewers) and the active sub-page in an
 * <Outlet/>. Desktop-first; the nav collapses to a horizontal scroll strip on
 * narrow widths. Each sub-page reads/writes its own mock.
 */
import { NavLink, Outlet } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useCanMutate, useIsInstanceAdmin } from '../../hooks/useRole';
import { cn } from '../../lib/cn';
import { SETTINGS_SECTIONS } from './nav';

export function SettingsLayout() {
  const { t } = useTranslation();
  const canMutate = useCanMutate();
  const isInstanceAdmin = useIsInstanceAdmin();
  const sections = SETTINGS_SECTIONS.filter((s) =>
    (!s.operatorOnly || canMutate) && (!s.instanceAdminOnly || isInstanceAdmin)
  );

  return (
    <div className="mx-auto max-w-250 px-6 py-6">
      <div className="flex flex-col gap-6 md:flex-row md:gap-8">
        <nav
          aria-label={t('settings.nav.label')}
          className="flex shrink-0 gap-1 overflow-x-auto md:w-52 md:flex-col md:overflow-visible"
        >
          {sections.map(({ segment, labelKey, labelDefault, Icon }) => (
            <NavLink
              key={segment}
              to={segment}
              className={({ isActive }) =>
                cn(
                  'group flex items-center gap-2.5 whitespace-nowrap rounded-md px-3 py-2 text-sm transition-colors',
                  isActive
                    ? 'bg-accent-primary-tint text-ink'
                    : 'text-ink-muted hover:bg-surface hover:text-ink'
                )
              }
            >
              {({ isActive }) => (
                <>
                  <span
                    className={cn(
                      'transition-colors',
                      isActive ? 'text-accent-primary-bright' : 'text-ink-faint group-hover:text-ink-muted'
                    )}
                  >
                    <Icon size={16} />
                  </span>
                  {t(labelKey, { defaultValue: labelDefault })}
                </>
              )}
            </NavLink>
          ))}
        </nav>

        <div className="min-w-0 flex-1">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
