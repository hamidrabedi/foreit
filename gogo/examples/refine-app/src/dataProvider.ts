/**
 * Gogo Data Provider for Refine
 * Copy this file to your Refine project
 */

import { DataProvider } from "@refinedev/core";

export interface GogoDataProviderOptions {
  apiUrl: string;
  getAccessToken?: () => string | null;
  headers?: Record<string, string>;
}

export function gogoDataProvider(
  options: GogoDataProviderOptions
): DataProvider {
  const { apiUrl, getAccessToken, headers = {} } = options;

  const fetchWithAuth = async (
    url: string,
    options: RequestInit = {}
  ): Promise<Response> => {
    const token = getAccessToken?.();
    const authHeaders: Record<string, string> = {
      "Content-Type": "application/json",
      ...headers,
    };

    if (token) {
      authHeaders.Authorization = `Bearer ${token}`;
    }

    return fetch(url, {
      ...options,
      headers: {
        ...authHeaders,
        ...options.headers,
      },
    });
  };

  return {
    getList: async ({ resource, pagination, filters, sorters, meta }) => {
      const url = new URL(`${apiUrl}/${resource}`);
      
      const current = pagination?.current || 1;
      const pageSize = pagination?.pageSize || 20;
      url.searchParams.set("page", current.toString());
      url.searchParams.set("page_size", pageSize.toString());

      if (sorters && sorters.length > 0) {
        const sorter = sorters[0];
        url.searchParams.set("sort_by", sorter.field);
        url.searchParams.set("sort_order", sorter.order === "asc" ? "asc" : "desc");
      }

      if (filters && filters.length > 0) {
        filters.forEach((filter) => {
          if ("field" in filter && "operator" in filter && "value" in filter) {
            const field = filter.field;
            const operator = filter.operator;
            const value = filter.value;

            switch (operator) {
              case "eq":
                url.searchParams.set(`filter_${field}`, String(value));
                break;
              case "ne":
                url.searchParams.set(`filter_${field}__ne`, String(value));
                break;
              case "gt":
                url.searchParams.set(`filter_${field}__gt`, String(value));
                break;
              case "gte":
                url.searchParams.set(`filter_${field}__gte`, String(value));
                break;
              case "lt":
                url.searchParams.set(`filter_${field}__lt`, String(value));
                break;
              case "lte":
                url.searchParams.set(`filter_${field}__lte`, String(value));
                break;
              case "contains":
                url.searchParams.set(`filter_${field}__contains`, String(value));
                break;
              case "in":
                if (Array.isArray(value)) {
                  url.searchParams.set(`filter_${field}__in`, value.join(","));
                }
                break;
            }
          }
        });
      }

      const response = await fetchWithAuth(url.toString());
      if (!response.ok) {
        throw new Error(`Failed to fetch list: ${response.statusText}`);
      }

      const data = await response.json();
      
      return {
        data: data.data || [],
        total: data.pagination?.total || 0,
      };
    },

    getOne: async ({ resource, id, meta }) => {
      const url = `${apiUrl}/${resource}/${id}`;
      const response = await fetchWithAuth(url);
      
      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(`Resource not found: ${resource}/${id}`);
        }
        throw new Error(`Failed to fetch resource: ${response.statusText}`);
      }

      const data = await response.json();
      return {
        data: data.data,
      };
    },

    create: async ({ resource, variables, meta }) => {
      const url = `${apiUrl}/${resource}`;
      const response = await fetchWithAuth(url, {
        method: "POST",
        body: JSON.stringify(variables),
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(
          error.details || error.error || `Failed to create resource: ${response.statusText}`
        );
      }

      const data = await response.json();
      return {
        data: data.data,
      };
    },

    update: async ({ resource, id, variables, meta }) => {
      const url = `${apiUrl}/${resource}/${id}`;
      const response = await fetchWithAuth(url, {
        method: "PUT",
        body: JSON.stringify(variables),
      });

      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(
          error.details || error.error || `Failed to update resource: ${response.statusText}`
        );
      }

      const data = await response.json();
      return {
        data: data.data,
      };
    },

    deleteOne: async ({ resource, id, meta }) => {
      const url = `${apiUrl}/${resource}/${id}`;
      const response = await fetchWithAuth(url, {
        method: "DELETE",
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error(`Resource not found: ${resource}/${id}`);
        }
        throw new Error(`Failed to delete resource: ${response.statusText}`);
      }

      return {
        data: { id },
      };
    },

    getApiUrl: () => {
      return apiUrl;
    },

    custom: async ({ url, method, filters, sorters, payload, query, headers, meta }) => {
      const requestUrl = url.startsWith("http") ? url : `${apiUrl}${url}`;
      const requestOptions: RequestInit = {
        method: method || "GET",
        headers: {
          "Content-Type": "application/json",
          ...headers,
        },
      };

      if (payload) {
        requestOptions.body = JSON.stringify(payload);
      }

      if (query) {
        const urlObj = new URL(requestUrl);
        Object.entries(query).forEach(([key, value]) => {
          urlObj.searchParams.set(key, String(value));
        });
        const response = await fetchWithAuth(urlObj.toString(), requestOptions);
        const data = await response.json();
        return { data };
      }

      const response = await fetchWithAuth(requestUrl, requestOptions);
      const data = await response.json();
      return { data };
    },
  };
}

