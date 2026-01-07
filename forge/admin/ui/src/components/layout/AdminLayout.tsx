import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from '@tanstack/react-router';
import { useModels } from '../../api/hooks/adminHooks';
import { Button } from '../ui/button';
import { LayoutDashboard, LogOut, Menu, Database, Bell } from 'lucide-react';
import { cn } from '../../lib/utils';
import { TabHeader } from './TabHeader';
import { useTabStore } from '../../store/tabStore';
import { clearAuth } from '../../lib/auth';
import { getAdminConfig } from '../../lib/config';
import { useAdminMetaStore } from '../../store/adminStore';
import { componentRegistry, CORE_COMPONENTS } from '../../lib/registry';

import { ThemeToggle } from './ThemeToggle';
import { GlobalSearch } from './GlobalSearch';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { data } = useModels();
  const navigate = useNavigate();
  const location = useLocation();
  const { basePath } = getAdminConfig();
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const { addTab } = useTabStore();
  const { setMetadata, customPages, menuEntries, plugins } = useAdminMetaStore();
  const SidebarTopSlot = componentRegistry.get(CORE_COMPONENTS.SLOTS.SIDEBAR_TOP);
  const SidebarBottomSlot = componentRegistry.get(CORE_COMPONENTS.SLOTS.SIDEBAR_BOTTOM);
  const TopbarLeftSlot = componentRegistry.get(CORE_COMPONENTS.SLOTS.TOPBAR_LEFT);
  const TopbarRightSlot = componentRegistry.get(CORE_COMPONENTS.SLOTS.TOPBAR_RIGHT);
  const FooterSlot = componentRegistry.get(CORE_COMPONENTS.SLOTS.FOOTER);
  const currentPath = location.pathname;
  const relativePath =
    basePath && currentPath.startsWith(basePath)
      ? currentPath.slice(basePath.length) || '/'
      : currentPath;

  const handleLogout = () => {
    clearAuth();
    navigate({ to: '/login' });
  };

  // Sync tab system with browser navigation
  useEffect(() => {
    // Auto-detect title from path
    let title = 'Dashboard';
    let id = 'dashboard';
    
    if (relativePath.startsWith('/')) {
      const parts = relativePath.split('/');
      id = parts.join('-');
      if (parts[2]) {
        title = parts[2].charAt(0).toUpperCase() + parts[2].slice(1);
        if (parts[3]) {
          title = `${title} - ${parts[3] === 'new' ? 'Create' : 'Edit'}`;
        }
      }
    } else if (relativePath === '/login') {
      return;
    }

    addTab({ id, title, path: relativePath });
  }, [location.pathname, addTab, basePath]);

  const models = data?.models || [];

  useEffect(() => {
    if (data) {
      setMetadata(data);
    }
  }, [data, setMetadata]);

  return (
    <div className="min-h-screen bg-background flex">
      {/* Sidebar */}
      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-50 w-64 bg-card border-r transition-transform duration-300 ease-in-out lg:relative lg:translate-x-0 premium-shadow",
          !sidebarOpen && "-translate-x-full lg:hidden"
        )}
      >
        <div className="h-16 flex items-center px-6 border-b">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center text-white font-bold">F</div>
            <span className="text-xl font-bold tracking-tight">Forge Admin</span>
          </div>
        </div>
        <div className="p-4 space-y-2 overflow-y-auto max-h-[calc(100vh-140px)]">
          {SidebarTopSlot && (
            <div className="pb-4">
              <SidebarTopSlot />
            </div>
          )}
          <Link
            to="/"
            data-testid="nav-dashboard"
            className={cn(
              "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground group",
              relativePath === '/' && "bg-primary text-primary-foreground hover:bg-primary/90"
            )}
          >
            <LayoutDashboard className="h-4 w-4" />
            <span className="font-medium text-sm">Dashboard</span>
          </Link>
          
          <div className="pt-4 pb-2">
            <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-[0.2em] px-3">
              Models
            </span>
          </div>
          
          {models.map((model: any) => (
            <div
              key={model.name}
              data-testid={`nav-${model.name}`}
              onClick={() => {
                const path = `/${model.name}`;
                addTab({ id: `admin-${model.name}`, title: model.verbose_name_plural, path });
                navigate({ to: path });
              }}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground cursor-pointer group",
                relativePath.startsWith(`/${model.name}`) && "bg-accent text-accent-foreground"
              )}
            >
              <Database className="h-4 w-4 text-muted-foreground group-hover:text-foreground" />
              <span className="text-sm font-medium">{model.verbose_name_plural}</span>
            </div>
          ))}

          {customPages.length > 0 && (
            <>
              <div className="pt-4 pb-2">
                <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-[0.2em] px-3">
                  Pages
                </span>
              </div>
              {customPages.map((page: any) => (
                <div
                  key={page.name}
                  onClick={() => {
                    const path = page.path || `/pages/${page.name}`;
                    addTab({ id: `page-${page.name}`, title: page.label, path });
                    navigate({ to: path });
                  }}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground cursor-pointer text-sm group",
                    relativePath === (page.path || `/pages/${page.name}`) && "bg-accent text-accent-foreground"
                  )}
                >
                  <div className="h-4 w-4 flex items-center justify-center">
                    <div className="w-1.5 h-1.5 rounded-full bg-primary" />
                  </div>
                  <span className="font-medium">{page.label}</span>
                </div>
              ))}
            </>
          )}

          {menuEntries.length > 0 && (
            <>
              <div className="pt-4 pb-2">
                <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-[0.2em] px-3">
                  Menu
                </span>
              </div>
              {menuEntries.map((entry: any, idx: number) => (
                <div
                  key={`menu-${idx}`}
                  onClick={() => {
                    const path = entry.path || '/';
                    addTab({ id: `menu-${idx}`, title: entry.label, path });
                    navigate({ to: path });
                  }}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground cursor-pointer text-sm group",
                    relativePath === entry.path && "bg-accent text-accent-foreground"
                  )}
                >
                  <div className="h-4 w-4 flex items-center justify-center">
                    <div className="w-1.5 h-1.5 rounded-full bg-primary" />
                  </div>
                  <span className="font-medium">{entry.label}</span>
                </div>
              ))}
            </>
          )}

          {plugins.map((plugin: any) => (
            <React.Fragment key={plugin.name}>
              {plugin.pages?.map((page: any, idx: number) => (
                <div
                  key={`${plugin.name}-page-${idx}`}
                  onClick={() => {
                    const path = page.path || `/pages/${page.name}`;
                    addTab({ id: `plugin-page-${plugin.name}-${idx}`, title: page.label, path });
                    navigate({ to: path });
                  }}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground cursor-pointer text-sm group",
                    relativePath === (page.path || `/pages/${page.name}`) && "bg-accent text-accent-foreground"
                  )}
                >
                  <div className="h-4 w-4 flex items-center justify-center">
                    <div className="w-1.5 h-1.5 rounded-full bg-primary" />
                  </div>
                  <span className="font-medium">{page.label}</span>
                </div>
              ))}
              {plugin.menuEntries?.map((entry: any, idx: number) => (
                <div
                  key={`${plugin.name}-menu-${idx}`}
                  onClick={() => {
                    const path = entry.path || '/';
                    addTab({ id: `plugin-${plugin.name}-${idx}`, title: entry.label, path });
                    navigate({ to: path });
                  }}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-all hover:bg-accent hover:text-accent-foreground cursor-pointer text-sm group",    
                    relativePath === entry.path && "bg-accent text-accent-foreground"
                  )}
                >
                  <div className="h-4 w-4 flex items-center justify-center">
                    <div className="w-1.5 h-1.5 rounded-full bg-primary" />
                  </div>
                  <span className="font-medium">{entry.label}</span>
                </div>
                  ))}
            </React.Fragment>
          ))}

          {SidebarBottomSlot && (
            <div className="pt-4">
              <SidebarBottomSlot />
            </div>
          )}
        </div>
        
        <div className="absolute bottom-4 left-4 right-4">
          <Button variant="ghost" className="w-full justify-start text-muted-foreground hover:text-destructive hover:bg-destructive/10" onClick={handleLogout}>
            <LogOut className="mr-3 h-4 w-4" />
            <span className="text-sm font-medium">Logout</span>
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col min-h-screen overflow-hidden">
        <header className="h-16 border-b flex items-center px-6 bg-card shrink-0 gap-4">
          <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setSidebarOpen(!sidebarOpen)}>
            <Menu className="h-5 w-5" />
          </Button>

          {TopbarLeftSlot ? <TopbarLeftSlot /> : <GlobalSearch />}
          <div className="ml-auto flex items-center gap-2">
            {TopbarRightSlot && <TopbarRightSlot />}
            <ThemeToggle />
            <Button variant="ghost" size="icon" className="rounded-full relative">
              <Bell className="h-4 w-4" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-destructive rounded-full border-2 border-card" />
            </Button>
            <div className="h-6 w-[1px] bg-border mx-2" />
            <div className="flex items-center gap-3 pl-2">
              <div className="text-right hidden sm:block">
                <p className="text-sm font-semibold leading-none">Admin User</p>
                <p className="text-[10px] text-muted-foreground mt-1 uppercase font-bold tracking-wider">Super Admin</p>
              </div>
              <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-primary to-indigo-400 border border-border shadow-sm flex items-center justify-center text-white text-xs font-bold">
                AU
              </div>
            </div>
          </div>
        </header>

        {/* Tab System */}
        <div className="bg-card/50 backdrop-blur-sm">
          <TabHeader />
        </div>
        <div className="flex-1 p-6 overflow-auto">
          {children}
        </div>
        {FooterSlot && (
          <div className="border-t border-border/50 bg-card/30 px-6 py-4">
            <FooterSlot />
          </div>
        )}
      </main>
    </div>
  );
}
