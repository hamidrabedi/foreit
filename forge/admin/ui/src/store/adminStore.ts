import { create } from "zustand";
import type { MetadataResponse, CustomPageMetadata, MenuEntryMetadata, DashboardConfig, PluginMetadata } from "../api/types";

interface AdminMetaState {
  metadata?: MetadataResponse;
  customPages: CustomPageMetadata[];
  menuEntries: MenuEntryMetadata[];
  dashboard?: DashboardConfig;
  plugins: PluginMetadata[];
  setMetadata: (metadata: MetadataResponse) => void;
}

export const useAdminMetaStore = create<AdminMetaState>((set: any) => ({
  metadata: undefined,
  customPages: [],
  menuEntries: [],
  dashboard: undefined,
  plugins: [],
  setMetadata: (metadata: MetadataResponse) =>
    set({
      metadata,
      customPages: metadata.custom_pages || [],
      menuEntries: metadata.menu_entries || [],
      dashboard: metadata.dashboard,
      plugins: metadata.plugins || [],
    }),
}));
