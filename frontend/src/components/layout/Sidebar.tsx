import { Link, useLocation } from 'react-router-dom';
import {
  LayoutDashboard,
  Activity,
  Network,
  Cable,
  Gauge,
  Settings2,
  Database,
  Settings,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useSidebarStore } from '@/store/sidebarStore';
import { ThemeToggle } from '@/components/common/ThemeToggle';
import { Button } from '@/components/ui/button';
import { APP_NAME } from '@/utils/constants';

// Menu configuration - v1.0.1
const menuItems = [
  { icon: LayoutDashboard, label: 'Dashboard', path: '/' },
  { icon: Activity, label: 'Monitoring', path: '/monitoring' },
  { icon: Network, label: 'Provisioning', path: '/provisioning' },
  { icon: Cable, label: 'VLAN Management', path: '/vlan' },
  { icon: Gauge, label: 'Traffic Control', path: '/traffic' },
  { icon: Settings2, label: 'ONU Management', path: '/onu-management' },
  { icon: Database, label: 'Config Backup', path: '/config-backup' },
  { icon: Settings, label: 'Settings', path: '/settings' },
];

export const Sidebar = () => {
  const location = useLocation();
  const { isCollapsed, toggle } = useSidebarStore();

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 z-40 h-screen border-r bg-card transition-all duration-300',
        isCollapsed ? 'w-16' : 'w-64'
      )}
    >
      <div className="flex h-full flex-col">
        {/* Logo */}
        <div className="flex h-16 items-center justify-between border-b px-4">
          {!isCollapsed && (
            <h1 className="text-lg font-bold text-primary">
              {APP_NAME.replace(' Dashboard', '')}
            </h1>
          )}
          <Button
            variant="ghost"
            size="icon"
            onClick={toggle}
            className={cn('ml-auto', isCollapsed && 'mx-auto')}
          >
            {isCollapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </Button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 space-y-1 overflow-y-auto p-2">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;

            return (
              <Link
                key={item.path}
                to={item.path}
                className={cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors',
                  isActive
                    ? 'bg-primary text-primary-foreground'
                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                  isCollapsed && 'justify-center'
                )}
                title={isCollapsed ? item.label : undefined}
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                {!isCollapsed && <span>{item.label}</span>}
              </Link>
            );
          })}
        </nav>

        {/* Footer */}
        <div className="border-t p-4">
          <div className={cn('flex items-center', isCollapsed && 'justify-center')}>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </aside>
  );
};
