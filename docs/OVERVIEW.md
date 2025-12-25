# Repository overview — ForgeGo + Gogo

This document summarizes the architecture, important modules, and conventions that callers and automation should know.

High-level architecture

- Two related projects live in the same repo:
  - `forgego/`: application framework built around a typed models layer and an admin UI.
  - `gogo/`: modular packages and documentation for reusable framework components.

Key components (forgego)

- `pkg/models` — typed model definitions, `ModelManager[T]`, `QuerySet`.
- `pkg/admin` — registry, viewset/service integrations, UI and API router built on `gofiber/fiber`.
- `cmd/` — CLI entrypoints and code generation helpers.

Important patterns & gotchas

- Registration-based wiring: model and viewset registration happens through `RegisterModel` / `RegisterViewSet` (see `forgego/pkg/admin/register.go`). Avoid editing `router.go` directly; add integrations through the registry.
- Generics + reflection: many runtime handlers create untyped handlers with reflection (see `router.go` generic handlers). When changing typed APIs, ensure any reflection-based call sites still find expected method names.
- Hooks & validation: `ModelManager` runs `ValidateModel` and `RunHooks` on CRUD operations. Use hook names like `BeforeCreate`, `AfterUpdate`.

Build, test, run

- Build: `go build ./...`
- Tests: `go test ./...` or, for a focused run, `cd forgego && go test ./...`
- Example app: `go run ./forgego/cmd/cli` or `go run ./forgego/examples/basic`

Files to inspect for common tasks

- Add a new model: `forgego/pkg/models/*`, then call `RegisterModel` in an integration point.
- Add a new REST resource: use `RegisterViewSet` or a ServiceIntegration; follow `modelNameToResource` for REST naming.
