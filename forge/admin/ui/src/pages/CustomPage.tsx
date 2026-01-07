import { useParams } from "@tanstack/react-router";
import { useAdminMetaStore } from "../store/adminStore";
import { componentRegistry } from "../lib/registry";
import AdminLayout from "../components/layout/AdminLayout";

export default function CustomPage() {
  const { page } = useParams({ strict: false }) as { page?: string };
  const { customPages, plugins } = useAdminMetaStore();

  const allPages = [
    ...customPages,
    ...plugins.flatMap((plugin: any) => plugin.pages || []),
  ];

  const match = allPages.find((entry) => {
    if (!page) return false;
    if (entry.name === page) return true;
    if (entry.path && entry.path.endsWith(`/${page}`)) return true;
    return false;
  });

  if (!match) {
    return (
      <AdminLayout>
        <div className="p-8 text-center text-muted-foreground">
          Unknown page.
        </div>
      </AdminLayout>
    );
  }

  const Component = componentRegistry.get(match.component);
  if (!Component) {
    return (
      <AdminLayout>
        <div className="p-8 text-center text-muted-foreground">
          Missing component: {match.component}
        </div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <Component />
    </AdminLayout>
  );
}
