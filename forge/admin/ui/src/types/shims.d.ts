declare module "@tanstack/react-router" {
  export interface Register {}
  export const RouterProvider: any;
  export const createRouter: any;
  export const createFileRoute: any;
  export const createRootRoute: any;
  export const redirect: any;
  export const Link: any;
  export const Outlet: any;
  export const useNavigate: any;
  export const useLocation: any;
  export const useParams: any;
}

declare module "@tanstack/react-query" {
  export const QueryClient: any;
  export const QueryClientProvider: any;
  export const useQuery: any;
  export const useMutation: any;
  export const useQueryClient: any;
}

declare module "zustand" {
  export function create<T>(initializer?: any): any;
}

declare module "zustand/middleware" {
  export const persist: any;
}

declare module "tailwind-merge" {
  export const twMerge: any;
}
