# foreit — ForgeGo + Gogo framework monorepo

This repository contains two related Go frameworks and examples:

- `forgego/` — application framework, admin UI, examples.
- `gogo/` — modular framework packages and documentation.

Quick commands

- Build everything: `go build ./...`
- Run tests: `go test ./...` (or `cd forgego && go test ./...` for submodule scope)
- Run example app: `go run ./forgego/cmd/cli` or `go run ./forgego/examples/basic`

Where to start

- Application runtime & examples: `forgego/cmd`, `forgego/examples`
- Model & DB layer: `forgego/pkg/models`
- Admin UI & routing: `forgego/pkg/admin`
- Framework modules & docs: `gogo/docs`, `gogo/pkg`

For AI agents and automation

See `AGENTS.md` for recommended API endpoints, authentication scopes, and example `curl`/`gh` commands to fetch repository content programmatically.
