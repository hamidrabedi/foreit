import { useModelMetadata } from '../api/hooks/adminHooks';
import ModelListPage from '../pages/ModelListPage';
import ModelUpsertPage from '../pages/ModelUpsertPage';
import { useParams } from '@tanstack/react-router';
import { componentRegistry, CORE_COMPONENTS } from '../lib/registry';

interface DynamicModelPageProps {
  mode: 'list' | 'create' | 'edit' | 'detail';
}

export default function DynamicModelPage({ mode }: DynamicModelPageProps) {
  const { model } = useParams({ strict: false }) as any;
  const { data: metadata, isLoading, error } = useModelMetadata(model);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (error || !metadata) {
    return (
      <div className="p-8 text-center text-red-500">
        Error loading metadata for "{model}"
      </div>
    );
  }

  const modelKey = String(model || '');
  const listOverride =
    componentRegistry.get(`forge.pages.model.${modelKey}.list`) ||
    componentRegistry.get(CORE_COMPONENTS.PAGES.LIST);
  const formOverride =
    componentRegistry.get(`forge.pages.model.${modelKey}.form`) ||
    componentRegistry.get(CORE_COMPONENTS.PAGES.FORM);

  if (mode === 'list') {
    const ListPage = (listOverride as any) || ModelListPage;
    return <ListPage mode={mode} model={modelKey} metadata={metadata} />;
  }

  if (mode === 'create' || mode === 'edit') {
    const FormPage = (formOverride as any) || ModelUpsertPage;
    return <FormPage mode={mode as any} model={modelKey} metadata={metadata} />;
  }

  // Detail view (TODO: implement ModelViewPage)
  return <ModelListPage />; 
}
