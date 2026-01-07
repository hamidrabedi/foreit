import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export interface Tab {
  id: string;
  title: string;
  path: string;
  closable?: boolean;
}

interface TabState {
  tabs: Tab[];
  activeTabId: string;
  addTab: (tab: Tab) => void;
  removeTab: (id: string) => void;
  setActiveTab: (id: string) => void;
  clearTabs: () => void;
}

export const useTabStore = create<TabState>()(
  persist(
    (set: any) => ({
      tabs: [{ id: 'dashboard', title: 'Dashboard', path: '/', closable: false }],
      activeTabId: 'dashboard',
      addTab: (tab: Tab) =>
        set((state: TabState) => {
          if (state.tabs.find((t) => t.id === tab.id)) {
            return { activeTabId: tab.id };
          }
          return {
            tabs: [...state.tabs, tab],
            activeTabId: tab.id,
          };
        }),
      removeTab: (id: string) =>
        set((state: TabState) => {
          const newTabs = state.tabs.filter((t: Tab) => t.id !== id);
          let newActiveTabId = state.activeTabId;
          if (state.activeTabId === id) {
            newActiveTabId = newTabs[newTabs.length - 1]?.id || 'dashboard';
          }
          return {
            tabs: newTabs,
            activeTabId: newActiveTabId,
          };
        }),
      setActiveTab: (id: string) => set({ activeTabId: id }),
      clearTabs: () =>
        set({
          tabs: [{ id: 'dashboard', title: 'Dashboard', path: '/', closable: false }],
          activeTabId: 'dashboard',
        }),
    }),
    {
      name: 'admin-tabs',
    }
  )
);
