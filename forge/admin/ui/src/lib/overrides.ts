import type React from "react";
import { componentRegistry } from "./registry";
import { widgetRegistry } from "./widgets";
import { getAdminConfig } from "./config";

export interface ForgeAdminAPI {
  registerComponent: (name: string, component: React.ComponentType<any>) => void;
  registerComponents: (components: Record<string, React.ComponentType<any>>) => void;
  registerWidget: (name: string, component: React.ComponentType<any>) => void;
  registerWidgets: (components: Record<string, React.ComponentType<any>>) => void;
  registry: {
    components: typeof componentRegistry;
    widgets: typeof widgetRegistry;
  };
}

declare global {
  interface Window {
    ForgeAdmin?: ForgeAdminAPI;
  }
}

export function exposeForgeAdmin() {
  if (typeof window === "undefined") {
    return;
  }
  if (window.ForgeAdmin) {
    return;
  }

  window.ForgeAdmin = {
    registerComponent: (name, component) => componentRegistry.register(name, component),
    registerComponents: (components) => componentRegistry.registerAll(components),
    registerWidget: (name, component) => widgetRegistry.register(name, component),
    registerWidgets: (components) => widgetRegistry.registerAll(components),
    registry: {
      components: componentRegistry,
      widgets: widgetRegistry,
    },
  };
}

export function loadOverrides(): Promise<void> {
  if (typeof window === "undefined") {
    return Promise.resolve();
  }
  const { overridesUrl } = getAdminConfig();
  if (!overridesUrl) {
    return Promise.resolve();
  }

  return new Promise((resolve) => {
    const script = document.createElement("script");
    script.src = overridesUrl;
    script.async = true;
    script.onload = () => resolve();
    script.onerror = () => resolve();
    document.head.appendChild(script);
  });
}
