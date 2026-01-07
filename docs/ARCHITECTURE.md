# forge Architecture (comprehensive, code-aligned)

## Table of contents

- [1. Scope and non-goals](#1-scope-and-non-goals)
- [2. Core philosophy](#2-core-philosophy)
- [3. Package map](#3-package-map)
- [4. End-to-end flows](#4-end-to-end-flows)
- [5. Schema system](#5-schema-system)
- [6. Code generation](#6-code-generation)
- [7. ORM and QuerySet](#7-orm-and-queryset)
- [8. Filtering system](#8-filtering-system)
- [9. Migration system](#9-migration-system)
- [10. Admin system](#10-admin-system)
- [11. API framework](#11-api-framework)
- [12. Identity system](#12-identity-system)
- [13. Server, middleware, security](#13-server-middleware-security)
- [14. Config, logging, validation](#14-config-logging-validation)
- [15. Extension points](#15-extension-points)
- [16. Testing architecture](#16-testing-architecture)
- [17. Alignment with archived docs](#17-alignment-with-archived-docs)

## 1. Scope and non-goals

This document describes how the current codebase is organized and how major runtime and generation flows work. It is not a line-by-line code reference. Historical long-form narratives and comparisons live under `docs/archive/`.

Non-goals:
- Reproducing every archived document verbatim.
- Promising features that are not present in `forge/` (we call those out as roadmap items instead).

## 2. Core philosophy

Forge is a Django-inspired framework for Go. The implemented philosophy, consistent with the archive, is:

- Type-safe first: generated models/managers/querysets provide compile-time safety.
- Dynamic when needed: runtime filter expressions and generic field references exist for edge cases.
- Convention over configuration: standard folder structure and CLI workflows.
- Extensible: registries and interfaces are preferred over hard-coded wiring.
- Security by default: middleware and identity system are built-in and test-covered.

## 3. Package map

The canonical implementation lives in the `forge/` module. The most important packages and why they exist:

- `forge/schema/`: the schema DSL (fields, relations, meta, hooks, registries).
- `forge/codegen/`: parses Go source using AST and generates strongly typed model helpers.
- `forge/orm/`: the type-safe ORM core (BaseQuerySet, expressions, managers, update builders, projections).
- `forge/filter/`: a reusable filter AST + filterset engine used by admin and API.
- `forge/db/` + `forge/db/migrate/`: database connection wrappers and the migration engine (state, diff, sql builder, execution).
- `forge/admin/`: admin registry, metadata builder, handlers, utilities, and UI helpers.
- `forge/api/`: DRF-like API framework (serializers, viewsets, auth, permissions, throttling, parsing/rendering, OpenAPI generator).
- `forge/identity/`: users/sessions/tokens + backends + middleware + services.
- `forge/server/`: server/router/middleware integration (chi router wrappers, security helpers).
- `forge/config/`, `forge/log/`, `forge/validate/`, `forge/utils/`: infrastructure.
- `forge/registry/`: registries/plugins/extensions.

Supporting modules:
- `examples/`: example projects used as references and by tests.
- `tests/`: integration and end-to-end tests (CLI, migrations, schema, ORM).

## 4. End-to-end flows

### 4.1 Model -> generated code

1. You implement a schema type in Go, usually embedding `schema.BaseSchema`.
2. You implement (some subset of) `Fields()`, `Relations()`, `Meta()`, `Hooks()`.
3. `forge generate` scans model directories and parses them using `forge/codegen/ast_parser.go`.
4. Codegen emits:
   - field expressions used for typed query building
   - a typed manager with CRUD helpers
   - a typed QuerySet wrapper embedding `forge/orm.BaseQuerySet`

The ecommerce example under `examples/ecommerce/app/*` and `examples/ecommerce/models/*` is used by migration and CLI tests.

### 4.2 HTTP request flow (server + middleware)

A typical request path:

1. Incoming request hits the server/router layer (`forge/server/router.go`).
2. Middleware stack runs (security, logging, identity, etc.).
3. The route handler is either:
   - an admin handler from `forge/admin/*` (list/detail/create/update/delete)
   - an API handler/viewset from `forge/api/*`
4. Handlers call into ORM/query layers to build SQL.
5. Database operations go through `forge/db` wrappers (transactions, connections).
6. Response is rendered using API renderers or admin templates/UI helpers.

### 4.3 Migration workflow

1. `forge makemigrations <name>` scans models and computes a schema diff.
2. Diffing uses:
   - AST model definitions from `forge/codegen`.
   - previous schema state loaded from existing migration SQL via `forge/db/migrate/state/loader.go`.
   - change detection in `forge/db/migrate/generate/detector.go`.
3. SQL generation uses `forge/db/migrate/sql/*` (driver-specific builders).
4. Execution uses `forge/db/migrations.go` (golang-migrate integration), plus recovery/status helpers under `forge/db/migrate/execute/*`.

### 4.4 Admin and identity integration

Admin is designed to work out-of-the-box with identity and permissions:

- Identity provides authentication middleware and session/token handling (`forge/identity/middleware/*`).
- Admin uses registry + metadata to render forms and enforce permissions.
- The design goal from the archive ("built-in auth but overridable") is met by ensuring admin and identity wire through registries/interfaces rather than requiring hard-coded calls.

## 5. Schema system

The schema DSL is implemented under `forge/schema/` and is the foundational contract for code generation and migrations.

### 5.1 Schema interface and BaseSchema

Schemas are plain Go types that implement some subset of these methods:

- `Fields() []schema.Field`
- `Relations() []schema.Relation`
- `Meta() schema.Meta`
- `Hooks() *schema.ModelHooks` (or equivalent; see `forge/schema/hooks.go`)

Most schemas embed `schema.BaseSchema` for conventions and for codegen discovery.

### 5.2 Field types

The type constructors in `forge/schema/typed_builders.go` include:

- Numeric: `schema.Int64`, `schema.Int32`, `schema.Float64`, `schema.Float32`, `schema.Decimal`
- Text: `schema.String`, `schema.Text`, `schema.Email`, `schema.URL`
- Boolean: `schema.Bool`
- Temporal: `schema.Time`, `schema.Date`, `schema.DateTime`
- Special: `schema.JSON`, `schema.Bytes`, `schema.UUID`

Field values are configured via chain methods on `schema.Field` (see `forge/schema/field_methods.go`).

### 5.3 Field options and how they map to DB and admin

Common options (non-exhaustive; see `forge/schema/field.go` and `forge/schema/field_methods.go`):

- Required / optional:
  - `WithRequired()` -> NOT NULL
  - `WithOptional()` -> NULL
- Key and identity:
  - `WithPrimary()` -> primary key
  - `WithAutoIncrement()` -> identity/serial behavior (driver-specific)
- Uniqueness and indexing:
  - `WithUnique()` -> unique constraint
  - `WithDBIndex()` / index helpers -> index creation
- Column metadata:
  - `WithDBColumn("...")` -> custom column name
  - `WithDBType("...")` -> custom database type
  - `WithDBDefault("...")` -> SQL default (driver-specific)
- Text constraints:
  - `WithMaxLength(n)`, `WithMinLength(n)`
- Temporal automation:
  - `WithAutoNowAdd()` implies DEFAULT now() and NOT NULL for created-at semantics
  - `WithAutoNow()` implies DEFAULT now() for updated-at semantics
- Generated columns:
  - `WithGeneratedColumn(expression, stored)` -> generated SQL (handled in migrations).

These options are consumed by:

- `forge/codegen/ast_parser.go` when producing `FieldDefinition.Options`.
- `forge/db/migrate/sql/base.go` when mapping field options into DDL.

### 5.4 Relations

Relations are defined via constructors in `forge/schema/relation.go`:

- `schema.ForeignKey(columnName, "TargetModel")`
- `schema.OneToOne(columnName, "TargetModel")`
- `schema.ManyToMany(name, "TargetModel")`
- `schema.OneToMany(name, "TargetModel", fkColumn)` (reverse relation helper)

Key relation options:

- `WithRelatedName("...")` and `WithRelatedQueryName("...")`
- `WithOnDelete(schema.CascadeCASCADE|CascadeSET_NULL|CascadePROTECT|...)`
- `WithOnUpdate(...)`
- `WithThroughTable("...")` (many-to-many join tables)
- Constraint tuning: `WithDBConstraint(bool)`, `WithConstraintName`, `WithDeferrable`, `WithMatch`

Migrations translate these into foreign keys (`ON DELETE`, `ON UPDATE`, constraint names, etc.) in `forge/db/migrate/sql/*`.

### 5.5 Meta options

The meta definition (see `forge/schema/meta.go`) includes table name, indexes, constraints, ordering hints, and naming helpers. Codegen and migrations read meta to ensure stable naming and to create secondary indexes/constraints.

### 5.6 Hooks

The hook system (see `forge/schema/hooks.go`) provides lifecycle entry points (before/after create/update/save/delete) to match Django-like behavior.

Design note: hooks are invoked in managers and/or model instance methods generated by codegen; changes to hook signatures should be coordinated with template updates.

### 5.7 Registries

Forge uses registries to avoid hard-coded global wiring:

- `forge/schema/registry.go` registers field and relation factories.
- `forge/registry/*` provides broader extension registries.

This design supports "auto discovery" and pluggability described in `docs/archive/AUTO_DISCOVERY.md`.

## 6. Code generation

Code generation is responsible for turning schema definitions into a type-safe developer API.

### 6.1 Inputs

- Model schema files in user apps and examples.
- Schema builders and relation definitions in `forge/schema`.

### 6.2 AST parser responsibilities

`forge/codegen/ast_parser.go`:

- Walks Go files and identifies schema types and their methods.
- Extracts:
  - Fields and their builder chains (options like `required`, `unique`, `primary`, `db_column`, `db_default`, `generated_column`).
  - Relations and relation chains (options like `on_delete`, `on_update`, `related_name`, `through`).
  - Meta (table name, indexes, constraints).

Important invariants:

- Primary keys must be treated as required (NOT NULL). The parser ensures this so migrations do not attempt to make identity primary keys nullable.
- Cascade constants are normalized and passed through as SQL cascade actions.

### 6.3 Templates and writer

`forge/codegen/writer.go` and templates under `forge/codegen/templates/`:

- Manage consistent file layout and imports.
- Emit:
  - Field expression types used for query building.
  - Manager types used for CRUD.
  - QuerySet wrapper types that embed `forge/orm.BaseQuerySet`.

### 6.4 Generated API surface

Generated code provides:

- `Model.Fields.<FieldName>` typed expressions (used for `Eq`, `Gt`, `Contains`, etc.).
- `Model.Objects` or manager instances with `All`, `Get`, `Create`, `Update`, `Delete`.
- `Model.Objects.Filter(...).OrderBy(...).Limit(...).All(ctx)` style chainable QuerySets.

See `forge/orm/TYPE_SAFE_API.md` for usage patterns and `forge/orm/queryset.go` for the canonical interface.

### 6.5 Codegen and migrations coupling

The migration generator loads model definitions via AST parser. If codegen and schema builders diverge, migrations become incorrect. Any new schema option must be represented in:

- AST parser extraction.
- migration SQL mapping (`forge/db/migrate/sql/*`).
- state conversion (`forge/db/migrate/state/*`) if state needs to preserve the option.

## 7. ORM and QuerySet

The ORM core lives in `forge/orm/`.

### 7.1 Core types

- `BaseQuerySet[T]` in `forge/orm/queryset.go` is the chainable query builder.
- Field expressions and query expressions live across `forge/orm/field_expr.go`, `query_expr.go`, and `expression.go`.
- Managers provide CRUD entry points and coordinate hook invocation (`forge/orm/manager.go`, plus generated managers).

### 7.2 QuerySet surface

The QuerySet interface supports:

- Filtering and exclusion: `Filter(...)`, `Exclude(...)`.
- Ordering: `OrderBy(...)`.
- Pagination: `Limit(...)`, `Offset(...)`.
- Projections: `Select(...)`, `Only(...)`, `Defer(...)`.
- Relation controls: `SelectRelated(...)`, `PrefetchRelated(...)`.
- Value projections: `Values(...)`, `ValuesList(...)`.

See:

- `forge/orm/queryset.go` for the authoritative interface and implementations.
- `forge/orm/queryset_test.go` for behavioral expectations (SelectRelated/PrefetchRelated/Values/ValuesList).

### 7.3 Field expressions

Generated FieldExprs provide typed conditions. They should map to:

- equality: `Eq`, `Ne`
- comparisons: `Gt`, `Gte`, `Lt`, `Lte`
- membership: `In`, `NotIn`
- null checks: `IsNull`, `IsNotNull`
- strings: `Contains`, `StartsWith`, `EndsWith`, case-insensitive variants

The archived API reference contained many examples; the up-to-date implementations live in `forge/orm/field_expr.go` and related helpers.

### 7.4 Runtime field references

When a field is determined at runtime, use the runtime field reference helpers documented in `docs/archive/MIGRATION_V1_TO_V2.md` (historical) and implemented in `forge/orm`.

Patterns:

- SQL-like: `Where("field", OpGreater, 18)`
- Django-like: `F("field").Gt(18)`

### 7.5 Relations and N+1 protection

`forge/orm/preload.go` defines errors and guardrails around accessing relations without preloading. This supports the design goal of preventing accidental N+1 query patterns.

### 7.6 Transactions

Transactions are provided by `forge/db/transaction.go` and used by managers/services. Long-running workflows should use transaction boundaries that match business logic rather than per-request implicit transactions.

### 7.7 Performance considerations

- QuerySet operations are lazy and only execute on terminal operations (`All`, `Get`, `Count`, etc.).
- Filters and projections should be pushed into SQL rather than applied in Go.
- Avoid deep dynamic filters without whitelisting (see filtering system security below).

## 8. Filtering system

The filtering system lives in `forge/filter/` and provides a shared filter AST and execution pipeline used by:

- Admin list pages
- REST API list endpoints
- Direct ORM usage

This subsystem was heavily described in the archive (`docs/archive/FILTERING_SYSTEM.md`). The current implementation includes:

### 8.1 Filter AST

- A serializable representation of filter trees (AND/OR/NOT) suitable for persistence and transport.
- Parsing helpers for query params and UI payloads.

See: `forge/filter/ast.go`, `forge/filter/parser.go`.

### 8.2 FilterSet

- `FilterSet[T]` binds a filter schema to a queryset and request context.
- Provides both an imperative builder API and declarative filters under `forge/filter/filters/*`.

See: `forge/filter/filterset.go`, `forge/filter/filters/*`.

### 8.3 Deep relation lookups

The system supports deep field paths like `author__company__country`. Implementations depend on the ORM schema registry and safe path resolution.

See: `forge/filter/relations.go`, `forge/filter/security.go`.

### 8.4 Security hardening

Filtering is a potential attack surface. The filter system defends with:

- field whitelisting / schema-based validation
- lookup whitelisting
- cost-based throttling and query planning

See: `forge/filter/security.go`, `forge/filter/optimizer.go`, `forge/filter/metrics.go`.

### 8.5 UI and tooling

Widgets for admin and debug tools for SQL previews are in `forge/filter/widgets/*`.

## 9. Migration system

Migrations are implemented in `forge/db/migrate/*` and orchestrated via `forge/db/migrations.go` and CLI commands.

### 9.1 Goals

- Generate deterministic up/down SQL migrations from schema diffs.
- Load previous state from existing migration files (not from database introspection) to support incremental generation in CI.
- Execute migrations via golang-migrate with recovery tools for dirty state.

### 9.2 Components

- Diff detection: `forge/db/migrate/generate/detector.go`
- State tracking: `forge/db/migrate/state/*`
  - File state loader parses existing `*.up.sql` migrations and reconstructs an in-memory schema model.
  - Converter turns schema state back into model definitions for diffing.
- SQL generation: `forge/db/migrate/sql/*`
  - Driver-specific DDL mapping, column type mapping, defaults, generated columns.
- Execution and verification: `forge/db/migrate/execute/*`, `forge/db/migrate/verify/*`
  - Status reporting reads `schema_migrations`.
  - Recovery handles dirty states.
  - Verification contains checksum/drift/safety checks.

### 9.3 Bookkeeping table

First migrations include a bookkeeping table (`schema_migrations`) so execution/status tools can determine applied versions. This is generated by `forge/db/migrate/generate/generator.go`.

Design guardrail: migration diffs must ignore internal bookkeeping tables when comparing schema state to current models.

### 9.4 Generated columns and DB-specific options

Field options like `db_type`, `db_default`, and generated columns are mapped in `forge/db/migrate/sql/base.go` and `forge/db/migrate/sql/common.go`.

### 9.5 Testing

The migration system is heavily tested under `tests/integration/migrate/*`. In particular:

- Incremental ecommerce migrations validate state loader correctness across multiple phases.
- Advanced fields tests validate generated/default columns.
- Recovery/status tests validate dirty state handling and status reporting.

## 10. Admin system

The admin system lives in `forge/admin/`. The archive contains extensive admin design notes and comparisons (Django admin comparison, system comparison, widget API specs). The current implementation contains:

### 10.1 Admin entry points

- `forge/admin/admin.go`: public entry points and wiring.
- `forge/admin/site.go`: Admin site abstraction, route binding, and high-level lifecycle.

### 10.2 Core registry and metadata

- Registry: `forge/admin/core/registry.go` stores model registrations and configuration.
- Metadata: `forge/admin/core/metadata.go` and `metadata_builder.go` build field lists, labels, widget hints, and permissions.

This registry is the intended override point: apps can replace or extend how models are presented without editing generated code.

### 10.3 API-first admin endpoints

`forge/admin/api/rest/router.go` implements a REST layer for admin UI interactions. This aligns with the archived redesign direction (headless admin) but is implemented inside the Go package rather than requiring a separate server.

### 10.4 UI layer

`forge/admin/ui/*` contains UI helpers, templates or template hooks, and widget definitions. The admin UI can be served as static assets or via server-rendered templates depending on configuration.

### 10.5 History, notifications, and plugins

- History/audit: `forge/admin/core/history.go`
- Notifications: `forge/admin/core/notifications.go`
- Plugins: `forge/admin/core/plugin.go`

### 10.6 Testing and examples

The ecommerce example (`examples/ecommerce/app/*/admin.go`) exercises how admin interacts with generated field instances and typed schemas.

## 11. API framework

The API framework lives in `forge/api/` and was originally described in archived docs as a DRF-like stack. The current implementation includes:

### 11.1 ViewSets

- Base and enhanced viewsets (`forge/api/viewset.go`, `viewset_enhanced*.go`) implement CRUD flows.
- Integration helpers (`forge/api/integration.go`) connect serializers, permissions, throttles, renderers, and parsers.

### 11.2 Serializers

- Core serializer interfaces in `forge/api/serializer.go`.
- Field-level serializers under `forge/api/serializers/fields/*`.
- Typed/enhanced serializers for structured output.

### 11.3 Authentication and permissions

- Auth backends under `forge/api/authentication/*` (API key, basic, JWT, session, token).
- Permission classes under `forge/api/permissions/*`.

Identity (forge/identity) provides the persistent user/session/token model layer; API auth focuses on request authentication and context population.

### 11.4 Throttling

- `forge/api/throttling/*` provides anon/user rate throttles, throttle parsing, and errors.

### 11.5 Content negotiation

- Parsers: `forge/api/parsers/*`
- Renderers: `forge/api/renderers/*`

### 11.6 Errors and exceptions

- Exceptions: `forge/api/exceptions/*`
- Problem+code mapping: `forge/api/errors/*`

### 11.7 OpenAPI

- `forge/api/docs/openapi.go` provides an OpenAPI 3.0 generator and HTTP handler.

### 11.8 Tests

The API framework is well-covered by unit tests in `forge/api/*_test.go`, including permissions, throttling, parsers/renderers, and viewset behavior.

## 12. Identity system

Identity is implemented under `forge/identity/` and follows the architecture described in `docs/archive/USER_SYSTEM_ARCHITECTURE.md` and `docs/archive/docs-new/features/identity-system.md`, updated to match the current package layout.

### 12.1 Layers

- Models: `forge/identity/models/*` (user, session, permission, group, tokens).
- Repository: `forge/identity/repository/*` (data access boundaries).
- Services: `forge/identity/service/*` (business logic, validations).
- Backends: `forge/identity/backends/*` (authentication strategy registry: password, token, etc.).
- Serializers: `forge/identity/serializers/*` (API-facing DTO conversion).
- Handlers: `forge/identity/handlers/*` (HTTP endpoints).
- Middleware: `forge/identity/middleware/*` (request authentication).

### 12.2 Extension points

- Add a new authentication mechanism by implementing the backend interface and registering it in `forge/identity/backends/registry.go`.
- Override routes by composing your own router that wraps `forge/identity/router.go`.

### 12.3 Testing

Identity behavior is covered by unit tests under `forge/identity/*/*_test.go`.

## 13. Server, middleware, security

`forge/server/` is responsible for HTTP integration:

- Router wiring and request context helpers.
- Security middleware in `forge/server/security.go`.
- Rate limiting and health endpoints.

When debugging 404s, static asset resolution, or middleware order, this is the first layer to inspect.

## 14. Config, logging, validation

- Config: `forge/config/*` (viper-based settings and environment support).
- Logging: `forge/log/*` (zap-based structured logging, exporters, hooks, middleware).
- Validation: `forge/validate/*` (schema validation tags and validators).

These packages are intentionally separate to keep core ORM/admin/API logic independent of specific logging or configuration frameworks.

## 15. Extension points

Forge is designed to be extended without editing core code. Common extension points:

- Schema: add new field/relation factories in `forge/schema/registry.go`.
- Codegen: extend templates and AST extraction (`forge/codegen/templates`, `forge/codegen/ast_parser.go`).
- ORM: add new expressions/operators in `forge/orm/expression.go`.
- Filter: register custom filters and widgets in `forge/filter/custom.go` and `forge/filter/widgets/*`.
- Admin: register plugins and override metadata/handlers via `forge/admin/core/plugin.go` and registry hooks.
- API: add auth, permissions, throttles, renderers, parsers using the existing interfaces.
- Identity: add backends and middleware.

## 16. Testing architecture

Tests are split across:

- `forge/` module tests: unit tests for API, filter, orm, identity.
- `examples/ecommerce`: example module tests.
- `tests/`: integration tests that exercise real CLI and migration flows.

If you change schema parsing or migration SQL generation, `tests/integration/migrate/*` is the critical suite.

## 17. Alignment with archived docs

The archive contains more detail than we keep in this canonical architecture file. The rule is:

- If a claim is implemented and test-covered, it belongs here.
- If a claim is conceptual or planned, it belongs in `docs/ROADMAP.md`.
- If a claim is historical context, comparisons, or a deep dive, it stays in `docs/archive/`.

Recommended archived references by topic:

- Deep architecture narrative: `docs/archive/FRAMEWORK_ARCHITECTURE.md`, `docs/archive/ARCHITECTURE.md`.
- Schema reference: `docs/archive/SCHEMA_REFERENCE.md`.
- API usage examples: `docs/archive/API_REFERENCE.md`, `docs/archive/REST_API.md`.
- Filtering deep dive: `docs/archive/FILTERING_SYSTEM.md`.
- Identity patterns: `docs/archive/USER_SYSTEM_ARCHITECTURE.md`.
- Admin deep dives: `docs/archive/ADMIN_*`.

