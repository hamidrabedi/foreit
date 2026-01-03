# Forge Admin Framework - Redesign Architecture

> **Comprehensive architecture plan for redesigning the Forge admin framework with modern React/TypeScript frontend and extensible Go backend**

## Executive Summary

This document outlines a complete redesign of the Forge admin framework following modern best practices for 2025. The new architecture will be:

- **Headless/API-First**: Complete decoupling of backend and frontend
- **Type-Safe**: Full TypeScript on frontend, generics in Go backend
- **Extensible**: Plugin system, widget marketplace, theme engine
- **Modern Stack**: React 18+, shadcn/ui, TanStack ecosystem, Tailwind CSS
- **Developer-Friendly**: Auto-generation, hot-reload, comprehensive docs

## Research Sources

Based on extensive research from:
- [Modern React Dashboard Best Practices](https://www.untitledui.com/blog/react-dashboards)
- [Headless CMS Architecture Patterns](https://www.contentful.com/headless-cms/)
- [Django-React Integration Patterns](https://www.digitalocean.com/community/tutorials/build-a-to-do-application-using-django-and-react)
- [shadcn/ui Admin Dashboard Components](https://www.shadcnblocks.com/admin-dashboard)
- [TanStack Ecosystem Guide](https://www.freecodecamp.org/news/build-an-admin-dashboard-with-shadcnui-and-tanstack-start/)
- Analyzed projects: CoreUI, ngx-admin, Gentelella, GoAdmin, meta-shell

---

## I. Overall Architecture

### 1.1 High-Level Design

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                              │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  React SPA (Vite + TypeScript)                             │ │
│  │  - shadcn/ui components                                     │ │
│  │  - TanStack Query (data fetching)                          │ │
│  │  - TanStack Router (routing)                               │ │
│  │  - TanStack Table (data grids)                             │ │
│  │  - Zustand (state management)                              │ │
│  │  - Recharts (data visualization)                           │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↕ REST/GraphQL API
┌─────────────────────────────────────────────────────────────────┐
│                        API LAYER                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Go Admin API Server                                        │ │
│  │  - RESTful endpoints (chi router)                          │ │
│  │  - GraphQL endpoint (gqlgen)                               │ │
│  │  - WebSocket (real-time updates)                           │ │
│  │  - OpenAPI/Swagger docs                                    │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                     BUSINESS LOGIC LAYER                         │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Admin Core System                                          │ │
│  │  - Model Registry                                           │ │
│  │  - Schema Discovery                                         │ │
│  │  - Permission System                                        │ │
│  │  - Action Framework                                         │ │
│  │  - Filter Engine                                            │ │
│  │  - Widget System                                            │ │
│  │  - Plugin Manager                                           │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────────┐
│                      DATA LAYER                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Forge ORM                                                  │ │
│  │  - Type-safe Manager[T]                                     │ │
│  │  - QuerySet[T]                                              │ │
│  │  - Schema System                                            │ │
│  │  - Migration System                                         │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Database (PostgreSQL/MySQL/SQLite)                        │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Design Principles

1. **API-First**: Everything exposed through well-documented APIs
2. **Headless**: Frontend completely decoupled, replaceable
3. **Type-Safe**: End-to-end type safety (Go generics ↔ TypeScript)
4. **Plugin-Based**: Core functionality extensible via plugins
5. **Auto-Generated**: Admin UI auto-generated from schema
6. **Customizable**: Every component overridable
7. **Real-Time**: WebSocket support for live updates
8. **Accessible**: WCAG 2.1 AA compliant
9. **Mobile-First**: Responsive design from ground up
10. **Developer Experience**: Hot reload, auto-complete, great docs

---

## II. Backend Architecture (Go)

### 2.1 Directory Structure

```
forge/
├── admin/
│   ├── core/
│   │   ├── registry.go          # Model registry
│   │   ├── admin.go              # Admin[T] type
│   │   ├── config.go             # Configuration
│   │   └── metadata.go           # Metadata extraction
│   ├── api/
│   │   ├── rest/
│   │   │   ├── router.go         # REST API router
│   │   │   ├── handlers/
│   │   │   │   ├── list.go
│   │   │   │   ├── detail.go
│   │   │   │   ├── create.go
│   │   │   │   ├── update.go
│   │   │   │   ├── delete.go
│   │   │   │   ├── bulk.go
│   │   │   │   └── actions.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   └── logger.go
│   │   │   └── serializers/
│   │   │       ├── model.go
│   │   │       ├── paginated.go
│   │   │       └── errors.go
│   │   ├── graphql/
│   │   │   ├── schema.graphql
│   │   │   ├── resolver.go
│   │   │   └── generated/
│   │   └── websocket/
│   │       ├── hub.go
│   │       ├── client.go
│   │       └── events.go
│   ├── widgets/
│   │   ├── registry.go
│   │   ├── base.go
│   │   └── builtin/
│   │       ├── stats_card.go
│   │       ├── chart.go
│   │       ├── table.go
│   │       └── timeline.go
│   ├── plugins/
│   │   ├── manager.go
│   │   ├── interface.go
│   │   ├── hooks.go
│   │   └── builtin/
│   │       ├── import_export/
│   │       ├── audit_log/
│   │       └── file_manager/
│   ├── permissions/
│   │   ├── checker.go
│   │   ├── rbac.go
│   │   └── rules.go
│   ├── actions/
│   │   ├── action.go
│   │   ├── bulk.go
│   │   └── builtin/
│   ├── filters/
│   │   ├── filter.go
│   │   ├── builder.go
│   │   └── types/
│   ├── search/
│   │   ├── engine.go
│   │   ├── indexer.go
│   │   └── highlighter.go
│   └── codegen/
│       ├── typescript.go         # Generate TS types from Go
│       ├── openapi.go            # Generate OpenAPI spec
│       └── client.go             # Generate API client
├── admin-ui/                      # Frontend (separate package)
└── examples/
    └── blog-admin/
```

### 2.2 Core API Endpoints

#### Admin Metadata API
```
GET  /api/admin/meta                    # List all registered models
GET  /api/admin/meta/:model              # Get model metadata
GET  /api/admin/meta/:model/schema       # Get JSON schema
GET  /api/admin/meta/:model/actions      # List available actions
GET  /api/admin/meta/:model/filters      # List available filters
```

#### CRUD API (per model)
```
GET    /api/admin/:model                 # List (with pagination, filtering, search)
GET    /api/admin/:model/:id             # Get single instance
POST   /api/admin/:model                 # Create new instance
PATCH  /api/admin/:model/:id             # Partial update
PUT    /api/admin/:model/:id             # Full update
DELETE /api/admin/:model/:id             # Delete instance
```

#### Bulk Operations
```
POST   /api/admin/:model/bulk-create     # Bulk create
PATCH  /api/admin/:model/bulk-update     # Bulk update
DELETE /api/admin/:model/bulk-delete     # Bulk delete
POST   /api/admin/:model/action/:name    # Execute bulk action
```

#### Related Data
```
GET    /api/admin/:model/:id/:relation   # Get related objects
POST   /api/admin/:model/:id/:relation   # Create related object
```

#### Search & Autocomplete
```
GET    /api/admin/search                 # Global search
GET    /api/admin/:model/autocomplete    # Field autocomplete
```

#### Dashboard & Widgets
```
GET    /api/admin/dashboard              # Dashboard config
GET    /api/admin/widgets                # Available widgets
GET    /api/admin/widgets/:id/data       # Widget data
```

#### File Management
```
POST   /api/admin/upload                 # File upload
GET    /api/admin/media/:path            # Serve media
DELETE /api/admin/media/:path            # Delete media
```

### 2.3 Go Types & Interfaces

```go
// Core Admin Interface
type AdminInterface interface {
    ModelName() string
    ModelType() reflect.Type
    GetMetadata() *Metadata
    GetListView() ListView
    GetDetailView() DetailView
    GetFormView() FormView
    GetPermissions() Permissions
    GetActions() []Action
    GetFilters() []Filter
    GetWidgets() []Widget
}

// Metadata for frontend consumption
type Metadata struct {
    Name              string                 `json:"name"`
    VerboseName       string                 `json:"verbose_name"`
    VerboseNamePlural string                 `json:"verbose_name_plural"`
    Fields            []FieldMetadata        `json:"fields"`
    Relations         []RelationMetadata     `json:"relations"`
    Permissions       PermissionMetadata     `json:"permissions"`
    Actions           []ActionMetadata       `json:"actions"`
    Filters           []FilterMetadata       `json:"filters"`
    ListDisplay       []string               `json:"list_display"`
    SearchFields      []string               `json:"search_fields"`
    Ordering          []string               `json:"ordering"`
    Pagination        PaginationConfig       `json:"pagination"`
}

// Field metadata with UI hints
type FieldMetadata struct {
    Name         string                 `json:"name"`
    Type         string                 `json:"type"`
    Label        string                 `json:"label"`
    HelpText     string                 `json:"help_text"`
    Required     bool                   `json:"required"`
    ReadOnly     bool                   `json:"read_only"`
    Choices      []Choice               `json:"choices,omitempty"`
    Widget       string                 `json:"widget"`
    Validators   []ValidatorMetadata    `json:"validators"`
    DefaultValue interface{}            `json:"default_value,omitempty"`
}

// Plugin system
type Plugin interface {
    Name() string
    Initialize(app *admin.App) error
    RegisterRoutes(router chi.Router)
    RegisterWidgets() []Widget
    RegisterActions() []Action
    Hooks() PluginHooks
}

type PluginHooks struct {
    BeforeList   []HookFunc
    AfterList    []HookFunc
    BeforeSave   []HookFunc
    AfterSave    []HookFunc
    BeforeDelete []HookFunc
    AfterDelete  []HookFunc
}
```

### 2.4 Code Generation System

Generate TypeScript types from Go schemas:

```go
// Generate TypeScript interfaces
func GenerateTypeScript(registry *Registry) (string, error) {
    var buf bytes.Buffer

    for _, admin := range registry.GetAll() {
        meta := admin.GetMetadata()

        // Generate interface
        fmt.Fprintf(&buf, "export interface %s {\n", meta.Name)
        for _, field := range meta.Fields {
            tsType := goTypeToTS(field.Type)
            optional := ""
            if !field.Required {
                optional = "?"
            }
            fmt.Fprintf(&buf, "  %s%s: %s;\n", field.Name, optional, tsType)
        }
        fmt.Fprintf(&buf, "}\n\n")

        // Generate API client methods
        generateAPIClient(&buf, meta)
    }

    return buf.String(), nil
}
```

---

## III. Frontend Architecture (React + TypeScript)

### 3.1 Directory Structure

```
admin-ui/
├── src/
│   ├── main.tsx                  # Entry point
│   ├── App.tsx                   # Root component
│   ├── vite-env.d.ts
│   ├── api/
│   │   ├── client.ts             # API client base
│   │   ├── admin.ts              # Admin API methods
│   │   ├── types.ts              # Generated types
│   │   └── hooks/                # React Query hooks
│   │       ├── useModels.ts
│   │       ├── useModelList.ts
│   │       ├── useModelDetail.ts
│   │       ├── useModelCreate.ts
│   │       └── useBulkAction.ts
│   ├── components/
│   │   ├── ui/                   # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── table.tsx
│   │   │   ├── form.tsx
│   │   │   ├── dialog.tsx
│   │   │   └── ...
│   │   ├── layout/
│   │   │   ├── AppShell.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── TopBar.tsx
│   │   │   ├── Breadcrumbs.tsx
│   │   │   └── Footer.tsx
│   │   ├── admin/
│   │   │   ├── ModelList.tsx
│   │   │   ├── ModelDetail.tsx
│   │   │   ├── ModelForm.tsx
│   │   │   ├── ModelTable.tsx
│   │   │   ├── FilterSidebar.tsx
│   │   │   ├── SearchBar.tsx
│   │   │   ├── BulkActions.tsx
│   │   │   └── Pagination.tsx
│   │   ├── fields/               # Form fields
│   │   │   ├── TextField.tsx
│   │   │   ├── TextArea.tsx
│   │   │   ├── NumberField.tsx
│   │   │   ├── DateField.tsx
│   │   │   ├── SelectField.tsx
│   │   │   ├── MultiSelect.tsx
│   │   │   ├── FileUpload.tsx
│   │   │   ├── RichText.tsx
│   │   │   ├── RelationField.tsx
│   │   │   └── FieldFactory.tsx  # Dynamic field renderer
│   │   ├── widgets/
│   │   │   ├── StatsCard.tsx
│   │   │   ├── ChartWidget.tsx
│   │   │   ├── TableWidget.tsx
│   │   │   ├── TimelineWidget.tsx
│   │   │   ├── ActivityWidget.tsx
│   │   │   └── WidgetRenderer.tsx
│   │   ├── data/
│   │   │   ├── DataTable.tsx     # TanStack Table wrapper
│   │   │   ├── VirtualList.tsx
│   │   │   ├── Charts.tsx
│   │   │   └── EmptyState.tsx
│   │   └── feedback/
│   │       ├── LoadingState.tsx
│   │       ├── ErrorBoundary.tsx
│   │       ├── Toast.tsx
│   │       └── CommandPalette.tsx
│   ├── features/                 # Feature-based organization
│   │   ├── dashboard/
│   │   │   ├── Dashboard.tsx
│   │   │   ├── DashboardConfig.tsx
│   │   │   └── hooks/
│   │   ├── auth/
│   │   │   ├── Login.tsx
│   │   │   ├── Logout.tsx
│   │   │   └── useAuth.ts
│   │   ├── settings/
│   │   │   └── Settings.tsx
│   │   └── admin/
│   │       ├── AdminRoot.tsx
│   │       ├── ListView.tsx
│   │       ├── DetailView.tsx
│   │       ├── CreateView.tsx
│   │       └── EditView.tsx
│   ├── lib/
│   │   ├── utils.ts
│   │   ├── validators.ts
│   │   ├── formatters.ts
│   │   └── queryClient.ts
│   ├── hooks/
│   │   ├── useTheme.ts
│   │   ├── useMediaQuery.ts
│   │   ├── useDebounce.ts
│   │   └── useLocalStorage.ts
│   ├── store/                    # Zustand stores
│   │   ├── authStore.ts
│   │   ├── uiStore.ts
│   │   └── adminStore.ts
│   ├── router/
│   │   ├── index.tsx             # TanStack Router config
│   │   └── routes/
│   │       ├── __root.tsx
│   │       ├── index.tsx
│   │       ├── dashboard.tsx
│   │       └── admin/
│   │           ├── $model.tsx
│   │           ├── $model.$id.tsx
│   │           └── $model.create.tsx
│   ├── styles/
│   │   ├── globals.css
│   │   └── themes/
│   │       ├── default.css
│   │       ├── dark.css
│   │       └── custom.css
│   └── types/
│       ├── admin.ts
│       ├── models.ts             # Generated from backend
│       └── api.ts
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── tailwind.config.js
└── components.json               # shadcn/ui config
```

### 3.2 Tech Stack

```json
{
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "@tanstack/react-query": "^5.83.0",
    "@tanstack/react-table": "^8.21.3",
    "@tanstack/react-router": "^1.80.0",
    "@tanstack/react-virtual": "^3.10.0",
    "zustand": "^5.0.9",
    "react-hook-form": "^7.61.1",
    "@hookform/resolvers": "^3.10.0",
    "zod": "^3.25.76",
    "axios": "^1.7.9",
    "@radix-ui/react-*": "latest",
    "lucide-react": "^0.462.0",
    "recharts": "^2.15.4",
    "cmdk": "^1.1.1",
    "sonner": "^1.7.4",
    "date-fns": "^3.6.0",
    "framer-motion": "^12.23.26",
    "tailwind-merge": "^2.6.0",
    "tailwindcss-animate": "^1.0.7",
    "class-variance-authority": "^0.7.1"
  },
  "devDependencies": {
    "typescript": "^5.8.3",
    "vite": "^5.4.19",
    "@vitejs/plugin-react-swc": "^3.11.0",
    "tailwindcss": "^3.4.17",
    "autoprefixer": "^8.5.6",
    "eslint": "^9.32.0",
    "prettier": "^3.4.2"
  }
}
```

### 3.3 Core Components

#### ModelList Component
```tsx
import { useModelList } from '@/api/hooks/useModelList';
import { DataTable } from '@/components/data/DataTable';
import { FilterSidebar } from '@/components/admin/FilterSidebar';
import { BulkActions } from '@/components/admin/BulkActions';

interface ModelListProps {
  modelName: string;
}

export function ModelList({ modelName }: ModelListProps) {
  const {
    data,
    isLoading,
    error,
    filters,
    setFilters,
    pagination,
    setPagination,
    sorting,
    setSorting,
  } = useModelList(modelName);

  const columns = useMemo(() =>
    generateColumns(data?.meta.fields),
    [data?.meta.fields]
  );

  return (
    <div className="flex h-full">
      <FilterSidebar
        filters={data?.meta.filters}
        values={filters}
        onChange={setFilters}
      />
      <div className="flex-1 flex flex-col">
        <BulkActions
          actions={data?.meta.actions}
          selectedRows={selectedRows}
        />
        <DataTable
          columns={columns}
          data={data?.results}
          pagination={pagination}
          onPaginationChange={setPagination}
          sorting={sorting}
          onSortingChange={setSorting}
        />
      </div>
    </div>
  );
}
```

#### Dynamic Form Builder
```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { FieldFactory } from '@/components/fields/FieldFactory';

interface ModelFormProps<T> {
  modelName: string;
  initialData?: T;
  onSubmit: (data: T) => void;
}

export function ModelForm<T>({ modelName, initialData, onSubmit }: ModelFormProps<T>) {
  const { data: meta } = useModelMetadata(modelName);

  const schema = useMemo(() =>
    generateZodSchema(meta?.fields),
    [meta?.fields]
  );

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: initialData,
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        {meta?.fields.map(field => (
          <FieldFactory
            key={field.name}
            field={field}
            control={form.control}
          />
        ))}
        <Button type="submit">Save</Button>
      </form>
    </Form>
  );
}
```

### 3.4 State Management Strategy

```typescript
// Auth Store (Zustand)
interface AuthStore {
  user: User | null;
  token: string | null;
  login: (credentials: Credentials) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  token: localStorage.getItem('token'),
  isAuthenticated: !!localStorage.getItem('token'),
  login: async (credentials) => {
    const { user, token } = await api.login(credentials);
    localStorage.setItem('token', token);
    set({ user, token, isAuthenticated: true });
  },
  logout: () => {
    localStorage.removeItem('token');
    set({ user: null, token: null, isAuthenticated: false });
  },
}));

// UI Store
interface UIStore {
  sidebarOpen: boolean;
  theme: 'light' | 'dark' | 'system';
  toggleSidebar: () => void;
  setTheme: (theme: UIStore['theme']) => void;
}

export const useUIStore = create<UIStore>((set) => ({
  sidebarOpen: true,
  theme: 'system',
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  setTheme: (theme) => set({ theme }),
}));
```

---

## IV. Widget System

### 4.1 Widget Architecture

Widgets are reusable dashboard components that can display data in various formats.

#### Backend Widget Interface
```go
type Widget interface {
    Name() string
    Type() WidgetType
    GetData(ctx context.Context, params map[string]interface{}) (interface{}, error)
    GetConfig() WidgetConfig
}

type WidgetType string

const (
    WidgetStatsCard  WidgetType = "stats_card"
    WidgetChart      WidgetType = "chart"
    WidgetTable      WidgetType = "table"
    WidgetTimeline   WidgetType = "timeline"
    WidgetActivity   WidgetType = "activity"
    WidgetCustom     WidgetType = "custom"
)

type WidgetConfig struct {
    ID          string                 `json:"id"`
    Type        WidgetType             `json:"type"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Size        WidgetSize             `json:"size"`
    RefreshRate int                    `json:"refresh_rate"` // seconds
    Params      map[string]interface{} `json:"params"`
}

// Example: Stats Card Widget
type StatsCardWidget struct {
    model   string
    field   string
    aggFunc string // count, sum, avg, min, max
}

func (w *StatsCardWidget) GetData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Query database using ORM
    result := performAggregation(w.model, w.field, w.aggFunc)

    return StatsCardData{
        Value:      result.Value,
        Label:      w.field,
        Trend:      result.Trend,
        Comparison: result.Comparison,
    }, nil
}
```

#### Frontend Widget System
```tsx
// Widget Registry
const widgetRegistry = {
  stats_card: StatsCardWidget,
  chart: ChartWidget,
  table: TableWidget,
  timeline: TimelineWidget,
  activity: ActivityWidget,
};

// Widget Renderer
interface WidgetRendererProps {
  config: WidgetConfig;
}

export function WidgetRenderer({ config }: WidgetRendererProps) {
  const { data, isLoading } = useWidgetData(config.id);

  const Widget = widgetRegistry[config.type];

  if (!Widget) {
    return <div>Unknown widget type: {config.type}</div>;
  }

  return (
    <Card className={getWidgetSize(config.size)}>
      <CardHeader>
        <CardTitle>{config.title}</CardTitle>
        {config.description && (
          <CardDescription>{config.description}</CardDescription>
        )}
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton />
        ) : (
          <Widget data={data} config={config} />
        )}
      </CardContent>
    </Card>
  );
}

// Dashboard with drag-and-drop grid
import { Responsive, WidthProvider } from 'react-grid-layout';

const ResponsiveGridLayout = WidthProvider(Responsive);

export function Dashboard() {
  const { data: widgets } = useDashboardWidgets();
  const [layout, setLayout] = useLocalStorage('dashboard-layout', []);

  return (
    <ResponsiveGridLayout
      layouts={{ lg: layout }}
      onLayoutChange={(newLayout) => setLayout(newLayout)}
      breakpoints={{ lg: 1200, md: 996, sm: 768 }}
      cols={{ lg: 12, md: 10, sm: 6 }}
    >
      {widgets.map(widget => (
        <div key={widget.id}>
          <WidgetRenderer config={widget} />
        </div>
      ))}
    </ResponsiveGridLayout>
  );
}
```

### 4.2 Built-in Widgets

1. **Stats Card**: Display single metric with trend
2. **Chart Widget**: Line, bar, pie, area charts
3. **Table Widget**: Mini data table
4. **Timeline Widget**: Activity timeline
5. **Activity Feed**: Recent actions
6. **Map Widget**: Geographic data
7. **Calendar Widget**: Events calendar
8. **Custom Widget**: Developer-defined

---

## V. Plugin System

### 5.1 Plugin Architecture

Plugins extend admin functionality without modifying core code.

```go
// Plugin interface
type Plugin interface {
    // Metadata
    Name() string
    Version() string
    Description() string
    Author() string

    // Lifecycle
    Initialize(app *AdminApp) error
    Start() error
    Stop() error

    // Extensibility
    RegisterRoutes(router chi.Router)
    RegisterWidgets() []Widget
    RegisterActions() []Action
    RegisterFieldTypes() []FieldType
    RegisterFilters() []Filter

    // Hooks
    Hooks() PluginHooks

    // Settings
    Settings() []Setting
    SettingsSchema() interface{}
}

// Plugin manager
type PluginManager struct {
    plugins map[string]Plugin
    hooks   *HookRegistry
}

func (pm *PluginManager) Register(plugin Plugin) error {
    if err := plugin.Initialize(pm.app); err != nil {
        return err
    }

    pm.plugins[plugin.Name()] = plugin
    pm.registerHooks(plugin.Hooks())

    return nil
}
```

### 5.2 Built-in Plugins

#### 1. Import/Export Plugin
```go
type ImportExportPlugin struct {
    supportedFormats []string // csv, xlsx, json, xml
}

func (p *ImportExportPlugin) RegisterRoutes(router chi.Router) {
    router.Post("/admin/{model}/import", p.HandleImport)
    router.Get("/admin/{model}/export", p.HandleExport)
}

func (p *ImportExportPlugin) HandleExport(w http.ResponseWriter, r *http.Request) {
    format := r.URL.Query().Get("format")
    modelName := chi.URLParam(r, "model")

    // Generate export file
    data := fetchModelData(modelName)
    file := exportToFormat(data, format)

    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", modelName, format))
    w.Write(file)
}
```

#### 2. Audit Log Plugin
Tracks all changes made through admin interface.

```go
type AuditLogPlugin struct {
    db *sql.DB
}

type AuditLog struct {
    ID         int64
    UserID     int64
    ModelName  string
    ObjectID   int64
    Action     string // create, update, delete
    Changes    JSONB
    Timestamp  time.Time
    IPAddress  string
    UserAgent  string
}

func (p *AuditLogPlugin) Hooks() PluginHooks {
    return PluginHooks{
        AfterSave: []HookFunc{p.LogSave},
        AfterDelete: []HookFunc{p.LogDelete},
    }
}
```

#### 3. File Manager Plugin
Advanced file/media management with thumbnails, cropping, CDN integration.

#### 4. Advanced Search Plugin
Full-text search with Elasticsearch/Meilisearch integration.

#### 5. Multi-tenancy Plugin
Support for multi-tenant applications.

---

## VI. API Design Patterns

### 6.1 REST API Patterns

#### Pagination
```json
{
  "count": 150,
  "next": "/api/admin/posts?page=2",
  "previous": null,
  "page_size": 25,
  "page": 1,
  "total_pages": 6,
  "results": [...]
}
```

#### Filtering
```
GET /api/admin/posts?status=published&created_after=2025-01-01&author=5
GET /api/admin/posts?search=django&ordering=-created_at
```

#### Field Selection (Sparse Fieldsets)
```
GET /api/admin/posts?fields=id,title,author.name
```

#### Including Relations
```
GET /api/admin/posts?include=author,tags,comments.count
```

#### Error Responses
```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": {
      "title": ["This field is required"],
      "email": ["Enter a valid email address"]
    }
  }
}
```

### 6.2 GraphQL Schema

```graphql
type Query {
  # Get model metadata
  adminModel(name: String!): AdminModel
  adminModels: [AdminModel!]!

  # CRUD operations
  adminObjects(
    model: String!
    page: Int
    pageSize: Int
    filters: JSON
    search: String
    ordering: [String!]
  ): PaginatedResult!

  adminObject(model: String!, id: ID!): AdminObject
}

type Mutation {
  adminCreate(model: String!, data: JSON!): AdminObject!
  adminUpdate(model: String!, id: ID!, data: JSON!): AdminObject!
  adminDelete(model: String!, id: ID!): Boolean!
  adminBulkAction(model: String!, action: String!, ids: [ID!]!): BulkActionResult!
}

type Subscription {
  adminObjectChanged(model: String!, id: ID!): AdminObject!
  adminListChanged(model: String!): AdminListEvent!
}

type AdminModel {
  name: String!
  verboseName: String!
  fields: [Field!]!
  actions: [Action!]!
  filters: [Filter!]!
}

type PaginatedResult {
  count: Int!
  pageSize: Int!
  page: Int!
  totalPages: Int!
  results: [AdminObject!]!
}
```

### 6.3 Real-Time Updates (WebSocket)

```typescript
// Frontend WebSocket client
class AdminWebSocket {
  private ws: WebSocket;
  private listeners: Map<string, Set<Function>>;

  connect(token: string) {
    this.ws = new WebSocket(`ws://localhost:8000/ws?token=${token}`);

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.dispatch(data.type, data.payload);
    };
  }

  subscribe(event: string, callback: Function) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set());
    }
    this.listeners.get(event)!.add(callback);

    // Send subscription message
    this.ws.send(JSON.stringify({
      type: 'subscribe',
      event,
    }));
  }

  private dispatch(event: string, payload: any) {
    const callbacks = this.listeners.get(event);
    if (callbacks) {
      callbacks.forEach(cb => cb(payload));
    }
  }
}

// Usage in React
function useAdminWebSocket(model: string, id?: string) {
  const queryClient = useQueryClient();

  useEffect(() => {
    const ws = adminWebSocket;

    ws.subscribe(`admin.${model}.changed`, (data) => {
      // Invalidate queries to refetch
      queryClient.invalidateQueries(['admin', model]);
    });

    return () => ws.unsubscribe(`admin.${model}.changed`);
  }, [model, queryClient]);
}
```

---

## VII. Implementation Roadmap

### Phase 1: Foundation (Weeks 1-4)

**Backend:**
- [ ] Restructure admin package
- [ ] Implement metadata extraction system
- [ ] Build REST API router with chi
- [ ] Create serializer system
- [ ] Add CORS middleware
- [ ] Implement authentication middleware
- [ ] Set up OpenAPI documentation

**Frontend:**
- [ ] Initialize Vite + React + TypeScript project
- [ ] Install and configure shadcn/ui
- [ ] Set up TanStack Query, Router, Table
- [ ] Create basic layout components
- [ ] Implement auth flow
- [ ] Build API client with axios

**DevOps:**
- [ ] Set up development environment
- [ ] Configure hot reload for both backend and frontend
- [ ] Create Docker compose for local development

### Phase 2: Core Admin Features (Weeks 5-8)

**Backend:**
- [ ] Implement List endpoint with filtering, search, pagination
- [ ] Implement Detail endpoint
- [ ] Implement Create/Update/Delete endpoints
- [ ] Add field validation
- [ ] Implement permission system
- [ ] Build action framework

**Frontend:**
- [ ] Create ModelList component with TanStack Table
- [ ] Build dynamic form generator
- [ ] Implement filter sidebar
- [ ] Create ModelDetail view
- [ ] Add create/edit forms
- [ ] Implement bulk actions UI

### Phase 3: Advanced Features (Weeks 9-12)

**Backend:**
- [ ] GraphQL endpoint with gqlgen
- [ ] WebSocket support for real-time updates
- [ ] Widget system implementation
- [ ] Plugin manager
- [ ] File upload handling
- [ ] Global search with indexing

**Frontend:**
- [ ] Dashboard with drag-and-drop widgets
- [ ] Command palette (Cmd+K)
- [ ] Advanced filtering UI
- [ ] File upload components
- [ ] Rich text editor integration
- [ ] Chart components with Recharts

### Phase 4: Extensibility (Weeks 13-16)

**Backend:**
- [ ] Plugin: Import/Export
- [ ] Plugin: Audit Log
- [ ] Plugin: File Manager
- [ ] Plugin: Advanced Search
- [ ] Hook system for plugins
- [ ] Code generation for TypeScript types

**Frontend:**
- [ ] Plugin system architecture
- [ ] Theme customization
- [ ] Widget marketplace UI
- [ ] Settings panel
- [ ] Keyboard shortcuts

### Phase 5: Polish & Documentation (Weeks 17-20)

- [ ] Comprehensive documentation
- [ ] Interactive tutorial
- [ ] Example projects
- [ ] Performance optimization
- [ ] Accessibility audit (WCAG 2.1 AA)
- [ ] Mobile responsiveness testing
- [ ] E2E testing with Playwright
- [ ] Storybook for components

---

## VIII. File Structure Overview

Final structure:

```
forge/
├── admin/                    # Go backend
│   ├── core/
│   ├── api/
│   ├── widgets/
│   ├── plugins/
│   ├── permissions/
│   ├── actions/
│   ├── filters/
│   ├── search/
│   └── codegen/
├── admin-ui/                 # React frontend
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── features/
│   │   ├── hooks/
│   │   ├── lib/
│   │   ├── router/
│   │   ├── store/
│   │   ├── styles/
│   │   └── types/
│   └── public/
├── examples/
│   ├── blog/
│   ├── ecommerce/
│   └── crm/
└── docs/
    ├── getting-started.md
    ├── api-reference.md
    ├── components.md
    ├── plugins.md
    └── deployment.md
```

---

## IX. Key Technologies Summary

### Backend Stack
- **Language**: Go 1.24+
- **Router**: chi/v5
- **ORM**: Forge ORM (existing)
- **GraphQL**: gqlgen
- **WebSocket**: gorilla/websocket
- **Auth**: JWT
- **API Docs**: swaggo/swag (OpenAPI)
- **Testing**: testify

### Frontend Stack
- **Framework**: React 18.3+
- **Language**: TypeScript 5.8+
- **Build Tool**: Vite 5.4+
- **UI Library**: shadcn/ui (Radix UI + Tailwind)
- **Data Fetching**: TanStack Query
- **Routing**: TanStack Router
- **Tables**: TanStack Table
- **Forms**: React Hook Form + Zod
- **State**: Zustand
- **Charts**: Recharts
- **Styling**: Tailwind CSS 3.4+
- **Animations**: Framer Motion
- **Icons**: Lucide React
- **Testing**: Vitest + Playwright

---

## X. Success Criteria

The redesigned admin framework will be considered successful when:

1. ✅ **Type Safety**: Full type safety from Go backend to TypeScript frontend
2. ✅ **Auto-Generation**: Admin UI auto-generated from schema with zero config
3. ✅ **Performance**: Initial load < 2s, list rendering > 60fps for 1000+ rows
4. ✅ **Extensibility**: Plugins can add widgets, actions, field types without core changes
5. ✅ **Developer Experience**: Hot reload < 500ms, comprehensive autocomplete
6. ✅ **Accessibility**: WCAG 2.1 AA compliant, keyboard navigation
7. ✅ **Mobile**: Fully responsive, touch-optimized
8. ✅ **Real-time**: WebSocket updates with < 100ms latency
9. ✅ **Documentation**: Every API endpoint and component documented
10. ✅ **Testing**: > 80% test coverage

---

## XI. Next Steps

1. **Review & Approval**: Review this architecture document
2. **Prototype**: Build small prototype to validate key decisions
3. **Team Formation**: Assign developers to backend/frontend
4. **Sprint Planning**: Break roadmap into 2-week sprints
5. **Start Development**: Begin Phase 1 implementation

---

**Document Version**: 1.0
**Last Updated**: 2026-01-01
**Author**: Architecture Team
**Status**: Pending Review
