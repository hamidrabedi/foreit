# forge Roadmap (comprehensive, code-aligned)

## Table of contents

- [1. Summary](#1-summary)
- [2. Current state (what is actually implemented)](#2-current-state-what-is-actually-implemented)
- [3. Guiding priorities](#3-guiding-priorities)
- [4. Near-term epics (0-3 months)](#4-near-term-epics-0-3-months)
- [5. Mid-term epics (3-6 months)](#5-mid-term-epics-3-6-months)
- [6. Longer-term epics (6-12 months)](#6-longer-term-epics-6-12-months)
- [7. Long-term bets (12+ months)](#7-long-term-bets-12-months)
- [8. Ongoing quality work](#8-ongoing-quality-work)
- [9. Release process and definitions of done](#9-release-process-and-definitions-of-done)
- [10. Cross-references to archived plans](#10-cross-references-to-archived-plans)

## 1. Summary

The archive contains multiple roadmaps and TODO lists (see `docs/archive/ROADMAP.md` and `docs/archive/docs-new/TODOS.md`). This roadmap reconciles them with the real implementation in `forge/` today.

Key themes:

- The "core stack" is present: schema, codegen, ORM, filters, migrations, admin, API, identity, CLI.
- The next work should focus on correctness, stability, and the developer experience of overrides and assets.
- Advanced ORM and advanced admin UI are the major growth areas.

## 2. Current state (what is actually implemented)

### 2.1 Core framework

Implemented:

- Schema DSL: `forge/schema/*`.
- AST-based code generation: `forge/codegen/*`.
- ORM core with QuerySet: `forge/orm/*`.
- Filter system: `forge/filter/*`.
- Migrations: `forge/db/migrate/*` + `forge/db/migrations.go`.
- CLI: `forge/cli/*`.

### 2.2 Admin

Implemented:

- Admin registry and site: `forge/admin/core/*`, `forge/admin/site.go`.
- Admin REST router: `forge/admin/api/rest/router.go`.
- UI scaffolding and widgets: `forge/admin/ui/*`.

### 2.3 API

Implemented:

- Serializers, viewsets, permissions, throttling, renderers/parsers.
- OpenAPI generator: `forge/api/docs/openapi.go`.

### 2.4 Identity

Implemented:

- Models, repos, services, backends, handlers, middleware.

### 2.5 Tests

Implemented:

- Migration integration tests: `tests/integration/migrate/*`.
- CLI e2e tests: `tests/cmd_forge/*`, `tests/e2e/cli/*`.
- ORM/schema integration tests: `tests/integration/schema`, `tests/integration/orm`.

## 3. Guiding priorities

1. Correctness and deterministic behavior (migrations, codegen, admin/API).
2. Security by default (filter whitelists, identity, middleware ordering).
3. Override-first extensibility (registry/plug-ins, fewer hard-coded calls).
4. Great DX: CLI workflows and examples stay in lockstep with code.

## 4. Near-term epics (0-3 months)

These are the highest ROI work items that improve reliability and unblock teams.

### Epic 4.1 Admin asset pipeline and static serving

Problem:

- Admin UI can produce hashed assets (e.g., Vite `index-<hash>.js`). 404s happen when server does not serve the correct dist directory or when build artifacts are missing.

Deliverables:

1. Define a single authoritative admin asset directory.
2. Ensure `forge runserver` serves it under `/admin/assets/`.
3. Add a build step (or auto-build in dev) that keeps hashed bundles in sync.
4. Add a regression test that requests `/admin/` and all referenced assets.

Implementation areas:

- `forge/admin/ui/*` (asset packaging)
- `forge/server/static.go` (static serving)
- `forge/cli/commands/server/runserver.go` (dev server behavior)
- tests (new integration test under `tests/` or `forge/admin/ui`)

Definition of done:

- No 404 on referenced admin assets in local dev.
- Tests cover at least the admin index page + asset fetch.

### Epic 4.2 Make overrides first-class (reduce main.go wiring)

Archived docs repeatedly emphasize: built-in systems should not require manual registration in `main.go`, but should remain overridable.

Deliverables:

1. Document and implement a default auto-registration mechanism.
2. Provide opt-out flags.
3. Provide clear extension APIs:
   - identity routes
   - admin routes
   - API routes

Implementation areas:

- `forge/registry/*`
- `forge/admin/core/registry.go`
- `forge/identity/router.go`
- `forge/server/router.go`

Definition of done:

- A sample app can start with near-zero wiring.
- Overrides can be applied by plugging into registries.

### Epic 4.3 Migrations safety and ergonomics

Deliverables:

1. Improve drift/lint output readability (`forge/db/migrate/verify/*`).
2. Improve status output (`forge/db/migrate/execute/status.go`).
3. Provide explicit migration dependency metadata.

Definition of done:

- Incremental ecommerce migrations remain stable.
- Developers can run `forge migrate status` and understand what is pending.

### Epic 4.4 CLI docs and UX

Deliverables:

1. Ensure CLI commands have consistent flags and help.
2. Add missing parity commands from archived roadmap that are now partially present:
   - startapp (alias for add app)
   - dbshell
   - collectstatic (if relevant)

Definition of done:

- CLI help is consistent.
- CLI e2e tests cover the core commands.

## 5. Mid-term epics (3-6 months)

### Epic 5.1 Advanced ORM features (close the gap with archived API reference)

The archive lists SelectRelated, PrefetchRelated, aggregates, annotations, values, and F expressions as major features. The codebase already contains scaffolding and partial implementations.

Deliverables:

1. Production-ready eager loading:
   - Implement join generation for `SelectRelated`.
   - Implement prefetch queries for `PrefetchRelated`.
   - Ensure safety errors in `forge/orm/preload.go` are only triggered when appropriate.
2. Aggregates and annotations:
   - Implement aggregate SQL in `forge/orm/aggregates.go`.
   - Implement annotation expressions in `forge/orm/annotations.go`.
3. Values/ValuesList:
   - Ensure projection SQL is correct.
   - Ensure type-safe variants (field expressions) work alongside string-based forms.

Tests:

- Expand `forge/orm/queryset_test.go`.
- Add integration tests under `tests/integration/orm` verifying SQL results.

Definition of done:

- ORM advanced methods work in real Postgres.
- No silent N+1 patterns.

### Epic 5.2 Admin UX (headless-first)

Archived redesign documents push a headless admin with a modern frontend.

Deliverables:

1. Solidify admin REST endpoints and models metadata format.
2. Provide stable TypeScript types generation (if desired) using OpenAPI or direct codegen.
3. Improve widget system:
   - consistent schema -> widget mapping
   - autocomplete endpoint
   - global search

Definition of done:

- Admin UI can list, edit, and filter ecommerce models reliably.

### Epic 5.3 REST auto-generation

Deliverables:

1. Auto-generate serializers from schema definitions.
2. Auto-generate viewsets from schema + serializer.
3. Auto-register routes with a registry-based approach.
4. Expand OpenAPI generator to include generated viewsets.

Risks:

- Over-generation can create APIs without proper security. Ensure permissions are explicit.

Definition of done:

- A project can expose CRUD REST APIs for selected models with minimal boilerplate.

## 6. Longer-term epics (6-12 months)

### Epic 6.1 Observability and operational maturity

Deliverables:

- Structured logs across all HTTP surfaces.
- Metrics hooks for request latency, errors, throttling events.
- Health/readiness endpoints.
- Tracing support (OpenTelemetry) as optional integration.

### Epic 6.2 Background tasks

Deliverables:

- Task runner + worker process.
- Admin actions can enqueue tasks.
- Retry policy and dead-letter behavior.

### Epic 6.3 Advanced migrations

Deliverables:

- Drift detection with DB introspection (optional).
- Schema checksums and verification.
- Safe rollback policies.

## 7. Long-term bets (12+ months)

### Epic 7.1 Multi-tenancy

Deliverables:

- tenant-aware query building
- tenant-aware migrations
- admin for tenant management

### Epic 7.2 GraphQL and realtime

Deliverables:

- GraphQL schema generation for models.
- WebSocket updates for admin dashboards.

### Epic 7.3 Plugin ecosystem

Deliverables:

- Stable plugin API surface.
- A curated plugin set for common needs (audit log, import/export, file manager).

## 8. Ongoing quality work

These items are continuous and should not be treated as one-time epics.

### 8.1 Code quality and cleanup

From `docs/archive/docs-new/TODOS.md`, the project tracks cleanup as "complete" but ongoing quality work remains:

- keep naming consistent
- keep error mapping consistent
- keep lint/staticcheck green

### 8.2 Examples consistency

Examples are part of the product contract.

- When schema/codegen changes, update `examples/ecommerce`.
- Keep example admin code using generated field instances (this has been a source of breakage historically).

### 8.3 Documentation

- Keep authoritative docs in `docs/` up to date.
- Keep archive for context.

## 9. Release process and definitions of done

### 9.1 Definitions of done

A feature is "done" when:

1. Code is implemented behind the intended extension points.
2. Unit tests cover the feature where appropriate.
3. Integration tests cover it if it touches migrations/CLI/admin.
4. Examples compile.
5. Docs are updated:
   - architecture/design if structural
   - PRD if product requirement
   - roadmap if status changed

### 9.2 Release checklist

- Run tests:
  - `go test ./forge/...`
  - `go test ./examples/ecommerce/...`
  - `go test ./tests/...`
- Run linters/checks:
  - `forge check` (or equivalent)
- Generate code and verify nothing unexpected changes:
  - `forge generate`
- Verify migrations:
  - `forge makemigrations <name> --auto`
  - review SQL
  - `forge migrate up` on test DB
- Verify admin assets:
  - build admin UI assets
  - run server and ensure no 404 on referenced bundles

## 10. Cross-references to archived plans

This roadmap is rooted in, but not identical to, archived planning.

- Archived long roadmap: `docs/archive/ROADMAP.md`.
- Archived TODO tracker: `docs/archive/docs-new/TODOS.md`.
- Archived design deep dives:
  - `docs/archive/FRAMEWORK_ARCHITECTURE.md`
  - `docs/archive/ADMIN_REDESIGN_ARCHITECTURE.md`
  - `docs/archive/USER_SYSTEM_ARCHITECTURE.md`

The rule is: archive documents are references; this file is the current plan.

## Appendix A: Backlog inventory (expanded)

This appendix enumerates backlog items from archived roadmaps and TODO lists, grouped by subsystem. It is intentionally exhaustive so future planning can prune rather than recreate.

### A.1 ORM backlog

- JOIN-based eager loading for `SelectRelated` (SQL builder + schema registry integration).
- Prefetch queries for `PrefetchRelated` (batch load related sets).
- Relation path resolution and validation.
- Aggregates: Count, Sum, Avg, Min, Max.
- Aggregation grouping.
- Annotations: computed expressions.
- Window functions.
- Subqueries.
- Raw SQL escape hatches with safe parameter binding.
- BulkCreate and BulkUpdate improvements.
- Values/ValuesList typed variants (field expressions strongly typed end-to-end).
- ValuesList flat mode improvements.
- Query plan introspection tooling.
- Query caching hooks.

### A.2 Filter system backlog

- Persisted filter storage implementation (DB-backed).
- Filter sharing/versioning.
- More widgets for numeric, date ranges.
- Better optimizer heuristics for EXISTS vs JOIN.
- Cost model calibration.
- Filter metadata schema improvements.
- Better error messages for invalid deep paths.

### A.3 Migration system backlog

- DB introspection drift detection (optional).
- Constraint diffing improvements (deferrable, match).
- Index diffing improvements (partial indexes).
- Safer down migrations.
- Squash migrations stability.
- Checksum storage and verification.
- Migration graph visualization.
- Migration lint rules for destructive changes.

### A.4 Admin backlog

- Robust SPA admin package integration (if SPA is chosen as default).
- Better list view performance.
- Autocomplete and global search endpoints.
- Rich text editor.
- File/image upload widget.
- Audit history UI.
- Model relationship navigation.
- Inline editing.
- Fieldsets and layout configuration.
- Theme system.
- Dashboard widgets library.

### A.5 API framework backlog

- Auto-generate serializers and viewsets from schema.
- Route auto-registration.
- API versioning strategy.
- Interactive API docs (Swagger UI) for OpenAPI generator.
- Better content negotiation customization.
- Consistent pagination metadata schema.
- Rate limiting policies and storage backends.
- CORS configuration improvements.
- API client generation.

### A.6 Identity backlog

- OAuth backend implementation.
- MFA support.
- Account lockout policies.
- Email provider abstraction.
- RBAC admin UI.
- Permissions caching.

### A.7 CLI/dev tooling backlog

- startapp alias and improved scaffolds.
- dbshell command.
- collectstatic command (if needed for asset packaging).
- createsuperuser improvements.
- hot reload.
- debug toolbar.
- improved test runner output.

### A.8 Observability backlog

- Prometheus metrics.
- request tracing.
- structured audit logs.
- log aggregation guidance.

### A.9 Product-level backlog

- Multi-tenancy.
- i18n.
- GraphQL.
- WebSocket support.
- Plugin marketplace.

## Appendix B: Timeline template (fill per release)

Use this template to plan releases without rewriting structure:

- Release name:
- Target date:
- Epics:
- Risks:
- Test gates:
- Docs updates required:
- Examples updates required:

## Appendix C: Milestones and checklists

### C.1 Migration safety milestone

- [ ] Improve drift detection output
- [ ] Improve safety lint rules
- [ ] Add destructive-change warnings
- [ ] Add schema checksum verification
- [ ] Add DB introspection mode (optional)
- [ ] Add migration plan printing
- [ ] Add migration dependency comments validation

### C.2 Admin asset milestone

- [ ] Define build directory contract
- [ ] Ensure `runserver` serves build directory
- [ ] Ensure `forge` CLI can build assets in dev
- [ ] Add test to fetch `/admin/` and verify bundles exist
- [ ] Add docs describing asset pipeline

### C.3 Advanced ORM milestone

- [ ] Implement join builder for select related
- [ ] Implement prefetch strategy
- [ ] Implement relation mapping and hydration
- [ ] Add integration tests for eager loading
- [ ] Expand aggregate SQL
- [ ] Expand annotation expressions

### C.4 API auto-generation milestone

- [ ] Generate serializer skeletons from schema
- [ ] Generate viewset skeletons
- [ ] Register routes via registry
- [ ] Expand OpenAPI to include generated APIs
- [ ] Add a tutorial and an example

### C.5 Identity milestone

- [ ] Add OAuth backend interface
- [ ] Add OAuth backend implementation
- [ ] Add MFA primitives
- [ ] Improve permissions caching

### C.6 Observability milestone

- [ ] Add request metrics
- [ ] Add tracing hooks
- [ ] Add admin audit events export

