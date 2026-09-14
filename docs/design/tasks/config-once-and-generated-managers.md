Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

# Load config once and pass it down; generated managers fail at init instead of hiding errors

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit: forge/db/migrate/generate/generator.go (+ its callers), forge/cli/core/context.go (+ callers that need a config), forge/codegen/templates/combined.tmpl, forge/orm/manager.go (only to add `MustNewManager`), and tests.

## Principles (mandatory)
- Tests FIRST. gofmt; functions < 40 lines; errors wrapped with %w; no new package-level mutable state; naming after behaviour.
- Keep existing exported constructors working for callers you do not update (add, don't break, unless every caller is in the workspace and updated).

## Facts (verified on master)
- `config.NewConfig()` is re-created deep in the stack: forge/db/migrate/generate/generator.go:35 (`NewMigrationGenerator`) and :51 (`NewMigrationGeneratorWithDefaults`, which reads `cfg.GetDriver()`), forge/cli/core/context.go:21 (`NewContext`). A config loaded by the CLI root (flags, env file) is ignored by these.
- Generated code (forge/codegen/templates/combined.tmpl:29) emits `var {{ .Name }}Objects, _ = orm.NewManager[...]` and discards the error, so a broken model schema produces a nil manager that panics later at the first query, far from the cause.

## Change
1. generator.go: `NewMigrationGenerator` and `NewMigrationGeneratorWithDefaults` take the driver (or a `*config.Config`) as a parameter instead of calling `config.NewConfig()`. Read how they are called (grep the workspace: cli/commands/migrations, tests/) and pass the config the caller already has (the CLI command context has `ctx.Config`). If a caller has no config, it creates one once at its entry point. Keep a thin `...FromEnv` wrapper only if an external test needs it.
2. cli/core/context.go: `NewContext(cfg *config.Config)`; the CLI root (find where `NewContext()` is called) loads the config once and passes it. Update all callers and tests.
3. orm/manager.go: add `func MustNewManager[T any](tableName string) *Manager[T]` that panics with a message naming the model type and the wrapped error. Test: an invalid model type panics with that message; a valid one returns a manager.
4. combined.tmpl: generate `var {{ .Name }}Objects = orm.MustNewManager[{{ .Name }}](...)`. Update any codegen golden/expected-output tests. Regenerate nothing under examples/ by hand; if examples/ecommerce contains generated files with the old pattern, update those lines to the new pattern so the example keeps compiling (it must still build).

## Acceptance (from the module root)
    gofmt -l .                                        # empty
    go build ./... && (cd ../examples/ecommerce && go build ./...)
    go vet ./db/... ./cli/... ./codegen/... ./orm/...
    staticcheck ./db/migrate/... ./cli/... ./codegen/...   # no output
    go test -race ./db/migrate/... ./cli/... ./codegen/... ./orm/... -count=1
    grep -rn 'config.NewConfig()' --include=*.go db cli | grep -v _test   # only the CLI entry point

Report (max 25 lines): signature changes and callers updated, each command result.
