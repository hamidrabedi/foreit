import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminAPI } from '../client';
import type {
  Metadata,
  PaginatedResponse,
  ListParams,
  ModelFormData,
  BulkActionRequest,
  BulkActionResponse,
  SearchRequest,
  SearchResponse,
  AutocompleteResponse,
  MetadataResponse,
} from '../types';

// Query keys factory
export const adminKeys = {
  all: ['admin'] as const,
  models: () => [...adminKeys.all, 'models'] as const,
  model: (model: string) => [...adminKeys.all, 'model', model] as const,
  modelMeta: (model: string) => [...adminKeys.model(model), 'meta'] as const,
  modelHistory: (model: string, id: string | number) =>
    [...adminKeys.model(model), 'history', id] as const,
  modelList: (model: string, params?: ListParams) =>
    [...adminKeys.model(model), 'list', params] as const,
  modelDetail: (model: string, id: string | number) =>
    [...adminKeys.model(model), 'detail', id] as const,
};

// Models and Plugins metadata hook
export function useMetadata(options?: any) {
  return useQuery({
    queryKey: adminKeys.models(),
    queryFn: () => adminAPI.getMetadata(),
    ...options,
  });
}

// Models list hook (backward compatibility)
export function useModels(options?: any) {
  return useMetadata(options);
}

// Model metadata hook
export function useModelMetadata(
  model: string,
  options?: any
) {
  return useQuery({
    queryKey: adminKeys.modelMeta(model),
    queryFn: () => adminAPI.getModelMetadata(model),
    enabled: !!model,
    staleTime: 5 * 60 * 1000, // 5 minutes
    ...options,
  });
}

// Model list hook
export function useModelList<T = any>(
  model: string,
  params?: ListParams,
  options?: any
) {
  return useQuery({
    queryKey: adminKeys.modelList(model, params),
    queryFn: () => adminAPI.listObjects<T>(model, params),
    enabled: !!model,
    ...options,
  });
}

// Model detail hook
export function useModelDetail<T = any>(
  model: string,
  id: string | number,
  options?: any
) {
  return useQuery({
    queryKey: adminKeys.modelDetail(model, id),
    queryFn: () => adminAPI.getObject<T>(model, id),
    enabled: !!model && !!id,
    ...options,
  });
}

export function useModelHistory(
  model: string,
  id: string | number,
  options?: any
) {
  return useQuery({
    queryKey: adminKeys.modelHistory(model, id),
    queryFn: () => adminAPI.getHistory(model, id),
    enabled: !!model && !!id,
    ...options,
  });
}

// Create object mutation
export function useCreateObject<T = any>(
  model: string,
  options?: any
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: ModelFormData) => adminAPI.createObject<T>(model, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.model(model) });
    },
    ...options,
  });
}

// Update object mutation
export function useUpdateObject<T = any>(
  model: string,
  options?: any
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string | number; data: Partial<ModelFormData> }) =>
      adminAPI.updateObject<T>(model, id, data),
    onSuccess: (_: T, variables: { id: string | number }) => {
      queryClient.invalidateQueries({ queryKey: adminKeys.modelDetail(model, variables.id) });
      queryClient.invalidateQueries({ queryKey: adminKeys.model(model) });
    },
    ...options,
  });
}

// Delete object mutation
export function useDeleteObject(
  model: string,
  options?: any
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string | number) => adminAPI.deleteObject(model, id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.model(model) });
    },
    ...options,
  });
}

// Bulk delete mutation
export function useBulkDelete(
  model: string,
  options?: any
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ids: (string | number)[]) => adminAPI.bulkDelete(model, ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.model(model) });
    },
    ...options,
  });
}

// Execute action mutation
export function useExecuteAction(
  model: string,
  action: string,
  options?: any
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: BulkActionRequest) =>
      adminAPI.executeAction(model, action, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: adminKeys.model(model) });
    },
    ...options,
  });
}

// Global search hook
export function useGlobalSearch(
  request: SearchRequest,
  options?: any
) {
  return useQuery({
    queryKey: [...adminKeys.all, 'search', request],
    queryFn: () => adminAPI.globalSearch(request),
    enabled: !!request.query && request.query.length > 0,
    ...options,
  });
}

// Autocomplete hook
export function useAutocomplete(
  model: string,
  field: string,
  query: string,
  options?: any
) {
  return useQuery({
    queryKey: [...adminKeys.model(model), 'autocomplete', field, query],
    queryFn: () => adminAPI.autocomplete(model, field, query),
    enabled: !!model && !!field && query.length > 0,
    staleTime: 30 * 1000, // 30 seconds
    ...options,
  });
}
