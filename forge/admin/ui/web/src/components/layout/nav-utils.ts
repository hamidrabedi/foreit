export const normalizeAdminPath = (path?: string) =>
  path?.startsWith("/admin/") ? path.substring(6) : path;

export const isEntryMatch = (entry: any, pathname: string): boolean => {
  const normalizedPath = normalizeAdminPath(entry.path);
  return Boolean(normalizedPath && normalizedPath === pathname);
};

export const isMenuEntryActive = (entry: any, pathname: string): boolean => {
  if (isEntryMatch(entry, pathname)) return true;
  return entry.children?.some((child: any) => isMenuEntryActive(child, pathname)) ?? false;
};

export const findActiveMenuEntry = (entries: any[], pathname: string): any | null => {
  for (const entry of entries) {
    if (isEntryMatch(entry, pathname)) return entry;
    if (entry.children?.length) {
      const nested = findActiveMenuEntry(entry.children, pathname);
      if (nested) return nested;
    }
  }
  return null;
};
