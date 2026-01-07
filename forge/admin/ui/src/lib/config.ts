export interface AdminUIConfig {
  basePath: string;
  apiBase: string;
  overridesUrl: string;
  uploadUrl: string;
}

const DEFAULT_BASE_PATH = "";

function normalizeBasePath(input: string): string {
  if (!input || input === "/") {
    return "";
  }

  let path = input.trim();
  if (!path.startsWith("/")) {
    path = `/${path}`;
  }
  if (path.length > 1 && path.endsWith("/")) {
    path = path.slice(0, -1);
  }
  return path === "/" ? "" : path;
}

function buildApiBase(basePath: string): string {
  const prefix = basePath || "";
  return `${prefix}/api`;
}

function inferBasePath(): string {
  if (typeof window === "undefined") {
    return DEFAULT_BASE_PATH;
  }
  if (window.location.pathname.startsWith("/admin")) {
    return "/admin";
  }
  return DEFAULT_BASE_PATH;
}

export function getAdminConfig(): AdminUIConfig {
  const raw =
    typeof window === "undefined"
      ? undefined
      : ((window as any).__FORGE_ADMIN__ as Partial<AdminUIConfig> | undefined);
  const basePath = normalizeBasePath(raw?.basePath ?? inferBasePath());

  return {
    basePath,
    apiBase: raw?.apiBase ?? buildApiBase(basePath),
    overridesUrl: raw?.overridesUrl ?? `${basePath || ""}/overrides.js`,
    uploadUrl: raw?.uploadUrl ?? `${basePath || ""}/uploads`,
  };
}
