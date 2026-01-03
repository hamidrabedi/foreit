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
    (set) => ({
      tabs: [{ id: 'dashboard', title: 'Dashboard', path: '/', closable: false }],
      activeTabId: 'dashboard',
      addTab: (tab) =>
        set((state) => {
          if (state.tabs.find((t) => t.id === tab.id)) {
            return { activeTabId: tab.id };
          }
          return {
            tabs: [...state.tabs, tab],
            activeTabId: tab.id,
          };
        }),
      removeTab: (id) =>
        set((state) => {
          const newTabs = state.tabs.filter((t) => t.id !== id);
          let newActiveTabId = state.activeTabId;
          if (state.activeTabId === id) {
            newActiveTabId = newTabs[newTabs.length - 1]?.id || 'dashboard';
          }
          return {
            tabs: newTabs,
            activeTabId: newActiveTabId,
          };
        }),
      setActiveTab: (id) => set({ activeTabId: id }),
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
