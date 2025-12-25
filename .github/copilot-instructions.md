# Copilot / AI Agent Instructions for ForgeGo / Gogo

Short, task-focused guidance to help an AI coding agent be productive in this repository.

- **Repo shape:** Two related projects live here: `forgego/` (application framework and examples) and `gogo/` (framework modules and docs). Start by scanning `forgego/cmd` and `gogo/docs` for runnable entry points and examples.

- **Big picture:**
  - `pkg/models` implements a Django-like typed model system (see [forgego/pkg/models/manager.go](forgego/pkg/models/manager.go#L1-L40)). Use `ModelDefinition[T]`, `ModelManager[T]` and `QuerySet` when modeling data.
  - `pkg/admin` provides an admin UI + routing layer built on Fiber; it wires model definitions into UI and API routes via registry/Integration patterns (see [forgego/pkg/admin/register.go](forgego/pkg/admin/register.go#L1-L80) and [forgego/pkg/admin/router.go](forgego/pkg/admin/router.go#L1-L80)). Prefer using `RegisterModel` / `RegisterViewSet` rather than manual route changes.
  - Code uses Go generics and reflection together: typed definitions (`T`) are common but some runtime wiring uses reflection (generic handlers in `router.go`). Be careful mixing compile-time generics changes with runtime reflection assumptions.

- **Key conventions & patterns**
  - Registration-based DI: modules register models and viewsets via `register.go` and `registry.go` patterns. Look for `Register*` helpers instead of editing global lists.
  - Model→resource naming: `modelNameToResource` converts CamelCase model names into REST resource names; use that function as ground truth for route names ([forgego/pkg/admin/router.go](forgego/pkg/admin/router.go#L1-L120)).
  - Read-only fields: serializers mark primary keys and timestamps read-only (see `RegisterModel` logic in [forgego/pkg/admin/register.go](forgego/pkg/admin/register.go#L1-L80)).
  - Hooks & validation: `ModelManager` runs `ValidateModel` and `RunHooks` in CRUD methods — follow the existing hook names (`BeforeCreate`, `AfterCreate`, etc.) when adding behavior ([forgego/pkg/models/manager.go](forgego/pkg/models/manager.go#L1-L140)).

- **Build / run / test workflows**
  - Build everything: `go build ./...` from repository root.
  - Run tests: `go test ./...` (module-aware). If tests target a submodule, run `cd forgego && go test ./...` etc.
  - Run an example app: `go run ./forgego/cmd/cli` or `go run ./forgego/examples/basic` depending on the entrypoint you want. Inspect `forgego/examples/*/README.md` for example-specific env vars.

- **Where to look first when asked to implement a feature**
  - For new models: `forgego/pkg/models/*` + `forgego/pkg/admin/register.go` for registration helpers.
  - For API/UI routes: `forgego/pkg/admin/router.go` — prefer adding ViewSet integrations; avoid directly editing generated route lists.
  - For project-wide docs and examples: `gogo/docs/QUICK_START.md` and `forgego/examples/`.

- **Examples to copy/pattern-match**
  - New model registration: follow the pattern in [forgego/pkg/admin/register.go](forgego/pkg/admin/register.go#L1-L80).
  - CRUD semantics and hooks: mirror `Create/Update/Delete` implementations in [forgego/pkg/models/manager.go](forgego/pkg/models/manager.go#L1-L120).

- **Integration points & external deps**
  - HTTP framework: uses `github.com/gofiber/fiber/v2` for routing and UI handlers.
  - Database layer: a DB abstraction is in `pkg/models` and `pkg/models/connection.go` (use `ModelManager` and `DB` helpers rather than raw SQL when possible).
  - Many modules are structured as reusable packages under `pkg/` and documented in `gogo/docs` — prefer reading module README before changing its public API.

- **AI agent behavior rules (project-specific)**
  - Do not change registry initialization patterns; instead add registration calls where other modules do (search for `RegisterModel` / `RegisterViewSet`).
  - Preserve generics/typed APIs: if adding types, follow the `T any` patterns and add explicit calls to typed helpers.
  - When adding routes, prefer creating a ViewSet or Service integration; avoid hand-editing `router.go` unless implementing new low-level behavior.

If anything here is unclear or you want more examples (e.g., a concrete model add/change PR), tell me which area to expand and I will iterate.
