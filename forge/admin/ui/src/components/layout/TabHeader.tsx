import { X } from 'lucide-react';
import { useTabStore } from '../../store/tabStore';
import type { Tab } from '../../store/tabStore';
import { useNavigate } from '@tanstack/react-router';
import { cn } from '../../lib/utils';

export function TabHeader() {
  const { tabs, activeTabId, setActiveTab, removeTab } = useTabStore();
  const navigate = useNavigate();

  const handleTabClick = (tab: Tab) => {
    setActiveTab(tab.id);
    navigate({ to: tab.path });
  };

  const handleClose = (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    removeTab(id);
    // After removal, navigate to the new active tab
    const currentActive = useTabStore.getState().activeTabId;
    const currentPath = tabs.find((t: Tab) => t.id === currentActive)?.path || '/';
    navigate({ to: currentPath });
  };

  return (
    <div className="flex bg-card border-b overflow-x-auto no-scrollbar h-10 items-end px-2 gap-1">
      {tabs.map((tab: Tab) => (
        <div
          key={tab.id}
          className={cn(
            "group flex items-center h-8 px-3 rounded-t-lg border-x border-t cursor-pointer transition-colors text-sm whitespace-nowrap gap-2",
            activeTabId === tab.id
              ? "bg-background border-border text-foreground font-medium"
              : "bg-muted/50 border-transparent text-muted-foreground hover:bg-muted"
          )}
          onClick={() => handleTabClick(tab)}
        >
          {tab.title}
          {tab.closable !== false && (
            <X
              className="h-3 w-3 hover:text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
              onClick={(e) => handleClose(e, tab.id)}
            />
          )}
        </div>
      ))}
    </div>
  );
}
