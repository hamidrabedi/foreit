# Admin Framework - Widget & API Specifications

> **Detailed specifications for widgets, APIs, and components**

## I. Widget Specifications

### 1. Stats Card Widget

**Purpose**: Display a single metric with optional trend and comparison

**Backend Implementation**:
```go
package widgets

type StatsCardWidget struct {
    id          string
    title       string
    model       string
    aggregation AggregationType
    field       string
    icon        string
    color       string
    trend       bool
    comparison  *ComparisonConfig
}

type AggregationType string

const (
    AggCount   AggregationType = "count"
    AggSum     AggregationType = "sum"
    AggAvg     AggregationType = "avg"
    AggMin     AggregationType = "min"
    AggMax     AggregationType = "max"
    AggDistinct AggregationType = "distinct"
)

type ComparisonConfig struct {
    Period      string // "day", "week", "month", "year"
    ShowPercent bool
}

type StatsCardData struct {
    Value       interface{} `json:"value"`
    Label       string      `json:"label"`
    Trend       *TrendData  `json:"trend,omitempty"`
    Comparison  *CompData   `json:"comparison,omitempty"`
    Icon        string      `json:"icon"`
    Color       string      `json:"color"`
}

type TrendData struct {
    Direction string  `json:"direction"` // "up", "down", "flat"
    Percent   float64 `json:"percent"`
    Data      []Point `json:"data"`
}

func (w *StatsCardWidget) GetData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Get current value
    currentValue := w.performAggregation(ctx, time.Now())

    data := &StatsCardData{
        Value: currentValue,
        Label: w.title,
        Icon:  w.icon,
        Color: w.color,
    }

    // Calculate trend if enabled
    if w.trend {
        previousValue := w.performAggregation(ctx, time.Now().AddDate(0, 0, -30))
        data.Trend = w.calculateTrend(currentValue, previousValue)
    }

    // Calculate comparison if configured
    if w.comparison != nil {
        compValue := w.performComparison(ctx, w.comparison.Period)
        data.Comparison = &CompData{
            Value:   compValue,
            Percent: calculatePercent(currentValue, compValue),
            Period:  w.comparison.Period,
        }
    }

    return data, nil
}
```

**Frontend Component**:
```tsx
interface StatsCardProps {
  data: StatsCardData;
  config: WidgetConfig;
}

export function StatsCard({ data, config }: StatsCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium">
          {data.label}
        </CardTitle>
        <div className={cn("h-4 w-4", `text-${data.color}-500`)}>
          <Icon name={data.icon} />
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">
          {formatValue(data.value)}
        </div>
        {data.trend && (
          <div className="flex items-center text-xs text-muted-foreground">
            <TrendIndicator trend={data.trend} />
            <span className="ml-2">
              {data.trend.percent}% from last month
            </span>
          </div>
        )}
        {data.comparison && (
          <p className="text-xs text-muted-foreground mt-2">
            {formatValue(data.comparison.value)} last {data.comparison.period}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
```

### 2. Chart Widget

**Purpose**: Display data visualization (line, bar, pie, area, etc.)

**Backend Implementation**:
```go
type ChartWidget struct {
    id        string
    title     string
    chartType ChartType
    query     ChartQuery
    options   ChartOptions
}

type ChartType string

const (
    ChartLine   ChartType = "line"
    ChartBar    ChartType = "bar"
    ChartPie    ChartType = "pie"
    ChartArea   ChartType = "area"
    ChartDonut  ChartType = "donut"
    ChartRadar  ChartType = "radar"
    ChartScatter ChartType = "scatter"
)

type ChartQuery struct {
    Model      string
    XAxis      string // field for X-axis
    YAxis      string // field for Y-axis
    Aggregation AggregationType
    GroupBy    string
    Filters    map[string]interface{}
    TimeRange  *TimeRange
}

type ChartOptions struct {
    XAxisLabel  string
    YAxisLabel  string
    ShowLegend  bool
    ShowGrid    bool
    Colors      []string
    Height      int
    Stacked     bool
    Smooth      bool // for line/area charts
}

type ChartData struct {
    Labels  []string                 `json:"labels"`
    Datasets []ChartDataset          `json:"datasets"`
    Options ChartOptions             `json:"options"`
}

type ChartDataset struct {
    Label string        `json:"label"`
    Data  []float64     `json:"data"`
    Color string        `json:"color"`
}

func (w *ChartWidget) GetData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Execute query
    results := w.executeChartQuery(ctx, w.query)

    // Transform to chart format
    data := &ChartData{
        Labels:   extractLabels(results, w.query.XAxis),
        Datasets: []ChartDataset{},
        Options:  w.options,
    }

    // If grouped, create multiple datasets
    if w.query.GroupBy != "" {
        groups := groupResults(results, w.query.GroupBy)
        for group, values := range groups {
            data.Datasets = append(data.Datasets, ChartDataset{
                Label: group,
                Data:  extractValues(values, w.query.YAxis),
                Color: w.getColorForGroup(group),
            })
        }
    } else {
        data.Datasets = append(data.Datasets, ChartDataset{
            Label: w.title,
            Data:  extractValues(results, w.query.YAxis),
            Color: w.options.Colors[0],
        })
    }

    return data, nil
}
```

**Frontend Component**:
```tsx
import { LineChart, Line, BarChart, Bar, PieChart, Pie, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface ChartWidgetProps {
  data: ChartData;
  config: WidgetConfig;
}

export function ChartWidget({ data, config }: ChartWidgetProps) {
  const ChartComponent = getChartComponent(config.params.chartType);

  return (
    <Card>
      <CardHeader>
        <CardTitle>{config.title}</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={data.options.height || 300}>
          <ChartComponent data={transformDataForRecharts(data)}>
            <XAxis
              dataKey="label"
              label={{ value: data.options.xAxisLabel, position: 'bottom' }}
            />
            <YAxis
              label={{ value: data.options.yAxisLabel, angle: -90, position: 'left' }}
            />
            {data.options.showGrid && <CartesianGrid strokeDasharray="3 3" />}
            {data.options.showLegend && <Legend />}
            <Tooltip />
            {data.datasets.map((dataset, idx) => (
              <Line
                key={idx}
                type={data.options.smooth ? "monotone" : "linear"}
                dataKey={dataset.label}
                stroke={dataset.color}
                strokeWidth={2}
              />
            ))}
          </ChartComponent>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
```

### 3. Data Table Widget

**Purpose**: Display mini data table with sorting/filtering

**Implementation**:
```tsx
import { useReactTable, getCoreRowModel, getSortedRowModel } from '@tanstack/react-table';

interface TableWidgetProps {
  data: TableData;
  config: WidgetConfig;
}

export function TableWidget({ data, config }: TableWidgetProps) {
  const table = useReactTable({
    data: data.rows,
    columns: data.columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>{config.title}</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map(headerGroup => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map(header => (
                  <TableHead
                    key={header.id}
                    onClick={header.column.getToggleSortingHandler()}
                    className="cursor-pointer"
                  >
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    {{
                      asc: ' ↑',
                      desc: ' ↓',
                    }[header.column.getIsSorted() as string] ?? null}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.map(row => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map(cell => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
```

### 4. Activity Feed Widget

**Purpose**: Display recent activities/audit log

```go
type ActivityWidget struct {
    id      string
    title   string
    limit   int
    models  []string // filter by models
    actions []string // filter by actions
}

type ActivityData struct {
    Activities []Activity `json:"activities"`
}

type Activity struct {
    ID        int64     `json:"id"`
    User      UserInfo  `json:"user"`
    Action    string    `json:"action"`
    Model     string    `json:"model"`
    ObjectID  int64     `json:"object_id"`
    ObjectRepr string   `json:"object_repr"`
    Timestamp time.Time `json:"timestamp"`
    Changes   map[string]Change `json:"changes,omitempty"`
}

type Change struct {
    Old interface{} `json:"old"`
    New interface{} `json:"new"`
}

func (w *ActivityWidget) GetData(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    activities := w.fetchRecentActivities(ctx, w.limit, w.models, w.actions)

    return &ActivityData{
        Activities: activities,
    }, nil
}
```

```tsx
export function ActivityWidget({ data }: { data: ActivityData }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[400px]">
          <div className="space-y-4">
            {data.activities.map(activity => (
              <div key={activity.id} className="flex items-start gap-3">
                <Avatar className="h-8 w-8">
                  <AvatarImage src={activity.user.avatar} />
                  <AvatarFallback>{activity.user.initials}</AvatarFallback>
                </Avatar>
                <div className="flex-1 space-y-1">
                  <p className="text-sm">
                    <span className="font-medium">{activity.user.name}</span>
                    {' '}
                    <span className="text-muted-foreground">
                      {getActivityText(activity)}
                    </span>
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {formatDistanceToNow(activity.timestamp)} ago
                  </p>
                  {activity.changes && (
                    <ChangesList changes={activity.changes} />
                  )}
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  );
}
```

### 5. Timeline Widget

**Purpose**: Display chronological events

### 6. Map Widget

**Purpose**: Display geographic data

### 7. Calendar Widget

**Purpose**: Display events calendar

---

## II. Complete API Specification

### Admin Metadata Endpoints

#### GET /api/admin/meta
List all registered admin models

**Response**:
```json
{
  "models": [
    {
      "name": "user",
      "verbose_name": "User",
      "verbose_name_plural": "Users",
      "icon": "user",
      "count": 150,
      "permissions": {
        "add": true,
        "change": true,
        "delete": true,
        "view": true
      }
    },
    {
      "name": "post",
      "verbose_name": "Post",
      "verbose_name_plural": "Posts",
      "icon": "file-text",
      "count": 532,
      "permissions": {
        "add": true,
        "change": true,
        "delete": false,
        "view": true
      }
    }
  ]
}
```

#### GET /api/admin/meta/:model
Get detailed metadata for a specific model

**Response**:
```json
{
  "name": "post",
  "verbose_name": "Post",
  "verbose_name_plural": "Posts",
  "description": "Blog posts in the system",
  "fields": [
    {
      "name": "id",
      "type": "integer",
      "label": "ID",
      "required": true,
      "read_only": true,
      "widget": "number",
      "help_text": "Unique identifier"
    },
    {
      "name": "title",
      "type": "string",
      "label": "Title",
      "required": true,
      "max_length": 200,
      "widget": "text",
      "validators": ["required", "max_length:200"]
    },
    {
      "name": "content",
      "type": "text",
      "label": "Content",
      "required": true,
      "widget": "rich_text"
    },
    {
      "name": "author",
      "type": "relation",
      "label": "Author",
      "required": true,
      "relation_type": "foreign_key",
      "related_model": "user",
      "widget": "select"
    },
    {
      "name": "status",
      "type": "string",
      "label": "Status",
      "required": true,
      "choices": [
        { "value": "draft", "label": "Draft" },
        { "value": "published", "label": "Published" },
        { "value": "archived", "label": "Archived" }
      ],
      "widget": "select",
      "default_value": "draft"
    },
    {
      "name": "tags",
      "type": "relation",
      "label": "Tags",
      "required": false,
      "relation_type": "many_to_many",
      "related_model": "tag",
      "widget": "multi_select"
    },
    {
      "name": "featured_image",
      "type": "file",
      "label": "Featured Image",
      "required": false,
      "widget": "image_upload",
      "accept": "image/*",
      "max_size": 5242880
    },
    {
      "name": "created_at",
      "type": "datetime",
      "label": "Created At",
      "read_only": true,
      "widget": "datetime"
    }
  ],
  "list_display": ["id", "title", "author", "status", "created_at"],
  "list_filter": ["status", "author", "created_at"],
  "search_fields": ["title", "content"],
  "ordering": ["-created_at"],
  "actions": [
    {
      "name": "publish",
      "label": "Publish",
      "description": "Publish selected posts",
      "confirmation": "Are you sure you want to publish {count} posts?",
      "permissions": ["change"]
    },
    {
      "name": "delete_selected",
      "label": "Delete",
      "description": "Delete selected posts",
      "confirmation": "Are you sure you want to delete {count} posts?",
      "permissions": ["delete"],
      "dangerous": true
    }
  ],
  "filters": [
    {
      "name": "status",
      "type": "choice",
      "label": "Status",
      "choices": [
        { "value": "draft", "label": "Draft" },
        { "value": "published", "label": "Published" },
        { "value": "archived", "label": "Archived" }
      ]
    },
    {
      "name": "author",
      "type": "relation",
      "label": "Author",
      "model": "user"
    },
    {
      "name": "created_at",
      "type": "date_range",
      "label": "Created Date"
    },
    {
      "name": "has_image",
      "type": "boolean",
      "label": "Has Featured Image"
    }
  ],
  "pagination": {
    "page_size": 25,
    "max_page_size": 100
  },
  "permissions": {
    "add": true,
    "change": true,
    "delete": true,
    "view": true
  }
}
```

### CRUD Endpoints

#### GET /api/admin/:model
List objects with pagination, filtering, search, ordering

**Query Parameters**:
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 25, max: 100)
- `search`: Search query
- `ordering`: Comma-separated fields, prefix with `-` for descending
- `{field}`: Filter by field value
- `{field}__gt`, `{field}__lt`, etc.: Field lookups

**Example**:
```
GET /api/admin/post?page=2&page_size=50&status=published&author=5&created_at__gte=2025-01-01&ordering=-created_at,title
```

**Response**:
```json
{
  "count": 150,
  "next": "/api/admin/post?page=3&page_size=50&...",
  "previous": "/api/admin/post?page=1&page_size=50&...",
  "page": 2,
  "page_size": 50,
  "total_pages": 3,
  "results": [
    {
      "id": 1,
      "title": "My First Post",
      "content": "...",
      "author": {
        "id": 5,
        "name": "John Doe",
        "email": "john@example.com"
      },
      "status": "published",
      "tags": [
        { "id": 1, "name": "Tech" },
        { "id": 2, "name": "Tutorial" }
      ],
      "featured_image": "/media/posts/image.jpg",
      "created_at": "2025-01-01T10:00:00Z",
      "updated_at": "2025-01-02T15:30:00Z"
    }
  ]
}
```

#### GET /api/admin/:model/:id
Get single object

**Response**:
```json
{
  "id": 1,
  "title": "My First Post",
  "content": "Full content here...",
  "author": {
    "id": 5,
    "name": "John Doe",
    "email": "john@example.com",
    "avatar": "/media/avatars/john.jpg"
  },
  "status": "published",
  "tags": [
    { "id": 1, "name": "Tech", "slug": "tech" },
    { "id": 2, "name": "Tutorial", "slug": "tutorial" }
  ],
  "featured_image": "/media/posts/image.jpg",
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-02T15:30:00Z",
  "_meta": {
    "can_edit": true,
    "can_delete": true,
    "view_url": "/blog/my-first-post"
  }
}
```

#### POST /api/admin/:model
Create new object

**Request**:
```json
{
  "title": "New Post",
  "content": "Content here",
  "author": 5,
  "status": "draft",
  "tags": [1, 2]
}
```

**Response** (201 Created):
```json
{
  "id": 150,
  "title": "New Post",
  "content": "Content here",
  "author": { "id": 5, "name": "John Doe" },
  "status": "draft",
  "tags": [
    { "id": 1, "name": "Tech" },
    { "id": 2, "name": "Tutorial" }
  ],
  "featured_image": null,
  "created_at": "2025-01-01T12:00:00Z",
  "updated_at": "2025-01-01T12:00:00Z"
}
```

#### PATCH /api/admin/:model/:id
Partial update

**Request**:
```json
{
  "status": "published"
}
```

**Response** (200 OK):
```json
{
  "id": 150,
  "title": "New Post",
  "status": "published",
  "updated_at": "2025-01-01T13:00:00Z"
}
```

#### DELETE /api/admin/:model/:id
Delete object

**Response** (204 No Content)

### Bulk Operations

#### POST /api/admin/:model/bulk-create
Create multiple objects

**Request**:
```json
{
  "objects": [
    { "title": "Post 1", "content": "...", "author": 5, "status": "draft" },
    { "title": "Post 2", "content": "...", "author": 5, "status": "draft" }
  ]
}
```

**Response** (201 Created):
```json
{
  "created": 2,
  "objects": [
    { "id": 151, "title": "Post 1", ... },
    { "id": 152, "title": "Post 2", ... }
  ]
}
```

#### POST /api/admin/:model/action/:action
Execute bulk action

**Request**:
```json
{
  "ids": [1, 2, 3, 5, 8],
  "params": {
    "confirm": true
  }
}
```

**Response**:
```json
{
  "success": true,
  "affected": 5,
  "message": "Successfully published 5 posts"
}
```

### Search & Autocomplete

#### GET /api/admin/search
Global search across all models

**Query**: `?q=django&models=post,user`

**Response**:
```json
{
  "results": [
    {
      "model": "post",
      "count": 5,
      "items": [
        {
          "id": 1,
          "title": "Introduction to Django",
          "highlight": "Introduction to <mark>Django</mark> framework",
          "url": "/admin/post/1"
        }
      ]
    },
    {
      "model": "user",
      "count": 2,
      "items": [
        {
          "id": 10,
          "name": "Django Expert",
          "highlight": "<mark>Django</mark> Expert",
          "url": "/admin/user/10"
        }
      ]
    }
  ]
}
```

#### GET /api/admin/:model/autocomplete
Field autocomplete for dropdowns

**Query**: `?field=author&q=john`

**Response**:
```json
{
  "results": [
    { "value": 5, "label": "John Doe" },
    { "value": 12, "label": "Johnny Smith" },
    { "value": 25, "label": "John Watson" }
  ]
}
```

### File Management

#### POST /api/admin/upload
Upload file

**Request** (multipart/form-data):
```
file: <binary data>
path: posts/images/
```

**Response**:
```json
{
  "url": "/media/posts/images/abc123.jpg",
  "filename": "abc123.jpg",
  "size": 245760,
  "mime_type": "image/jpeg",
  "width": 1920,
  "height": 1080
}
```

---

## III. Field Widget Specifications

### Text Input
```tsx
<Input
  type="text"
  placeholder={field.help_text}
  maxLength={field.max_length}
  required={field.required}
  {...register(field.name)}
/>
```

### Rich Text Editor
```tsx
<RichTextEditor
  value={value}
  onChange={onChange}
  toolbar={['bold', 'italic', 'link', 'heading', 'list']}
  placeholder={field.help_text}
/>
```

### Date Picker
```tsx
<DatePicker
  selected={value}
  onChange={onChange}
  showTimeSelect={field.type === 'datetime'}
  dateFormat={field.type === 'datetime' ? 'yyyy-MM-dd HH:mm' : 'yyyy-MM-dd'}
/>
```

### File Upload
```tsx
<FileUpload
  accept={field.accept}
  maxSize={field.max_size}
  multiple={field.multiple}
  onUpload={handleUpload}
  preview={field.widget === 'image_upload'}
/>
```

### Relation Select
```tsx
<RelationSelect
  model={field.related_model}
  value={value}
  onChange={onChange}
  multiple={field.relation_type === 'many_to_many'}
  searchable
  createable={field.allow_create}
/>
```

---

## IV. Permission System

### Permission Types
- `view`: Can view list and detail
- `add`: Can create new objects
- `change`: Can edit existing objects
- `delete`: Can delete objects

### Permission Check Flow
```
Request → Authentication → Permission Check → Execute → Response
```

### Custom Permissions
```go
type PermissionChecker interface {
    HasPermission(ctx context.Context, user interface{}, permission string) bool
    HasObjectPermission(ctx context.Context, user interface{}, obj interface{}, permission string) bool
}

// Example: Owner-only permission
func (p *OwnerPermission) HasObjectPermission(ctx context.Context, user interface{}, obj interface{}, permission string) bool {
    if permission == "change" || permission == "delete" {
        return obj.GetOwnerID() == user.GetID()
    }
    return true
}
```

---

**Document Version**: 1.0
**Companion to**: ADMIN_REDESIGN_ARCHITECTURE.md
**Last Updated**: 2026-01-01
